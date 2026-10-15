<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import { useToast } from '@nuxt/ui/composables'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'
import PlaylistEntryItem from '@/components/molecules/PlaylistEntryItem.vue'
import { downloadAll, downloadEntry, downloadEntryFile, downloadPlaylistZip } from '@/api/downloadService'
import { DownloadStatus } from '@/models/enums'
import type { FilterOptions } from '@/models/filterOptions'
import type { Playlist } from '@/models/playlist'

interface SongListProps {
  playlist: Playlist
  filters?: FilterOptions
  allComplete?: boolean
}

const props = withDefaults(defineProps<SongListProps>(), {
  allComplete: false,
})

const emit = defineEmits<{
  'update:playlist': [playlist: Playlist]
}>()

const toast = useToast()

const downloadingAll = ref(false)
const downloadingEntryIds = ref<number[]>([])
// When true, we are waiting for downloads to finish before fetching the zip.
const downloadAllZipPending = ref(false)

const entries = computed(() => props.playlist.entries ?? [])
const songCount = computed(() => entries.value.length)

function isEntryDownloading(entryId: number): boolean {
  return downloadingEntryIds.value.includes(entryId)
}

function isEntryCompleted(entryId: number): boolean {
  return entries.value.some(
    (entry) => entry.id === entryId && entry.status === DownloadStatus.Completed,
  )
}

// "Download all": download nothing (already cached) or kick off background
// downloads, then (once all entries complete) fetch the zip to the browser.
async function handleDownloadAll() {
  if (downloadingAll.value) return

  downloadingAll.value = true
  try {
    const updated = await downloadAll(props.playlist.id, props.filters)
    emit('update:playlist', { ...props.playlist, ...updated })

    // If everything is already downloaded (all cached), auto-zip immediately.
    // Otherwise set a pending flag and let the watcher on allComplete fire.
    downloadAllZipPending.value = true
    if (props.allComplete) {
      await fetchZip()
    }
  } catch (error) {
    console.error('Failed to start playlist download:', error)
    toast.add({
      title: 'Download failed',
      description: 'Something went wrong while starting the playlist download.',
      color: 'error',
    })
  } finally {
    downloadingAll.value = false
  }
}

async function handleDownloadEntry(entryId: number) {
  if (isEntryDownloading(entryId)) return

  // Reuse the per-row button: for an already cached file, save it to the
  // browser; otherwise trigger the background download.
  if (isEntryCompleted(entryId)) {
    await saveEntryFile(entryId)
    return
  }

  // Optimistically flip the status to "downloading" so the row updates
  // immediately; the SSE stream then brings it to completed/failed.
  emit('update:playlist', {
    ...props.playlist,
    entries: entries.value.map((e) =>
      e.id === entryId ? { ...e, status: DownloadStatus.Downloading, file_path: '' } : e,
    ),
  })
  downloadingEntryIds.value = [...downloadingEntryIds.value, entryId]
  try {
    // The 202 response reflects the pre-download state, so don't merge it back
    // (it would downgrade the optimistic status). Progress comes via SSE.
    await downloadEntry(entryId, props.filters)
  } catch (error) {
    console.error('Failed to start entry download:', error)
    toast.add({
      title: 'Download failed',
      description: 'Something went wrong while starting this song download.',
      color: 'error',
    })
  } finally {
    downloadingEntryIds.value = downloadingEntryIds.value.filter((id) => id !== entryId)
  }
}

async function saveEntryFile(entryId: number) {
  const entry = entries.value.find((e) => e.id === entryId)
  if (!entry) return

  const ok = await downloadEntryFile(entry)
  if (!ok) {
    toast.add({
      title: 'File unavailable',
      description: 'This file is no longer available on the server. Try downloading again.',
      color: 'warning',
    })
  }
}

async function fetchZip() {
  const ok = await downloadPlaylistZip(
    props.playlist.id,
    `${props.playlist.title || 'playlist'}.zip`,
  )
  downloadAllZipPending.value = false
  if (!ok) {
    toast.add({
      title: 'No files to zip',
      description: 'There are no downloaded files available to bundle yet.',
      color: 'warning',
    })
  }
}

// When all downloads complete after a "Download all" click, fetch the zip.
watch(
  () => props.allComplete,
  async (complete) => {
    if (complete && downloadAllZipPending.value) {
      await fetchZip()
    }
  },
)
</script>

<template>
  <div>
    <div v-if="songCount > 0" class="flex items-center justify-end px-4 py-3 sm:px-6">
      <DownloadButton
        label="Download all"
        :loading="downloadingAll || downloadAllZipPending"
        @click="handleDownloadAll"
      />
    </div>

    <ul v-if="songCount > 0" class="divide-y divide-default overflow-y-auto max-h-[60vh]">
      <PlaylistEntryItem
        v-for="entry in entries"
        :key="entry.id"
        :entry="entry"
        :downloading="isEntryDownloading(entry.id)"
        :on-download="() => handleDownloadEntry(entry.id)"
      />
    </ul>

    <div v-else class="flex flex-col items-center gap-2 py-8 text-center">
      <UIcon name="i-lucide-music" class="text-dimmed size-8" />
      <p class="text-muted text-sm">No songs were found for this playlist.</p>
    </div>
  </div>
</template>
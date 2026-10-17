<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { refDebounced } from '@vueuse/core'
import UButton from '@nuxt/ui/components/Button.vue'
import UProgress from '@nuxt/ui/components/Progress.vue'
import { useToast } from '@nuxt/ui/composables'
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'
import CancelButton from '@/components/molecules/buttons/CancelButton.vue'
import PlaylistEntryItem from '@/components/molecules/PlaylistEntryItem.vue'
import { apiErrorStatus, getApiErrorMessage } from '@/api/apiErrors'
import {
  FileRequestError,
  cancelEntryDownload,
  cancelPlaylistDownload,
  downloadAll,
  downloadEntry,
  playlistZipUrl,
  saveEntryFile,
  saveFromUrl,
} from '@/api/downloadService'
import type { Playlist } from '@/models/playlist'
import type { PlaylistEntry } from '@/models/playlistEntry'
import type { ButtonColor, ButtonVariant } from '@/types/buttonTypes'
import type { EntryAction } from '@/types/entryAction'
import type { EntryProgress } from '@/types/progress'
import { playlistStats } from '@/utils/playlistStats'

interface SongListProps {
  playlist: Playlist
  progress?: Record<number, EntryProgress>
}

const props = withDefaults(defineProps<SongListProps>(), {
  progress: () => ({}),
})

const ACTION_TOAST_DURATION_MS = 10_000

const toast = useToast()

const startingRun = ref(false)
const cancellingRun = ref(false)
const busyEntryIds = ref<number[]>([])

const entries = computed(() => props.playlist.entries ?? [])
const stats = computed(() => playlistStats(entries.value))
const zipUrl = computed(() => playlistZipUrl(props.playlist.id))

function songs(count: number): string {
  return `${count} ${count === 1 ? 'song' : 'songs'}`
}

function missedParts(): string[] {
  const { failed, unavailable } = stats.value
  return [
    failed > 0 ? `${failed} failed` : '',
    unavailable > 0 ? `${unavailable} unavailable` : '',
  ].filter(Boolean)
}

const progressSummary = computed(() => {
  const { completed, total } = stats.value
  return [`${completed} of ${total} downloaded`, ...missedParts()].join(' · ')
})

const overallPercent = computed(() => {
  const { completed, total } = stats.value
  if (total === 0) return 0
  const partial = Object.values(props.progress).reduce(
    (sum, progress) => sum + (progress.percent ?? 0) / 100,
    0,
  )
  return Math.min(100, ((completed + partial) / total) * 100)
})

const liveText = computed(() =>
  stats.value.active > 0 ? progressSummary.value.split(' · ').join(', ') : '',
)
const announcement = refDebounced(liveText, 1000)

interface PrimaryAction {
  label: string
  icon: string
  color: ButtonColor
  variant: ButtonVariant
  disabled: boolean
}

const primaryAction = computed<PrimaryAction>(() => {
  const { completed, failed, unavailable, fresh, remaining } = stats.value
  if (remaining === 0) {
    return {
      label: unavailable > 0 ? 'Nothing left to download' : 'All downloaded',
      icon: 'i-lucide-circle-check',
      color: 'neutral',
      variant: 'soft',
      disabled: true,
    }
  }
  if (fresh === 0) {
    return {
      label: `Retry ${failed} failed`,
      icon: 'i-lucide-rotate-ccw',
      color: 'warning',
      variant: 'solid',
      disabled: false,
    }
  }
  return {
    label:
      completed + failed + unavailable === 0 ? 'Download all' : `Download ${remaining} remaining`,
    icon: 'i-lucide-download',
    color: 'primary',
    variant: 'solid',
    disabled: false,
  }
})

const zipCount = computed(() => {
  const { completed, total } = stats.value
  return completed < total ? `${completed} of ${total}` : ''
})
const zipLabel = computed(() =>
  zipCount.value ? `Save as zip (${zipCount.value})` : 'Save as zip',
)

function isConflict(error: unknown): boolean {
  return apiErrorStatus(error) === 409
}

async function handleDownloadAll() {
  if (startingRun.value) return

  startingRun.value = true
  try {
    await downloadAll(props.playlist.id)
  } catch (error) {
    if (isConflict(error)) {
      toast.add({
        title: 'Already downloading',
        description: 'This playlist is already downloading.',
        color: 'info',
      })
      return
    }
    console.error('Failed to start playlist download:', error)
    toast.add({
      title: 'Download failed',
      description: getApiErrorMessage(error, "Couldn't start the download."),
      color: 'error',
      duration: 0,
    })
  } finally {
    startingRun.value = false
  }
}

let skipNextRunEndToast = false

async function handleCancelAll() {
  if (cancellingRun.value) return

  cancellingRun.value = true
  skipNextRunEndToast = true
  try {
    await cancelPlaylistDownload(props.playlist.id)
    const { completed } = stats.value
    toast.add({
      title: completed > 0 ? `Stopped. ${songs(completed)} kept.` : 'Download stopped.',
      color: 'neutral',
      icon: 'i-lucide-circle-stop',
    })
  } catch (error) {
    skipNextRunEndToast = false
    if (isConflict(error)) return
    console.error('Failed to cancel playlist download:', error)
    toast.add({
      title: "Couldn't stop the download",
      description: getApiErrorMessage(error, 'Something went wrong while cancelling. Try again.'),
      color: 'error',
      duration: 0,
    })
  } finally {
    cancellingRun.value = false
  }
}

function isEntryBusy(entryId: number): boolean {
  return busyEntryIds.value.includes(entryId)
}

async function handleEntryAction(entry: PlaylistEntry, action: EntryAction) {
  if (isEntryBusy(entry.id)) return

  busyEntryIds.value = [...busyEntryIds.value, entry.id]
  try {
    if (action === 'save') await saveEntry(entry)
    else if (action === 'cancel') await cancelEntry(entry)
    else await startEntry(entry)
  } finally {
    busyEntryIds.value = busyEntryIds.value.filter((id) => id !== entry.id)
  }
}

async function startEntry(entry: PlaylistEntry) {
  try {
    await downloadEntry(entry.id)
  } catch (error) {
    if (isConflict(error)) {
      toast.add({
        title: 'Already downloading',
        description: `${entry.title} is already downloading.`,
        color: 'info',
      })
      return
    }
    console.error('Failed to start entry download:', error)
    toast.add({
      title: 'Download failed',
      description: getApiErrorMessage(error, "Couldn't start this song download."),
      color: 'error',
      duration: 0,
    })
  }
}

async function cancelEntry(entry: PlaylistEntry) {
  try {
    await cancelEntryDownload(entry.id)
  } catch (error) {
    if (isConflict(error)) return
    console.error('Failed to cancel entry download:', error)
    toast.add({
      title: "Couldn't cancel",
      description: getApiErrorMessage(
        error,
        `Something went wrong while cancelling ${entry.title}.`,
      ),
      color: 'error',
      duration: 0,
    })
  }
}

async function saveEntry(entry: PlaylistEntry) {
  try {
    await saveEntryFile(entry)
  } catch (error) {
    if (error instanceof FileRequestError && (error.status === 404 || error.status === 410)) {
      toast.add({
        title: 'File unavailable',
        description: `${entry.title} was removed from the server.`,
        color: 'warning',
        duration: ACTION_TOAST_DURATION_MS,
        actions: [
          {
            label: 'Download again',
            icon: 'i-lucide-download',
            color: 'neutral',
            variant: 'outline',
            onClick: () => handleEntryAction(entry, 'download'),
          },
        ],
      })
      return
    }
    console.error('Failed to save entry file:', error)
    toast.add({
      title: "Couldn't save the file",
      description:
        error instanceof FileRequestError
          ? 'The server could not send this file. Try again.'
          : 'Check your connection and try again.',
      color: 'error',
      duration: 0,
    })
  }
}

let runPeak = 0

watch(
  () => stats.value.active,
  (active, previous) => {
    if (active > 0) {
      runPeak = Math.max(runPeak, active)
      return
    }
    if (!previous) return

    const peak = runPeak
    runPeak = 0
    if (skipNextRunEndToast) {
      skipNextRunEndToast = false
      return
    }
    if (peak > 1) announceRunEnd()
  },
  { immediate: true },
)

function announceRunEnd() {
  const { completed, failed, unavailable } = stats.value
  if (completed === 0) {
    toast.add({
      title: 'Download failed',
      description: `${songs(failed + unavailable)} couldn't be downloaded.`,
      color: 'error',
      duration: 0,
    })
    return
  }

  const missed = missedParts()
  toast.add({
    title: [`${songs(completed)} ready`, ...missed].join(' · '),
    color: missed.length > 0 ? 'warning' : 'success',
    duration: ACTION_TOAST_DURATION_MS,
    actions: [
      {
        label: 'Save as zip',
        icon: 'i-lucide-folder-archive',
        color: 'neutral',
        variant: 'outline',
        onClick: () => saveFromUrl(zipUrl.value),
      },
    ],
  })
}
</script>

<template>
  <div>
    <div class="bg-default/90 sticky top-0 z-10 border-b border-default backdrop-blur">
      <div class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2 px-4 py-3 sm:px-6">
        <div class="flex min-w-0 flex-wrap items-baseline gap-x-2">
          <h2 class="text-highlighted font-semibold">Songs</h2>
          <p v-if="stats.total > 0" class="text-muted text-sm tabular-nums">
            {{ progressSummary }}
          </p>
        </div>

        <div v-if="stats.total > 0" class="flex flex-wrap items-center gap-2">
          <UButton
            v-if="stats.completed > 0"
            :to="zipUrl"
            external
            download
            icon="i-lucide-folder-archive"
            color="neutral"
            variant="outline"
            :aria-label="zipLabel"
          >
            Save as zip
            <span v-if="zipCount" class="hidden sm:inline">({{ zipCount }})</span>
          </UButton>
          <CancelButton
            v-if="stats.active > 0"
            :label="`Cancel (${stats.active} left)`"
            :loading="cancellingRun"
            @click="handleCancelAll"
          />
          <BaseButton
            v-else
            :label="primaryAction.label"
            :icon="primaryAction.icon"
            :color="primaryAction.color"
            :variant="primaryAction.variant"
            :disabled="primaryAction.disabled"
            :loading="startingRun"
            @click="handleDownloadAll"
          />
        </div>
      </div>

      <UProgress
        v-if="stats.active > 0"
        :model-value="overallPercent"
        size="xs"
        :aria-label="progressSummary"
        class="absolute inset-x-0 bottom-0 translate-y-1/2"
      />
      <p class="sr-only" role="status" aria-live="polite">{{ announcement }}</p>
    </div>

    <ul class="divide-y divide-default">
      <PlaylistEntryItem
        v-for="entry in entries"
        :key="entry.id"
        :entry="entry"
        :progress="progress[entry.id]"
        :busy="isEntryBusy(entry.id)"
        @action="(action) => handleEntryAction(entry, action)"
      />
    </ul>
  </div>
</template>

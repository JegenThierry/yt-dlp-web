<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import SongList from '@/components/molecules/SongList.vue'
import { fetchPlaylist } from '@/api/playlistService'
import { subscribePlaylistEvents } from '@/api/downloadService'
import { DownloadStatus } from '@/models/enums'
import type { Playlist } from '@/models/playlist'
import type { ProgressEvent } from '@/models/download'
import { formatDate } from '@/utils/format'

const route = useRoute()

const playlist = ref<Playlist | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const unsubscribe = ref<(() => void) | null>(null)

const playlistId = computed(() => Number((route.params as { id?: string }).id))

async function loadPlaylist() {
  const id = playlistId.value

  if (!Number.isInteger(id) || id < 1) {
    error.value = 'This playlist url is not valid.'
    loading.value = false
    return
  }

  loading.value = true
  error.value = null

  try {
    playlist.value = await fetchPlaylist(id)
    subscribeToProgress(id)
  } catch (err) {
    if (axios.isAxiosError(err) && err.response?.status === 404) {
      error.value = 'This playlist does not exist.'
    } else {
      console.error('Failed to fetch playlist:', err)
      error.value = 'Something went wrong while loading the playlist.'
    }
  } finally {
    loading.value = false
  }
}

function subscribeToProgress(id: number) {
  unsubscribe.value?.()

  unsubscribe.value = subscribePlaylistEvents(id, (event: ProgressEvent) => {
    const current = playlist.value
    if (!current || current.id !== event.playlist_id) return

    const entries = (current.entries ?? []).map((entry) =>
      entry.id === event.entry_id
        ? { ...entry, status: event.status, file_path: event.file_path ?? entry.file_path }
        : entry,
    )

    playlist.value = { ...current, entries }
    maybeStopListening()
  })
}

function maybeStopListening() {
  const entries = playlist.value?.entries ?? []
  if (entries.length === 0) return

  const stillRunning = entries.some(
    (entry) => entry.status !== DownloadStatus.Completed && entry.status !== DownloadStatus.Failed,
  )

  if (!stillRunning) unsubscribe.value?.()
}

function handlePlaylistUpdate(updated: Playlist) {
  playlist.value = updated
  if (!unsubscribe.value) subscribeToProgress(updated.id)
  else maybeStopListening()
}

onMounted(loadPlaylist)
onBeforeUnmount(() => unsubscribe.value?.())

const songCount = computed(() => playlist.value?.entries?.length ?? 0)

// True once every entry has reached a terminal state (completed or failed).
const allComplete = computed(() => {
  const entries = playlist.value?.entries ?? []
  if (entries.length === 0) return false
  return entries.every(
    (entry) => entry.status === DownloadStatus.Completed || entry.status === DownloadStatus.Failed,
  )
})
</script>

<template>
  <PageLayout :title="playlist?.title">
    <template #navigation>
      <UButton
        to="/playlists"
        icon="i-lucide-arrow-left"
        color="neutral"
        variant="ghost"
        size="sm"
        label="Back to playlists"
        class="-ms-2.5 self-start"
      />
    </template>

    <template v-if="!error" #title>
      <USkeleton v-if="loading" class="h-8 w-2/3 sm:max-w-md" />
      <template v-else-if="playlist">{{ playlist.title }}</template>
    </template>

    <template v-if="playlist" #meta>
      <UBadge
        color="neutral"
        variant="subtle"
        icon="i-lucide-list-music"
        :label="`${songCount} ${songCount === 1 ? 'song' : 'songs'}`"
      />
      <UBadge
        color="neutral"
        variant="subtle"
        icon="i-lucide-clock"
        :label="`Updated ${formatDate(playlist.updated_at)}`"
      />
    </template>

    <template v-if="playlist" #actions>
      <UButton
        :to="`/playlists/${playlist.id}/edit`"
        icon="i-lucide-pen-line"
        color="neutral"
        variant="outline"
        label="Clean up metadata"
      />
      <UButton
        :to="playlist.url"
        target="_blank"
        external
        icon="i-lucide-external-link"
        color="neutral"
        variant="outline"
        label="Open source"
      />
    </template>

    <UCard v-if="loading" variant="outline">
      <template #header>
        <USkeleton class="h-5 w-16" />
      </template>

      <div class="space-y-4">
        <USkeleton v-for="n in 5" :key="n" class="h-10" />
      </div>
    </UCard>

    <UEmpty
      v-else-if="error"
      variant="naked"
      icon="i-lucide-search-x"
      title="Playlist unavailable"
      :description="error"
      class="py-16"
    >
      <template #actions>
        <UButton
          icon="i-lucide-refresh-cw"
          color="neutral"
          variant="outline"
          label="Try again"
          @click="loadPlaylist"
        />
        <UButton to="/playlists" color="neutral" variant="ghost" label="All playlists" />
      </template>
    </UEmpty>

    <UCard v-else-if="playlist" variant="outline" :ui="{ body: 'p-0' }">
      <template #header>
        <h2 class="font-semibold">Songs</h2>
      </template>

      <SongList :playlist="playlist" :all-complete="allComplete" @update:playlist="handlePlaylistUpdate" />
    </UCard>
  </PageLayout>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import { useToast } from '@nuxt/ui/composables'
import { createPlaylist } from '@/api/playlistService'
import { subscribePlaylistEvents } from '@/api/downloadService'
import InputText from '@/components/atoms/inputs/InputText.vue'
import InputSelect from '@/components/atoms/inputs/InputSelect.vue'
import InputSwitch from '@/components/atoms/inputs/InputSwitch.vue'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'
import SongList from '@/components/molecules/SongList.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
import { AudioFormat, AudioQuality, DownloadStatus } from '@/models/enums'
import type { Playlist, CreatePlaylistRequest } from '@/models/playlist'
import type { FilterOptions } from '@/models/filterOptions'
import type { ProgressEvent } from '@/models/download'

interface DownloadForm {
  url: string
  format: AudioFormat
  quality: AudioQuality
  embedMetadata: boolean
}

const toast = useToast()

const createdPlaylist = ref<Playlist | null>(null)
const activeFilters = ref<FilterOptions | undefined>()
const unsubscribe = ref<(() => void) | null>(null)

// True once every entry has reached a terminal state (completed or failed).
const allComplete = computed(() => {
  const entries = createdPlaylist.value?.entries ?? []
  if (entries.length === 0) return false
  return entries.every(
    (entry) => entry.status === DownloadStatus.Completed || entry.status === DownloadStatus.Failed,
  )
})

const formatOptions = Object.values(AudioFormat).map((format) => ({
  label: format.toUpperCase(),
  value: format,
}))

const qualityOptions = Object.values(AudioQuality).map((quality) => ({
  label: quality.charAt(0).toUpperCase() + quality.slice(1),
  value: quality,
}))

async function handleDownload(values: DownloadForm) {
  try {
    const response = await createPlaylist(buildCreatePlaylistRequest(values))
    activeFilters.value = buildFilters(values)
    createdPlaylist.value = response
    subscribeToProgress(response.id)
    toast.add({
      title: 'Playlist ready',
      description: `Found ${response.entry_count} song${response.entry_count === 1 ? '' : 's'}.`,
      color: 'success',
    })
  } catch (error) {
    console.error('Error creating playlist:', error)
    throw error
  }
}

function buildCreatePlaylistRequest(values: DownloadForm): CreatePlaylistRequest {
  return {
    url: values.url,
    target_format: values.format || AudioFormat.MP3,
  }
}

function buildFilters(values: DownloadForm): FilterOptions {
  return {
    format: values.format || AudioFormat.MP3,
    quality: values.quality || AudioQuality.Best,
    embed_metadata: values.embedMetadata,
  }
}

function subscribeToProgress(playlistId: number) {
  unsubscribe.value?.()

  unsubscribe.value = subscribePlaylistEvents(playlistId, (event: ProgressEvent) => {
    const playlist = createdPlaylist.value
    if (!playlist || playlist.id !== event.playlist_id) return

    const entries = (playlist.entries ?? []).map((entry) =>
      entry.id === event.entry_id
        ? { ...entry, status: event.status, file_path: event.file_path ?? entry.file_path }
        : entry,
    )

    createdPlaylist.value = { ...playlist, entries }
    maybeStopListening()
  })
}

function maybeStopListening() {
  const entries = createdPlaylist.value?.entries ?? []
  if (entries.length === 0) return

  const stillRunning = entries.some(
    (entry) => entry.status !== DownloadStatus.Completed && entry.status !== DownloadStatus.Failed,
  )

  if (!stillRunning) unsubscribe.value?.()
}

function handlePlaylistUpdate(playlist: Playlist) {
  createdPlaylist.value = playlist
  if (!unsubscribe.value) subscribeToProgress(playlist.id)
  else maybeStopListening()
}

onBeforeUnmount(() => unsubscribe.value?.())
</script>

<template>
  <PageLayout title="Download">
    <UCard variant="outline">
      <FormWrapper
        :initial-values="{
          url: '',
          format: AudioFormat.MP3,
          quality: AudioQuality.Best,
          embedMetadata: true,
        }"
        :on-submit="handleDownload"
        v-slot="{ state, loading }"
      >
        <InputText
          v-model="state.url"
          label="Video URL"
          placeholder="youtube.com/video-url"
          icon="i-lucide-link"
          required
        />

        <div class="grid gap-4 sm:grid-cols-2">
          <InputSelect
            v-model="state.format"
            label="Audio format"
            :options="formatOptions"
            icon="i-lucide-file-audio"
            required
          />
          <InputSelect
            v-model="state.quality"
            label="Quality"
            :options="qualityOptions"
            icon="i-lucide-gauge"
            required
          />
        </div>

        <InputSwitch
          v-model="state.embedMetadata"
          label="Embed metadata"
          description="Write title and artist tags into the downloaded file."
        />

        <DownloadButton :loading="loading" :disabled="!state.url" type="submit" />
      </FormWrapper>
    </UCard>

    <UCard v-if="createdPlaylist" variant="outline" :ui="{ body: 'p-0' }">
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h2 class="truncate font-semibold">{{ createdPlaylist.title }}</h2>
          <UButton
            :to="`/playlists/${createdPlaylist.id}`"
            color="neutral"
            variant="ghost"
            size="sm"
            icon="i-lucide-external-link"
            label="Open playlist"
          />
        </div>
      </template>

      <SongList
        :playlist="createdPlaylist"
        :filters="activeFilters"
        :all-complete="allComplete"
        @update:playlist="handlePlaylistUpdate"
      />
    </UCard>
  </PageLayout>
</template>

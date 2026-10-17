<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import { useToast } from '@nuxt/ui/composables'
import { apiErrorStatus, getApiErrorMessage } from '@/api/apiErrors'
import { downloadAll } from '@/api/downloadService'
import { createPlaylist } from '@/api/playlistService'
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'
import InputSelect from '@/components/atoms/inputs/InputSelect.vue'
import InputSwitch from '@/components/atoms/inputs/InputSwitch.vue'
import InputText from '@/components/atoms/inputs/InputText.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import PlaylistRow from '@/components/molecules/PlaylistRow.vue'
import { useLibraryRefresh } from '@/composables/useLibraryRefresh'
import { AudioFormat, AudioQuality } from '@/models/enums'
import type { PlaylistImport } from '@/models/playlist'
import { formatQuality } from '@/utils/format'

interface AddPlaylistForm extends Record<string, unknown> {
  url: string
  format: AudioFormat
  quality: AudioQuality
  autoStart: boolean
}

const RECENT_COUNT = 5
const SLOW_HINT_DELAY_MS = 3_000
const URL_ERROR_STATUSES = [400, 403, 404, 422]
const URL_PLACEHOLDER = 'https://www.youtube.com/playlist?list=…'

const router = useRouter()
const toast = useToast()
const store = useLibraryRefresh()
const { playlists } = storeToRefs(store)

const recentPlaylists = computed(() => playlists.value.slice(0, RECENT_COUNT))

const slow = ref(false)
let slowTimer: ReturnType<typeof setTimeout> | undefined

onBeforeUnmount(() => clearTimeout(slowTimer))

const formatOptions = Object.values(AudioFormat).map((format) => ({
  label: format.toUpperCase(),
  value: format,
}))

const qualityOptions = Object.values(AudioQuality).map((quality) => ({
  label: formatQuality(quality),
  value: quality,
}))

function isHttpUrl(value: string): boolean {
  if (!URL.canParse(value)) return false
  const url = new URL(value)
  return (url.protocol === 'http:' || url.protocol === 'https:') && url.hostname !== ''
}

function validate(values: AddPlaylistForm) {
  const url = values.url.trim()
  if (!url) return [{ field: 'url', message: 'Paste a playlist or video link.' }]
  if (isHttpUrl(url)) return []
  return [{ field: 'url', message: 'Enter a full link that starts with https://' }]
}

function urlFieldErrors(error: unknown) {
  const status = apiErrorStatus(error)
  if (status === undefined || !URL_ERROR_STATUSES.includes(status)) return []
  return [{ field: 'url', message: getApiErrorMessage(error, "This link can't be added.") }]
}

function count(n: number, singular: string, plural: string): string {
  return `${n} ${n === 1 ? singular : plural}`
}

function announce(result: PlaylistImport) {
  const parts = result.existing
    ? [
        result.added_count > 0
          ? `${count(result.added_count, 'new song', 'new songs')} added`
          : 'No new songs',
      ]
    : [`${count(result.added_count, 'song', 'songs')} added`]
  if (result.skipped_count > 0) {
    parts.push(`${count(result.skipped_count, 'unavailable video', 'unavailable videos')} skipped`)
  }

  toast.add({
    title: result.existing ? 'Already in your library' : 'Playlist added',
    description: parts.join(' · '),
    color: result.existing ? 'info' : 'success',
  })
}

async function startDownloads(playlistId: number) {
  try {
    await downloadAll(playlistId)
  } catch (error) {
    if (apiErrorStatus(error) === 409) return
    console.error('Failed to start playlist download:', error)
    toast.add({
      title: "Couldn't start downloading",
      description: getApiErrorMessage(error, 'Use Download all on the playlist page to try again.'),
      color: 'error',
      duration: 0,
    })
  }
}

async function handleAdd(values: AddPlaylistForm) {
  slowTimer = setTimeout(() => (slow.value = true), SLOW_HINT_DELAY_MS)
  try {
    const result = await createPlaylist({
      url: values.url.trim(),
      target_format: values.format || AudioFormat.MP3,
      quality: values.quality || AudioQuality.Best,
    })
    announce(result)
    if (values.autoStart) await startDownloads(result.id)
    await router.push(`/playlists/${result.id}`)
  } finally {
    clearTimeout(slowTimer)
    slow.value = false
  }
}
</script>

<template>
  <PageLayout title="Add a playlist">
    <UCard variant="outline">
      <FormWrapper
        :initial-values="{
          url: '',
          format: AudioFormat.MP3,
          quality: AudioQuality.Best,
          autoStart: true,
        }"
        :on-submit="handleAdd"
        :on-validate="validate"
        :field-errors="urlFieldErrors"
        error-title="Couldn't add the playlist"
        v-slot="{ state, loading }"
      >
        <InputText
          v-model="state.url"
          name="url"
          type="url"
          label="Playlist or video link"
          :placeholder="URL_PLACEHOLDER"
          icon="i-lucide-link"
          inputmode="url"
          autocomplete="off"
          spellcheck="false"
          required
        />

        <div class="grid gap-4 sm:grid-cols-2">
          <InputSelect
            v-model="state.format"
            name="format"
            label="Audio format"
            :options="formatOptions"
            icon="i-lucide-file-audio"
          />
          <InputSelect
            v-model="state.quality"
            name="quality"
            label="Quality"
            :options="qualityOptions"
            icon="i-lucide-gauge"
          />
        </div>

        <InputSwitch
          v-model="state.autoStart"
          label="Start downloading right away"
          description="Songs download on the server as soon as the playlist is added."
        />

        <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
          <BaseButton
            type="submit"
            label="Add playlist"
            icon="i-lucide-list-plus"
            :loading="loading"
          />
          <p v-if="slow" class="text-muted text-sm" role="status">
            Reading the playlist from YouTube… large playlists can take a minute.
          </p>
        </div>
      </FormWrapper>
    </UCard>

    <UCard
      v-if="recentPlaylists.length > 0"
      variant="outline"
      :ui="{ header: 'py-3 sm:py-3', body: 'p-0 sm:p-0' }"
    >
      <template #header>
        <div class="flex items-center justify-between gap-3">
          <h2 class="text-highlighted font-semibold">Recently added</h2>
          <UButton
            to="/playlists"
            color="neutral"
            variant="ghost"
            size="sm"
            trailing-icon="i-lucide-arrow-right"
            label="View all"
            class="-me-2.5"
          />
        </div>
      </template>

      <ul class="divide-y divide-default">
        <PlaylistRow
          v-for="playlist in recentPlaylists"
          :key="playlist.id"
          :playlist="playlist"
          heading="h3"
          :actions="false"
        />
      </ul>
    </UCard>
  </PageLayout>
</template>

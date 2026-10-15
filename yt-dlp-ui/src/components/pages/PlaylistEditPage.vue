<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import { useToast } from '@nuxt/ui/composables'
import PageLayout from '@/components/layouts/PageLayout.vue'
import InputText from '@/components/atoms/inputs/InputText.vue'
import { fetchPlaylist, updatePlaylistEntryMetadata } from '@/api/playlistService'
import type { Playlist } from '@/models/playlist'
import type { PlaylistEntry } from '@/models/playlistEntry'

interface EntryDraft {
  id: number
  title: string
  artist: string
  genre: string
  year: string
}

const MIN_YEAR = 1900
const MAX_YEAR = 2100

const route = useRoute('/playlists/[id]/edit')
const router = useRouter()
const toast = useToast()

const playlist = ref<Playlist | null>(null)
const drafts = ref<EntryDraft[]>([])
const loading = ref(true)
const saving = ref(false)
const error = ref<string | null>(null)

const playlistId = Number(route.params.id)

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const data = await fetchPlaylist(playlistId)
    playlist.value = data
    drafts.value = (data.entries ?? []).map(toDraft)
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
})

function toDraft(entry: PlaylistEntry): EntryDraft {
  return {
    id: entry.id,
    title: entry.title,
    artist: entry.artist ?? '',
    genre: entry.genre ?? '',
    year: entry.year ? String(entry.year) : '',
  }
}

function buildChanges(): {
  id: number
  title: string
  artist: string
  genre: string
  year: string
}[] {
  const original = (playlist.value?.entries ?? []).reduce<Record<number, PlaylistEntry>>(
    (map, entry) => {
      map[entry.id] = entry
      return map
    },
    {},
  )

  return drafts.value
    .filter((draft) => {
      const prev = original[draft.id]
      if (!prev) return true
      return (
        draft.title !== prev.title ||
        draft.artist !== (prev.artist ?? '') ||
        draft.genre !== (prev.genre ?? '') ||
        draft.year !== (prev.year ? String(prev.year) : '')
      )
    })
    .map((draft) => ({
      id: draft.id,
      title: draft.title,
      artist: draft.artist,
      genre: draft.genre,
      year: draft.year,
    }))
}

async function handleSave() {
  const changes = buildChanges()
  if (changes.length === 0) {
    toast.add({ title: 'Nothing to save', color: 'neutral' })
    return
  }

  if (
    changes.some((c) => c.year !== '' && (Number(c.year) < MIN_YEAR || Number(c.year) > MAX_YEAR))
  ) {
    toast.add({
      title: 'Invalid year',
      description: 'Year must be between 1900 and 2100.',
      color: 'error',
    })
    return
  }

  saving.value = true
  try {
    await Promise.all(
      changes.map((c) =>
        updatePlaylistEntryMetadata(c.id, {
          title: c.title,
          artist: c.artist,
          genre: c.genre,
          year: c.year === '' ? 0 : Number(c.year),
        }),
      ),
    )
    toast.add({
      title: 'Metadata saved',
      description: `Updated ${changes.length} song${changes.length === 1 ? '' : 's'}.`,
      color: 'success',
    })
    router.push(`/playlists/${playlistId}`)
  } catch (err) {
    console.error('Failed to save metadata:', err)
    toast.add({
      title: 'Save failed',
      description: 'Something went wrong while saving the metadata.',
      color: 'error',
    })
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <PageLayout :title="playlist?.title">
    <template #navigation>
      <UButton
        :to="`/playlists/${playlistId}`"
        icon="i-lucide-arrow-left"
        color="neutral"
        variant="ghost"
        size="sm"
        label="Back to playlist"
        class="-ms-2.5 self-start"
      />
    </template>

    <template v-if="!error" #title>
      <USkeleton v-if="loading" class="h-8 w-2/3 sm:max-w-md" />
      <template v-else-if="playlist">Edit metadata</template>
    </template>

    <template v-if="playlist && !error" #meta>
      <p class="text-muted text-sm">
        Correct inaccurate source metadata. Values you enter are embedded when a song is encoded;
        leave a field empty to keep the value pulled from the source.
      </p>
    </template>

    <template v-if="playlist && !error" #actions>
      <UButton
        icon="i-lucide-save"
        color="primary"
        variant="solid"
        label="Save changes"
        :loading="saving"
        :disabled="saving || drafts.length === 0"
        @click="handleSave"
      />
    </template>

    <UCard v-if="loading" variant="outline">
      <template #header>
        <USkeleton class="h-5 w-16" />
      </template>
      <div class="space-y-4">
        <USkeleton v-for="n in 5" :key="n" class="h-12" />
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
        <UButton to="/playlists" color="neutral" variant="ghost" label="All playlists" />
      </template>
    </UEmpty>

    <UCard v-else-if="playlist && drafts.length === 0" variant="outline">
      <p class="text-muted p-6 text-center text-sm">No songs in this playlist.</p>
    </UCard>

    <UCard
      v-else-if="playlist"
      variant="outline"
      :ui="{ body: 'max-h-[70vh] overflow-y-auto p-0 divide-y divide-default' }"
    >
      <div
        v-for="draft in drafts"
        :key="draft.id"
        class="grid gap-3 px-4 py-3 sm:px-6 sm:grid-cols-2 lg:grid-cols-4"
      >
        <InputText v-model="draft.title" label="Title" :required="false" placeholder="Song title" />
        <InputText v-model="draft.artist" label="Artist" :required="false" placeholder="Artist" />
        <InputText v-model="draft.genre" label="Genre" :required="false" placeholder="Genre" />
        <InputText
          v-model="draft.year"
          label="Year"
          :required="false"
          placeholder="1900–2100"
          inputmode="numeric"
        />
      </div>
    </UCard>
  </PageLayout>
</template>

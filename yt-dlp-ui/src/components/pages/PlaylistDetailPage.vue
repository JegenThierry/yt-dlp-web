<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import PlaylistEntryItem from '@/components/molecules/PlaylistEntryItem.vue'
import { fetchPlaylist } from '@/api/playlistService'
import type { Playlist } from '@/models/playlist'
import { formatDate } from '@/utils/format'

const route = useRoute('/playlists/[id]')

const playlist = ref<Playlist | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const playlistId = computed(() => Number(route.params.id))

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

onMounted(loadPlaylist)

const songCount = computed(() => playlist.value?.entries?.length ?? 0)
</script>

<template>
  <PageLayout>
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

      <ul v-if="songCount > 0" class="divide-y divide-default overflow-y-auto max-h-[60vh]">
        <PlaylistEntryItem v-for="entry in playlist.entries" :key="entry.id" :entry="entry" />
      </ul>

      <div v-else class="flex flex-col items-center gap-2 py-8 text-center">
        <UIcon name="i-lucide-music" class="text-dimmed size-8" />
        <p class="text-muted text-sm">No songs were found for this playlist.</p>
      </div>
    </UCard>
  </PageLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import PlaylistRow from '@/components/molecules/PlaylistRow.vue'
import { fetchPlaylists } from '@/api/playlistService'
import type { Playlist } from '@/models/playlist'

const playlists = ref<Playlist[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

async function loadPlaylists() {
  loading.value = true
  error.value = null

  try {
    playlists.value = await fetchPlaylists()
  } catch (err) {
    console.error('Failed to fetch playlists:', err)
    error.value = 'Something went wrong while loading your playlists.'
  } finally {
    loading.value = false
  }
}

onMounted(loadPlaylists)
</script>

<template>
  <PageLayout title="Playlists">
    <template v-if="playlists.length > 0" #actions>
      <UBadge color="neutral" variant="subtle" :label="playlists.length" />
    </template>

    <UCard v-if="loading" variant="outline">
      <div class="space-y-4">
        <USkeleton v-for="n in 5" :key="n" class="h-12" />
      </div>
    </UCard>

    <UEmpty
      v-else-if="error"
      variant="naked"
      icon="i-lucide-circle-x"
      title="Failed to load playlists"
      :description="error"
      class="py-16"
    >
      <template #actions>
        <UButton
          icon="i-lucide-refresh-cw"
          color="neutral"
          variant="outline"
          label="Try again"
          @click="loadPlaylists"
        />
      </template>
    </UEmpty>

    <UEmpty
      v-else-if="playlists.length === 0"
      variant="naked"
      icon="i-lucide-list-music"
      title="No playlists yet"
      description="Download a playlist to see it show up here."
      class="py-16"
    >
      <template #actions>
        <UButton to="/" icon="i-lucide-download" label="Go to download" variant="link" />
      </template>
    </UEmpty>

    <UCard v-else variant="outline" :ui="{ body: 'p-0' }">
      <ul class="divide-y divide-default">
        <PlaylistRow v-for="playlist in playlists" :key="playlist.id" :playlist="playlist" />
      </ul>
    </UCard>
  </PageLayout>
</template>

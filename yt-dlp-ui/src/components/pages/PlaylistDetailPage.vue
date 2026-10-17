<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { useMediaQuery, useTitle } from '@vueuse/core'
import type { DropdownMenuItem } from '@nuxt/ui'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UDropdownMenu from '@nuxt/ui/components/DropdownMenu.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import DeletePlaylistModal from '@/components/molecules/DeletePlaylistModal.vue'
import SongList from '@/components/molecules/SongList.vue'
import { getApiErrorMessage } from '@/api/apiErrors'
import { fetchPlaylist } from '@/api/playlistService'
import { useCheckForNewSongs } from '@/composables/useCheckForNewSongs'
import { usePlaylistProgress } from '@/composables/usePlaylistProgress'
import type { Playlist } from '@/models/playlist'
import { formatDate, formatQuality } from '@/utils/format'
import { playlistStats } from '@/utils/playlistStats'
import { sourceLabel } from '@/utils/source'
import { pageTitle } from '@/utils/title'

const route = useRoute()
const router = useRouter()
const smUp = useMediaQuery('(min-width: 640px)')

const playlist = ref<Playlist | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const errorRetryable = ref(false)
const deleteOpen = ref(false)

const { progress } = usePlaylistProgress(playlist)
const { refreshing, checkForNewSongs } = useCheckForNewSongs()

const playlistId = computed(() => Number((route.params as { id?: string }).id))

async function loadPlaylist() {
  const id = playlistId.value

  if (!Number.isInteger(id) || id < 1) {
    error.value = 'This playlist url is not valid.'
    errorRetryable.value = false
    loading.value = false
    return
  }

  loading.value = true
  error.value = null
  errorRetryable.value = false

  try {
    playlist.value = await fetchPlaylist(id)
  } catch (err) {
    if (axios.isAxiosError(err) && err.response?.status === 404) {
      error.value = 'This playlist does not exist.'
    } else {
      console.error('Failed to fetch playlist:', err)
      error.value = getApiErrorMessage(err, 'Something went wrong while loading the playlist.')
      errorRetryable.value = true
    }
  } finally {
    loading.value = false
  }
}

onMounted(loadPlaylist)

async function checkSongs() {
  if (!playlist.value) return
  const result = await checkForNewSongs(playlist.value.id)
  if (result && playlist.value?.id === result.id) playlist.value = result
}

function onDeleted() {
  void router.push('/playlists')
}

const stats = computed(() => playlistStats(playlist.value?.entries ?? []))
const songCount = computed(() => stats.value.total)

useTitle(
  () => {
    if (!playlist.value) return pageTitle('Playlist')
    const { active, completed, total } = stats.value
    const progressPrefix = active > 0 ? `(${completed}/${total}) ` : ''
    return progressPrefix + pageTitle(playlist.value.title)
  },
  { restoreOnUnmount: false },
)

const settingsLabel = computed(() => {
  if (!playlist.value) return ''
  const formats = new Set((playlist.value.entries ?? []).map((entry) => entry.format))
  const formatLabels = [...formats].filter(Boolean).map((format) => format.toUpperCase())
  return [...formatLabels, `${formatQuality(playlist.value.quality)} quality`].join(' · ')
})

const openSourceLabel = computed(() => sourceLabel(playlist.value?.url ?? ''))

const editPath = computed(() => `/playlists/${playlistId.value}/edit`)

const headerMenuItems = computed<DropdownMenuItem[][]>(() => {
  if (!playlist.value) return []

  const manage: DropdownMenuItem[] = []
  if (songCount.value > 0 && !smUp.value) {
    manage.push({ label: 'Edit metadata', icon: 'i-lucide-pen-line', to: editPath.value })
  }
  if (songCount.value > 0) {
    manage.push({
      label: 'Check for new songs',
      icon: 'i-lucide-refresh-cw',
      disabled: refreshing.value,
      onSelect: () => void checkSongs(),
    })
  }
  if (!smUp.value) {
    manage.push({
      label: openSourceLabel.value,
      icon: 'i-lucide-external-link',
      to: playlist.value.url,
      target: '_blank',
    })
  }

  const danger: DropdownMenuItem[] = [
    {
      label: 'Delete…',
      icon: 'i-lucide-trash-2',
      color: 'error',
      onSelect: () => (deleteOpen.value = true),
    },
  ]

  return [manage, danger].filter((group) => group.length > 0)
})
</script>

<template>
  <PageLayout :title="error ? 'Playlist unavailable' : playlist?.title" :loading="loading">
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

    <template v-if="playlist" #meta>
      <UBadge
        color="neutral"
        variant="subtle"
        icon="i-lucide-list-music"
        :label="`${songCount} ${songCount === 1 ? 'song' : 'songs'}`"
      />
      <UBadge color="neutral" variant="subtle" icon="i-lucide-file-audio" :label="settingsLabel" />
      <UBadge
        color="neutral"
        variant="subtle"
        icon="i-lucide-clock"
        :label="`Updated ${formatDate(playlist.updated_at)}`"
      />
    </template>

    <template v-if="playlist" #actions>
      <UButton
        v-if="songCount > 0"
        :to="editPath"
        icon="i-lucide-pen-line"
        color="neutral"
        variant="outline"
        label="Edit metadata"
        class="hidden sm:inline-flex"
      />
      <UButton
        :to="playlist.url"
        target="_blank"
        external
        icon="i-lucide-external-link"
        color="neutral"
        variant="outline"
        :label="openSourceLabel"
        class="hidden sm:inline-flex"
      />
      <UDropdownMenu :items="headerMenuItems" :content="{ align: 'end' }">
        <UButton
          icon="i-lucide-ellipsis"
          color="neutral"
          variant="outline"
          :loading="refreshing"
          aria-label="More actions"
        />
      </UDropdownMenu>
      <DeletePlaylistModal v-model:open="deleteOpen" :playlist="playlist" @deleted="onDeleted" />
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
      :description="error"
      class="py-16"
    >
      <template #actions>
        <template v-if="errorRetryable">
          <UButton
            icon="i-lucide-refresh-cw"
            color="neutral"
            variant="outline"
            label="Try again"
            @click="loadPlaylist"
          />
          <UButton to="/playlists" color="neutral" variant="ghost" label="All playlists" />
        </template>
        <UButton
          v-else
          to="/playlists"
          icon="i-lucide-arrow-left"
          color="neutral"
          variant="outline"
          label="All playlists"
        />
      </template>
    </UEmpty>

    <UEmpty
      v-else-if="playlist && songCount === 0"
      variant="naked"
      icon="i-lucide-music"
      title="No songs in this playlist"
      description="Nothing downloadable was found when this link was added. Check again to pick up single videos and new songs, or delete it."
      class="py-16"
    >
      <template #actions>
        <UButton
          icon="i-lucide-refresh-cw"
          label="Check for new songs"
          :loading="refreshing"
          @click="checkSongs"
        />
        <UButton
          :to="playlist.url"
          target="_blank"
          external
          icon="i-lucide-external-link"
          color="neutral"
          variant="outline"
          :label="openSourceLabel"
        />
        <UButton
          icon="i-lucide-trash-2"
          color="error"
          variant="outline"
          label="Delete playlist"
          @click="deleteOpen = true"
        />
      </template>
    </UEmpty>

    <UCard
      v-else-if="playlist"
      variant="outline"
      :ui="{ root: 'overflow-clip', body: 'p-0 sm:p-0' }"
    >
      <SongList :playlist="playlist" :progress="progress" />
    </UCard>
  </PageLayout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useMediaQuery } from '@vueuse/core'
import type { DropdownMenuItem } from '@nuxt/ui'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UDropdownMenu from '@nuxt/ui/components/DropdownMenu.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import UProgress from '@nuxt/ui/components/Progress.vue'
import { useToast } from '@nuxt/ui/composables'
import DeletePlaylistModal from '@/components/molecules/DeletePlaylistModal.vue'
import { apiErrorStatus, getApiErrorMessage } from '@/api/apiErrors'
import { downloadAll, playlistZipUrl, saveFromUrl } from '@/api/downloadService'
import { useCheckForNewSongs } from '@/composables/useCheckForNewSongs'
import type { Playlist } from '@/models/playlist'
import { usePlaylistsStore } from '@/stores/playlists'
import type { BadgeColor } from '@/types/badgeTypes'
import { countLabel, formatDateTime, formatRelative, formatTotalDuration } from '@/utils/format'
import { libraryStats } from '@/utils/playlistStats'
import { sourceLabel } from '@/utils/source'

interface PlaylistRowProps {
  playlist: Playlist
  heading?: 'h2' | 'h3'
  actions?: boolean
}

const props = withDefaults(defineProps<PlaylistRowProps>(), {
  heading: 'h2',
  actions: true,
})

interface RowBadge {
  label: string
  color: BadgeColor
  icon: string
}

interface RowAction {
  kind: 'download' | 'zip'
  label: string
  icon: string
}

const toast = useToast()
const store = usePlaylistsStore()
const { refreshing, checkForNewSongs } = useCheckForNewSongs()
const smUp = useMediaQuery('(min-width: 640px)')

const starting = ref(false)
const deleteOpen = ref(false)

const to = computed(() => `/playlists/${props.playlist.id}`)
const stats = computed(() => libraryStats(props.playlist))
const zipUrl = computed(() => playlistZipUrl(props.playlist.id))

const songsLabel = computed(() => countLabel(stats.value.total, 'song', 'songs'))
const durationLabel = computed(() => formatTotalDuration(props.playlist.total_duration))
const addedLabel = computed(() => `added ${formatRelative(props.playlist.created_at)}`)
const addedTitle = computed(() => `Added ${formatDateTime(props.playlist.created_at)}`)

const downloadedCount = computed(() => {
  const { completed, active, total } = stats.value
  return completed === 0 && active === 0 ? '' : `${completed} of ${total}`
})

const showProgress = computed(() => {
  const { completed, active, total } = stats.value
  return active > 0 || (completed > 0 && completed < total)
})
const percent = computed(() =>
  stats.value.total === 0 ? 0 : (stats.value.completed / stats.value.total) * 100,
)

const badges = computed<RowBadge[]>(() => {
  const { total, active, failed, unavailable } = stats.value
  const list: RowBadge[] = []
  if (total === 0) {
    list.push({ label: 'No songs', color: 'warning', icon: 'i-lucide-triangle-alert' })
  }
  if (active > 0) {
    list.push({ label: 'Downloading…', color: 'info', icon: 'i-lucide-download' })
  }
  if (failed > 0) {
    list.push({ label: `${failed} failed`, color: 'error', icon: 'i-lucide-circle-x' })
  }
  if (unavailable > 0) {
    list.push({ label: `${unavailable} unavailable`, color: 'neutral', icon: 'i-lucide-ban' })
  }
  return list
})

const primaryAction = computed<RowAction | null>(() => {
  const { total, active, remaining, fresh, completed, failed, unavailable } = stats.value
  if (total === 0 || active > 0) return null
  if (remaining > 0) {
    if (fresh === 0) {
      return { kind: 'download', label: `Retry ${failed} failed`, icon: 'i-lucide-rotate-ccw' }
    }
    return {
      kind: 'download',
      label:
        completed + failed + unavailable === 0 ? 'Download' : `Download ${remaining} remaining`,
      icon: 'i-lucide-download',
    }
  }
  if (completed > 0) return { kind: 'zip', label: 'Save zip', icon: 'i-lucide-folder-archive' }
  return null
})

const zipMenuLabel = computed(() => {
  const { completed, total } = stats.value
  return completed < total ? `Save as zip (${completed} of ${total})` : 'Save as zip'
})

function runAction(action: RowAction) {
  if (action.kind === 'zip') saveFromUrl(zipUrl.value)
  else void startDownload()
}

const menuItems = computed<DropdownMenuItem[][]>(() => {
  const { total, completed } = stats.value
  const action = primaryAction.value

  const downloads: DropdownMenuItem[] = []
  if (action && !smUp.value) {
    downloads.push({ label: action.label, icon: action.icon, onSelect: () => runAction(action) })
  }
  if (completed > 0 && action?.kind !== 'zip') {
    downloads.push({
      label: zipMenuLabel.value,
      icon: 'i-lucide-folder-archive',
      onSelect: () => saveFromUrl(zipUrl.value),
    })
  }

  const manage: DropdownMenuItem[] = [
    {
      label: sourceLabel(props.playlist.url),
      icon: 'i-lucide-external-link',
      to: props.playlist.url,
      target: '_blank',
    },
  ]
  if (total > 0) {
    manage.push({
      label: 'Edit metadata',
      icon: 'i-lucide-pen-line',
      to: `/playlists/${props.playlist.id}/edit`,
    })
  }
  manage.push({
    label: 'Check for new songs',
    icon: 'i-lucide-refresh-cw',
    disabled: refreshing.value,
    onSelect: () => void refresh(),
  })

  const danger: DropdownMenuItem[] = [
    {
      label: 'Delete…',
      icon: 'i-lucide-trash-2',
      color: 'error',
      onSelect: () => (deleteOpen.value = true),
    },
  ]

  return [downloads, manage, danger].filter((group) => group.length > 0)
})

async function startDownload() {
  if (starting.value) return

  starting.value = true
  try {
    await downloadAll(props.playlist.id)
  } catch (error) {
    const conflict = apiErrorStatus(error) === 409
    if (!conflict) console.error('Failed to start playlist download:', error)
    toast.add({
      title: conflict ? 'Nothing new to download' : "Couldn't start downloading",
      description: getApiErrorMessage(error, 'Open the playlist and try again.'),
      color: conflict ? 'info' : 'error',
      duration: conflict ? undefined : 0,
    })
  } finally {
    starting.value = false
    void store.load()
  }
}

async function refresh() {
  const result = await checkForNewSongs(props.playlist.id)
  if (result) void store.load()
}

function focusPageHeading() {
  requestAnimationFrame(() => {
    document.querySelector<HTMLElement>('#main-content h1[tabindex]')?.focus()
  })
}
</script>

<template>
  <li
    class="has-focus-visible:bg-elevated/50 hover:bg-elevated/50 relative grid grid-cols-[minmax(0,1fr)_auto] items-center gap-x-3 px-4 py-3 transition-colors sm:grid-cols-[2.5rem_minmax(0,1fr)_auto] sm:px-6"
  >
    <div class="hidden size-10 items-center justify-center rounded-md bg-elevated sm:flex">
      <UIcon name="i-lucide-list-music" class="size-5 text-muted" />
    </div>

    <div class="min-w-0">
      <div class="flex min-w-0 items-center gap-2">
        <component
          :is="heading"
          class="text-highlighted min-w-0 text-base font-semibold break-words line-clamp-2 sm:line-clamp-1"
        >
          <RouterLink
            :to="to"
            :title="playlist.title"
            class="outline-none after:absolute after:inset-0 after:rounded-md after:outline-primary focus-visible:after:outline-2 focus-visible:after:-outline-offset-2"
          >
            {{ playlist.title }}
          </RouterLink>
        </component>
        <UBadge
          v-for="badge in badges"
          :key="badge.label"
          :color="badge.color"
          variant="subtle"
          size="sm"
          :icon="badge.icon"
          :label="badge.label"
          class="hidden shrink-0 sm:inline-flex"
        />
      </div>

      <div
        class="text-muted mt-1 flex flex-wrap items-center gap-x-1.5 gap-y-1 text-xs tabular-nums sm:text-sm"
      >
        <template v-if="stats.total > 0">
          <span>{{ songsLabel }}</span>
          <span aria-hidden="true" class="hidden sm:inline">·</span>
          <span class="hidden sm:inline">{{ durationLabel }}</span>
          <span v-if="!showProgress" aria-hidden="true">·</span>
          <span v-else aria-hidden="true" class="sm:hidden">·</span>
          <UProgress
            v-if="showProgress"
            :model-value="percent"
            size="xs"
            aria-hidden="true"
            class="hidden w-28 sm:block"
          />
          <span v-if="downloadedCount">
            {{ downloadedCount }}<span class="hidden sm:inline"> downloaded</span>
          </span>
          <span v-else>not downloaded</span>
        </template>
        <span v-else class="hidden sm:inline">No songs were found at this link</span>
        <span aria-hidden="true" class="hidden sm:inline">·</span>
        <time
          :datetime="new Date(playlist.created_at).toISOString()"
          :title="addedTitle"
          class="hidden sm:inline"
        >
          {{ addedLabel }}
        </time>
        <UBadge
          v-for="badge in badges"
          :key="badge.label"
          :color="badge.color"
          variant="subtle"
          size="sm"
          :icon="badge.icon"
          :label="badge.label"
          class="sm:hidden"
        />
      </div>

      <UProgress
        v-if="showProgress"
        :model-value="percent"
        size="xs"
        aria-hidden="true"
        class="mt-2 sm:hidden"
      />
    </div>

    <div v-if="actions" class="relative z-10 flex items-center gap-1">
      <UButton
        v-if="primaryAction?.kind === 'zip'"
        :to="zipUrl"
        external
        download
        :icon="primaryAction.icon"
        :label="primaryAction.label"
        color="neutral"
        variant="outline"
        size="sm"
        class="hidden sm:inline-flex"
      />
      <UButton
        v-else-if="primaryAction"
        :icon="primaryAction.icon"
        :label="primaryAction.label"
        :loading="starting"
        color="neutral"
        variant="outline"
        size="sm"
        class="hidden sm:inline-flex"
        @click="startDownload"
      />
      <UDropdownMenu :items="menuItems" :content="{ align: 'end' }">
        <UButton
          icon="i-lucide-ellipsis-vertical"
          color="neutral"
          variant="ghost"
          size="sm"
          :loading="refreshing || (starting && !smUp)"
          :aria-label="`More actions for ${playlist.title}`"
        />
      </UDropdownMenu>

      <DeletePlaylistModal
        v-model:open="deleteOpen"
        :playlist="playlist"
        @deleted="focusPageHeading"
      />
    </div>
    <UIcon
      v-else
      name="i-lucide-chevron-right"
      class="text-dimmed size-5 shrink-0"
      aria-hidden="true"
    />
  </li>
</template>

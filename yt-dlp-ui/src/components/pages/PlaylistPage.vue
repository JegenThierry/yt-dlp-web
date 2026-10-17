<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter, type LocationQueryValue } from 'vue-router'
import { storeToRefs } from 'pinia'
import { refDebounced } from '@vueuse/core'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import UInput from '@nuxt/ui/components/Input.vue'
import USelect from '@nuxt/ui/components/Select.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import PlaylistRow from '@/components/molecules/PlaylistRow.vue'
import { useLibraryRefresh } from '@/composables/useLibraryRefresh'
import type { Playlist } from '@/models/playlist'
import { countLabel } from '@/utils/format'
import { libraryStats } from '@/utils/playlistStats'

type SortKey = 'added' | 'activity' | 'title' | 'attention'

const DEFAULT_SORT: SortKey = 'added'
const SEARCH_DEBOUNCE_MS = 300

const sortOptions: { label: string; value: SortKey }[] = [
  { label: 'Recently added', value: 'added' },
  { label: 'Recent activity', value: 'activity' },
  { label: 'Title A–Z', value: 'title' },
  { label: 'Needs attention', value: 'attention' },
]

const route = useRoute()
const router = useRouter()
const store = useLibraryRefresh()
const { playlists, loaded, error } = storeToRefs(store)

function queryText(value: LocationQueryValue | LocationQueryValue[] | undefined): string {
  return typeof value === 'string' ? value : ''
}

function querySort(value: LocationQueryValue | LocationQueryValue[] | undefined): SortKey {
  const sort = queryText(value)
  return sortOptions.some((option) => option.value === sort) ? (sort as SortKey) : DEFAULT_SORT
}

const search = ref(queryText(route.query.q))
const sort = ref<SortKey>(querySort(route.query.sort))
const debouncedSearch = refDebounced(search, SEARCH_DEBOUNCE_MS)

watch(
  () => [route.query.q, route.query.sort],
  ([q, sortValue]) => {
    if (queryText(q) !== debouncedSearch.value.trim()) search.value = queryText(q)
    sort.value = querySort(sortValue)
  },
)

watch([debouncedSearch, sort], ([q, sortValue]) => {
  const nextQ = q.trim() || undefined
  const nextSort = sortValue === DEFAULT_SORT ? undefined : sortValue
  if (nextQ === (route.query.q ?? undefined) && nextSort === (route.query.sort ?? undefined)) return
  void router.replace({ query: { ...route.query, q: nextQ, sort: nextSort } })
})

const collator = new Intl.Collator('en', { sensitivity: 'base', numeric: true })

function time(date: Date | string): number {
  return new Date(date).getTime()
}

function attentionRank(playlist: Playlist): number {
  const { total, failed, remaining, active } = libraryStats(playlist)
  if (total === 0 || failed > 0) return 0
  if (remaining > 0) return 1
  if (active > 0) return 2
  return 3
}

const comparators: Record<SortKey, (a: Playlist, b: Playlist) => number> = {
  added: (a, b) => time(b.created_at) - time(a.created_at),
  activity: (a, b) => time(b.last_activity_at) - time(a.last_activity_at),
  title: (a, b) => collator.compare(a.title, b.title),
  attention: (a, b) =>
    attentionRank(a) - attentionRank(b) || time(b.last_activity_at) - time(a.last_activity_at),
}

const searchTerms = computed(() => search.value.toLocaleLowerCase().split(/\s+/).filter(Boolean))

const visiblePlaylists = computed(() => {
  const terms = searchTerms.value
  const matches =
    terms.length === 0
      ? playlists.value
      : playlists.value.filter((playlist) => {
          const haystack = `${playlist.title} ${playlist.url}`.toLocaleLowerCase()
          return terms.every((term) => haystack.includes(term))
        })
  return [...matches].sort(comparators[sort.value])
})

const countText = computed(() => {
  const total = playlists.value.length
  const shown = visiblePlaylists.value.length
  const label = countLabel(total, 'playlist', 'playlists')
  return shown < total ? `${shown} of ${label}` : label
})

function clearSearch() {
  search.value = ''
}
</script>

<template>
  <PageLayout title="Playlists">
    <template v-if="playlists.length > 0" #actions>
      <UButton to="/" icon="i-lucide-list-plus" label="Add playlist" />
    </template>

    <template v-if="playlists.length > 0" #meta>
      <UBadge color="neutral" variant="subtle" icon="i-lucide-list-music" :label="countText" />
    </template>

    <UCard v-if="!loaded && !error" variant="outline">
      <div class="space-y-4">
        <USkeleton v-for="n in 5" :key="n" class="h-12" />
      </div>
    </UCard>

    <UEmpty
      v-else-if="!loaded"
      variant="naked"
      icon="i-lucide-circle-x"
      title="Failed to load playlists"
      :description="error ?? undefined"
      class="py-16"
    >
      <template #actions>
        <UButton
          icon="i-lucide-refresh-cw"
          color="neutral"
          variant="outline"
          label="Try again"
          @click="store.load"
        />
      </template>
    </UEmpty>

    <UEmpty
      v-else-if="playlists.length === 0"
      variant="naked"
      icon="i-lucide-list-music"
      title="No playlists yet"
      description="Add a playlist or video link to see it show up here."
      class="py-16"
    >
      <template #actions>
        <UButton to="/" icon="i-lucide-list-plus" label="Add a playlist" />
      </template>
    </UEmpty>

    <template v-else>
      <div class="flex items-center gap-2" role="search">
        <UInput
          v-model="search"
          icon="i-lucide-search"
          placeholder="Search playlists"
          aria-label="Search playlists"
          autocomplete="off"
          class="min-w-0 flex-1"
          :ui="{ trailing: 'pe-1' }"
          @keydown.esc="clearSearch"
        >
          <template v-if="search" #trailing>
            <UButton
              color="neutral"
              variant="link"
              size="sm"
              icon="i-lucide-x"
              aria-label="Clear search"
              @click="clearSearch"
            />
          </template>
        </UInput>
        <USelect
          v-model="sort"
          :items="sortOptions"
          icon="i-lucide-arrow-down-up"
          aria-label="Sort playlists"
          class="w-48 shrink-0 sm:w-52"
        />
      </div>

      <UEmpty
        v-if="visiblePlaylists.length === 0"
        variant="naked"
        icon="i-lucide-search-x"
        title="No playlists match"
        :description="`Nothing in your library matches “${search.trim()}”.`"
        class="py-16"
      >
        <template #actions>
          <UButton
            icon="i-lucide-x"
            color="neutral"
            variant="outline"
            label="Clear search"
            @click="clearSearch"
          />
        </template>
      </UEmpty>

      <UCard v-else variant="outline" :ui="{ body: 'p-0 sm:p-0' }">
        <ul class="divide-y divide-default">
          <PlaylistRow
            v-for="playlist in visiblePlaylists"
            :key="playlist.id"
            :playlist="playlist"
          />
        </ul>
      </UCard>
    </template>
  </PageLayout>
</template>

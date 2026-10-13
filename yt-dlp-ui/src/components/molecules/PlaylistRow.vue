<script setup lang="ts">
import { computed } from 'vue'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import ULink from '@nuxt/ui/components/Link.vue'
import UTooltip from '@nuxt/ui/components/Tooltip.vue'
import type { Playlist } from '@/models/playlist'
import { formatDate } from '@/utils/format'

interface PlaylistRowProps {
  playlist: Playlist
}

const props = defineProps<PlaylistRowProps>()

const to = computed(() => ({
  name: '/playlists/[id]' as const,
  params: { id: props.playlist.id },
}))

const songLabel = computed(
  () => `${props.playlist.entry_count} ${props.playlist.entry_count === 1 ? 'song' : 'songs'}`,
)
</script>

<template>
  <li>
    <ULink
      :to="to"
      raw
      :aria-label="`Open playlist ${playlist.title}`"
      class="group flex items-center gap-3 rounded-md px-4 py-3 transition-colors outline-primary/25 hover:bg-elevated/50 focus-visible:outline-3 sm:px-6"
    >
      <div class="flex size-10 shrink-0 items-center justify-center rounded-md bg-elevated">
        <UIcon name="i-lucide-list-music" class="size-5 text-muted" />
      </div>

      <UTooltip :text="playlist.title" class="min-w-0 max-w-full flex-1">
        <h2 class="text-highlighted truncate text-base font-semibold">
          {{ playlist.title }}
        </h2>
      </UTooltip>

      <UBadge
        color="neutral"
        variant="subtle"
        size="sm"
        icon="i-lucide-music"
        :label="songLabel"
        class="hidden shrink-0 sm:inline-flex"
      />

      <UBadge
        color="neutral"
        variant="subtle"
        size="sm"
        icon="i-lucide-clock"
        :label="formatDate(playlist.updated_at)"
        class="hidden shrink-0 lg:inline-flex"
      />

      <UIcon
        name="i-lucide-chevron-right"
        class="text-dimmed group-hover:text-muted size-5 shrink-0 transition group-hover:translate-x-0.5"
      />
    </ULink>
  </li>
</template>

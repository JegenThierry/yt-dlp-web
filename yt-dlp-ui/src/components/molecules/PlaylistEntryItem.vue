<script setup lang="ts">
import UBadge from '@nuxt/ui/components/Badge.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import StatusBadge from '@/components/atoms/badges/StatusBadge.vue'
import type { PlaylistEntry } from '@/models/playlistEntry'
import { formatDuration } from '@/utils/format'

interface PlaylistEntryItemProps {
  entry: PlaylistEntry
}

defineProps<PlaylistEntryItemProps>()
</script>

<template>
  <li class="flex items-center gap-4 px-4 py-3 sm:px-6">
    <img
      v-if="entry.thumbnail"
      :src="entry.thumbnail"
      alt=""
      loading="lazy"
      class="size-10 shrink-0 rounded-md object-cover"
    />
    <div v-else class="flex size-10 shrink-0 items-center justify-center rounded-md bg-elevated">
      <UIcon name="i-lucide-music" class="size-5 text-muted" />
    </div>

    <div class="min-w-0 flex-1">
      <p class="truncate text-highlighted text-sm font-medium">{{ entry.title }}</p>
      <p v-if="entry.artist" class="truncate text-muted text-xs">{{ entry.artist }}</p>
    </div>

    <p class="text-muted text-xs tabular-nums">{{ formatDuration(entry.duration) }}</p>

    <UBadge
      v-if="entry.format"
      color="neutral"
      variant="outline"
      size="sm"
      :label="entry.format"
      class="uppercase"
    />

    <StatusBadge :status="entry.status" />
  </li>
</template>

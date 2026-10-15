<script setup lang="ts">
import { computed } from 'vue'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import StatusBadge from '@/components/atoms/badges/StatusBadge.vue'
import IconButton from '@/components/molecules/buttons/IconButton.vue'
import { DownloadStatus } from '@/models/enums'
import type { PlaylistEntry } from '@/models/playlistEntry'
import { formatDuration } from '@/utils/format'

interface PlaylistEntryItemProps {
  entry: PlaylistEntry
  downloading?: boolean
  onDownload?: () => void
}

const props = withDefaults(defineProps<PlaylistEntryItemProps>(), {
  downloading: false,
})

const completed = computed(() => props.entry.status === DownloadStatus.Completed)
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

    <IconButton
      v-if="onDownload"
      :icon="completed ? 'i-lucide-folder-down' : 'i-lucide-download'"
      :ariaLabel="completed ? `Save ${entry.title}` : `Download ${entry.title}`"
      color="primary"
      variant="ghost"
      size="sm"
      :loading="downloading"
      :disabled="downloading"
      @click="onDownload"
    />
  </li>
</template>

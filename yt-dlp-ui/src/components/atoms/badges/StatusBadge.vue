<script setup lang="ts">
import { computed } from 'vue'
import UBadge from '@nuxt/ui/components/Badge.vue'
import type { BadgeColor, BadgeSize } from '@/types/badgeTypes.ts'
import { DownloadStatus } from '@/models/enums'

interface StatusBadgeProps {
  status: DownloadStatus
  size?: BadgeSize
}

const props = withDefaults(defineProps<StatusBadgeProps>(), {
  size: 'sm',
})

const statusStyles: Record<DownloadStatus, { color: BadgeColor; icon: string }> = {
  [DownloadStatus.New]: { color: 'neutral', icon: 'i-lucide-sparkles' },
  [DownloadStatus.Pending]: { color: 'warning', icon: 'i-lucide-hourglass' },
  [DownloadStatus.Downloading]: { color: 'info', icon: 'i-lucide-download' },
  [DownloadStatus.Completed]: { color: 'success', icon: 'i-lucide-circle-check' },
  [DownloadStatus.Failed]: { color: 'error', icon: 'i-lucide-circle-x' },
}

const statusStyle = computed(() => statusStyles[props.status] ?? statusStyles[DownloadStatus.New])
</script>

<template>
  <UBadge
    :color="statusStyle.color"
    variant="subtle"
    :icon="statusStyle.icon"
    :label="status"
    :size="size"
    class="capitalize"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import type { BadgeColor, BadgeSize } from '@/types/badgeTypes.ts'
import { DownloadStatus } from '@/models/enums'

interface StatusBadgeProps {
  status: DownloadStatus
  size?: BadgeSize
  compact?: boolean
  unavailable?: boolean
}

const props = withDefaults(defineProps<StatusBadgeProps>(), {
  size: 'sm',
  compact: false,
  unavailable: false,
})

interface StatusStyle {
  label: string
  color: BadgeColor
  icon: string
  text: string
}

const statusStyles: Record<DownloadStatus, StatusStyle> = {
  [DownloadStatus.New]: {
    label: 'New',
    color: 'neutral',
    icon: 'i-lucide-sparkles',
    text: 'text-muted',
  },
  [DownloadStatus.Pending]: {
    label: 'Queued',
    color: 'warning',
    icon: 'i-lucide-hourglass',
    text: 'text-warning',
  },
  [DownloadStatus.Downloading]: {
    label: 'Downloading',
    color: 'info',
    icon: 'i-lucide-download',
    text: 'text-info',
  },
  [DownloadStatus.Completed]: {
    label: 'Completed',
    color: 'success',
    icon: 'i-lucide-circle-check',
    text: 'text-success',
  },
  [DownloadStatus.Failed]: {
    label: 'Failed',
    color: 'error',
    icon: 'i-lucide-circle-x',
    text: 'text-error',
  },
}

const unavailableStyle: StatusStyle = {
  label: 'Unavailable',
  color: 'neutral',
  icon: 'i-lucide-ban',
  text: 'text-muted',
}

const statusStyle = computed(() =>
  props.unavailable
    ? unavailableStyle
    : (statusStyles[props.status] ?? statusStyles[DownloadStatus.New]),
)
</script>

<template>
  <span v-if="compact" class="inline-flex shrink-0">
    <UIcon :name="statusStyle.icon" :class="['size-3.5', statusStyle.text]" />
    <span class="sr-only">{{ statusStyle.label }}</span>
  </span>
  <UBadge
    v-else
    :color="statusStyle.color"
    variant="subtle"
    :icon="statusStyle.icon"
    :label="statusStyle.label"
    :size="size"
  />
</template>

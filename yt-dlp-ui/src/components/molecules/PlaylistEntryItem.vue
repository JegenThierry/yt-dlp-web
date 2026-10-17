<script setup lang="ts">
import { computed } from 'vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import StatusBadge from '@/components/atoms/badges/StatusBadge.vue'
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'
import IconButton from '@/components/molecules/buttons/IconButton.vue'
import { DownloadPhase, DownloadStatus } from '@/models/enums'
import type { PlaylistEntry } from '@/models/playlistEntry'
import type { ButtonColor, ButtonVariant } from '@/types/buttonTypes'
import type { EntryAction } from '@/types/entryAction'
import type { EntryProgress } from '@/types/progress'
import { failureText, isUnavailable } from '@/utils/entryFailure'
import { formatDuration } from '@/utils/format'

interface PlaylistEntryItemProps {
  entry: PlaylistEntry
  progress?: EntryProgress
  busy?: boolean
}

const props = withDefaults(defineProps<PlaylistEntryItemProps>(), {
  busy: false,
})

const emit = defineEmits<{
  action: [action: EntryAction]
}>()

interface ActionStyle {
  label: string
  icon: string
  color: ButtonColor
  variant: ButtonVariant
}

const actionStyles: Record<EntryAction, ActionStyle> = {
  download: { label: 'Download', icon: 'i-lucide-download', color: 'primary', variant: 'soft' },
  cancel: { label: 'Cancel', icon: 'i-lucide-x', color: 'neutral', variant: 'ghost' },
  save: { label: 'Save', icon: 'i-lucide-hard-drive-download', color: 'neutral', variant: 'soft' },
  retry: { label: 'Retry', icon: 'i-lucide-rotate-ccw', color: 'warning', variant: 'soft' },
}

const unavailable = computed(() => isUnavailable(props.entry))
const failure = computed(() => failureText(props.entry))

const action = computed<EntryAction | null>(() => {
  switch (props.entry.status) {
    case DownloadStatus.Pending:
    case DownloadStatus.Downloading:
      return 'cancel'
    case DownloadStatus.Completed:
      return 'save'
    case DownloadStatus.Failed:
      return unavailable.value ? null : 'retry'
    default:
      return 'download'
  }
})

const actionStyle = computed(() => (action.value ? actionStyles[action.value] : null))
const actionLabel = computed(() =>
  actionStyle.value ? `${actionStyle.value.label} ${props.entry.title}` : '',
)

function runAction() {
  if (action.value) emit('action', action.value)
}

const duration = computed(() => formatDuration(props.entry.duration))
const subtitle = computed(() =>
  [props.entry.artist, props.entry.genre, props.entry.year || ''].filter(Boolean).join(' · '),
)

const progressText = computed(() => {
  const progress = props.progress
  if (!progress || props.entry.status !== DownloadStatus.Downloading) return ''
  if (progress.phase === DownloadPhase.Converting) return 'Converting…'
  if (progress.percent === undefined) return ''

  const percent = `${Math.floor(progress.percent)}%`
  return progress.eta === undefined ? percent : `${percent} · ${formatDuration(progress.eta)} left`
})
</script>

<template>
  <li
    class="grid grid-cols-[2.5rem_minmax(0,1fr)_auto] items-center gap-x-3 px-4 py-3 sm:grid-cols-[2.5rem_minmax(0,1fr)_3.5rem_7rem_6.5rem] sm:gap-x-4 sm:px-6"
  >
    <img
      v-if="entry.thumbnail"
      :src="entry.thumbnail"
      alt=""
      loading="lazy"
      :class="['size-10 rounded-md object-cover', { 'opacity-50 grayscale': unavailable }]"
    />
    <div v-else class="flex size-10 items-center justify-center rounded-md bg-elevated">
      <UIcon name="i-lucide-music" class="size-5 text-muted" />
    </div>

    <div class="min-w-0">
      <p
        :class="[
          'line-clamp-2 text-sm font-medium break-words sm:line-clamp-1',
          unavailable ? 'text-muted' : 'text-highlighted',
        ]"
      >
        {{ entry.title }}
      </p>
      <p class="text-muted mt-0.5 flex min-w-0 items-center gap-1.5 text-xs sm:hidden">
        <StatusBadge :status="entry.status" :unavailable="unavailable" compact />
        <span v-if="progressText" class="text-info truncate tabular-nums">{{ progressText }}</span>
        <span
          v-else-if="failure"
          class="text-error line-clamp-2 min-w-0 break-words"
          :title="entry.error_message"
        >
          {{ failure }}
        </span>
        <template v-else>
          <span v-if="subtitle" class="truncate">{{ subtitle }}</span>
          <span v-if="subtitle" aria-hidden="true">·</span>
          <span class="shrink-0 tabular-nums">{{ duration }}</span>
        </template>
      </p>
      <p v-if="progressText" class="text-info hidden truncate text-xs tabular-nums sm:block">
        {{ progressText }}
      </p>
      <p
        v-else-if="failure"
        class="text-error hidden text-xs break-words sm:line-clamp-2"
        :title="entry.error_message"
      >
        {{ failure }}
      </p>
      <p v-else-if="subtitle" class="text-muted hidden truncate text-xs sm:block">
        {{ subtitle }}
      </p>
    </div>

    <p class="text-muted hidden text-right text-xs tabular-nums sm:block">{{ duration }}</p>

    <div class="hidden sm:block">
      <StatusBadge :status="entry.status" :unavailable="unavailable" />
    </div>

    <div v-if="actionStyle" class="flex justify-end">
      <IconButton
        :icon="actionStyle.icon"
        :ariaLabel="actionLabel"
        :tooltip="actionStyle.label"
        :color="actionStyle.color"
        :variant="actionStyle.variant"
        size="sm"
        :loading="busy"
        :disabled="busy"
        class="sm:hidden"
        @click="runAction"
      />
      <BaseButton
        :icon="actionStyle.icon"
        :label="actionStyle.label"
        :ariaLabel="actionLabel"
        :color="actionStyle.color"
        :variant="actionStyle.variant"
        size="sm"
        :loading="busy"
        :disabled="busy"
        class="hidden sm:inline-flex"
        @click="runAction"
      />
    </div>
  </li>
</template>

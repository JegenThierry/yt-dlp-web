<script setup lang="ts">
import { computed } from 'vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import InputText from '@/components/atoms/inputs/InputText.vue'
import IconButton from '@/components/molecules/buttons/IconButton.vue'
import { DownloadStatus } from '@/models/enums'
import type { PlaylistEntry } from '@/models/playlistEntry'
import type { DraftErrors, DraftField, EntryDraft } from '@/types/entryDraft'
import { changedFields } from '@/utils/entryDraft'
import { isUnavailable } from '@/utils/entryFailure'
import { formatDuration } from '@/utils/format'

interface EntryMetadataRowProps {
  entry: PlaylistEntry
  draft: EntryDraft
  position: number
  total: number
  errors?: DraftErrors
  sourceUrl?: string
}

const props = withDefaults(defineProps<EntryMetadataRowProps>(), {
  errors: () => ({}),
  sourceUrl: undefined,
})

const emit = defineEmits<{
  update: [field: DraftField, value: string]
  blur: [field: DraftField]
}>()

const FIELD_UI = { label: 'sm:sr-only', help: 'truncate' }

const statusText: Record<DownloadStatus, string> = {
  [DownloadStatus.New]: 'Not downloaded',
  [DownloadStatus.Pending]: 'Queued',
  [DownloadStatus.Downloading]: 'Downloading',
  [DownloadStatus.Completed]: 'Downloaded',
  [DownloadStatus.Failed]: 'Failed',
}

const unavailable = computed(() => isUnavailable(props.entry))
const changed = computed(() => new Set(changedFields(props.draft, props.entry)))
const songName = computed(
  () => props.entry.title || props.entry.source_title || `Song ${props.position}`,
)
const meta = computed(() => {
  const status = unavailable.value ? 'Unavailable' : statusText[props.entry.status]
  return `${formatDuration(props.entry.duration)} · ${status}`
})

const originals = computed(() => {
  const { source_title: title, source_artist: artist } = props.entry
  return {
    title: title && props.draft.title.trim() !== title ? title : undefined,
    artist: props.draft.artist.trim() !== artist ? artist || 'not set' : undefined,
  }
})

function update(field: DraftField, value: string | undefined) {
  emit('update', field, value ?? '')
}

function reset(field: 'title' | 'artist') {
  emit('update', field, field === 'title' ? props.entry.source_title : props.entry.source_artist)
}

function highlight(field: DraftField) {
  return changed.value.has(field) ? { highlight: true, color: 'warning' as const } : {}
}
</script>

<template>
  <fieldset
    :id="`entry-${entry.id}`"
    class="grid min-w-0 scroll-mt-28 scroll-mb-20 grid-cols-2 gap-3 px-4 py-4 sm:grid-cols-[4rem_minmax(0,3fr)_minmax(0,2fr)_minmax(0,1.25fr)_6.5rem] sm:gap-y-1.5 sm:px-6 sm:py-3"
  >
    <legend class="sr-only">Song {{ position }} of {{ total }}: {{ songName }}</legend>

    <div class="col-span-2 flex min-w-0 items-center gap-2 sm:col-span-1 sm:mt-1 sm:h-8">
      <span
        class="text-muted hidden w-6 shrink-0 text-right text-xs tabular-nums sm:block"
        aria-hidden="true"
      >
        {{ position }}
      </span>
      <img
        v-if="entry.thumbnail"
        :src="entry.thumbnail"
        alt=""
        loading="lazy"
        :class="['size-8 shrink-0 rounded object-cover', { 'opacity-50 grayscale': unavailable }]"
      />
      <div v-else class="bg-elevated flex size-8 shrink-0 items-center justify-center rounded">
        <UIcon name="i-lucide-music" class="text-muted size-4" />
      </div>
      <p class="text-muted flex min-w-0 flex-wrap items-center gap-x-1.5 text-xs sm:hidden">
        <span>Song {{ position }} of {{ total }}</span>
        <span aria-hidden="true">·</span>
        <span class="tabular-nums">{{ meta }}</span>
        <template v-if="sourceUrl">
          <span aria-hidden="true">·</span>
          <a
            :href="sourceUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hover:text-default inline-flex items-center gap-1 underline-offset-2 hover:underline"
          >
            YouTube<span class="sr-only">: {{ songName }}</span>
            <UIcon name="i-lucide-external-link" class="size-3" />
          </a>
        </template>
      </p>
    </div>

    <div class="col-span-2 min-w-0 sm:col-span-1">
      <InputText
        :model-value="draft.title"
        label="Title"
        placeholder="Song title"
        :field-ui="FIELD_UI"
        :help="originals.title && `Original: ${originals.title}`"
        :error="errors.title"
        v-bind="highlight('title')"
        data-field="title"
        autocomplete="off"
        @update:model-value="update('title', $event)"
        @blur="emit('blur', 'title')"
      >
        <template v-if="originals.title" #trailing>
          <IconButton
            icon="i-lucide-undo-2"
            :ariaLabel="`Reset title to the original: ${originals.title}`"
            tooltip="Reset to original"
            size="xs"
            @click="reset('title')"
          />
        </template>
      </InputText>
    </div>

    <div class="col-span-2 min-w-0 sm:col-span-1">
      <InputText
        :model-value="draft.artist"
        label="Artist"
        placeholder="Not set"
        :field-ui="FIELD_UI"
        :help="originals.artist && `Original: ${originals.artist}`"
        :error="errors.artist"
        v-bind="highlight('artist')"
        data-field="artist"
        autocomplete="off"
        @update:model-value="update('artist', $event)"
        @blur="emit('blur', 'artist')"
      >
        <template v-if="originals.artist" #trailing>
          <IconButton
            icon="i-lucide-undo-2"
            :ariaLabel="`Reset artist to the original: ${originals.artist}`"
            tooltip="Reset to original"
            size="xs"
            @click="reset('artist')"
          />
        </template>
      </InputText>
    </div>

    <div class="min-w-0">
      <InputText
        :model-value="draft.genre"
        label="Genre"
        placeholder="Not set"
        :field-ui="FIELD_UI"
        :error="errors.genre"
        v-bind="highlight('genre')"
        data-field="genre"
        autocomplete="off"
        @update:model-value="update('genre', $event)"
        @blur="emit('blur', 'genre')"
      />
    </div>

    <div class="min-w-0">
      <InputText
        :model-value="draft.year"
        label="Year"
        placeholder="Not set"
        :field-ui="FIELD_UI"
        :error="errors.year"
        v-bind="highlight('year')"
        data-field="year"
        inputmode="numeric"
        autocomplete="off"
        @update:model-value="update('year', $event)"
        @blur="emit('blur', 'year')"
      />
    </div>

    <p
      class="text-muted col-start-2 col-end-6 hidden min-w-0 items-center gap-x-1.5 text-xs sm:flex"
    >
      <span class="tabular-nums">{{ meta }}</span>
      <template v-if="sourceUrl">
        <span aria-hidden="true">·</span>
        <a
          :href="sourceUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="hover:text-default inline-flex items-center gap-1 underline-offset-2 hover:underline"
        >
          Open on YouTube<span class="sr-only">: {{ songName }}</span>
          <UIcon name="i-lucide-external-link" class="size-3" />
        </a>
      </template>
    </p>
  </fieldset>
</template>

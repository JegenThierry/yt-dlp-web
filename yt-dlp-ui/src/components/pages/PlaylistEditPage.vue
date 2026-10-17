<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import { onBeforeRouteLeave, useRoute } from 'vue-router'
import axios from 'axios'
import { useEventListener, useTitle } from '@vueuse/core'
import UBadge from '@nuxt/ui/components/Badge.vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import UIcon from '@nuxt/ui/components/Icon.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import { defineShortcuts, useOverlay, useToast } from '@nuxt/ui/composables'
import PageLayout from '@/components/layouts/PageLayout.vue'
import EntryMetadataRow from '@/components/molecules/EntryMetadataRow.vue'
import MetadataBulkActions from '@/components/molecules/MetadataBulkActions.vue'
import UnsavedChangesModal from '@/components/molecules/UnsavedChangesModal.vue'
import { apiFieldErrors, getApiErrorMessage } from '@/api/apiErrors'
import { fetchPlaylist, updatePlaylistMetadata } from '@/api/playlistService'
import type { Playlist } from '@/models/playlist'
import type { PlaylistEntry, PlaylistEntryMetadataUpdate } from '@/models/playlistEntry'
import type { DraftEdit, DraftErrors, DraftField, EntryDraft } from '@/types/entryDraft'
import {
  DRAFT_FIELDS,
  applyEdits,
  changedFields,
  draftErrors,
  revertEdits,
  savedValue,
  toDraft,
  toUpdate,
} from '@/utils/entryDraft'
import { countLabel } from '@/utils/format'
import { entrySourceUrl } from '@/utils/source'
import { pageTitle } from '@/utils/title'

interface FieldError {
  key: string
  id: number
  field: DraftField
  message: string
  server: boolean
}

interface ServerError {
  value: string
  message: string
}

const UNDO_TOAST_DURATION_MS = 10_000

const route = useRoute('/playlists/[id]/edit')
const toast = useToast()
const leaveModal = useOverlay().create(UnsavedChangesModal)

const playlist = ref<Playlist | null>(null)
const drafts = ref<EntryDraft[]>([])
const loading = ref(true)
const saving = ref(false)
const saved = ref(false)
const error = ref<string | null>(null)
const errorRetryable = ref(false)

const touched = reactive(new Set<string>())
const showAllErrors = ref(false)
const serverErrors = reactive(new Map<string, ServerError>())

const playlistId = Number(route.params.id)

useTitle(() => pageTitle(playlist.value ? `Edit: ${playlist.value.title}` : 'Edit metadata'), {
  restoreOnUnmount: false,
})

async function loadPlaylist() {
  if (!Number.isInteger(playlistId) || playlistId < 1) {
    error.value = 'This playlist url is not valid.'
    errorRetryable.value = false
    loading.value = false
    return
  }

  loading.value = true
  error.value = null
  errorRetryable.value = false
  try {
    const data = await fetchPlaylist(playlistId)
    playlist.value = data
    drafts.value = (data.entries ?? []).map(toDraft)
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

function fieldKey(id: number, field: DraftField): string {
  return `${id}:${field}`
}

function isDraftField(field: string): field is DraftField {
  return (DRAFT_FIELDS as string[]).includes(field)
}

const entriesById = computed(
  () => new Map((playlist.value?.entries ?? []).map((entry) => [entry.id, entry])),
)

const rows = computed(() =>
  drafts.value.flatMap((draft) => {
    const entry = entriesById.value.get(draft.id)
    return entry ? [{ draft, entry }] : []
  }),
)

const changes = computed(() =>
  rows.value.flatMap(({ draft, entry }) => {
    const fields = changedFields(draft, entry)
    return fields.length > 0 ? [{ draft, entry, fields }] : []
  }),
)

const changeCount = computed(() =>
  changes.value.reduce((sum, { fields }) => sum + fields.length, 0),
)

const fieldErrors = computed<FieldError[]>(() =>
  changes.value.flatMap(({ draft }) => {
    const local = draftErrors(draft)
    return DRAFT_FIELDS.flatMap((field) => {
      const key = fieldKey(draft.id, field)
      const server = serverErrors.get(key)
      const serverMessage = server?.value === draft[field] ? server.message : undefined
      const message = local[field] ?? serverMessage
      if (!message) return []
      return [{ key, id: draft.id, field, message, server: !local[field] }]
    })
  }),
)

const shownErrors = computed(() =>
  fieldErrors.value.filter((err) => err.server || showAllErrors.value || touched.has(err.key)),
)

const errorsById = computed(() => {
  const result: Record<number, DraftErrors> = {}
  for (const err of shownErrors.value) (result[err.id] ??= {})[err.field] = err.message
  return result
})

const changeSummary = computed(
  () =>
    `${countLabel(changeCount.value, 'unsaved change', 'unsaved changes')} in ${countLabel(changes.value.length, 'song', 'songs')}`,
)
const saveLabel = computed(() =>
  changeCount.value > 0 ? `Save ${countLabel(changeCount.value, 'change', 'changes')}` : 'Save',
)
const canSave = computed(
  () => !saving.value && changeCount.value > 0 && shownErrors.value.length === 0,
)

function updateField(draft: EntryDraft, field: DraftField, value: string) {
  draft[field] = value
}

async function focusField(id: number, field: DraftField) {
  await nextTick()
  const input = document.querySelector<HTMLInputElement>(`#entry-${id} [data-field="${field}"]`)
  if (!input) return
  input.scrollIntoView({ block: 'center' })
  input.focus({ preventScroll: true })
}

function focusFirstError() {
  const first = shownErrors.value[0] ?? fieldErrors.value[0]
  if (first) focusField(first.id, first.field)
}

async function handleSave(): Promise<boolean> {
  if (saving.value) return false
  if (changes.value.length === 0) return true

  showAllErrors.value = true
  if (fieldErrors.value.length > 0) {
    focusFirstError()
    return false
  }

  const updates = changes.value.map(({ draft }) => toUpdate(draft))
  saving.value = true
  try {
    const updated = await updatePlaylistMetadata(playlistId, updates)
    applySaved(updated, updates)
    showAllErrors.value = false
    touched.clear()
    serverErrors.clear()
    saved.value = true
    toast.add({
      title: 'Changes saved',
      description: `Updated ${countLabel(updated.length, 'song', 'songs')}. Files you save from now on carry the new tags.`,
      color: 'success',
      icon: 'i-lucide-circle-check',
    })
    return true
  } catch (err) {
    showSaveError(err)
    return false
  } finally {
    saving.value = false
  }
}

function applySaved(updated: PlaylistEntry[], sent: PlaylistEntryMetadataUpdate[]) {
  if (!playlist.value) return

  const updatedById = new Map(updated.map((entry) => [entry.id, entry]))
  const sentById = new Map(sent.map((update) => [update.id, update]))
  playlist.value = {
    ...playlist.value,
    entries: (playlist.value.entries ?? []).map((entry) => updatedById.get(entry.id) ?? entry),
  }

  for (const draft of drafts.value) {
    const entry = updatedById.get(draft.id)
    const update = sentById.get(draft.id)
    if (entry && update && sameUpdate(toUpdate(draft), update)) Object.assign(draft, toDraft(entry))
  }
}

function sameUpdate(a: PlaylistEntryMetadataUpdate, b: PlaylistEntryMetadataUpdate): boolean {
  return a.title === b.title && a.artist === b.artist && a.genre === b.genre && a.year === b.year
}

function showSaveError(err: unknown) {
  const fields = apiFieldErrors(err).filter((field) => isDraftField(field.field))
  if (fields.length > 0) {
    for (const { entry_id, field, message } of fields) {
      const draft = drafts.value.find((candidate) => candidate.id === entry_id)
      if (draft && isDraftField(field)) {
        serverErrors.set(fieldKey(entry_id, field), { value: draft[field], message })
      }
    }
    focusFirstError()
    return
  }

  console.error('Failed to save metadata:', err)
  toast.add({
    title: "Couldn't save your changes",
    description: getApiErrorMessage(
      err,
      'Something went wrong while saving. Your changes are still here.',
    ),
    color: 'error',
    duration: 0,
  })
}

function applyDraftEdits(edits: DraftEdit[], title: string, description: string) {
  applyEdits(drafts.value, edits)
  toast.add({
    title,
    description,
    icon: 'i-lucide-wand-sparkles',
    duration: UNDO_TOAST_DURATION_MS,
    actions: [
      {
        label: 'Undo',
        icon: 'i-lucide-undo-2',
        color: 'neutral',
        variant: 'outline',
        onClick: () => revertEdits(drafts.value, edits),
      },
    ],
  })
}

function onBulkApply(edits: DraftEdit[], summary: string) {
  if (edits.length === 0) {
    toast.add({
      title: 'Nothing to change',
      description: 'No song needed this change.',
      color: 'neutral',
      icon: 'i-lucide-info',
    })
    return
  }

  const songs = new Set(edits.map((edit) => edit.id)).size
  applyDraftEdits(
    edits,
    `Changed ${countLabel(songs, 'song', 'songs')}`,
    `${summary}. Not saved yet.`,
  )
}

function discardChanges() {
  const edits = changes.value.flatMap(({ draft, entry, fields }) =>
    fields.map((field) => ({
      id: draft.id,
      field,
      from: draft[field],
      to: savedValue(entry, field),
    })),
  )
  const count = changeCount.value
  showAllErrors.value = false
  touched.clear()
  applyDraftEdits(
    edits,
    `Discarded ${countLabel(count, 'change', 'changes')}`,
    'The fields are back to their saved values.',
  )
}

onBeforeRouteLeave(async () => {
  if (changeCount.value === 0) return true

  const hasErrors = fieldErrors.value.length > 0
  const choice = await leaveModal.open({
    changes: changeCount.value,
    songs: changes.value.length,
    canSave: !hasErrors,
  })
  if (choice === 'discard') return true
  if (choice === 'save') return await handleSave()
  if (hasErrors) showAllErrors.value = true
  return false
})

useEventListener(window, 'beforeunload', (event: BeforeUnloadEvent) => {
  if (changeCount.value > 0) event.preventDefault()
})

defineShortcuts({
  meta_s: {
    usingInput: true,
    handler: () => {
      if (playlist.value) void handleSave()
    },
  },
})
</script>

<template>
  <PageLayout :title="error ? 'Playlist unavailable' : 'Edit metadata'" :loading="loading">
    <template #navigation>
      <UButton
        :to="`/playlists/${playlistId}`"
        icon="i-lucide-arrow-left"
        color="neutral"
        variant="ghost"
        size="sm"
        label="Back to playlist"
        class="-ms-2.5 self-start"
      />
    </template>

    <template v-if="playlist && !error" #meta>
      <div class="flex w-full min-w-0 flex-col items-start gap-2">
        <UBadge
          color="neutral"
          variant="subtle"
          icon="i-lucide-list-music"
          :label="playlist.title"
          class="max-w-full"
        />
        <p class="text-muted text-sm">
          Fix titles and artists. Changes are written into the tags and file names of songs you
          save, and the original values stay available to reset to.
        </p>
      </div>
    </template>

    <UCard v-if="loading" variant="outline">
      <template #header>
        <USkeleton class="h-5 w-16" />
      </template>
      <div class="space-y-4">
        <USkeleton v-for="n in 5" :key="n" class="h-12" />
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
          <UButton
            :to="`/playlists/${playlistId}`"
            color="neutral"
            variant="ghost"
            label="Back to playlist"
          />
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

    <UCard v-else-if="playlist && drafts.length === 0" variant="outline">
      <p class="text-muted p-6 text-center text-sm">No songs in this playlist.</p>
    </UCard>

    <form v-else-if="playlist" novalidate @submit.prevent="handleSave">
      <UCard
        variant="outline"
        :ui="{
          root: 'overflow-clip',
          header: 'bg-default/95 p-0 backdrop-blur sm:sticky sm:top-0 sm:z-10 sm:p-0',
          body: 'p-0 sm:p-0 divide-y divide-default',
          footer: 'bg-default/95 sticky bottom-0 z-10 -mt-px border-t border-default backdrop-blur',
        }"
      >
        <template #header>
          <div class="flex items-center justify-between gap-3 px-4 py-3 sm:px-6">
            <h2 class="text-highlighted text-sm font-semibold">
              {{ countLabel(drafts.length, 'song', 'songs') }}
            </h2>
            <MetadataBulkActions :drafts="drafts" @apply="onBulkApply" />
          </div>
          <div
            class="text-muted border-default hidden grid-cols-[4rem_minmax(0,3fr)_minmax(0,2fr)_minmax(0,1.25fr)_6.5rem] gap-x-3 border-t px-6 py-2 text-xs font-medium sm:grid"
            aria-hidden="true"
          >
            <span class="w-6 text-right">#</span>
            <span>Title</span>
            <span>Artist</span>
            <span>Genre</span>
            <span>Year</span>
          </div>
        </template>

        <EntryMetadataRow
          v-for="({ draft, entry }, index) in rows"
          :key="draft.id"
          :entry="entry"
          :draft="draft"
          :position="index + 1"
          :total="rows.length"
          :errors="errorsById[draft.id]"
          :source-url="entrySourceUrl(playlist.url, entry.video_id)"
          @update="(field, value) => updateField(draft, field, value)"
          @blur="(field) => touched.add(fieldKey(draft.id, field))"
        />

        <template #footer>
          <div class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2">
            <div class="flex min-w-0 flex-wrap items-center gap-x-3 gap-y-1 text-sm">
              <p aria-live="polite">
                <span v-if="changeCount > 0" class="text-default">{{ changeSummary }}</span>
                <span v-else-if="saved" class="text-success inline-flex items-center gap-1.5">
                  <UIcon name="i-lucide-circle-check" class="size-4" />
                  All changes saved
                </span>
                <span v-else class="text-muted">No unsaved changes</span>
              </p>
              <UButton
                v-if="shownErrors.length > 0"
                color="error"
                variant="link"
                icon="i-lucide-circle-alert"
                class="p-0"
                :label="`Fix ${countLabel(shownErrors.length, 'error', 'errors')} to save`"
                @click="focusFirstError"
              />
            </div>
            <div class="ms-auto flex items-center gap-2">
              <UButton
                v-if="changeCount > 0"
                color="neutral"
                variant="ghost"
                label="Discard"
                :disabled="saving"
                @click="discardChanges"
              />
              <UButton
                type="submit"
                icon="i-lucide-save"
                :color="changeCount > 0 ? 'primary' : 'neutral'"
                :variant="changeCount > 0 ? 'solid' : 'outline'"
                :label="saveLabel"
                :loading="saving"
                :disabled="!canSave"
              />
            </div>
          </div>
        </template>
      </UCard>
    </form>
  </PageLayout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { DropdownMenuItem } from '@nuxt/ui'
import UButton from '@nuxt/ui/components/Button.vue'
import UDropdownMenu from '@nuxt/ui/components/DropdownMenu.vue'
import UModal from '@nuxt/ui/components/Modal.vue'
import InputSelect from '@/components/atoms/inputs/InputSelect.vue'
import InputText from '@/components/atoms/inputs/InputText.vue'
import type { DraftEdit, EntryDraft } from '@/types/entryDraft'
import { yearError } from '@/utils/entryDraft'
import { countLabel } from '@/utils/format'
import { type CleanupAction, cleanupEdits, setAllEdits } from '@/utils/metadataCleanup'

interface MetadataBulkActionsProps {
  drafts: EntryDraft[]
}

const props = defineProps<MetadataBulkActionsProps>()

const emit = defineEmits<{
  apply: [edits: DraftEdit[], summary: string]
}>()

type SetField = 'artist' | 'genre' | 'year'

const fieldOptions: { label: string; value: SetField }[] = [
  { label: 'Artist', value: 'artist' },
  { label: 'Genre', value: 'genre' },
  { label: 'Year', value: 'year' },
]

const setAllOpen = ref(false)
const setField = ref<string>('genre')
const setValue = ref('')
const setError = ref<string>()

const cleanups: { action: CleanupAction; label: string; icon: string; summary: string }[] = [
  {
    action: 'remove-topic',
    label: 'Remove " - Topic" from artists',
    icon: 'i-lucide-user-round-x',
    summary: 'Removed " - Topic" from artists',
  },
  {
    action: 'split-artist-title',
    label: 'Split "Artist - Title" titles',
    icon: 'i-lucide-split',
    summary: 'Split "Artist - Title" titles',
  },
  {
    action: 'strip-suffixes',
    label: 'Strip "(Official Video)" and "[Lyrics]"',
    icon: 'i-lucide-eraser',
    summary: 'Stripped video suffixes from titles',
  },
  {
    action: 'trim-spaces',
    label: 'Trim extra spaces',
    icon: 'i-lucide-space',
    summary: 'Trimmed extra spaces',
  },
]

const items = computed<DropdownMenuItem[][]>(() => [
  [
    {
      label: 'Set artist, genre or year for all songs…',
      icon: 'i-lucide-pen-line',
      onSelect: openSetAll,
    },
  ],
  cleanups.map(({ action, label, icon, summary }) => ({
    label,
    icon,
    onSelect: () => emit('apply', cleanupEdits(props.drafts, action), summary),
  })),
])

const applyLabel = computed(() => `Apply to ${countLabel(props.drafts.length, 'song', 'songs')}`)

function openSetAll() {
  setValue.value = ''
  setError.value = undefined
  setAllOpen.value = true
}

function applySetAll() {
  const field = fieldOptions.find((option) => option.value === setField.value)?.value ?? 'genre'
  const value = setValue.value.trim()
  if (field === 'year') {
    setError.value = yearError(value)
    if (setError.value) return
  }

  const summary = value === '' ? `Cleared the ${field} of every song` : `Set ${field} to "${value}"`
  emit('apply', setAllEdits(props.drafts, field, value), summary)
  setAllOpen.value = false
}
</script>

<template>
  <UDropdownMenu :items="items" :content="{ align: 'end' }">
    <UButton
      icon="i-lucide-wand-sparkles"
      trailing-icon="i-lucide-chevron-down"
      color="neutral"
      variant="outline"
      size="sm"
      label="Bulk edit"
    />
  </UDropdownMenu>

  <UModal
    v-model:open="setAllOpen"
    title="Set a field for every song"
    description="This changes your unsaved edits only. Nothing is saved until you press Save."
  >
    <template #body>
      <form
        id="set-all-form"
        class="grid gap-4 sm:grid-cols-[10rem_minmax(0,1fr)]"
        novalidate
        @submit.prevent="applySetAll"
      >
        <InputSelect
          v-model="setField"
          label="Field"
          :options="fieldOptions"
          @update:model-value="setError = undefined"
        />
        <InputText
          v-model="setValue"
          label="New value"
          placeholder="Leave empty to clear"
          :error="setError"
          :inputmode="setField === 'year' ? 'numeric' : undefined"
          autocomplete="off"
          @update:model-value="setError = undefined"
        />
      </form>
    </template>
    <template #footer>
      <div class="flex w-full justify-end gap-2">
        <UButton color="neutral" variant="outline" label="Cancel" @click="setAllOpen = false" />
        <UButton type="submit" form="set-all-form" icon="i-lucide-check" :label="applyLabel" />
      </div>
    </template>
  </UModal>
</template>

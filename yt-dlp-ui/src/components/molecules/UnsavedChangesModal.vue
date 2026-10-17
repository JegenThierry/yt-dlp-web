<script setup lang="ts">
import { computed } from 'vue'
import UButton from '@nuxt/ui/components/Button.vue'
import UModal from '@nuxt/ui/components/Modal.vue'
import { countLabel } from '@/utils/format'

interface UnsavedChangesModalProps {
  changes: number
  songs: number
  canSave: boolean
}

const props = defineProps<UnsavedChangesModalProps>()

const emit = defineEmits<{
  close: [choice: 'save' | 'discard' | 'stay']
}>()

let leaving = false

const description = computed(() => {
  const summary = `You have ${countLabel(props.changes, 'unsaved change', 'unsaved changes')} in ${countLabel(props.songs, 'song', 'songs')}.`
  return props.canSave
    ? summary
    : `${summary} Some fields need fixing before they can be saved, so keep editing or discard them.`
})

function choose(choice: 'save' | 'discard' | 'stay') {
  leaving = choice !== 'stay'
  emit('close', choice)
}

function preventFocusRestore(event: Event) {
  if (leaving) event.preventDefault()
}
</script>

<template>
  <UModal
    title="Leave without saving?"
    :description="description"
    :close="false"
    :content="{ onCloseAutoFocus: preventFocusRestore }"
    :ui="{ footer: 'flex-wrap justify-end' }"
  >
    <template #footer>
      <UButton color="neutral" variant="outline" label="Keep editing" @click="choose('stay')" />
      <UButton color="error" variant="soft" label="Discard changes" @click="choose('discard')" />
      <UButton v-if="canSave" icon="i-lucide-save" label="Save and leave" @click="choose('save')" />
    </template>
  </UModal>
</template>

<script setup lang="ts">
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'
import type { ButtonType } from '@/types/buttonTypes.ts'

interface DownloadButtonProps {
  disabled?: boolean
  loading?: boolean
  type?: ButtonType
}

withDefaults(defineProps<DownloadButtonProps>(), {
  type: 'button',
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

function onClick(event: MouseEvent) {
  emit('click', event)
}
</script>

<template>
  <BaseButton
    color="primary"
    variant="solid"
    label="Download"
    icon="i-lucide-download"
    size="md"
    :type="type"
    :loading="loading"
    :disabled="disabled"
    @click="onClick"
  />
</template>

<style scoped></style>

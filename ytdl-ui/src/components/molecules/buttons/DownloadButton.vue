<script setup lang="ts">
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'

interface DownloadButtonProps {
  disabled?: boolean
  loading?: boolean
}

defineProps<DownloadButtonProps>()

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
    :loading="loading"
    :disabled="disabled"
    @click="onClick"
  />
</template>

<style scoped></style>

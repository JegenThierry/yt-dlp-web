<script setup lang="ts">
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'

interface CancelButtonProps {
  disabled?: boolean
  loading?: boolean
}

defineProps<CancelButtonProps>()

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

function onClick(event: MouseEvent) {
  emit('click', event)
}
</script>

<template>
  <BaseButton
    color="error"
    variant="soft"
    label="Cancel"
    icon="i-lucide-x"
    size="md"
    :disabled="disabled"
    :loading="loading"
    @click="onClick"
  />
</template>

<style scoped></style>

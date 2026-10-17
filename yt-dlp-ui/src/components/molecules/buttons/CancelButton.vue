<script setup lang="ts">
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'

interface CancelButtonProps {
  disabled?: boolean
  loading?: boolean
  label?: string
}

withDefaults(defineProps<CancelButtonProps>(), {
  label: 'Cancel',
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
    color="error"
    variant="soft"
    :label="label"
    icon="i-lucide-x"
    size="md"
    :disabled="disabled"
    :loading="loading"
    @click="onClick"
  />
</template>

<style scoped></style>

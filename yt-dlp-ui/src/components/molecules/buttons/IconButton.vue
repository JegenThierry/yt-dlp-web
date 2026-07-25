<script setup lang="ts">
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'
import type { ButtonColor, ButtonSize, ButtonVariant } from '@/types/buttonTypes.ts'

interface IconButtonProps {
  icon: string
  ariaLabel: string
  color?: ButtonColor
  variant?: ButtonVariant
  size?: ButtonSize
  disabled?: boolean
  loading?: boolean
}

withDefaults(defineProps<IconButtonProps>(), {
  color: 'neutral',
  variant: 'ghost',
  size: 'md',
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
    :icon="icon"
    :color="color"
    :variant="variant"
    :size="size"
    :disabled="disabled"
    :loading="loading"
    :aria-label="ariaLabel"
    @click="onClick"
  />
</template>

<style scoped></style>

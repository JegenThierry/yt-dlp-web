<script setup lang="ts">
import UButton from '@nuxt/ui/components/Button.vue'
import type { ButtonColor, ButtonSize, ButtonType, ButtonVariant } from '@/types/buttonTypes.ts'

interface BaseButtonProps {
  color?: ButtonColor
  label?: string
  icon?: string
  disabled?: boolean
  size?: ButtonSize
  variant?: ButtonVariant
  loading?: boolean
  ariaLabel?: string
  ariaExpanded?: boolean
  type?: ButtonType
}

withDefaults(defineProps<BaseButtonProps>(), {
  color: 'primary',
  size: 'md',
  variant: 'solid',
  ariaExpanded: undefined,
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
  <UButton
    :variant="variant"
    :color="color"
    :label="label"
    :loading="loading"
    :size="size"
    :disabled="disabled"
    :icon="icon"
    :aria-label="ariaLabel"
    :aria-expanded="ariaExpanded"
    :type="type"
    @click="onClick"
  >
    <slot v-if="!label" />
  </UButton>
</template>

<style scoped></style>

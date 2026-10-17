<script setup lang="ts">
import UTooltip from '@nuxt/ui/components/Tooltip.vue'
import BaseButton from '@/components/atoms/buttons/BaseButton.vue'
import type { ButtonColor, ButtonSize, ButtonVariant } from '@/types/buttonTypes.ts'

interface IconButtonProps {
  icon: string
  ariaLabel: string
  tooltip?: string
  ariaExpanded?: boolean
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
  ariaExpanded: undefined,
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

function onClick(event: MouseEvent) {
  emit('click', event)
}
</script>

<template>
  <UTooltip :text="tooltip ?? ariaLabel">
    <BaseButton
      :icon="icon"
      :color="color"
      :variant="variant"
      :size="size"
      :disabled="disabled"
      :loading="loading"
      :aria-label="ariaLabel"
      :aria-expanded="ariaExpanded"
      @click="onClick"
    />
  </UTooltip>
</template>

<style scoped></style>

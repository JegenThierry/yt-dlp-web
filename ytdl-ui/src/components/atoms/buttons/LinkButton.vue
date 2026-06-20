<script setup lang="ts">
import { NButton, type ButtonSize, type ButtonType } from 'naive-ui'

defineProps<{
  href: string
  type?: ButtonType
  size?: ButtonSize
}>();

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()

function onClick(event: MouseEvent) {
  emit('click', event)
}
</script>

<template>
  <n-button tag="a"
            text
            :type="type"
            :size="size"
            :href="href"
            @click="onClick">
    <slot />
  </n-button>
</template>

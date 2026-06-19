<script setup lang="ts">
import { NText } from 'naive-ui'

defineProps<{
  size?: 'small' | 'medium' | 'large'
}>()

const sizeMap = {
  small: '1rem',
  medium: '1.25rem',
  large: '1.5rem',
} as const
</script>

<template>
  <RouterLink to="/" class="app-logo">
    <n-text
      :style="{ fontSize: sizeMap[size ?? 'medium'], fontWeight: 700 }"
    >
      YT-DL UI
    </n-text>
  </RouterLink>
</template>

<style scoped>
.app-logo {
  text-decoration: none;
}
</style>

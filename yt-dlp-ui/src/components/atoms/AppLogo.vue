<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  size?: 'small' | 'medium' | 'large'
}>()

const textClass = computed(() =>
  props.size === 'small' ? 'text-base' : props.size === 'large' ? 'text-2xl' : 'text-xl',
)
</script>

<template>
  <RouterLink to="/" class="no-underline font-bold" :class="textClass">
    YT-DL UI
  </RouterLink>
</template>

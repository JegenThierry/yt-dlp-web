<script setup lang="ts">
import { computed } from 'vue'
import { APP_NAME } from '@/utils/title'

const props = defineProps<{
  size?: 'small' | 'medium' | 'large'
}>()

const textClass = computed(() =>
  props.size === 'small' ? 'text-base' : props.size === 'large' ? 'text-2xl' : 'text-xl',
)
</script>

<template>
  <RouterLink to="/" class="no-underline font-bold" :class="textClass">{{ APP_NAME }}</RouterLink>
</template>

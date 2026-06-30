<script setup lang="ts">
import { useColorMode } from '@vueuse/core'
import UButton from '@nuxt/ui/components/Button.vue'

const colorMode = useColorMode()

function toggle() {
  colorMode.value = colorMode.value === 'dark' ? 'light' : 'dark'
}
</script>

<template>
  <UButton
    :icon="colorMode === 'dark' ? 'i-lucide-sun' : 'i-lucide-moon'"
    color="neutral"
    variant="ghost"
    aria-label="Toggle theme"
    @click="toggle"
  />
</template>

<script setup lang="ts">
import { useColorMode } from '@vueuse/core'
import IconButton from '@/components/molecules/buttons/IconButton.vue'

const colorMode = useColorMode()

function toggle() {
  colorMode.value = colorMode.value === 'dark' ? 'light' : 'dark'
}
</script>

<template>
  <IconButton
    :icon="colorMode === 'dark' ? 'i-lucide-sun' : 'i-lucide-moon'"
    color="neutral"
    variant="ghost"
    ariaLabel="Toggle theme"
    @click="toggle"
  />
</template>

<script setup lang="ts">
import UCard from '@nuxt/ui/components/Card.vue'

interface ContentContainerProps {
  variant?: 'subtle' | 'outline' | 'solid' | 'soft'
}

withDefaults(defineProps<ContentContainerProps>(), {
  variant: 'subtle',
})

defineSlots<{
  default(): unknown
  header(): unknown
  footer(): unknown
}>()
</script>

<template>
  <div class="mx-auto w-full max-w-7xl">
    <UCard :variant="variant">
      <template v-if="$slots.header" #header>
        <slot name="header" />
      </template>

      <slot />

      <template v-if="$slots.footer" #footer>
        <slot name="footer" />
      </template>
    </UCard>
  </div>
</template>

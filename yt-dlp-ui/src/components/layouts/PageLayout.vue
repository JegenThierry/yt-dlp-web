<script setup lang="ts">
import UContainer from '@nuxt/ui/components/Container.vue'

interface PageLayoutProps {
  title?: string
}

defineProps<PageLayoutProps>()

defineSlots<{
  navigation(): unknown
  title(): unknown
  meta(): unknown
  actions(): unknown
  default(): unknown
}>()
</script>

<template>
  <UContainer class="max-w-5xl flex flex-col gap-6">
    <slot name="navigation" />

    <header
      v-if="title || $slots.title"
      class="flex flex-wrap items-start justify-between gap-x-4 gap-y-3"
    >
      <div class="min-w-0">
        <h1 class="text-highlighted text-pretty text-2xl font-semibold">
          <slot name="title">{{ title }}</slot>
        </h1>
        <div v-if="$slots.meta" class="mt-3 flex flex-wrap items-center gap-2">
          <slot name="meta" />
        </div>
      </div>

      <div v-if="$slots.actions" class="flex shrink-0 items-center gap-2">
        <slot name="actions" />
      </div>
    </header>

    <slot />
  </UContainer>
</template>

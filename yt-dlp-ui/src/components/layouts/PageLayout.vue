<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useResizeObserver } from '@vueuse/core'
import UContainer from '@nuxt/ui/components/Container.vue'
import USkeleton from '@nuxt/ui/components/Skeleton.vue'
import UTooltip from '@nuxt/ui/components/Tooltip.vue'

interface PageLayoutProps {
  title?: string
  loading?: boolean
}

const props = defineProps<PageLayoutProps>()

defineSlots<{
  navigation(): unknown
  meta(): unknown
  actions(): unknown
  default(): unknown
}>()

const titleEl = ref<HTMLElement | null>(null)
const titleClamped = ref(false)

function measureTitle() {
  const el = titleEl.value
  titleClamped.value =
    !!el && (el.scrollHeight > el.clientHeight || el.scrollWidth > el.clientWidth)
}

useResizeObserver(titleEl, measureTitle)
watch(
  () => [props.title, props.loading],
  () => nextTick(measureTitle),
)
</script>

<template>
  <UContainer class="max-w-5xl flex flex-col gap-6">
    <slot name="navigation" />

    <header v-if="title || loading" class="flex flex-col gap-3">
      <div class="flex flex-wrap items-start justify-between gap-x-4 gap-y-3">
        <h1
          ref="titleEl"
          tabindex="-1"
          class="text-highlighted line-clamp-2 min-w-0 flex-[1_1_16rem] text-2xl font-semibold break-words outline-none sm:line-clamp-1"
        >
          <USkeleton v-if="loading" as="span" class="block h-8 w-2/3 sm:max-w-md" />
          <UTooltip v-else :text="title" :disabled="!titleClamped">
            <span>{{ title }}</span>
          </UTooltip>
        </h1>

        <div v-if="$slots.actions" class="flex shrink-0 flex-wrap items-center gap-2">
          <slot name="actions" />
        </div>
      </div>

      <div v-if="$slots.meta" class="flex flex-wrap items-center gap-2">
        <slot name="meta" />
      </div>
    </header>

    <slot />
  </UContainer>
</template>

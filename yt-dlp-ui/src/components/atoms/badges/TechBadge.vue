<script setup lang="ts">
import UBadge from '@nuxt/ui/components/Badge.vue'
import type { BadgeSize } from '@/types/badgeTypes'
import type { Technology } from '@/types/techStack'

interface TechBadgeProps {
  tech: Technology
  size?: BadgeSize
}

withDefaults(defineProps<TechBadgeProps>(), {
  size: 'sm',
})
</script>

<template>
  <UBadge :icon="tech.icon" :label="tech.name" :size="size" color="neutral" variant="outline" />
</template>

<script setup lang="ts">
import TechBadge from '@/components/atoms/badges/TechBadge.vue'
import type { TechGroup } from '@/types/techStack'

defineProps<{
  groups: TechGroup[]
}>()
</script>

<template>
  <dl class="flex flex-col gap-4 sm:gap-3">
    <div
      v-for="group in groups"
      :key="group.title"
      class="grid gap-1.5 sm:grid-cols-[7rem_minmax(0,1fr)] sm:gap-6"
    >
      <dt class="text-muted text-sm sm:pt-0.5">{{ group.title }}</dt>
      <dd>
        <ul class="flex flex-wrap gap-1.5">
          <li v-for="tech in group.items" :key="tech.name">
            <TechBadge :tech="tech" size="md" />
          </li>
        </ul>
      </dd>
    </div>
  </dl>
</template>

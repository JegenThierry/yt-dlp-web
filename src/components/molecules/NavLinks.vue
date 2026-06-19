<script setup lang="ts">
import { NSpace } from 'naive-ui'
import NavLink from '@/components/atoms/NavLink.vue'

interface NavItem {
  to: string
  label: string
}

defineProps<{
  items: NavItem[]
}>()
</script>

<template>
  <n-space align="center" :size="0">
    <NavLink v-for="item in items"
             :key="item.to"
             :to="item.to"
             :label="item.label" />
  </n-space>
</template>

<script setup lang="ts">
import { NLayoutHeader, NSpace, NButton } from 'naive-ui'
import { SunnyOutline, MoonOutline } from '@vicons/ionicons5'
import AppLogo from '@/components/atoms/AppLogo.vue'
import NavLinks from '@/components/molecules/NavLinks.vue'
import { useThemeStore } from '@/stores/theme'

interface NavItem {
  to: string
  label: string
}

defineProps<{
  navItems: NavItem[]
}>()

const themeStore = useThemeStore()
</script>

<template>
  <n-layout-header bordered class="top-navbar">
    <n-space justify="space-between" align="center">
      <AppLogo />
      <n-space align="center">
        <NavLinks :items="navItems" />
        <n-button quaternary circle @click="themeStore.toggle">
          <template #icon>
            <SunnyOutline v-if="themeStore.state.isDark" />
            <MoonOutline v-else />
          </template>
        </n-button>
      </n-space>
    </n-space>
  </n-layout-header>
</template>

<style scoped>
.top-navbar {
  padding: 1rem;
}
</style>

<script setup lang="ts">
import { ref } from 'vue'
import type { NavigationMenuItem } from '@nuxt/ui'
import UNavigationMenu from '@nuxt/ui/components/NavigationMenu.vue'
import USidebar from '@nuxt/ui/components/Sidebar.vue'
import AppLogo from '@/components/atoms/AppLogo.vue'
import ThemeToggleButton from '@/components/molecules/buttons/ThemeToggleButton.vue'
import IconButton from '@/components/molecules/buttons/IconButton.vue'

const sidebarOpen = ref(true)

const navItems: NavigationMenuItem[] = [
  { label: 'Home', icon: 'i-lucide-house', to: '/' },
  { label: 'About', icon: 'i-lucide-info', to: '/about' },
  { label: 'Playlists', icon: 'i-lucide-info', to: '/playlists' },
]
</script>

<template>
  <div class="flex h-dvh bg-neutral-50 dark:bg-neutral-950">
    <USidebar
      v-model:open="sidebarOpen"
      variant="inset"
      collapsible="offcanvas"
      :ui="{ container: 'h-full' }"
    >
      <template #header>
        <AppLogo size="small" />
      </template>

      <UNavigationMenu
        :items="navItems"
        orientation="vertical"
        :ui="{ link: 'p-1.5 overflow-hidden' }"
      />

      <template #footer>
        <div class="flex justify-center p-2">
          <ThemeToggleButton />
        </div>
      </template>
    </USidebar>

    <div
      class="flex-1 flex flex-col overflow-hidden peer-data-[variant=inset]:m-4 peer-data-[variant=inset]:rounded-xl peer-data-[variant=inset]:shadow-sm peer-data-[variant=inset]:ring peer-data-[variant=inset]:ring-default bg-default"
    >
      <div class="h-(--ui-header-height) shrink-0 flex items-center px-4 border-b border-default">
        <IconButton
          icon="i-lucide-panel-left"
          color="neutral"
          variant="ghost"
          ariaLabel="Toggle sidebar"
          @click="sidebarOpen = !sidebarOpen"
        />
      </div>

      <main class="flex-1 p-4">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMediaQuery } from '@vueuse/core'
import type { NavigationMenuItem } from '@nuxt/ui'
import UNavigationMenu from '@nuxt/ui/components/NavigationMenu.vue'
import USidebar from '@nuxt/ui/components/Sidebar.vue'
import AppLogo from '@/components/atoms/AppLogo.vue'
import ThemeToggleButton from '@/components/molecules/buttons/ThemeToggleButton.vue'
import IconButton from '@/components/molecules/buttons/IconButton.vue'

const route = useRoute()

const sidebarOpen = ref(true)

const isMobile = useMediaQuery('(max-width: 1023px)')

watch(
  () => route.fullPath,
  () => {
    if (isMobile.value) sidebarOpen.value = false
  },
)

function skipToContent() {
  document.getElementById('main-content')?.focus()
}

const navItems: NavigationMenuItem[] = [
  { label: 'Add playlist', icon: 'i-lucide-list-plus', to: '/' },
  { label: 'Playlists', icon: 'i-lucide-list-music', to: '/playlists' },
]

const footerItems: NavigationMenuItem[] = [{ label: 'About', icon: 'i-lucide-info', to: '/about' }]
</script>

<template>
  <div class="flex h-dvh bg-neutral-50 dark:bg-neutral-950">
    <a
      href="#main-content"
      class="bg-default text-highlighted ring-primary fixed top-2 left-2 z-50 -translate-y-20 rounded-md px-3 py-2 text-sm font-medium ring-2 focus:translate-y-0"
      @click.prevent="skipToContent"
    >
      Skip to content
    </a>

    <USidebar
      v-model:open="sidebarOpen"
      variant="inset"
      collapsible="offcanvas"
      :menu="{ ui: { content: 'max-w-[18rem]' } }"
      :ui="{ container: 'h-full' }"
    >
      <template #header="{ close }">
        <div class="flex flex-1 items-center justify-between gap-2">
          <AppLogo size="small" />
          <IconButton
            icon="i-lucide-x"
            color="neutral"
            variant="ghost"
            ariaLabel="Close menu"
            class="lg:hidden"
            @click="close"
          />
        </div>
      </template>

      <UNavigationMenu
        :items="navItems"
        orientation="vertical"
        :ui="{ link: 'p-1.5 overflow-hidden' }"
      />

      <template #footer>
        <div class="flex flex-1 items-center justify-between gap-2">
          <UNavigationMenu
            :items="footerItems"
            orientation="vertical"
            class="min-w-0 flex-1"
            :ui="{ link: 'p-1.5 overflow-hidden' }"
          />
          <ThemeToggleButton />
        </div>
      </template>
    </USidebar>

    <div
      class="flex-1 flex flex-col overflow-hidden lg:peer-data-[variant=inset]:m-4 lg:peer-data-[variant=inset]:rounded-xl lg:peer-data-[variant=inset]:shadow-sm lg:peer-data-[variant=inset]:ring lg:peer-data-[variant=inset]:ring-default bg-default"
    >
      <div class="h-(--ui-header-height) shrink-0 flex items-center px-4 border-b border-default">
        <IconButton
          icon="i-lucide-panel-left"
          color="neutral"
          variant="ghost"
          ariaLabel="Toggle sidebar"
          :ariaExpanded="sidebarOpen"
          @click="sidebarOpen = !sidebarOpen"
        />
      </div>

      <main id="main-content" tabindex="-1" class="flex-1 overflow-y-auto outline-none">
        <div class="py-4 sm:p-4">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  NConfigProvider,
  NGlobalStyle,
  NLayout,
  NLayoutContent,
  NMessageProvider,
  NDialogProvider,
  NNotificationProvider,
  NLoadingBarProvider,
  darkTheme,
} from 'naive-ui'
import { useThemeStore } from '@/stores/theme.ts'
import { computed } from 'vue'
import TopNavbar from '@/components/organisms/TopNavbar.vue'

const themeStore = useThemeStore()
const theme = computed(() => (themeStore.state.isDark ? darkTheme : null))
</script>

<template>
  <n-config-provider :theme="theme">
    <n-global-style />
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <n-layout>
              <TopNavbar
                :nav-items="[
                  { to: '/', label: 'Home' },
                  { to: '/about', label: 'About' },
                ]"
              />
              <n-layout-content class="layout-padding">
                <RouterView />
              </n-layout-content>
            </n-layout>
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<style scoped>
.layout-padding {
  padding: 2rem;
}
</style>

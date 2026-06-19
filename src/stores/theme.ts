import { defineStore } from 'pinia'
import { useStorage } from '@vueuse/core'

export const useThemeStore = defineStore('theme', () => {
  const state = useStorage('yt-dl-theme', {
    isDark: window.matchMedia('(prefers-color-scheme: dark)').matches,
  })

  function toggle() {
    state.value.isDark = !state.value.isDark
  }

  return { state, toggle }
})

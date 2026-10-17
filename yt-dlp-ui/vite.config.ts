import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import VueRouter from 'vue-router/vite'
import ui from '@nuxt/ui/vite'

export default defineConfig({
  plugins: [
    VueRouter({}),
    vue(),
    vueJsx(),
    vueDevTools(),
    ui({ autoImport: false, components: false }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})

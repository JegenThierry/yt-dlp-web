import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import VueRouter from 'vue-router/vite'
import ui from '@nuxt/ui/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
    VueRouter({}),
    // Disable Nuxt UI auto-import: components (`U*`) and composables are NOT
    // registered globally, so they must be imported explicitly with direct deep
    // imports, e.g. `import UButton from '@nuxt/ui/components/Button.vue'`.
    // Note: <UIcon>/<ULink> must come from the Vue-compatible runtime paths
    // (`@nuxt/ui/runtime/vue/components/Icon.vue`,
    // `@nuxt/ui/runtime/vue/overrides/vue-router/Link.vue`) instead of
    // `@nuxt/ui/components/*` to avoid pulling in Nuxt-only dependencies.
    ui({ autoImport: false, components: false }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})

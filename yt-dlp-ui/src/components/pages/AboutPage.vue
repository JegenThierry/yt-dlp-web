<script setup lang="ts">
import UButton from '@nuxt/ui/components/Button.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import ULink from '@nuxt/ui/components/Link.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import TechStack from '@/components/molecules/TechStack.vue'
import type { TechGroup } from '@/types/techStack'
import { APP_NAME } from '@/utils/title'

const SOURCE_URL = 'https://codeberg.org/thierryjegen/yt-dlp-web'
const YT_DLP_URL = 'https://github.com/yt-dlp/yt-dlp'

const groups: TechGroup[] = [
  {
    title: 'Interface',
    items: [
      { name: 'Vue', icon: 'simple-icons:vuedotjs' },
      { name: 'TypeScript', icon: 'simple-icons:typescript' },
      { name: 'Vite', icon: 'simple-icons:vite' },
      { name: 'Nuxt UI', icon: 'simple-icons:nuxt' },
      { name: 'Tailwind CSS', icon: 'simple-icons:tailwindcss' },
      { name: 'Pinia', icon: 'simple-icons:pinia' },
    ],
  },
  {
    title: 'Server',
    items: [
      { name: 'Go', icon: 'simple-icons:go' },
      { name: 'PostgreSQL', icon: 'simple-icons:postgresql' },
      { name: 'sqlc', icon: 'i-lucide-file-code' },
    ],
  },
  {
    title: 'Downloads',
    items: [
      { name: 'yt-dlp', icon: 'i-lucide-square-terminal' },
      { name: 'FFmpeg', icon: 'simple-icons:ffmpeg' },
      { name: 'Deno', icon: 'simple-icons:deno' },
    ],
  },
  {
    title: 'Hosting',
    items: [
      { name: 'Docker', icon: 'simple-icons:docker' },
      { name: 'nginx', icon: 'simple-icons:nginx' },
      { name: 'Traefik', icon: 'simple-icons:traefikproxy' },
    ],
  },
]
</script>

<template>
  <PageLayout title="About">
    <UCard variant="outline">
      <div class="flex flex-col gap-4">
        <p class="text-default text-pretty">
          {{ APP_NAME }} saves YouTube playlists and videos as audio files. Paste a link, pick a
          format, and the server downloads every song for you. It's a web interface for
          <ULink
            :to="YT_DLP_URL"
            target="_blank"
            class="text-primary hover:text-primary decoration-primary/40 hover:decoration-primary font-medium underline underline-offset-4"
            >yt-dlp<span class="sr-only"> (opens in a new tab)</span></ULink
          >.
        </p>

        <div>
          <UButton
            :to="SOURCE_URL"
            target="_blank"
            icon="simple-icons:codeberg"
            trailing-icon="i-lucide-arrow-up-right"
            label="Source code"
            color="neutral"
            variant="outline"
            size="sm"
          />
        </div>
      </div>
    </UCard>

    <UCard variant="outline" :ui="{ header: 'py-3 sm:py-3' }">
      <template #header>
        <h2 class="text-highlighted font-semibold">Built with</h2>
      </template>

      <TechStack :groups="groups" />
    </UCard>
  </PageLayout>
</template>

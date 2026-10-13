<script setup lang="ts">
import UCard from '@nuxt/ui/components/Card.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
</script>

<template>
  <PageLayout title="About">
    <UCard variant="outline">
      <p class="text-muted text-sm leading-6">
        A web frontend for the command-line tool
        <a
          href="https://github.com/yt-dlp/yt-dlp"
          target="_blank"
          rel="noopener noreferrer"
          class="text-primary underline-offset-4 hover:underline"
        >
          yt-dlp
        </a>
        .
      </p>
    </UCard>
  </PageLayout>
</template>
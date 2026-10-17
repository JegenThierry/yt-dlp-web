<script setup lang="ts">
import { useRoute } from 'vue-router'
import UButton from '@nuxt/ui/components/Button.vue'
import UEmpty from '@nuxt/ui/components/Empty.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'

const route = useRoute()
</script>

<template>
  <PageLayout title="Page not found">
    <UEmpty
      variant="naked"
      icon="i-lucide-map-pin-off"
      :description="`There is no page at ${route.path}. The link may be mistyped or out of date.`"
      class="py-16"
    >
      <template #actions>
        <UButton to="/playlists" icon="i-lucide-list-music" label="Playlists" />
        <UButton
          to="/"
          icon="i-lucide-plus"
          color="neutral"
          variant="outline"
          label="Add a playlist"
        />
      </template>
    </UEmpty>
  </PageLayout>
</template>

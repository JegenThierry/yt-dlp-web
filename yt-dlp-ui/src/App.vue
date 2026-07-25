<script setup lang="ts">
import UApp from '@nuxt/ui/components/App.vue'
import MainLayout from '@/components/layouts/MainLayout.vue'
</script>

<template>
  <UApp>
    <MainLayout>
      <RouterView />
    </MainLayout>
  </UApp>
</template>

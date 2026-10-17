<script setup lang="ts">
import UApp from '@nuxt/ui/components/App.vue'
import MainLayout from '@/components/layouts/MainLayout.vue'
</script>

<template>
  <UApp :toaster="{ position: 'top-center' }">
    <MainLayout>
      <RouterView />
    </MainLayout>
  </UApp>
</template>

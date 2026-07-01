<script setup lang="ts">
import InputText from '@/components/atoms/inputs/InputText.vue'
import { ref } from 'vue'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'

const link = ref<string>();
</script>

<template>
  <div class="max-w-180 mx-auto p-4 border border-default rounded">
    <h3 class="text-lg font-semibold mb-4">Download</h3>
    <form>
      <div class="flex flex-col gap-1 mb-4">
        <InputText class="w-full" v-model="link" label="Video URL" placeholder="https://youtube.com/..."></InputText>
      </div>
      <DownloadButton :disabled="!link" />
    </form>
  </div>
</template>

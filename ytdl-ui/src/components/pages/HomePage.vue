<script setup lang="ts">
import InputText from '@/components/atoms/inputs/InputText.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
import SubmitButton from '@/components/molecules/buttons/SubmitButton.vue'

interface DownloadForm {
  url: string
}

async function handleDownload(values: DownloadForm) {
  console.log('Downloading:', values.url)
}
</script>

<template>
  <div class="max-w-180 mx-auto p-4 border border-default rounded">
    <h3 class="text-lg font-semibold mb-4">Download</h3>
    <FormWrapper
      :initial-values="{ url: '' }"
      :on-submit="handleDownload"
      v-slot="{ state, loading }"
    >
      <InputText v-model="state.url" label="Video URL" placeholder="https://youtube.com/..." required />
      <SubmitButton :loading="loading" :disabled="!state.url" />
    </FormWrapper>
  </div>
</template>

<script setup lang="ts">
import InputText from '@/components/atoms/inputs/InputText.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'
interface DownloadForm {
  url: string
}

async function handleDownload(values: DownloadForm) {
  await new Promise<void>((resolve) => {
    setTimeout(() => {
      try {
        console.log('Downloading:', values.url)
      } catch (error) {
        console.error('Error downloading:', error)
      } finally {
        resolve()
      }
    }, 5000)
  })
}
</script>

<template>
  <div class="max-w-180 mx-auto p-4 space-y-4">
    <div class="mx-auto p-4 border border-default rounded">
      <h3 class="text-lg font-semibold mb-4">Download</h3>
      <FormWrapper
        :initial-values="{ url: '' }"
        :on-submit="handleDownload"
        v-slot="{ state, loading }"
      >
        <InputText
          v-model="state.url"
          label="Video URL"
          placeholder="https://youtube.com/..."
          required
        />
        <DownloadButton :loading="loading" :disabled="!state.url" type="submit" />
      </FormWrapper>
    </div>
  </div>
</template>

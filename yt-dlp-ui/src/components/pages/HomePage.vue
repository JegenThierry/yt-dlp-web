<script setup lang="ts">
import { createPlaylist } from '@/api/playlistService'
import InputText from '@/components/atoms/inputs/InputText.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'
import { AudioFormat } from '@/models/enums'
import type { CreatePlaylistRequest } from '@/models/playlist'

interface DownloadForm {
  url: string
  format?: AudioFormat
}

async function handleDownload(values: DownloadForm) {
  try {
    const response = await createPlaylist(buildCreatePlaylistRequest(values))
    console.log('Playlist created successfully:', response)
  } catch (error) {
    console.error('Error creating playlist:', error)
  }
}

function buildCreatePlaylistRequest(values: DownloadForm): CreatePlaylistRequest {
  return {
    url: values.url,
    target_format: values.format || AudioFormat.MP3,
  }
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

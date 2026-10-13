<script setup lang="ts">
import { createPlaylist } from '@/api/playlistService'
import InputText from '@/components/atoms/inputs/InputText.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'
import { AudioFormat } from '@/models/enums'
import type { CreatePlaylistRequest } from '@/models/playlist'
import ContentContainer from '../layouts/ContentContainer.vue'

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
  <ContentContainer>
    <template #header>
      <h3 class="text-lg font-semibold">Download</h3>
    </template>

    <FormWrapper
      :initial-values="{ url: '' }"
      :on-submit="handleDownload"
      v-slot="{ state, loading }"
    >
      <InputText
        v-model="state.url"
        label="Video URL"
        placeholder="youtube.com/video-url"
        required
      />
      <DownloadButton :loading="loading" :disabled="!state.url" type="submit" />
    </FormWrapper>
  </ContentContainer>
</template>

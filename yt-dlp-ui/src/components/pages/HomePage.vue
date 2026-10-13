<script setup lang="ts">
import UCard from '@nuxt/ui/components/Card.vue'
import { createPlaylist } from '@/api/playlistService'
import InputText from '@/components/atoms/inputs/InputText.vue'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'
import PageLayout from '@/components/layouts/PageLayout.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
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
  <PageLayout title="Download">
    <UCard variant="outline">
      <FormWrapper
        :initial-values="{ url: '' }"
        :on-submit="handleDownload"
        v-slot="{ state, loading }"
      >
        <InputText
          v-model="state.url"
          label="Video URL"
          placeholder="youtube.com/video-url"
          icon="i-lucide-link"
          required
        />
        <DownloadButton :loading="loading" :disabled="!state.url" type="submit" />
      </FormWrapper>
    </UCard>
  </PageLayout>
</template>

<script setup lang="ts">
import { createPlaylist } from '@/api/playlistService'
import InputText from '@/components/atoms/inputs/InputText.vue'
import FormWrapper from '@/components/layouts/FormWrapper.vue'
import DownloadButton from '@/components/molecules/buttons/DownloadButton.vue'
import { AudioFormat } from '@/models/enums'
import type { CreatePlaylistRequest } from '@/models/playlist'
import UPageList from "@nuxt/ui/components/PageList.vue"
import UPageCard from "@nuxt/ui/components/PageCard.vue"
import UUser from "@nuxt/ui/components/User.vue"
import { ref } from 'vue'

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
const users = ref([
  {
    name: 'Benjamin Canac',
    description: 'benjamincanac',
    to: 'https://github.com/benjamincanac',
    target: '_blank',
    avatar: {
      src: 'https://github.com/benjamincanac.png',
      alt: 'benjamincanac',
      loading: 'lazy' as const,
    },
  },
  {
    name: 'Romain Hamel',
    description: 'romhml',
    to: 'https://github.com/romhml',
    target: '_blank',
    avatar: {
      src: 'https://github.com/romhml.png',
      alt: 'romhml',
      loading: 'lazy' as const,
    },
  },
  {
    name: 'Sébastien Chopin',
    description: 'atinux',
    to: 'https://github.com/atinux',
    target: '_blank',
    avatar: {
      src: 'https://github.com/atinux.png',
      alt: 'atinux',
      loading: 'lazy' as const,
    },
  },
  {
    name: 'Hugo Richard',
    description: 'HugoRCD',
    to: 'https://github.com/HugoRCD',
    target: '_blank',
    avatar: {
      src: 'https://github.com/HugoRCD.png',
      alt: 'HugoRCD',
      loading: 'lazy' as const,
    },
  },
  {
    name: 'Sandro Circi',
    description: 'sandros94',
    to: 'https://github.com/sandros94',
    target: '_blank',
    avatar: {
      src: 'https://github.com/sandros94.png',
      alt: 'sandros94',
      loading: 'lazy' as const,
    },
  },
  {
    name: 'Daniel Roe',
    description: 'danielroe',
    to: 'https://github.com/danielroe',
    target: '_blank',
    avatar: {
      src: 'https://github.com/danielroe.png',
      alt: 'danielroe',
      loading: 'lazy' as const,
    },
  },
  {
    name: 'Jakub Michálek',
    description: 'J-Michalek',
    to: 'https://github.com/J-Michalek',
    target: '_blank',
    avatar: {
      src: 'https://github.com/J-Michalek.png',
      alt: 'J-Michalek',
      loading: 'lazy' as const,
    },
  },
  {
    name: 'Eugen Istoc',
    description: 'genu',
    to: 'https://github.com/genu',
    target: '_blank',
    avatar: {
      src: 'https://github.com/genu.png',
      alt: 'genu',
      loading: 'lazy' as const,
    },
  },
])
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
    <div class="mx-auto p-4 border border-default rounded">
      <h3 class="text-lg font-semibold mb-4">Playlist</h3>
      <UPageList divide>
        <UPageCard
          v-for="(user, index) in users"
          :key="index"
          variant="ghost"
          :to="user.to"
          :target="user.target"
        >
          <template #body>
            <UUser
              :name="user.name"
              :description="user.description"
              :avatar="user.avatar"
              size="xl"
            />
          </template>
        </UPageCard>
      </UPageList>
    </div>
  </div>
</template>

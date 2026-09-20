<script setup lang="ts">
import { onMounted, ref } from 'vue';
import UTable from '@nuxt/ui/components/Table.vue'
import UCard from '@nuxt/ui/components/Card.vue'
import { fetchPlaylists } from '@/api/playlistService';
import type { Playlist } from '@/models/playlist';
import ContentContainer from '../layouts/ContentContainer.vue';

const data = ref();

async function updatePlaylist() {
  try {
    const playlists = await fetchPlaylists();
    data.value = playlists.map((x: Playlist) => ({
      id: x.id,
      title: x.title,
      createdAt: x.created_at,
      upatedAt: x.updated_at,
    }));
  } catch (error) {
    console.error(error);
  }
}

onMounted(() => {
  updatePlaylist();
});

</script>

<template>
  <ContentContainer>
      <UTable :data="data" sticky class="max-h-[70vh]" />
  </ContentContainer>
  <div class="mx-auto w-full max-w-4xl px-4 py-6">
    <UCard :ui="{ body: 'p-0 sm:p-0' }">
    </UCard>
  </div>
</template>

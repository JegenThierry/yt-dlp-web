import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { getApiErrorMessage } from '@/api/apiErrors'
import { fetchPlaylists } from '@/api/playlistService'
import type { Playlist } from '@/models/playlist'

export const usePlaylistsStore = defineStore('playlists', () => {
  const playlists = ref<Playlist[]>([])
  const loaded = ref(false)
  const error = ref<string | null>(null)

  const hasActive = computed(() => playlists.value.some((playlist) => playlist.active_count > 0))

  let pending: Promise<void> | null = null

  function load(): Promise<void> {
    pending ??= fetchAll().finally(() => {
      pending = null
    })
    return pending
  }

  async function fetchAll() {
    error.value = null
    try {
      playlists.value = await fetchPlaylists()
      loaded.value = true
    } catch (err) {
      console.error('Failed to fetch playlists:', err)
      error.value = getApiErrorMessage(err, 'Something went wrong while loading your playlists.')
    }
  }

  function remove(id: number) {
    playlists.value = playlists.value.filter((playlist) => playlist.id !== id)
  }

  return { playlists, loaded, error, hasActive, load, remove }
})

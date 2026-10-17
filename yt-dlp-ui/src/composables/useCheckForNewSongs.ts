import { ref } from 'vue'
import { useToast } from '@nuxt/ui/composables'
import { getApiErrorMessage } from '@/api/apiErrors'
import { refreshPlaylist } from '@/api/playlistService'
import type { PlaylistImport } from '@/models/playlist'
import { countLabel } from '@/utils/format'

export function useCheckForNewSongs() {
  const toast = useToast()
  const refreshing = ref(false)

  async function checkForNewSongs(playlistId: number): Promise<PlaylistImport | undefined> {
    if (refreshing.value) return

    refreshing.value = true
    try {
      const result = await refreshPlaylist(playlistId)
      const added = result.added_count > 0
      toast.add({
        title: added
          ? `${countLabel(result.added_count, 'new song', 'new songs')} added`
          : 'No new songs',
        description:
          result.skipped_count > 0
            ? `${countLabel(result.skipped_count, 'unavailable video', 'unavailable videos')} skipped.`
            : undefined,
        color: added ? 'success' : 'neutral',
        icon: added ? 'i-lucide-list-plus' : 'i-lucide-list-checks',
      })
      return result
    } catch (err) {
      console.error('Failed to refresh playlist:', err)
      toast.add({
        title: "Couldn't check for new songs",
        description: getApiErrorMessage(err, 'Something went wrong while reading the playlist.'),
        color: 'error',
        duration: 0,
      })
    } finally {
      refreshing.value = false
    }
  }

  return { refreshing, checkForNewSongs }
}

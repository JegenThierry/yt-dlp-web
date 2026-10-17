import { onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useDocumentVisibility, useIntervalFn } from '@vueuse/core'
import { usePlaylistsStore } from '@/stores/playlists'

const POLL_INTERVAL_MS = 4_000

export function useLibraryRefresh() {
  const store = usePlaylistsStore()
  const { hasActive } = storeToRefs(store)
  const visibility = useDocumentVisibility()

  const { pause, resume } = useIntervalFn(store.load, POLL_INTERVAL_MS, { immediate: false })

  watch(
    () => hasActive.value && visibility.value === 'visible',
    (polling) => (polling ? resume() : pause()),
    { immediate: true },
  )

  watch(visibility, (state) => {
    if (state === 'visible' && hasActive.value) void store.load()
  })

  onMounted(store.load)

  return store
}

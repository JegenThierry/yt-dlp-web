import { onScopeDispose, ref, watch, type Ref } from 'vue'
import { fetchPlaylist } from '@/api/playlistService'
import { subscribePlaylistEvents } from '@/api/downloadService'
import { DownloadStatus } from '@/models/enums'
import type { ProgressEvent } from '@/models/download'
import type { Playlist } from '@/models/playlist'
import type { EntryProgress } from '@/types/progress'

const RECONNECT_DELAYS_MS = [1_000, 2_000, 5_000, 10_000, 30_000]

export function usePlaylistProgress(playlist: Ref<Playlist | null>) {
  const progress = ref<Record<number, EntryProgress>>({})

  let closeStream: (() => void) | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | undefined
  let reconnectAttempts = 0

  function stop() {
    clearTimeout(reconnectTimer)
    closeStream?.()
    closeStream = null
  }

  function start(playlistId: number) {
    stop()
    closeStream = subscribePlaylistEvents(playlistId, {
      onOpen: () => {
        reconnectAttempts = 0
        void resync(playlistId)
      },
      onEvent: applyEvent,
      onClosed: () => scheduleReconnect(playlistId),
    })
  }

  function scheduleReconnect(playlistId: number) {
    stop()
    const delay = RECONNECT_DELAYS_MS[Math.min(reconnectAttempts, RECONNECT_DELAYS_MS.length - 1)]
    reconnectAttempts++
    reconnectTimer = setTimeout(() => start(playlistId), delay)
  }

  async function resync(playlistId: number) {
    try {
      const fresh = await fetchPlaylist(playlistId)
      if (playlist.value?.id !== playlistId) return
      playlist.value = fresh
      pruneProgress()
    } catch (err) {
      console.error('Failed to refresh playlist:', err)
    }
  }

  function applyEvent(event: ProgressEvent) {
    const current = playlist.value
    if (!current || current.id !== event.playlist_id) return

    const entries = current.entries ?? []
    const entry = entries.find((e) => e.id === event.entry_id)
    if (!entry) return

    if (event.phase) {
      progress.value = {
        ...progress.value,
        [event.entry_id]: { phase: event.phase, percent: event.percent, eta: event.eta },
      }
      if (entry.status === event.status) return
    } else {
      clearProgress(event.entry_id)
    }

    playlist.value = {
      ...current,
      entries: entries.map((e) =>
        e.id === event.entry_id
          ? {
              ...e,
              status: event.status,
              file_path: event.phase ? e.file_path : (event.file_path ?? ''),
              error_message: event.error_message,
              error_code: event.error_code,
            }
          : e,
      ),
    }
  }

  function clearProgress(entryId: number) {
    if (!(entryId in progress.value)) return
    const next = { ...progress.value }
    delete next[entryId]
    progress.value = next
  }

  function pruneProgress() {
    const downloading = new Set(
      (playlist.value?.entries ?? [])
        .filter((entry) => entry.status === DownloadStatus.Downloading)
        .map((entry) => entry.id),
    )
    progress.value = Object.fromEntries(
      Object.entries(progress.value).filter(([id]) => downloading.has(Number(id))),
    )
  }

  watch(
    () => playlist.value?.id,
    (playlistId) => {
      progress.value = {}
      if (playlistId) start(playlistId)
      else stop()
    },
    { immediate: true },
  )

  onScopeDispose(stop)

  return { progress }
}

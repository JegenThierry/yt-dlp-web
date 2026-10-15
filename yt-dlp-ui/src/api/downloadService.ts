import axiosInstance from './apiClient'
import { config } from '@/config/config'
import type { FilterOptions } from '@/models/filterOptions'
import type { ProgressEvent } from '@/models/download'
import type { Playlist } from '@/models/playlist'
import type { PlaylistEntry } from '@/models/playlistEntry'

export async function downloadAll(playlistId: number, filters?: FilterOptions): Promise<Playlist> {
  const { data } = await axiosInstance.post<Playlist>(
    `/playlists/${playlistId}/download`,
    filters ? { filters } : {},
  )
  return data
}

export async function downloadEntry(
  entryId: number,
  filters?: FilterOptions,
): Promise<PlaylistEntry> {
  const { data } = await axiosInstance.post<PlaylistEntry>(
    `/playlists/entries/${entryId}/download`,
    filters ? { filters } : {},
  )
  return data
}

/**
 * Downloads a single cached entry file to the browser. Returns true on success,
 * false when no file is available for the entry yet.
 */
export async function downloadEntryFile(entry: PlaylistEntry): Promise<boolean> {
  if (!entry.file_path) return false

  const dataUrl = axiosInstance.getUri({ url: `/playlists/entries/${entry.id}/file` })
  const response = await fetch(dataUrl)

  if (!response.ok) return false

  const blob = await response.blob()
  const filename = filenameFromDisposition(response.headers.get('Content-Disposition'))
  triggerBlobDownload(blob, filename || basenameFromPath(entry.file_path))
  return true
}

/**
 * Downloads a zip of every cached file in the playlist to the browser. Returns
 * true on success, false when the playlist has no downloadable files yet.
 */
export async function downloadPlaylistZip(playlistId: number, fallbackName: string): Promise<boolean> {
  const dataUrl = axiosInstance.getUri({ url: `/playlists/${playlistId}/download/zip` })
  const response = await fetch(dataUrl)

  if (!response.ok) return false

  const blob = await response.blob()
  const filename = filenameFromDisposition(response.headers.get('Content-Disposition'))
  triggerBlobDownload(blob, filename || fallbackName)
  return true
}

/**
 * Parses the filename out of a Content-Disposition header value, or returns
 * empty string when absent.
 */
function filenameFromDisposition(value: string | null): string {
  if (!value) return ''
  const match = /filename="?([^";]+)"?/.exec(value)
  return match?.[1] ?? ''
}

function basenameFromPath(path: string): string {
  const parts = path.split(/[\\/]/)
  return parts[parts.length - 1] ?? 'download'
}

/** Triggers a client-side file download for a blob using a temporary anchor. */
function triggerBlobDownload(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Subscribes to a playlist's download progress stream via Server-Sent Events.
 * Returns a cleanup function that closes the stream and detaches the listener.
 */
export function subscribePlaylistEvents(
  playlistId: number,
  onMessage: (event: ProgressEvent) => void,
): () => void {
  const source = new EventSource(
    `${config.apiBaseUrl.replace(/\/$/, '')}/playlists/${playlistId}/events`,
  )

  const handler = (event: MessageEvent) => {
    try {
      onMessage(JSON.parse(event.data) as ProgressEvent)
    } catch (err) {
      console.error('Failed to parse progress event:', err)
    }
  }

  source.addEventListener('message', handler)

  return () => {
    source.removeEventListener('message', handler)
    source.close()
  }
}

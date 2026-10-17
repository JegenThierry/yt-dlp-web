import axiosInstance from './apiClient'
import { config } from '@/config/config'
import type { FilterOptions } from '@/models/filterOptions'
import type { ProgressEvent } from '@/models/download'
import type { PlaylistEntry } from '@/models/playlistEntry'

export class FileRequestError extends Error {
  constructor(readonly status: number) {
    super(`File request failed with status ${status}`)
    this.name = 'FileRequestError'
  }
}

export async function downloadAll(playlistId: number, filters?: FilterOptions): Promise<void> {
  await axiosInstance.post(`/playlists/${playlistId}/download`, filters ? { filters } : {})
}

export async function downloadEntry(entryId: number, filters?: FilterOptions): Promise<void> {
  await axiosInstance.post(`/playlists/entries/${entryId}/download`, filters ? { filters } : {})
}

export async function cancelPlaylistDownload(playlistId: number): Promise<void> {
  await axiosInstance.post(`/playlists/${playlistId}/download/cancel`)
}

export async function cancelEntryDownload(entryId: number): Promise<void> {
  await axiosInstance.post(`/playlists/entries/${entryId}/download/cancel`)
}

export function playlistZipUrl(playlistId: number): string {
  return axiosInstance.getUri({ url: `/playlists/${playlistId}/download/zip` })
}

export function saveFromUrl(url: string) {
  const link = document.createElement('a')
  link.href = url
  link.download = ''
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

export async function saveEntryFile(entry: PlaylistEntry): Promise<void> {
  const response = await fetch(axiosInstance.getUri({ url: `/playlists/entries/${entry.id}/file` }))
  if (!response.ok) throw new FileRequestError(response.status)

  const blob = await response.blob()
  const filename = filenameFromDisposition(response.headers.get('Content-Disposition'))
  triggerBlobDownload(blob, filename || basenameFromPath(entry.file_path))
}

function filenameFromDisposition(value: string | null): string {
  if (!value) return ''

  const encoded = /filename\*\s*=\s*UTF-8''([^;]+)/i.exec(value)?.[1]
  const decoded = encoded ? safeDecode(encoded.trim()) : ''
  if (decoded) return decoded

  return /filename\s*=\s*"?([^";]+)"?/i.exec(value)?.[1] ?? ''
}

function safeDecode(value: string): string {
  try {
    return decodeURIComponent(value)
  } catch {
    return ''
  }
}

function basenameFromPath(path: string): string {
  const parts = path.split(/[\\/]/)
  return parts[parts.length - 1] || 'download'
}

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

export interface PlaylistEventHandlers {
  onEvent: (event: ProgressEvent) => void
  onOpen?: () => void
  onClosed?: () => void
}

export function subscribePlaylistEvents(
  playlistId: number,
  handlers: PlaylistEventHandlers,
): () => void {
  const source = new EventSource(
    `${config.apiBaseUrl.replace(/\/$/, '')}/playlists/${playlistId}/events`,
  )

  const handleMessage = (event: MessageEvent) => {
    try {
      handlers.onEvent(JSON.parse(event.data) as ProgressEvent)
    } catch (err) {
      console.error('Failed to parse progress event:', err)
    }
  }

  const handleOpen = () => handlers.onOpen?.()

  const handleError = () => {
    if (source.readyState === EventSource.CLOSED) handlers.onClosed?.()
  }

  source.addEventListener('message', handleMessage)
  source.addEventListener('open', handleOpen)
  source.addEventListener('error', handleError)

  return () => {
    source.removeEventListener('message', handleMessage)
    source.removeEventListener('open', handleOpen)
    source.removeEventListener('error', handleError)
    source.close()
  }
}

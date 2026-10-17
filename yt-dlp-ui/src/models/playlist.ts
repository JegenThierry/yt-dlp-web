import type { AudioFormat, AudioQuality } from './enums'
import type { FilterOptions } from './filterOptions'
import type { PlaylistEntry } from './playlistEntry'

export interface CreatePlaylistRequest {
  url: string
  target_format: AudioFormat
  quality: AudioQuality
}

export interface DownloadPlaylistRequest {
  filters?: FilterOptions
  retry_unavailable?: boolean
}

export interface DownloadEntryRequest {
  filters?: FilterOptions
}

export interface Playlist {
  id: number
  url: string
  title: string
  quality: AudioQuality
  entries?: PlaylistEntry[]
  entry_count: number
  completed_count: number
  failed_count: number
  unavailable_count: number
  active_count: number
  total_duration: number
  created_at: Date
  updated_at: Date
  last_activity_at: Date
}

export interface PlaylistImport extends Playlist {
  existing: boolean
  added_count: number
  skipped_count: number
}

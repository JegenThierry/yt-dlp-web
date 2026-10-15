import type { AudioFormat } from './enums'
import type { FilterOptions } from './filterOptions'
import type { PlaylistEntry } from './playlistEntry'

export interface CreatePlaylistRequest {
  url: string
  target_format: AudioFormat
}

export interface DownloadPlaylistRequest {
  filters?: FilterOptions
}

export interface DownloadEntryRequest {
  filters?: FilterOptions
}

export interface Playlist {
  id: number
  url: string
  title: string
  entries?: PlaylistEntry[]
  entry_count: number
  created_at: Date
  updated_at: Date
}

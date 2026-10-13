import type { AudioFormat } from './enums'
import type { PlaylistEntry } from './playlistEntry'

export interface CreatePlaylistRequest {
  url: string
  target_format: AudioFormat
}

export interface Playlist {
  id: number
  url: string
  title: string
  entries?: PlaylistEntry[]
  created_at: Date
  updated_at: Date
}

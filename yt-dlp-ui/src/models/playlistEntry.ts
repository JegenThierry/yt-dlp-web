import type { DownloadStatus, AudioFormat } from './enums'

export interface PlaylistEntry {
  id: number
  playlist_id: number
  video_id: string
  title: string
  artist: string
  genre: string
  year: number
  duration: number
  thumbnail: string
  status: DownloadStatus
  file_path: string
  format: AudioFormat
  created_at: Date
  updated_at: Date
}

export interface PlaylistEntryMetadataUpdate {
  title: string
  artist: string
  genre: string
  year: number
}

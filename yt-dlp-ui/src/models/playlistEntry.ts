import type { DownloadStatus, AudioFormat, EntryErrorCode } from './enums'

export interface PlaylistEntry {
  id: number
  playlist_id: number
  video_id: string
  title: string
  artist: string
  source_title: string
  source_artist: string
  genre: string
  year: number
  duration: number
  thumbnail: string
  status: DownloadStatus
  file_path: string
  format: AudioFormat
  error_message?: string
  error_code?: EntryErrorCode
  created_at: Date
  updated_at: Date
}

export interface PlaylistEntryMetadataUpdate {
  id: number
  title: string
  artist: string
  genre: string
  year: number
}

export interface MetadataFieldError {
  entry_id: number
  field: string
  message: string
}

export interface MetadataValidationError {
  error: string
  fields: MetadataFieldError[]
}

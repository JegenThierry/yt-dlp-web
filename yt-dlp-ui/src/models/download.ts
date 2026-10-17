import type { DownloadPhase, DownloadStatus, EntryErrorCode } from './enums'

export interface ProgressEvent {
  playlist_id: number
  entry_id: number
  status: DownloadStatus
  file_path?: string
  error_message?: string
  error_code?: EntryErrorCode
  phase?: DownloadPhase
  percent?: number
  eta?: number
}

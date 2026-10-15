import type { DownloadStatus } from './enums'

export interface ProgressEvent {
  playlist_id: number
  entry_id: number
  status: DownloadStatus
  file_path?: string
}

import type { DownloadPhase } from '@/models/enums'

export interface EntryProgress {
  phase: DownloadPhase
  percent?: number
  eta?: number
}

export interface Download {
  id: number
  url: string
  status: DownloadStatus
}

export enum DownloadStatus {
  Pending = 'pending',
  Downloading = 'downloading',
  Completed = 'completed',
  Failed = 'failed',
}

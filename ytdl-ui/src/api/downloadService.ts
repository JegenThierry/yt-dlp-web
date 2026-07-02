import type { Download } from '@/models/download'
import { api } from './apiClient'

export async function submitDownloadRequest(url: string): Promise<Download[]> {
  return await api.post<Download[]>('/download', { url })
}

export async function fetchDownloads(): Promise<Download[]> {
  return await api.get<Download[]>('/downloads')
}

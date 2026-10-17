import { DownloadStatus } from '@/models/enums'
import type { Playlist } from '@/models/playlist'
import type { PlaylistEntry } from '@/models/playlistEntry'
import { isUnavailable } from '@/utils/entryFailure'

export interface PlaylistStats {
  total: number
  fresh: number
  queued: number
  downloading: number
  completed: number
  failed: number
  unavailable: number
  active: number
  remaining: number
}

export function playlistStats(entries: PlaylistEntry[]): PlaylistStats {
  const counts: Record<DownloadStatus, number> = {
    [DownloadStatus.New]: 0,
    [DownloadStatus.Pending]: 0,
    [DownloadStatus.Downloading]: 0,
    [DownloadStatus.Completed]: 0,
    [DownloadStatus.Failed]: 0,
  }
  let unavailable = 0
  for (const entry of entries) {
    counts[entry.status] = (counts[entry.status] ?? 0) + 1
    if (isUnavailable(entry)) unavailable++
  }

  const failed = counts[DownloadStatus.Failed] - unavailable

  return {
    total: entries.length,
    fresh: counts[DownloadStatus.New],
    queued: counts[DownloadStatus.Pending],
    downloading: counts[DownloadStatus.Downloading],
    completed: counts[DownloadStatus.Completed],
    failed,
    unavailable,
    active: counts[DownloadStatus.Pending] + counts[DownloadStatus.Downloading],
    remaining: counts[DownloadStatus.New] + failed,
  }
}

export type LibraryStats = Pick<
  PlaylistStats,
  'total' | 'fresh' | 'completed' | 'failed' | 'unavailable' | 'active' | 'remaining'
>

export function libraryStats(playlist: Playlist): LibraryStats {
  const total = playlist.entry_count
  const completed = playlist.completed_count
  const failed = playlist.failed_count
  const unavailable = playlist.unavailable_count
  const active = playlist.active_count
  const fresh = Math.max(0, total - completed - failed - unavailable - active)

  return { total, fresh, completed, failed, unavailable, active, remaining: fresh + failed }
}

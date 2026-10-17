import { DownloadStatus, EntryErrorCode } from '@/models/enums'
import type { PlaylistEntry } from '@/models/playlistEntry'

const failureHints: Partial<Record<EntryErrorCode, string>> = {
  [EntryErrorCode.Blocked]: 'YouTube is blocking or rate-limiting the server. Try again later.',
  [EntryErrorCode.SignInRequired]:
    "YouTube asked the server to sign in. Set up cookies in the API's .env, then retry.",
  [EntryErrorCode.Network]: "Couldn't reach YouTube. Check the server's connection, then retry.",
  [EntryErrorCode.Postprocessing]:
    'Converting the audio failed. Check that ffmpeg is installed on the server.',
}

export function isUnavailable(entry: PlaylistEntry): boolean {
  return entry.status === DownloadStatus.Failed && entry.error_code === EntryErrorCode.Unavailable
}

export function failureText(entry: PlaylistEntry): string {
  if (entry.status !== DownloadStatus.Failed) return ''
  const hint = entry.error_code ? failureHints[entry.error_code] : undefined
  return hint ?? entry.error_message ?? ''
}

export enum DownloadStatus {
  New = 'new',
  Pending = 'pending',
  Downloading = 'downloading',
  Completed = 'completed',
  Failed = 'failed',
}

export enum AudioFormat {
  MP3 = 'mp3',
  AAC = 'aac',
  FLAC = 'flac',
  WAV = 'wav',
  M4A = 'm4a',
  OPUS = 'opus',
  VORBIS = 'vorbis',
  ALAC = 'alac',
}

export enum AudioQuality {
  Best = 'best',
  High = 'high',
  Standard = 'standard',
  Low = 'low',
}

export enum DownloadPhase {
  Downloading = 'downloading',
  Converting = 'converting',
}

export enum EntryErrorCode {
  Unavailable = 'unavailable',
  SignInRequired = 'sign_in_required',
  Blocked = 'blocked',
  Network = 'network',
  Postprocessing = 'postprocessing',
  Unknown = 'unknown',
}

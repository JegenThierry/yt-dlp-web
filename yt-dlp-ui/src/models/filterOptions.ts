import type { AudioFormat, AudioQuality } from './enums'

export interface FilterOptions {
  format: AudioFormat
  quality: AudioQuality
  embed_metadata: boolean
}

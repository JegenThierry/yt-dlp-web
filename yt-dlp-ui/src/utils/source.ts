function isYouTube(url: string): boolean {
  if (!URL.canParse(url)) return false
  return /(^|\.)(youtube\.com|youtu\.be)$/.test(new URL(url).hostname)
}

export function sourceLabel(url: string): string {
  return isYouTube(url) ? 'Open on YouTube' : 'Open original'
}

export function entrySourceUrl(playlistUrl: string, videoId: string): string | undefined {
  if (!isYouTube(playlistUrl) || !videoId) return undefined
  return `https://www.youtube.com/watch?v=${encodeURIComponent(videoId)}`
}

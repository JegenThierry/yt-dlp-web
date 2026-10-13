const dateFormatter = new Intl.DateTimeFormat('en', { dateStyle: 'medium' })
const dateTimeFormatter = new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' })

/** Formats a date as a compact date, e.g. "Feb 12, 2026". */
export function formatDate(date: Date | string): string {
  return dateFormatter.format(new Date(date))
}

/** Formats a date including its time, e.g. "Feb 12, 2026, 3:24 PM". */
export function formatDateTime(date: Date | string): string {
  return dateTimeFormatter.format(new Date(date))
}

/** Formats a duration in seconds as "m:ss" or "h:mm:ss". */
export function formatDuration(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds <= 0) return '0:00'

  const total = Math.round(seconds)
  const hours = Math.floor(total / 3600)
  const minutes = Math.floor((total % 3600) / 60)
  const secs = total % 60

  return hours > 0
    ? `${hours}:${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
    : `${minutes}:${String(secs).padStart(2, '0')}`
}

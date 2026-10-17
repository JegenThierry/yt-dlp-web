import type { AudioQuality } from '@/models/enums'

const DAY_MS = 86_400_000
const RECENT_DAYS = 7

const dateFormatter = new Intl.DateTimeFormat('en', { dateStyle: 'medium' })
const dateTimeFormatter = new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' })
const monthDayFormatter = new Intl.DateTimeFormat('en', { month: 'short', day: 'numeric' })
const relativeDayFormatter = new Intl.RelativeTimeFormat('en', { numeric: 'auto' })

export function formatDate(date: Date | string): string {
  return dateFormatter.format(new Date(date))
}

export function formatDateTime(date: Date | string): string {
  return dateTimeFormatter.format(new Date(date))
}

function startOfDay(date: Date): number {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate()).getTime()
}

export function formatRelative(date: Date | string, now = new Date()): string {
  const value = new Date(date)
  const days = Math.round((startOfDay(now) - startOfDay(value)) / DAY_MS)
  if (days >= 0 && days < RECENT_DAYS) return relativeDayFormatter.format(-days, 'day')
  return value.getFullYear() === now.getFullYear()
    ? monthDayFormatter.format(value)
    : dateFormatter.format(value)
}

export function countLabel(count: number, singular: string, plural: string): string {
  return `${count} ${count === 1 ? singular : plural}`
}

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

export function formatTotalDuration(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds < 60) return formatDuration(seconds)

  const totalMinutes = Math.round(seconds / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours === 0) return `${minutes} min`
  return minutes === 0 ? `${hours} h` : `${hours} h ${minutes} min`
}

export function formatQuality(quality: AudioQuality): string {
  return quality.charAt(0).toUpperCase() + quality.slice(1)
}

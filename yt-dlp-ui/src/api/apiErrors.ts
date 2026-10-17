import axios from 'axios'
import type { MetadataFieldError } from '@/models/playlistEntry'

export const API_UNREACHABLE_MESSAGE = "Can't reach the API. Is yt-dlp-api running?"

export function apiErrorStatus(error: unknown): number | undefined {
  return axios.isAxiosError(error) ? error.response?.status : undefined
}

export function getApiErrorMessage(error: unknown, fallback: string): string {
  if (!axios.isAxiosError(error)) return fallback

  if (error.code === axios.AxiosError.ECONNABORTED || error.code === axios.AxiosError.ETIMEDOUT) {
    return 'The API took too long to answer. Try again.'
  }

  const response = error.response
  if (!response) return API_UNREACHABLE_MESSAGE

  const data: unknown = response.data
  if (data && typeof data === 'object' && 'error' in data) {
    const message = (data as { error: unknown }).error
    if (typeof message === 'string' && message.trim()) return message
  }

  if (response.status >= 500 && !data) return API_UNREACHABLE_MESSAGE

  return fallback
}

export function apiFieldErrors(error: unknown): MetadataFieldError[] {
  if (apiErrorStatus(error) !== 400 || !axios.isAxiosError(error)) return []

  const data: unknown = error.response?.data
  if (!data || typeof data !== 'object' || !('fields' in data)) return []

  const fields = (data as { fields: unknown }).fields
  return Array.isArray(fields) ? (fields as MetadataFieldError[]) : []
}

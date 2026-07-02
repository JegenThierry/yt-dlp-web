import { config } from '@/config/config'

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

async function request<T>(endpointUrl: string, options: RequestInit): Promise<T> {
  const url = `${config.apiBaseUrl}/${endpointUrl}`
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })

  if (!response.ok) {
    const text = await response.text()
    let message: string
    try {
      const body = JSON.parse(text)
      message = body.error || body.message || response.statusText
    } catch {
      message = text || response.statusText
    }
    throw new ApiError(message, response.status)
  }

  return response.json()
}

export const api = {
  get: async <T>(endpointUrl: string): Promise<T> => {
    return request<T>(endpointUrl, { method: 'GET' })
  },

  post: async <T>(endpointUrl: string, data: unknown): Promise<T> => {
    return request<T>(endpointUrl, { method: 'POST', body: JSON.stringify(data) })
  },

  put: async <T>(endpointUrl: string, data: unknown): Promise<T> => {
    return request<T>(endpointUrl, { method: 'PUT', body: JSON.stringify(data) })
  },

  delete: async <T>(endpointUrl: string): Promise<T> => {
    return request<T>(endpointUrl, { method: 'DELETE' })
  },
}

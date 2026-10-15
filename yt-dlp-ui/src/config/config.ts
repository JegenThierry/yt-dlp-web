function getBaseUrl(): string {
  return import.meta.env.VITE_API_BASE_URL ?? '/api'
}

export const config = {
  apiBaseUrl: getBaseUrl(),
}

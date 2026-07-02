function getBaseUrl(): string {
  return import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:3000/api'
}

export const config = {
  apiBaseUrl: getBaseUrl(),
}

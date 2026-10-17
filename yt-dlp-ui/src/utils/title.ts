export const APP_NAME = 'YT-DL UI'

export function pageTitle(page?: string): string {
  return page ? `${page} · ${APP_NAME}` : APP_NAME
}

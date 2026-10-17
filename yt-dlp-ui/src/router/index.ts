import { nextTick } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from 'vue-router/auto-routes'
import { pageTitle } from '@/utils/title'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
  }
}

const scrollOffsets = new Map<number, number>()
let currentPosition: number | undefined
let poppedState = false
let initialNavigation = true

window.addEventListener('popstate', () => {
  poppedState = true
})

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

function mainElement(): HTMLElement | null {
  return document.getElementById('main-content')
}

router.beforeEach(() => {
  const main = mainElement()
  if (main && currentPosition !== undefined) scrollOffsets.set(currentPosition, main.scrollTop)
})

router.afterEach(async (to, from, failure) => {
  const restoring = poppedState
  poppedState = false
  if (failure) return

  document.title = pageTitle(to.meta.title)
  currentPosition = history.state?.position

  if (initialNavigation) {
    initialNavigation = false
    return
  }
  if (to.path === from.path) return

  await nextTick()
  const main = mainElement()
  if (!main) return

  main.scrollTop =
    restoring && currentPosition !== undefined ? (scrollOffsets.get(currentPosition) ?? 0) : 0

  const heading = main.querySelector<HTMLElement>('h1[tabindex]') ?? main
  heading.focus({ preventScroll: true })
})

export default router

import type { C311Runtime } from '../plugins/c311'
import { canAccessC311Route, type C311ErrorLike, type C311RouteMeta } from './c311'

type RouterLike = {
  beforeEach: (handler: (to: any, from: any, next: (location?: any) => void) => any) => void
  afterEach: (handler: (to: any, from: any) => any) => void
}

type RouteLike = {
  meta?: C311RouteMeta
}

export function installC311RouterGuards (router: RouterLike, getRuntime: () => C311Runtime): void {
  router.beforeEach(async (to, _from, next) => {
    const meta = (to.matched || []).reduce((result: C311RouteMeta['c311'], route: RouteLike) => ({
      ...(result || {}),
      ...(route.meta?.c311 || {}),
    }), {})

    if (!meta || Object.keys(meta).length === 0) return next()

    const runtime = getRuntime()
    const session = await runtime.loadSession()
    const returnTo = encodeURIComponent(to.fullPath || to.path)
    const sessionError = runtime.error as C311ErrorLike | null
    if (sessionError && sessionError.status !== 401 && sessionError.status !== 403 && sessionError.code !== 'UNAUTHENTICATED' && sessionError.code !== 'FORBIDDEN') {
      // A session endpoint outage is not an authentication failure. Let the
      // page render its retryable/terminal state instead of redirecting to 401.
      return next()
    }
    if ((meta.requiresAuth && (!session?.authenticated || runtime.isExpired()))) {
      return next({ name: 'c311.unauthorized', query: { returnTo } })
    }
    if (!canAccessC311Route(session, meta)) {
      return next({ name: 'c311.forbidden', query: { returnTo } })
    }
    return next()
  })

  router.afterEach(() => {
    if (typeof document === 'undefined') return
    setTimeout(() => {
      const target = document.querySelector('[data-c311-main] h1, [data-c311-main]') as HTMLElement | null
      if (target && typeof target.focus === 'function') target.focus()
    }, 0)
  })
}

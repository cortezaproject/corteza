import type { Router } from 'vue-router'

/**
 * Returns the local path (pathname + search + hash) when `target` points to a
 * route the given router actually hosts internally, or null when it's external
 * / belongs to another app / can't be resolved.
 *
 * This is what lets the unified app navigate client-side between its own
 * sections while still doing a full page load for not-yet-merged apps: the
 * router only knows its own routes, so anything else falls back to the
 * catch-all and is treated as external.
 */
export function resolveInternalPath(router: Router | undefined, target: string): string | null {
  if (!router || !target) return null

  let local: string
  try {
    const url = new URL(target, window.location.origin)
    if (url.origin !== window.location.origin) return null // different origin → external
    local = url.pathname + url.search + url.hash
  } catch {
    return null
  }

  const resolved = router.resolve(local)

  // The catch-all fallback has a `:pathMatch` segment — that means "not ours".
  if (resolved.matched.some(record => /:pathMatch/.test(record.path))) return null

  // Guard against a redirect fallback resolving to a different path than asked.
  const pathname = local.split('?')[0].split('#')[0]
  if (resolved.path !== pathname) return null

  return local
}

/**
 * Anchor click handler: client-side navigation when the href is an internal
 * route, native browser navigation otherwise. Modified clicks (new tab, etc.)
 * are always left to the browser.
 */
export function handleInternalAnchorClick(
  router: Router | undefined,
  event: MouseEvent,
  target: string,
): void {
  if (
    event.button !== 0 ||
    event.metaKey ||
    event.ctrlKey ||
    event.shiftKey ||
    event.altKey
  ) {
    return
  }

  const internal = resolveInternalPath(router, target)
  if (internal && router) {
    event.preventDefault()
    router.push(internal)
  }
}

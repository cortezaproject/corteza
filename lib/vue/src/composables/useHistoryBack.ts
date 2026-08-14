/**
 * A Back button that always lands somewhere.
 *
 * `router.back()` alone is a dead end whenever the screen is the first entry in
 * the session history — a deep link, a bookmark, a fresh tab, or a redirect that
 * replaced the entry it came from. There it either does nothing or walks the
 * user out of the app entirely.
 *
 * `history.length` cannot answer this: it counts the whole tab's history,
 * including entries that predate the app, and never decreases. The router's own
 * `state.back` is the honest signal — it is null exactly when there is nothing
 * of ours to go back to.
 *
 * @example
 * const goBack = useHistoryBack()
 * // previous screen, or the module list when the editor was opened cold
 * goBack({ name: 'admin.modules', params: { slug } })
 */
import { useRouter, type RouteLocationRaw } from 'vue-router'

export function useHistoryBack() {
  const router = useRouter()

  return function goBack(fallback: RouteLocationRaw) {
    if (router.options.history.state?.back) {
      router.back()
      return
    }
    router.push(fallback)
  }
}

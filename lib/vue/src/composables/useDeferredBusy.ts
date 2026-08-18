import { onScopeDispose, ref, watch, toValue, type MaybeRefOrGetter, type Ref } from 'vue'

/**
 * A busy flag that is late to appear and slow to leave.
 *
 * A spinner that renders the moment work starts is the wrong instrument for
 * work that usually finishes in under a tenth of a second: it appears and
 * vanishes inside a couple of frames, which reads as a glitch rather than as
 * progress. One that leaves the moment work ends is the same glitch when the
 * work took barely longer than the delay.
 *
 * So the flag waits `delay` before turning on, and once on it stays on for at
 * least `minimum`. Work that beats the delay never shows a spinner at all —
 * the common case, and the point.
 *
 * The minimum is the price of the delay: a cover that stays up holds back what
 * is underneath it, so `minimum` trades a later reveal for one that does not
 * flicker. Keep it short enough that the trade is worth making.
 *
 * @example
 * const swapping = ref(false)
 * const cover = useDeferredBusy(swapping)
 * // swapping true for 80ms  → cover never turns on
 * // swapping true for 900ms → cover on at 150ms, off at 900ms
 * // swapping true for 200ms → cover on at 150ms, off at 650ms
 */
export function useDeferredBusy(
  source: MaybeRefOrGetter<boolean>,
  { delay = 150, minimum = 500 }: { delay?: number; minimum?: number } = {},
): Ref<boolean> {
  const shown = ref(false)

  let showTimer: ReturnType<typeof setTimeout> | null = null
  let hideTimer: ReturnType<typeof setTimeout> | null = null
  let shownAt = 0

  function clear() {
    if (showTimer) clearTimeout(showTimer)
    if (hideTimer) clearTimeout(hideTimer)
    showTimer = null
    hideTimer = null
  }

  watch(
    () => toValue(source),
    busy => {
      clear()

      if (busy) {
        // Already on: the work that turned it off is back before it left, so it
        // stays on and its minimum runs from when it first appeared.
        if (shown.value) return
        showTimer = setTimeout(() => {
          shown.value = true
          shownAt = Date.now()
        }, delay)
        return
      }

      if (!shown.value) return

      const left = minimum - (Date.now() - shownAt)
      if (left <= 0) {
        shown.value = false
        return
      }
      hideTimer = setTimeout(() => {
        shown.value = false
      }, left)
    },
    { immediate: true },
  )

  onScopeDispose(clear)

  return shown
}

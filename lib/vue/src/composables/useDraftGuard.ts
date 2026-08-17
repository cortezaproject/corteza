import { cloneDeep, isEqual } from 'lodash-es'
import { toValue, type MaybeRefOrGetter } from 'vue'
import { useUnsavedGuard } from './useUnsavedGuard'

interface DraftGuardOptions<T> {
  /** The resource being edited, live. */
  draft: MaybeRefOrGetter<T>
  /** True while a save, delete or clone is in flight — nothing to warn about. */
  busy?: MaybeRefOrGetter<boolean>
  /** State edited alongside the draft but held apart from it: member lists, a
   *  raw JSON pane, a layout set. Whatever it returns is compared the same way. */
  extra?: () => unknown
  messageKey?: string
  tabClose?: boolean
}

/**
 * The unsaved-changes guard for a screen that edits one resource: it owns the
 * baseline, the comparison and the suppression, so the screen says only what
 * its draft is and when that draft became clean.
 *
 *   const { capture, markSaved } = useDraftGuard({ draft: chart, busy: saving })
 *   // after load, and after a save that stays on the screen:
 *   capture()
 *   // before navigating away from a save:
 *   markSaved()
 *
 * `capture()` is explicit on purpose. Only the screen knows when its draft is
 * the thing the user was actually shown: data arrives asynchronously, and field
 * editors resolve presets on their own schedule, so any moment this composable
 * picked for itself would be a guess that reads as a false warning on load.
 */
export function useDraftGuard<T>(options: DraftGuardOptions<T>) {
  const { draft, busy, extra, messageKey = 'general.editor.unsavedChanges', tabClose } = options

  // Null until captured: an editor whose resource has not loaded has nothing to
  // lose, and comparing against a baseline it never took would warn on load.
  let baseline: { draft: unknown; extra: unknown } | null = null

  function snapshot() {
    return { draft: cloneDeep(toValue(draft)), extra: cloneDeep(extra?.()) }
  }

  /** Mark the current state as saved — the point comparisons are made against. */
  function capture() {
    baseline = snapshot()
  }

  /** Forget the baseline, so the guard stays quiet until the next capture(). */
  function reset() {
    baseline = null
  }

  // A getter rather than a computed: it is read only when the user leaves, and
  // a computed would cache a result whose inputs it cannot all see — a draft
  // mutated through a raw reference registers no dependency and the cached
  // "clean" would stand.
  function isDirty() {
    if (toValue(busy)) return false
    if (!baseline || toValue(draft) == null) return false

    const now = snapshot()
    return !isEqual(now.draft, baseline.draft) || !isEqual(now.extra, baseline.extra)
  }

  const { markSaved } = useUnsavedGuard({ isDirty, messageKey, tabClose })

  return { isDirty, capture, reset, markSaved }
}

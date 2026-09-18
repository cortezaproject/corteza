import { onScopeDispose, toValue, watch, type MaybeRefOrGetter, type WatchSource } from 'vue'

/**
 * Column widths for an auto-layout DataTable that only grow while it is
 * mounted.
 *
 * Such a table sizes each column to the widest value on the rows it shows, so a
 * sort or page that brings narrower rows shrinks columns the next one widens
 * again. Around every change of `sources` each column's natural width is read
 * and the column is held at the widest seen, as its header's min-width.
 *
 * Columns are the header cells carrying `data-field`. One the user drags to a
 * width of their own is theirs from then on, and a click on a resize grip is
 * never taken for a sort.
 *
 * @example
 * const table = ref()
 * useGrowOnlyColumns(() => table.value?.$el, [() => props.items, () => props.fields])
 */
export function useGrowOnlyColumns(
  root: MaybeRefOrGetter<HTMLElement | null | undefined>,
  sources: WatchSource | WatchSource[],
) {
  const widest = new Map<string, number>()
  const userSized = new Set<string>()

  function hold() {
    const table = toValue(root)?.querySelector?.('table')
    if (!table) return

    const headers = [...table.querySelectorAll<HTMLElement>('thead th[data-field]')]
    if (!headers.length) return

    // The width each column asks for on its own, before the table spreads its
    // spare room across them
    const { width, minWidth } = table.style
    table.style.width = 'max-content'
    table.style.minWidth = '0'
    const natural = headers.map(th => th.getBoundingClientRect().width)
    table.style.width = width
    table.style.minWidth = minWidth

    headers.forEach((th, i) => {
      const field = th.dataset.field as string
      if (userSized.has(field)) return

      const held = Math.max(widest.get(field) || 0, Math.ceil(natural[i]))
      widest.set(field, held)
      th.style.minWidth = `${held}px`
    })
  }

  function onResizeStart(event: Event) {
    const target = event.target as HTMLElement | null
    if (!target?.closest?.('.p-datatable-column-resizer')) return

    const th = target.closest<HTMLElement>('th[data-field]')
    if (!th?.dataset.field) return

    userSized.add(th.dataset.field)
    th.style.minWidth = ''
  }

  function onGripClick(event: Event) {
    const target = event.target as HTMLElement | null
    if (target?.closest?.('.p-datatable-column-resizer')) event.stopPropagation()
  }

  function listen(el: HTMLElement | null | undefined, on: boolean) {
    const method = on ? 'addEventListener' : 'removeEventListener'
    el?.[method]?.('mousedown', onResizeStart, true)
    el?.[method]?.('click', onGripClick, true)
  }

  watch(
    () => toValue(root),
    (el, prev) => {
      listen(prev, false)
      listen(el, true)
    },
    { immediate: true, flush: 'post' },
  )

  // Before the swap as well as after it: cells that filled in since the last
  // swap (async field viewers) are only seen here
  watch(sources, hold, { flush: 'pre' })
  watch(sources, hold, { flush: 'post' })

  onScopeDispose(() => listen(toValue(root), false))

  return { hold }
}

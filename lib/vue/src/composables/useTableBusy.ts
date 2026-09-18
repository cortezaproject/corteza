import { computed, nextTick, ref, toValue, watch, type MaybeRefOrGetter } from 'vue'
import { useDeferredBusy } from './useDeferredBusy'

/**
 * Loading feedback for a lazy DataTable: a sort shows a spinner in the header
 * of the column it was asked for and leaves the rows readable; any other
 * refetch gets the table's mask. Both are deferred (see useDeferredBusy), so
 * work that beats the delay shows nothing but the new rows.
 *
 * Call `sortStarted(field)` from the table's `@sort` handler, after starting
 * the fetch it asks for.
 *
 * @example
 * const { masked, sortingColumn, sortStarted } = useTableBusy(loading)
 * // <DataTable :loading="masked" …>
 * //   <Column …><template v-if="sortingColumn === col.name" #sorticon="{ class: c }">
 * //     <i :class="[c, 'pi pi-spin pi-spinner']" />
 */
export function useTableBusy(loading: MaybeRefOrGetter<boolean>) {
  const busy = useDeferredBusy(loading)
  const sortField = ref<string | null>(null)

  // Held until the spinner has gone, so it keeps its minimum like the mask does
  watch([() => toValue(loading), busy], ([isLoading, isBusy]) => {
    if (!isLoading && !isBusy) sortField.value = null
  })

  function sortStarted(field?: string | null) {
    sortField.value = field || null
    // A sort that started no fetch has nothing to wait for. A `loading` prop
    // reaches the child only once the parent re-renders, hence the tick.
    nextTick(() => {
      if (!toValue(loading)) sortField.value = null
    })
  }

  return {
    masked: computed(() => busy.value && !sortField.value),
    sortingColumn: computed(() => (busy.value ? sortField.value : null)),
    sortStarted,
  }
}

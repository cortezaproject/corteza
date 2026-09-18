import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { effectScope, nextTick, ref, watch } from 'vue'
import { useTableBusy } from './useTableBusy'

// What the tests pin is where the feedback goes: the sorted column's header
// for a sort, the mask for anything else, and nowhere for work that is quick.

describe('useTableBusy', () => {
  let scope: ReturnType<typeof effectScope>

  beforeEach(() => {
    vi.useFakeTimers()
    scope = effectScope()
  })

  afterEach(() => {
    scope.stop()
    vi.useRealTimers()
  })

  function setup(loading: ReturnType<typeof ref<boolean>>) {
    return scope.run(() => useTableBusy(loading))!
  }

  async function sort(loading, busy, field) {
    loading.value = true
    busy.sortStarted(field)
    await nextTick()
  }

  it('shows nothing for a sort that beats the delay', async () => {
    const loading = ref(false)
    const busy = setup(loading)
    const seen: unknown[] = []
    scope.run(() => watch([busy.masked, busy.sortingColumn], v => seen.push(v)))

    await sort(loading, busy, 'name')
    vi.advanceTimersByTime(80)
    loading.value = false
    await nextTick()
    vi.advanceTimersByTime(1000)
    await nextTick()

    expect(seen).toEqual([])
  })

  it('puts a slow sort on its column and never masks the rows', async () => {
    const loading = ref(false)
    const busy = setup(loading)
    const masked: boolean[] = []
    scope.run(() => watch(busy.masked, v => masked.push(v)))

    await sort(loading, busy, 'name')
    vi.advanceTimersByTime(150)
    await nextTick()
    expect(busy.sortingColumn.value).toBe('name')

    vi.advanceTimersByTime(600)
    loading.value = false
    await nextTick()
    expect(busy.sortingColumn.value).toBe(null)
    expect(masked).toEqual([])
  })

  it('keeps the header spinner for the minimum when the sort ends just after it appeared', async () => {
    const loading = ref(false)
    const busy = setup(loading)

    await sort(loading, busy, 'name')
    vi.advanceTimersByTime(200)
    loading.value = false
    await nextTick()
    expect(busy.sortingColumn.value).toBe('name')

    vi.advanceTimersByTime(450)
    await nextTick()
    expect(busy.sortingColumn.value).toBe(null)
  })

  it('masks a slow refetch that is not a sort', async () => {
    const loading = ref(false)
    const busy = setup(loading)

    loading.value = true
    await nextTick()
    vi.advanceTimersByTime(150)
    await nextTick()

    expect(busy.masked.value).toBe(true)
    expect(busy.sortingColumn.value).toBe(null)
  })

  it('moves a mask already up to the header once a sort takes over', async () => {
    const loading = ref(false)
    const busy = setup(loading)

    loading.value = true
    await nextTick()
    vi.advanceTimersByTime(300)
    await nextTick()
    expect(busy.masked.value).toBe(true)

    busy.sortStarted('name')
    await nextTick()
    expect(busy.masked.value).toBe(false)
    expect(busy.sortingColumn.value).toBe('name')
  })

  it('forgets a sort that started no fetch, so a later refetch is masked', async () => {
    const loading = ref(false)
    const busy = setup(loading)

    busy.sortStarted('name')
    await nextTick()

    loading.value = true
    await nextTick()
    vi.advanceTimersByTime(150)
    await nextTick()

    expect(busy.masked.value).toBe(true)
    expect(busy.sortingColumn.value).toBe(null)
  })
})

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { effectScope, nextTick, ref, watch } from 'vue'
import { useDeferredBusy } from './useDeferredBusy'

// A flag that turned on and off inside two frames is the glitch this replaces,
// so what the tests pin is when it turns on, and how long it stays.

describe('useDeferredBusy', () => {
  let scope: ReturnType<typeof effectScope>

  beforeEach(() => {
    vi.useFakeTimers()
    scope = effectScope()
  })

  afterEach(() => {
    scope.stop()
    vi.useRealTimers()
  })

  function setup(busy: ReturnType<typeof ref<boolean>>, opts = {}) {
    return scope.run(() => useDeferredBusy(busy, opts))!
  }

  it('never turns on for work that beats the delay', async () => {
    const busy = ref(false)
    const shown = setup(busy)
    // Every transition, not just where it ends up: turning on and off again is
    // the flicker this exists to prevent, and a final `false` cannot see it
    const seen: boolean[] = []
    scope.run(() => watch(shown, v => seen.push(v)))

    busy.value = true
    await nextTick()
    vi.advanceTimersByTime(80)
    busy.value = false
    await nextTick()

    vi.advanceTimersByTime(1000)
    await nextTick()
    expect(seen).toEqual([])
    expect(shown.value).toBe(false)
  })

  it('turns on once the work outlasts the delay', async () => {
    const busy = ref(false)
    const shown = setup(busy)

    busy.value = true
    await nextTick()

    vi.advanceTimersByTime(149)
    expect(shown.value).toBe(false)

    vi.advanceTimersByTime(1)
    expect(shown.value).toBe(true)
  })

  it('stays on for the minimum when the work ends just after it appeared', async () => {
    const busy = ref(false)
    const shown = setup(busy)

    busy.value = true
    await nextTick()
    vi.advanceTimersByTime(200) // on at 150
    busy.value = false
    await nextTick()

    // 50ms of the 500ms minimum spent
    vi.advanceTimersByTime(449)
    expect(shown.value).toBe(true)

    vi.advanceTimersByTime(1)
    expect(shown.value).toBe(false)
  })

  it('adds no delay of its own to work that outlasted the minimum', async () => {
    const busy = ref(false)
    const shown = setup(busy)

    busy.value = true
    await nextTick()
    vi.advanceTimersByTime(1200) // on at 150, minimum spent by 650
    busy.value = false
    await nextTick()

    expect(shown.value).toBe(false)
  })

  it('stays on when work resumes before the minimum is up', async () => {
    const busy = ref(false)
    const shown = setup(busy)

    busy.value = true
    await nextTick()
    vi.advanceTimersByTime(200)
    busy.value = false
    await nextTick()
    vi.advanceTimersByTime(50)

    busy.value = true
    await nextTick()
    expect(shown.value).toBe(true)

    // and the minimum still runs from when it first appeared, not from here
    vi.advanceTimersByTime(50)
    busy.value = false
    await nextTick()
    vi.advanceTimersByTime(1000)
    expect(shown.value).toBe(false)
  })

  it('drops a pending turn-on when the scope goes', async () => {
    const busy = ref(false)
    const shown = setup(busy)

    busy.value = true
    await nextTick()
    scope.stop()

    vi.advanceTimersByTime(1000)
    expect(shown.value).toBe(false)
  })
})

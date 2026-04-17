import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { withMinDuration, useMinDuration } from './useMinDuration'

describe('withMinDuration', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('resolves with the promise value', async () => {
    const promise = Promise.resolve('hello')
    const result = withMinDuration(promise, 100)
    vi.advanceTimersByTime(100)
    expect(await result).toBe('hello')
  })

  it('waits at least minMs when promise is faster', async () => {
    let settled = false
    const promise = Promise.resolve('fast')
    const result = withMinDuration(promise, 300).then(v => {
      settled = true
      return v
    })

    vi.advanceTimersByTime(299)
    await Promise.resolve() // flush microtasks
    expect(settled).toBe(false)

    vi.advanceTimersByTime(1)
    expect(await result).toBe('fast')
    expect(settled).toBe(true)
  })

  it('does not add extra delay when promise is slower', async () => {
    let resolved = false
    const slowPromise = new Promise<string>(resolve =>
      setTimeout(() => {
        resolved = true
        resolve('slow')
      }, 500),
    )

    const result = withMinDuration(slowPromise, 100)
    vi.advanceTimersByTime(500)
    expect(await result).toBe('slow')
    expect(resolved).toBe(true)
  })
})

describe('useMinDuration', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('starts with loading false', () => {
    const { loading } = useMinDuration(300)
    expect(loading.value).toBe(false)
  })

  it('sets loading true during run and false after', async () => {
    const { loading, run } = useMinDuration(100)

    const runPromise = run(() => Promise.resolve('done'))
    expect(loading.value).toBe(true)

    vi.advanceTimersByTime(100)
    await runPromise

    expect(loading.value).toBe(false)
  })

  it('returns the value from the fn', async () => {
    const { run } = useMinDuration(0)
    const result = run(() => Promise.resolve(42))
    vi.advanceTimersByTime(0)
    expect(await result).toBe(42)
  })

  it('resets loading to false even when fn throws', async () => {
    const { loading, run } = useMinDuration(0)
    await expect(
      run(() => Promise.reject(new Error('oops'))).finally(() => vi.advanceTimersByTime(0)),
    ).rejects.toThrow('oops')
    expect(loading.value).toBe(false)
  })
})

import { describe, expect, it } from 'vitest'
import { statusFilter, statusOf } from './useResourceStatus'

describe('resource status', () => {
  const states = ['suspended', 'deleted']

  it('reads the state set to only, else active', () => {
    expect(statusOf({ suspended: '0', deleted: '0' }, states)).toBe('active')
    expect(statusOf({ suspended: '2', deleted: '0' }, states)).toBe('suspended')
    expect(statusOf({ suspended: '1', deleted: '2' }, states)).toBe('deleted')
    expect(statusOf({ suspended: '1', deleted: '1' }, states)).toBe('active')
  })

  it('shows one status at a time', () => {
    expect(statusFilter('active', states)).toEqual({ suspended: '0', deleted: '0' })
    expect(statusFilter('suspended', states)).toEqual({ suspended: '2', deleted: '0' })
  })

  it('keeps deleted rows whatever else they carry', () => {
    expect(statusFilter('deleted', states)).toEqual({ suspended: '1', deleted: '2' })
  })

  it('round-trips every status', () => {
    for (const s of ['active', ...states]) expect(statusOf(statusFilter(s, states), states)).toBe(s)
  })
})

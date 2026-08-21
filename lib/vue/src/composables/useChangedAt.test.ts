import { describe, it, expect } from 'vitest'
import {
  CHANGED_AT_KEY,
  CHANGED_AT_SORT,
  changedAt,
  changedAtField,
  changedAtText,
} from './useChangedAt'

describe('changedAt', () => {
  it('prefers deletedAt, then updatedAt, then createdAt', () => {
    expect(
      changedAt({ deletedAt: '2026-03-03', updatedAt: '2026-02-02', createdAt: '2026-01-01' }),
    ).toBe('2026-03-03')
    expect(changedAt({ updatedAt: '2026-02-02', createdAt: '2026-01-01' })).toBe('2026-02-02')
    expect(changedAt({ createdAt: '2026-01-01' })).toBe('2026-01-01')
  })

  it('has no value when the resource carries no timestamp', () => {
    expect(changedAt({})).toBeUndefined()
    expect(changedAt(undefined)).toBeUndefined()
  })

  // A catalog connection is assembled by the server rather than stored, so its
  // createdAt arrives as Go's zero time — truthy, and a date nobody wants shown.
  // A lib/js model casts the same fields to Date in its constructor, so the
  // column has to read both shapes — ProjectList maps rows to system.Project.
  it('reads a model row whose timestamps are Date objects', () => {
    const d = new Date('2026-02-02T10:00:00Z')
    expect(changedAt({ updatedAt: d, createdAt: new Date('2026-01-01T10:00:00Z') })).toBe(d)
    expect(changedAtText({ createdAt: d })).not.toBe('')
  })

  it('reads a zero-time Date as no timestamp', () => {
    expect(changedAt({ createdAt: new Date('0001-01-01T00:00:00Z') })).toBeUndefined()
    expect(changedAt({ createdAt: new Date('nope') })).toBeUndefined()
  })

  it('reads the zero time as no timestamp, in any offset', () => {
    expect(changedAt({ createdAt: '0001-01-01T00:00:00Z' })).toBeUndefined()
    expect(changedAt({ createdAt: '0001-01-01T00:58:04+00:58' })).toBeUndefined()
    expect(changedAt({ updatedAt: '0001-01-01T00:00:00Z', createdAt: '2026-01-01' })).toBe(
      '2026-01-01',
    )
  })
})

describe('changedAtText', () => {
  it('renders the resolved timestamp', () => {
    expect(changedAtText({ createdAt: '2026-01-01T10:00:00Z' })).not.toBe('')
  })

  // A missing date reaches moment as undefined, which it reads as *now* — an
  // unguarded call renders today for a resource that was never written.
  it('is empty rather than today when there is no timestamp', () => {
    expect(changedAtText({})).toBe('')
    expect(changedAtText({ deletedAt: null, updatedAt: null, createdAt: null })).toBe('')
    expect(changedAtText({ createdAt: '0001-01-01T00:00:00Z' })).toBe('')
    expect(changedAtText({ createdAt: '0001-01-01T00:58:04+00:58' })).toBe('')
  })
})

describe('changedAtField', () => {
  it('carries the key useResourceList maps to the COALESCE sort', () => {
    expect(changedAtField('Last change').key).toBe(CHANGED_AT_KEY)
    expect(CHANGED_AT_SORT).toBe('coalesce(deletedAt, updatedAt, createdAt)')
  })

  it('is sortable and right-aligned by default', () => {
    const f = changedAtField('Last change')
    expect(f.header).toBe('Last change')
    expect(f.sortable).toBe(true)
    expect(f.class).toBe('text-right')
  })

  it('takes overrides for lists the server cannot sort', () => {
    expect(changedAtField('Last change', { sortable: false }).sortable).toBe(false)
  })
})

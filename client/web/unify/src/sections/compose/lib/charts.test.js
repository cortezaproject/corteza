import { describe, it, expect } from 'vitest'
import { wrapLabel } from './charts'

// One unit per character keeps the widths readable.
const measure = text => text.length

describe('wrapLabel', () => {
  it('leaves a label that fits on one line', () => {
    expect(wrapLabel('Closed Won', 10, measure)).toBe('Closed Won')
  })

  it('breaks between words once a line would run past the width', () => {
    expect(wrapLabel('Europe, Middle East & Africa', 14, measure)).toBe(
      'Europe, Middle\nEast & Africa',
    )
  })

  it('never splits a word wider than the width', () => {
    expect(wrapLabel('Qualification', 5, measure)).toBe('Qualification')
    expect(wrapLabel('2026-01-01', 4, measure)).toBe('2026-01-01')
  })

  it('takes numbers', () => {
    expect(wrapLabel(2026, 2, measure)).toBe('2026')
  })
})

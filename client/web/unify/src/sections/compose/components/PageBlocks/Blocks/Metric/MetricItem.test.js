import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'

import MetricItem from './MetricItem.vue'

// The value is sized in container units (cqw/cqh), which scale one glyph against
// the card and know nothing about how many glyphs there are. A one-character
// count fits; a formatted currency sum is an order of magnitude wider than the
// card and clips at its edge. --metric-chars is the divisor that makes the rule
// length-aware, so it has to count exactly what the template renders.

const mountItem = (metric, value) => mount(MetricItem, { props: { metric, value: { value } } })

const chars = wrapper =>
  Number(
    wrapper
      .get('.metric-value')
      .attributes('style')
      .match(/--metric-chars:\s*(\d+)/)[1],
  )

describe('MetricItem', () => {
  it('counts a bare value', () => {
    expect(chars(mountItem({}, '18'))).toBe(2)
  })

  it('counts prefix and suffix along with the value', () => {
    // '$ ' + ' ' + '145,873.92' — the template puts a space between the parts
    expect(chars(mountItem({ prefix: '$ ', suffix: '' }, '145,873.92'))).toBe(13)
    expect(chars(mountItem({ suffix: ' %' }, '42.8'))).toBe(7)
  })

  it('never divides by zero', () => {
    expect(chars(mountItem({}, ''))).toBe(1)
    expect(chars(mountItem({}, undefined))).toBe(1)
  })

  it('gives a long value a smaller share of the width than a short one', () => {
    expect(chars(mountItem({ prefix: '$ ' }, '145,873.92'))).toBeGreaterThan(
      chars(mountItem({}, '18')),
    )
  })
})

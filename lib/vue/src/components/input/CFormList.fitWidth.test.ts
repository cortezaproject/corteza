import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CFormList from './CFormList.vue'

// A row of inputs has no meaningful content width — an `<input>` reports the
// browser's default `size`, so a grid sized to its content is wider than any
// container and scrolls sideways from the first row, hiding whatever columns
// sit on the right. `fitWidth` lets the tracks share the width the list is
// given instead; the column minimums still hold and the list still scrolls
// under them.

const rows = [{ name: 'first' }, { name: 'second' }]

const list = (props = {}) =>
  mount(CFormList, {
    props: { modelValue: rows, columns: [{ label: 'Name' }, { label: 'Kind' }], ...props },
    global: { stubs: { CInputDelete: true } },
  })

const inner = (w: ReturnType<typeof list>) => w.findComponent({ name: 'CDraggableList' })

describe('CFormList fitWidth', () => {
  it('sizes to its content by default', () => {
    expect(inner(list()).classes()).toContain('min-w-max')
    expect(inner(list({ stickyHeader: true })).classes()).toContain('min-w-max')
  })

  it('shares the given width when asked to', () => {
    expect(inner(list({ fitWidth: true })).classes()).not.toContain('min-w-max')
    expect(inner(list({ fitWidth: true, stickyHeader: true })).classes()).not.toContain('min-w-max')
  })

  it('keeps the frame and the scroll region either way', () => {
    // Plain list: the border sits on the inner element, the root scrolls.
    const plain = list({ fitWidth: true })
    expect(inner(plain).classes()).toContain('border')
    expect(plain.classes()).toContain('overflow-x-auto')

    // Sticky list: the frame moves out to the root so the header can pin to it.
    const sticky = list({ fitWidth: true, stickyHeader: true })
    expect(inner(sticky).classes()).not.toContain('border')
    expect(sticky.classes()).toEqual(
      expect.arrayContaining(['overflow-auto', 'border', 'rounded-border']),
    )
  })
})

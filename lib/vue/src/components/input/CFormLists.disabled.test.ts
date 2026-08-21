import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CFormList from './CFormList.vue'
import CFormItemList from './CFormItemList.vue'

// Both lists carry their own controls — a remove button and a drag handle — so
// an editor cannot make them read-only by disabling its own fields. Without
// `disabled` here, a user with no update permission still gets a live Remove
// on every row of a form whose Save button is hidden.

const rows = [
  { id: 'a', name: 'first' },
  { id: 'b', name: 'second' },
]

const removeButtons = (w: ReturnType<typeof mount>) =>
  w.findAll('button').filter(b => b.find('.pi-trash').exists())

const draggable = (w: ReturnType<typeof mount>) => w.findComponent({ name: 'CDraggableList' })

describe('CFormList disabled', () => {
  const mountList = (props = {}) =>
    mount(CFormList, {
      props: { modelValue: rows, columns: [{ field: 'name' }], ...props },
      global: { stubs: { CInputDelete: true } },
    })

  it('offers remove and reordering by default', () => {
    const w = mountList({ draggable: true })
    expect(removeButtons(w).length).toBe(rows.length)
    expect(w.find('.c-drag-handle').exists()).toBe(true)
    expect(draggable(w).props('disabled')).toBe(false)
  })

  it('withdraws remove, the drag handle and reordering when disabled', () => {
    const w = mountList({ draggable: true, disabled: true })
    expect(removeButtons(w).length).toBe(0)
    expect(w.find('.c-drag-handle').exists()).toBe(false)
    expect(draggable(w).props('disabled')).toBe(true)
  })
})

describe('CFormItemList disabled', () => {
  const mountList = (props = {}) =>
    mount(CFormItemList, { props: { items: rows, itemKey: 'id', ...props } })

  it('offers remove and reordering by default', () => {
    const w = mountList({ draggable: true })
    expect(removeButtons(w).length).toBe(rows.length)
    expect(w.find('.c-drag-handle').exists()).toBe(true)
    expect(draggable(w).props('disabled')).toBe(false)
  })

  it('withdraws remove, the drag handle and reordering when disabled', () => {
    const w = mountList({ draggable: true, disabled: true })
    expect(removeButtons(w).length).toBe(0)
    expect(w.find('.c-drag-handle').exists()).toBe(false)
    expect(draggable(w).props('disabled')).toBe(true)
  })

  // The rows themselves stay: a list you may read but not change is still worth
  // showing.
  it('still renders every row when disabled', () => {
    const w = mount(CFormItemList, {
      props: { items: rows, itemKey: 'id', disabled: true },
      slots: { default: '<template #default="{ item }">{{ item.name }}</template>' },
    })
    expect(w.text()).toContain('first')
    expect(w.text()).toContain('second')
  })
})

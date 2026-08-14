import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CFormList from './CFormList.vue'
import CFormItemList from './CFormItemList.vue'

// Both lists reorder through the shared draggable now, and a real drag needs a
// pointer jsdom cannot provide. What is worth pinning is the seam either side
// of it: the rows the sortable is given to address, and what each list does
// with the new order it hands back.
//
// CFormItemList matters most — its `items` is a prop, so `reorder` is the only
// way a move reaches the owner. Two consumers act on that payload (chatbot
// Scenarios, admin ApiGateway); dropping it would lose their reorder silently.

const draggable = (w: ReturnType<typeof mount>) => w.findComponent({ name: 'CDraggableList' })

const rows = [
  { id: 'a', name: 'first' },
  { id: 'b', name: 'second' },
  { id: 'c', name: 'third' },
]

describe('CFormList reorder', () => {
  const mountList = (props = {}) =>
    mount(CFormList, {
      props: { modelValue: [...rows], draggable: true, ...props },
      global: { stubs: { Button: true, CInputDelete: true }, directives: { tooltip: {} } },
    })

  it('marks every row as an item the sortable can address', () => {
    expect(mountList().findAll('[data-drag-item]')).toHaveLength(rows.length)
  })

  it('drags by the grip, so an input in a row keeps its own pointer gesture', () => {
    const w = mountList()

    expect(draggable(w).props('handle')).toBe('.c-drag-handle')
    expect(w.findAll('.c-drag-handle')).toHaveLength(rows.length)
  })

  it('renders no grip and disables the sortable when not draggable', () => {
    const w = mountList({ draggable: false })

    expect(draggable(w).props('disabled')).toBe(true)
    expect(w.find('.c-drag-handle').exists()).toBe(false)
  })

  it('reports a new order through the model and the reorder event', async () => {
    const w = mountList()
    const moved = [rows[2], rows[0], rows[1]]

    await draggable(w).vm.$emit('update:modelValue', moved)
    await draggable(w).vm.$emit('update')

    expect(w.emitted('update:modelValue')?.at(-1)?.[0]).toEqual(moved)
    expect(w.emitted('reorder')?.at(-1)?.[0]).toEqual(moved)
  })
})

describe('CFormItemList reorder', () => {
  const mountList = (props = {}) =>
    mount(CFormItemList, {
      props: { items: [...rows], draggable: true, ...props },
      global: { stubs: { Button: true }, directives: { tooltip: {} } },
    })

  it('marks every row as an item the sortable can address', () => {
    expect(mountList().findAll('[data-drag-item]')).toHaveLength(rows.length)
  })

  it('drags by the grip', () => {
    const w = mountList()

    expect(draggable(w).props('handle')).toBe('.c-drag-handle')
    expect(w.findAll('.c-drag-handle')).toHaveLength(rows.length)
  })

  it('reports the new order to its owner, which owns the list', async () => {
    const w = mountList()
    const moved = [rows[1], rows[2], rows[0]]

    await draggable(w).vm.$emit('update:modelValue', moved)

    expect(w.emitted('reorder')?.at(-1)?.[0]).toEqual(moved)
  })

  it('leaves the rows alone until the owner accepts the move', async () => {
    const w = mountList()

    await draggable(w).vm.$emit('update:modelValue', [rows[2], rows[1], rows[0]])

    // The prop never changed, so the list still shows what its owner holds.
    expect(draggable(w).props('modelValue')).toEqual(rows)
  })
})

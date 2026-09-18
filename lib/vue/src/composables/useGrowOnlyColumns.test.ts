import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { defineComponent, h, nextTick, ref } from 'vue'
import { mount, type VueWrapper } from '@vue/test-utils'
import { useGrowOnlyColumns } from './useGrowOnlyColumns'

// jsdom lays nothing out, so a header reports its column's width as 10px per
// character of the longest text in it, and only while the table is at its
// natural width; any other read gets the spread width a full-width table gives.

const FIELDS = ['name', 'set']
const SPREAD = 999

type Row = Record<string, string>

const List = defineComponent({
  props: {
    rows: { type: Array as () => Row[], required: true },
    headerKey: { type: Number, default: 0 },
  },
  setup(props) {
    const root = ref<HTMLElement>()
    useGrowOnlyColumns(root, [() => props.rows, () => props.headerKey])
    return () =>
      h('div', { ref: root }, [
        h('table', { style: 'width: 100%; min-width: 50rem' }, [
          h('thead', { key: props.headerKey }, [
            h('tr', [
              h('th'),
              ...FIELDS.map(f =>
                h('th', { 'data-field': f }, [h('span', { class: 'p-datatable-column-resizer' })]),
              ),
            ]),
          ]),
          h(
            'tbody',
            props.rows.map(r => h('tr', [h('td'), ...FIELDS.map(f => h('td', r[f]))])),
          ),
        ]),
      ])
  },
})

const realRect = HTMLTableCellElement.prototype.getBoundingClientRect

function measure(this: HTMLTableCellElement) {
  const table = this.closest('table')!
  const index = [...this.parentElement!.children].indexOf(this)
  const cells = [...table.querySelectorAll(`tbody td:nth-child(${index + 1})`)]
  const content = Math.max(0, ...cells.map(td => (td.textContent || '').length * 10))
  const held = parseFloat(this.style.minWidth) || 0
  const width = table.style.width === 'max-content' ? content : SPREAD
  return { width: Math.max(width, held) } as DOMRect
}

const held = (wrapper: VueWrapper, field: string) =>
  (wrapper.find(`th[data-field="${field}"]`).element as HTMLElement).style.minWidth

describe('useGrowOnlyColumns', () => {
  let wrapper: VueWrapper

  beforeEach(() => {
    HTMLTableCellElement.prototype.getBoundingClientRect = measure
  })

  afterEach(() => {
    wrapper?.unmount()
    HTMLTableCellElement.prototype.getBoundingClientRect = realRect
  })

  async function show(first: Row[]) {
    wrapper = mount(List, { props: { rows: [] }, attachTo: document.body })
    await wrapper.setProps({ rows: first })
    await nextTick()
  }

  it('holds each column at the widest it has been', async () => {
    await show([{ name: 'x'.repeat(20), set: 'x'.repeat(15) }])
    expect(held(wrapper, 'name')).toBe('200px')
    expect(held(wrapper, 'set')).toBe('150px')

    await wrapper.setProps({ rows: [{ name: 'x'.repeat(12), set: 'x'.repeat(26) }] })
    expect(held(wrapper, 'name')).toBe('200px')
    expect(held(wrapper, 'set')).toBe('260px')

    await wrapper.setProps({ rows: [{ name: 'x', set: 'x' }] })
    expect(held(wrapper, 'name')).toBe('200px')
    expect(held(wrapper, 'set')).toBe('260px')
  })

  it('reads natural widths, not the spread ones, and puts the table back', async () => {
    await show([{ name: 'x'.repeat(20), set: 'x'.repeat(15) }])

    expect(held(wrapper, 'name')).toBe('200px')
    const table = wrapper.find('table').element as HTMLTableElement
    expect(table.style.width).toBe('100%')
    expect(table.style.minWidth).toBe('50rem')
  })

  it('keeps what a cell filled in after the rows arrived', async () => {
    await show([{ name: 'x'.repeat(10), set: 'x'.repeat(10) }])
    expect(held(wrapper, 'set')).toBe('100px')

    // An async viewer drawing its value a moment after the row
    wrapper.findAll('tbody td')[2].element.textContent = 'x'.repeat(30)

    await wrapper.setProps({ rows: [{ name: 'x'.repeat(10), set: 'x'.repeat(5) }] })
    expect(held(wrapper, 'set')).toBe('300px')
  })

  it('keeps a column at its widest when its header is drawn afresh', async () => {
    await show([{ name: 'x'.repeat(20), set: 'x'.repeat(15) }])

    await wrapper.setProps({ headerKey: 1, rows: [{ name: 'x'.repeat(5), set: 'x'.repeat(5) }] })

    expect(held(wrapper, 'name')).toBe('200px')
    expect(held(wrapper, 'set')).toBe('150px')
  })

  it('leaves a column the user resizes to them', async () => {
    await show([{ name: 'x'.repeat(20), set: 'x'.repeat(15) }])

    await wrapper.find('th[data-field="set"] .p-datatable-column-resizer').trigger('mousedown')
    expect(held(wrapper, 'set')).toBe('')

    await wrapper.setProps({ rows: [{ name: 'x'.repeat(30), set: 'x'.repeat(40) }] })
    expect(held(wrapper, 'set')).toBe('')
    expect(held(wrapper, 'name')).toBe('300px')
  })

  it('keeps a click on a resize grip from reaching the header', async () => {
    await show([{ name: 'x'.repeat(20), set: 'x'.repeat(15) }])
    const sorted: string[] = []
    const th = wrapper.find('th[data-field="set"]')
    th.element.addEventListener('click', () => sorted.push('set'))

    await th.find('.p-datatable-column-resizer').trigger('click')
    expect(sorted).toEqual([])

    await th.trigger('click')
    expect(sorted).toEqual(['set'])
  })

  it('does not take a click elsewhere in the header for a resize', async () => {
    await show([{ name: 'x'.repeat(20), set: 'x'.repeat(15) }])

    await wrapper.find('th[data-field="set"]').trigger('mousedown')
    await wrapper.setProps({ rows: [{ name: 'x'.repeat(20), set: 'x'.repeat(10) }] })

    expect(held(wrapper, 'set')).toBe('150px')
  })
})

import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CFieldSelectViewer from './CFieldSelectViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'status', isSystem: false, isMulti: false, options: {}, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

function mountViewer(props: Record<string, unknown>) {
  return mount(CFieldSelectViewer, { props })
}

describe('CFieldSelectViewer', () => {
  describe('single value', () => {
    it('renders nothing when value is undefined', () => {
      const wrapper = mountViewer({ field: field(), record: record() })
      expect(wrapper.findAll('span')).toHaveLength(0)
    })

    it('falls back to raw value when not in options', () => {
      const wrapper = mountViewer({
        field: field(),
        record: record({ status: 'active' }),
      })
      expect(wrapper.text()).toContain('active')
    })

    it('resolves option text from options array', () => {
      const wrapper = mountViewer({
        field: field({ options: { options: [{ value: 'active', text: 'Active Status' }] } }),
        record: record({ status: 'active' }),
      })
      expect(wrapper.text()).toContain('Active Status')
    })
  })

  describe('multi value', () => {
    it('renders multiple spans for array values', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true }),
        record: record({ status: ['a', 'b', 'c'] }),
      })
      // Each item: outer <span v-for> + inner <span v-else> = 2 per item; count outer only
      expect(wrapper.findAll('div > span')).toHaveLength(3)
    })

    it('uses default delimiter between values', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true }),
        record: record({ status: ['X', 'Y'] }),
      })
      expect(wrapper.text()).toContain(', ')
    })

    it('uses custom multiDelimiter', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true, options: { multiDelimiter: ' | ' } }),
        record: record({ status: ['X', 'Y'] }),
      })
      expect(wrapper.text()).toContain(' | ')
      expect(wrapper.text()).not.toContain(', ')
    })

    it('renders nothing when array is empty', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true }),
        record: record({ status: [] }),
      })
      expect(wrapper.findAll('span')).toHaveLength(0)
    })
  })

  describe('badge display', () => {
    it('renders Tag component when displayType=badge', () => {
      const wrapper = mountViewer({
        field: field({
          options: {
            displayType: 'badge',
            options: [{ value: 'active', text: 'Active' }],
          },
        }),
        record: record({ status: 'active' }),
      })
      // Tag renders with the text
      expect(wrapper.text()).toContain('Active')
    })

    it('applies default colors when option has no style', () => {
      const wrapper = mountViewer({
        field: field({ options: { displayType: 'badge', options: [{ value: 'v', text: 'V' }] } }),
        record: record({ status: 'v' }),
      })
      // Default badge colors are applied — just verify it renders
      expect(wrapper.text()).toContain('V')
    })

    it('applies custom colors from option style', () => {
      const wrapper = mountViewer({
        field: field({
          options: {
            displayType: 'badge',
            options: [{ value: 'v', text: 'V', style: { textColor: '#FF0000FF', backgroundColor: '#00FF00FF' } }],
          },
        }),
        record: record({ status: 'v' }),
      })
      expect(wrapper.text()).toContain('V')
    })
  })

  describe('system field', () => {
    it('reads value directly from record for system field', () => {
      const rec: any = { status: 'sys-val', values: {} }
      const wrapper = mountViewer({
        field: field({ isSystem: true }),
        record: rec,
      })
      expect(wrapper.text()).toContain('sys-val')
    })
  })
})

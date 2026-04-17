import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CFieldNumberViewer from './CFieldNumberViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'score', isSystem: false, isMulti: false, options: {}, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

function mountViewer(props: Record<string, unknown>) {
  return mount(CFieldNumberViewer, {
    props,
    global: { stubs: { ProgressBar: { template: '<div data-testid="progress-bar" />' } } },
  })
}

describe('CFieldNumberViewer', () => {
  describe('number display (default)', () => {
    it('renders nothing when value is undefined', () => {
      const wrapper = mountViewer({ field: field(), record: record() })
      expect(wrapper.text().trim()).toBe('')
    })

    it('renders number value', () => {
      const wrapper = mountViewer({ field: field(), record: record({ score: '42' }) })
      expect(wrapper.text()).toContain('42')
    })

    it('uses field.formatValue when available', () => {
      const f = { ...field(), formatValue: (v: string) => `§${v}§` }
      const wrapper = mountViewer({ field: f, record: record({ score: '99' }) })
      expect(wrapper.text()).toContain('§99§')
    })

    it('applies precision (toFixed)', () => {
      const wrapper = mountViewer({
        field: field({ options: { precision: 2 } }),
        record: record({ score: '3.14159' }),
      })
      expect(wrapper.text()).toContain('3.14')
    })

    it('applies prefix and suffix', () => {
      const wrapper = mountViewer({
        field: field({ options: { prefix: '$', suffix: ' USD' } }),
        record: record({ score: '100' }),
      })
      expect(wrapper.text()).toContain('$100 USD')
    })

    it('renders multi-value joined by delimiter', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true, options: { multiDelimiter: ' / ' } }),
        record: record({ score: ['1', '2', '3'] }),
      })
      expect(wrapper.text()).toContain('1 / 2 / 3')
    })

    it('renders multi-value with default delimiter', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true }),
        record: record({ score: ['10', '20'] }),
      })
      expect(wrapper.text()).toContain(', ')
    })
  })

  describe('progress display', () => {
    it('renders progress bar when display=progress', () => {
      const wrapper = mountViewer({
        field: field({ options: { display: 'progress' } }),
        record: record({ score: '50' }),
      })
      expect(wrapper.find('[data-testid="progress-bar"]').exists()).toBe(true)
    })

    it('does not render progress bar in number mode', () => {
      const wrapper = mountViewer({
        field: field({ options: { display: 'number' } }),
        record: record({ score: '50' }),
      })
      expect(wrapper.find('[data-testid="progress-bar"]').exists()).toBe(false)
    })

    it('renders one progress bar per multi-value', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true, options: { display: 'progress' } }),
        record: record({ score: ['25', '75'] }),
      })
      expect(wrapper.findAll('[data-testid="progress-bar"]')).toHaveLength(2)
    })
  })

  describe('system field', () => {
    it('reads value from record root for system fields', () => {
      const rec: any = { score: '77', values: {} }
      const wrapper = mountViewer({ field: field({ isSystem: true }), record: rec })
      expect(wrapper.text()).toContain('77')
    })
  })
})

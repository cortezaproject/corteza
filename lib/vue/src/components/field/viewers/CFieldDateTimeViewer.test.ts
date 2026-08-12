import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import CFieldDateTimeViewer from './CFieldDateTimeViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'createdAt', isSystem: false, isMulti: false, options: {}, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

function mountViewer(props: Record<string, unknown>) {
  return mount(CFieldDateTimeViewer, { props })
}

describe('CFieldDateTimeViewer', () => {
  describe('empty / falsy values', () => {
    it('renders nothing when value is undefined', () => {
      const wrapper = mountViewer({ field: field(), record: record() })
      expect(wrapper.find('span').exists()).toBe(false)
    })

    it('renders nothing when value is empty string', () => {
      const wrapper = mountViewer({ field: field(), record: record({ createdAt: '' }) })
      expect(wrapper.find('span').exists()).toBe(false)
    })
  })

  describe('field.formatValue delegation', () => {
    it('uses formatValue when available', () => {
      const f = { ...field(), formatValue: (_v: string) => 'FORMATTED' }
      const wrapper = mountViewer({ field: f, record: record({ createdAt: '2024-01-01' }) })
      expect(wrapper.text()).toContain('FORMATTED')
    })

    it('joins multi-values with delimiter via formatValue', () => {
      const f = { ...field({ isMulti: true }), formatValue: (v: string) => `[${v}]` }
      const wrapper = mountViewer({
        field: f,
        record: record({ createdAt: ['2024-01-01', '2024-02-01'] }),
      })
      expect(wrapper.text()).toContain('[2024-01-01]')
      expect(wrapper.text()).toContain('[2024-02-01]')
    })
  })

  describe('fallback formatting', () => {
    it('renders a date string via toLocaleString', () => {
      const wrapper = mountViewer({
        field: field(),
        record: record({ createdAt: '2024-06-15T12:00:00Z' }),
      })
      const span = wrapper.find('span')
      expect(span.exists()).toBe(true)
      expect(span.text()).toBeTruthy()
    })

    it('onlyDate uses toLocaleDateString', () => {
      const dateStr = '2024-06-15T12:00:00Z'
      const expected = new Date(dateStr).toLocaleDateString()
      const wrapper = mountViewer({
        field: field({ options: { onlyDate: true } }),
        record: record({ createdAt: dateStr }),
      })
      expect(wrapper.text()).toContain(expected)
    })

    it('onlyTime uses toLocaleTimeString', () => {
      const dateStr = '2024-06-15T12:00:00Z'
      const expected = new Date(dateStr).toLocaleTimeString()
      const wrapper = mountViewer({
        field: field({ options: { onlyTime: true } }),
        record: record({ createdAt: dateStr }),
      })
      expect(wrapper.text()).toContain(expected)
    })

    it('returns raw string for invalid date', () => {
      const wrapper = mountViewer({ field: field(), record: record({ createdAt: 'not-a-date' }) })
      expect(wrapper.text()).toContain('not-a-date')
    })
  })

  describe('multi-value', () => {
    it('joins multiple dates with default delimiter', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true }),
        record: record({ createdAt: ['2024-01-01', '2024-02-01'] }),
      })
      expect(wrapper.text()).toContain(', ')
    })

    it('joins with custom delimiter', () => {
      const wrapper = mountViewer({
        field: field({ isMulti: true, options: { multiDelimiter: ' — ' } }),
        record: record({ createdAt: ['2024-01-01', '2024-02-01'] }),
      })
      expect(wrapper.text()).toContain(' — ')
    })
  })

  describe('outputRelative', () => {
    it('shows relative text for a past date', () => {
      const pastDate = new Date(Date.now() - 5 * 60 * 1000).toISOString() // 5 mins ago
      const wrapper = mountViewer({
        field: field({ options: { outputRelative: true } }),
        record: record({ createdAt: pastDate }),
      })
      // With no i18n messages, vue-i18n returns the key itself
      expect(wrapper.find('span').exists()).toBe(true)
      expect(wrapper.text()).toBeTruthy()
    })

    it('shows relative text for a future date', () => {
      const futureDate = new Date(Date.now() + 5 * 60 * 1000).toISOString() // 5 mins future
      const wrapper = mountViewer({
        field: field({ options: { outputRelative: true } }),
        record: record({ createdAt: futureDate }),
      })
      expect(wrapper.find('span').exists()).toBe(true)
      expect(wrapper.text()).toBeTruthy()
    })
  })

  describe('system field', () => {
    it('reads value from record root for system fields', () => {
      const rec: any = { createdAt: '2024-06-15T12:00:00Z', values: {} }
      const wrapper = mountViewer({ field: field({ isSystem: true }), record: rec })
      expect(wrapper.find('span').exists()).toBe(true)
    })
  })
})

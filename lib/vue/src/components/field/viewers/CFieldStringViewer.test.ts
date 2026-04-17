import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CFieldStringViewer from './CFieldStringViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'title', isMulti: false, isSystem: false, options: {}, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

describe('CFieldStringViewer', () => {
  it('renders nothing when value is undefined', () => {
    const wrapper = mount(CFieldStringViewer, {
      props: { field: field(), record: record() },
    })
    expect(wrapper.find('p').exists()).toBe(false)
  })

  it('renders the string value', () => {
    const wrapper = mount(CFieldStringViewer, {
      props: { field: field(), record: record({ title: 'Hello World' }) },
    })
    expect(wrapper.find('p').text()).toBe('Hello World')
  })

  it('joins multi-value array with default delimiter', () => {
    const wrapper = mount(CFieldStringViewer, {
      props: {
        field: field({ isMulti: true }),
        record: record({ title: ['A', 'B', 'C'] }),
      },
    })
    expect(wrapper.find('p').text()).toBe('A, B, C')
  })

  it('joins multi-value array with custom delimiter', () => {
    const wrapper = mount(CFieldStringViewer, {
      props: {
        field: field({ isMulti: true, options: { multiDelimiter: ' | ' } }),
        record: record({ title: ['X', 'Y'] }),
      },
    })
    expect(wrapper.find('p').text()).toBe('X | Y')
  })

  it('reads system field from record root (not values)', () => {
    const wrapper = mount(CFieldStringViewer, {
      props: {
        field: field({ name: 'createdAt', isSystem: true }),
        record: { createdAt: '2024-01-01', values: {} },
      },
    })
    expect(wrapper.find('p').text()).toBe('2024-01-01')
  })

  it('adds rt-content class when useRichTextEditor is set', () => {
    const wrapper = mount(CFieldStringViewer, {
      props: {
        field: field({ options: { useRichTextEditor: true } }),
        record: record({ title: 'rich' }),
      },
    })
    expect(wrapper.find('p').classes()).toContain('rt-content')
  })

  it('adds multiline class for multi fields', () => {
    const wrapper = mount(CFieldStringViewer, {
      props: {
        field: field({ isMulti: true }),
        record: record({ title: ['a', 'b'] }),
      },
    })
    expect(wrapper.find('p').classes()).toContain('multiline')
  })
})

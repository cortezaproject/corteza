import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import CFieldBoolViewer from './CFieldBoolViewer.vue'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      general: { label: { yes: 'Yes', no: 'No' } },
    },
  },
})

function field(options: Record<string, unknown> = {}) {
  return { name: 'active', isSystem: false, options }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

function mountViewer(fieldOpts: Record<string, unknown>, value: unknown) {
  return mount(CFieldBoolViewer, {
    props: { field: field(fieldOpts), record: record({ active: value }) },
    global: { plugins: [i18n] },
  })
}

describe('CFieldBoolViewer', () => {
  it('shows "Yes" for value "1"', () => {
    expect(mountViewer({}, '1').find('span').text()).toBe('Yes')
  })

  it('shows "Yes" for boolean true', () => {
    expect(mountViewer({}, true).find('span').text()).toBe('Yes')
  })

  it('shows "No" for value "0"', () => {
    expect(mountViewer({}, '0').find('span').text()).toBe('No')
  })

  it('shows "No" for undefined value', () => {
    expect(mountViewer({}, undefined).find('span').text()).toBe('No')
  })

  it('uses custom trueLabel', () => {
    expect(mountViewer({ trueLabel: 'Active' }, '1').find('span').text()).toBe('Active')
  })

  it('uses custom falseLabel', () => {
    expect(mountViewer({ falseLabel: 'Inactive' }, '0').find('span').text()).toBe('Inactive')
  })
})

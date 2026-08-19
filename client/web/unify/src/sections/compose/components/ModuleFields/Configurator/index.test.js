import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

import CFieldConfigurator from './index.vue'

const toastDanger = vi.fn()

// The dialog and the tabs are chrome; the test needs the footer's Save button
// and nothing else, so both stubs just render their slots.
const slotHost = names => ({
  template: `<div>${names.map(n => `<slot name="${n}" />`).join('')}<slot /></div>`,
})

function mountConfigurator(expressions) {
  return mount(CFieldConfigurator, {
    props: {
      visible: true,
      field: { name: 'title', kind: 'RecordID', label: 'Title', expressions },
    },
    shallow: true,
    global: {
      mocks: { $t: k => k },
      provide: { $toast: { toastDanger } },
      stubs: {
        Dialog: slotHost(['footer']),
        Tabs: slotHost([]),
        TabList: slotHost([]),
        TabPanels: slotHost([]),
        Tab: true,
        TabPanel: true,
        // No explicit emit: the parent's @click lands on the root element as a
        // fallthrough attribute, and emitting as well would fire it twice.
        Button: { template: '<button>{{ label }}</button>', props: ['label'] },
      },
    },
  })
}

async function clickSave(wrapper) {
  const save = wrapper.findAll('button').at(-1)
  await save.trigger('click')
}

describe('CFieldConfigurator save guard', () => {
  it('refuses to save a validator with no expression', async () => {
    toastDanger.mockClear()
    const wrapper = mountConfigurator({ validators: [{ test: '', error: 'Enter a value' }] })

    await clickSave(wrapper)

    expect(wrapper.emitted('save')).toBeUndefined()
    expect(toastDanger).toHaveBeenCalledWith('field.expressions.incomplete')
  })

  it('refuses to save a validator with no error message', async () => {
    toastDanger.mockClear()
    const wrapper = mountConfigurator({ validators: [{ test: 'value == ""', error: '' }] })

    await clickSave(wrapper)

    expect(wrapper.emitted('save')).toBeUndefined()
  })

  it('refuses to save a sanitizer with no expression', async () => {
    toastDanger.mockClear()
    const wrapper = mountConfigurator({ sanitizers: [''] })

    await clickSave(wrapper)

    expect(wrapper.emitted('save')).toBeUndefined()
  })

  it('saves once every row is filled in', async () => {
    toastDanger.mockClear()
    const wrapper = mountConfigurator({
      sanitizers: ['trim(value)'],
      validators: [{ test: 'value == ""', error: 'Enter a value' }],
    })

    await clickSave(wrapper)

    expect(wrapper.emitted('save')).toHaveLength(1)
    expect(toastDanger).not.toHaveBeenCalled()
  })

  it('saves a field that has no expressions at all', async () => {
    toastDanger.mockClear()
    const wrapper = mountConfigurator(undefined)

    await clickSave(wrapper)

    expect(wrapper.emitted('save')).toHaveLength(1)
  })
})

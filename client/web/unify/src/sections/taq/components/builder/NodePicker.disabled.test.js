import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'

// A construct the server lists with `disabled: true` exists in the catalog but
// cannot run on this server: the picker shows it, explains why, and never adds it.

const store = { functions: [], triggers: [] }

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@planetcrust/human-vue', () => ({
  useAutomationStore: () => store,
}))

import NodePicker from './NodePicker.vue'

const corredorExec = {
  ref: 'corredorExec',
  kind: 'function',
  groups: ['Corredor'],
  meta: { short: 'Run Corredor script', description: 'Runs a Corredor server script' },
}

const otherOff = {
  ref: 'someOtherStep',
  kind: 'function',
  groups: ['Corredor'],
  meta: { short: 'Some other step', description: 'Does something else' },
  disabled: true,
}

const logStep = {
  ref: 'logInfo',
  kind: 'function',
  groups: ['Corredor'],
  meta: { short: 'Log message' },
}

let wrapper

// Mounts the picker and opens the Corredor group (Branches is listed first).
async function mountPicker(functions) {
  store.functions = functions
  wrapper = mount(NodePicker, {
    global: {
      mocks: { $t: k => k },
      stubs: { IconField: true, InputIcon: true, InputText: true, TaqIcon: true },
    },
  })
  await wrapper
    .findAll('button')
    .find(b => b.text() === 'Corredor')
    .trigger('click')
  return wrapper
}

function nodeItem(label) {
  return wrapper.findAll('.node-item').find(n => n.text().includes(label))
}

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('NodePicker — a construct the server has turned off', () => {
  it('shows it greyed, without hover emphasis, marked disabled', async () => {
    await mountPicker([{ ...corredorExec, disabled: true }, logStep])

    const item = nodeItem('Run Corredor script')
    expect(item.classes()).toEqual(expect.arrayContaining(['opacity-60', 'cursor-not-allowed']))
    expect(item.classes()).not.toContain('hover:bg-emphasis')
    expect(item.classes()).not.toContain('cursor-pointer')
    expect(item.attributes('aria-disabled')).toBe('true')
  })

  it('names Corredor as the reason for the Corredor step', async () => {
    await mountPicker([{ ...corredorExec, disabled: true }])

    expect(nodeItem('Run Corredor script').find('.node-disabled-hint').text()).toBe(
      'builder.nodePicker.disabled.corredor',
    )
  })

  it('uses the generic reason for any other turned-off construct', async () => {
    await mountPicker([otherOff])

    expect(nodeItem('Some other step').find('.node-disabled-hint').text()).toBe(
      'builder.nodePicker.disabled.generic',
    )
  })

  it('ignores a click on it', async () => {
    await mountPicker([{ ...corredorExec, disabled: true }])

    await nodeItem('Run Corredor script').trigger('click')

    expect(wrapper.emitted('select')).toBeUndefined()
    expect(wrapper.emitted('close')).toBeUndefined()
  })

  it('adds the same construct as usual when the server runs it', async () => {
    await mountPicker([corredorExec])

    const item = nodeItem('Run Corredor script')
    expect(item.classes()).toEqual(expect.arrayContaining(['cursor-pointer', 'hover:bg-emphasis']))
    expect(item.attributes('aria-disabled')).toBe('false')
    expect(item.find('.node-disabled-hint').exists()).toBe(false)

    await item.trigger('click')

    expect(wrapper.emitted('select')[0][0]).toMatchObject({
      ref: 'corredorExec',
      label: 'Run Corredor script',
      disabled: false,
    })
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('leaves the enabled steps beside it addable', async () => {
    await mountPicker([{ ...corredorExec, disabled: true }, logStep])

    await nodeItem('Log message').trigger('click')

    expect(wrapper.emitted('select')[0][0]).toMatchObject({ ref: 'logInfo' })
  })
})

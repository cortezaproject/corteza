import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

import AutomationButtonsEditor from './AutomationButtonsEditor.vue'

// What GET /compose/automation/ answers with: one client and one server script
// for compose, and one belonging to another app entirely.
const SCRIPTS = [
  {
    name: '/client-scripts/compose/agent-sandbox/ContactGreet.js:default',
    label: 'Greet contact (client)',
    description: 'Shows a toast in the browser',
    triggers: [
      {
        eventTypes: ['onManual'],
        resourceTypes: ['compose:record'],
        uiProps: [{ name: 'app', value: 'compose' }],
      },
    ],
  },
  {
    name: '/server-scripts/agent-sandbox/ContactActivate.js:default',
    label: 'Activate contact (server)',
    triggers: [
      {
        eventTypes: ['onManual'],
        resourceTypes: ['compose:record'],
        uiProps: [{ name: 'app', value: 'compose' }],
      },
    ],
  },
  {
    name: '/client-scripts/admin/UserGreet.js:default',
    label: 'Greet user (client)',
    triggers: [
      {
        eventTypes: ['onManual'],
        resourceTypes: ['system:user'],
        uiProps: [{ name: 'app', value: 'admin' }],
      },
    ],
  },
]

const passThrough = { template: '<div><slot /></div>' }
const itemList = { name: 'CFormItemList', props: ['items'], template: '<div />' }

let automationList

function mountEditor(buttons = []) {
  automationList = vi.fn(() => Promise.resolve({ set: SCRIPTS }))

  return mount(AutomationButtonsEditor, {
    props: { buttons },
    global: {
      stubs: {
        CFormGroup: passThrough,
        CFormList: passThrough,
        CFormItemList: itemList,
        CInputExpression: true,
        CExpressionHint: true,
        Tabs: passThrough,
        TabList: passThrough,
        Tab: passThrough,
        TabPanels: passThrough,
        TabPanel: passThrough,
        InputText: true,
        Select: true,
        Button: true,
        Divider: true,
        ProgressSpinner: true,
        Tag: true,
      },
      mocks: { $t: k => k },
      provide: {
        $AutomationAPI: {
          triggerList: () => Promise.resolve({ set: [] }),
          ngAutomationListCancellable: () => ({ response: () => Promise.resolve({ set: [] }) }),
        },
        $ComposeAPI: { automationList },
      },
    },
  })
}

// The three tabs, in template order
const scriptsTab = wrapper => wrapper.findAllComponents(itemList)[2]

beforeEach(() => {
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

describe('AutomationButtonsEditor scripts tab', () => {
  it('lists the manual Corredor scripts this app may offer', async () => {
    const wrapper = mountEditor()
    await flushPromises()

    expect(automationList).toHaveBeenCalledWith({
      eventTypes: ['onManual'],
      excludeInvalid: true,
    })

    expect(scriptsTab(wrapper).props('items')).toEqual([
      {
        script: '/client-scripts/compose/agent-sandbox/ContactGreet.js:default',
        label: 'Greet contact (client)',
        resourceType: 'compose:record',
        description: 'Shows a toast in the browser',
        isScript: true,
      },
      {
        script: '/server-scripts/agent-sandbox/ContactActivate.js:default',
        label: 'Activate contact (server)',
        resourceType: 'compose:record',
        description: undefined,
        isScript: true,
      },
    ])
  })

  it('leaves out a script already configured as a button', async () => {
    const wrapper = mountEditor([
      { label: 'Greet', script: '/client-scripts/compose/agent-sandbox/ContactGreet.js:default' },
    ])
    await flushPromises()

    expect(
      scriptsTab(wrapper)
        .props('items')
        .map(i => i.script),
    ).toEqual(['/server-scripts/agent-sandbox/ContactActivate.js:default'])
  })

  it('adds the picked script as a button that names it', async () => {
    const wrapper = mountEditor()
    await flushPromises()

    scriptsTab(wrapper).vm.$emit('select', scriptsTab(wrapper).props('items')[0])

    expect(wrapper.emitted('update:buttons')[0][0]).toEqual([
      {
        label: 'Greet contact (client)',
        variant: 'primary',
        resourceType: 'compose:record',
        script: '/client-scripts/compose/agent-sandbox/ContactGreet.js:default',
        scriptType: 'script',
      },
    ])
  })

  it('still lists workflows and TAQs when no script can be read', async () => {
    automationList = vi.fn(() => Promise.reject(new Error('no corredor')))

    const wrapper = mount(AutomationButtonsEditor, {
      props: { buttons: [] },
      global: {
        stubs: {
          CFormGroup: passThrough,
          CFormList: passThrough,
          CFormItemList: itemList,
          Tabs: passThrough,
          TabList: passThrough,
          Tab: passThrough,
          TabPanels: passThrough,
          TabPanel: passThrough,
          CInputExpression: true,
          CExpressionHint: true,
          InputText: true,
          Select: true,
          Button: true,
          Divider: true,
          ProgressSpinner: true,
          Tag: true,
        },
        mocks: { $t: k => k },
        provide: {
          $AutomationAPI: {
            triggerList: () => Promise.resolve({ set: [] }),
            ngAutomationListCancellable: () => ({ response: () => Promise.resolve({ set: [] }) }),
          },
          $ComposeAPI: { automationList },
        },
      },
    })
    await flushPromises()

    expect(scriptsTab(wrapper).props('items')).toEqual([])
  })
})

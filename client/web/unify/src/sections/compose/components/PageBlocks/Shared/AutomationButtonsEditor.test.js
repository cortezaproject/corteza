import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

import AutomationButtonsEditor from './AutomationButtonsEditor.vue'

// What GET /compose/automation/ answers with: one client and one server script
// for compose, and one belonging to another app entirely.
const CONTACT_CONSTRAINTS = [
  { name: 'module', value: ['agent-contact'] },
  { name: 'namespace', value: ['agent-sandbox'] },
]

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
        constraints: CONTACT_CONSTRAINTS,
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
        constraints: CONTACT_CONSTRAINTS,
      },
    ],
  },
  {
    name: '/server-scripts/agent-sandbox/RebuildSpace.js:default',
    label: 'Rebuild the space (server)',
    triggers: [
      {
        eventTypes: ['onManual'],
        resourceTypes: ['compose:namespace'],
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

const NAMESPACE = { slug: 'agent-sandbox', name: 'Agent sandbox' }
const CONTACT_MODULE = { handle: 'agent-contact', name: 'Contact' }
const RECORD_PAGE = { pageID: '1', moduleID: '2', isRecordPage: true }
const DASHBOARD_PAGE = { pageID: '3', moduleID: '0', isRecordPage: false }

const passThrough = { template: '<div><slot /></div>' }
const itemList = { name: 'CFormItemList', props: ['items'], template: '<div />' }

let automationList

function mountEditor(buttons = [], props = {}) {
  automationList = vi.fn(() => Promise.resolve({ set: SCRIPTS }))

  return mount(AutomationButtonsEditor, {
    props: { buttons, ...props },
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
        constraints: CONTACT_CONSTRAINTS,
        constraintChips: ['module = agent-contact', 'namespace = agent-sandbox'],
        applies: true,
      },
      {
        script: '/server-scripts/agent-sandbox/ContactActivate.js:default',
        label: 'Activate contact (server)',
        resourceType: 'compose:record',
        description: undefined,
        isScript: true,
        constraints: CONTACT_CONSTRAINTS,
        constraintChips: ['module = agent-contact', 'namespace = agent-sandbox'],
        applies: true,
      },
      {
        script: '/server-scripts/agent-sandbox/RebuildSpace.js:default',
        label: 'Rebuild the space (server)',
        resourceType: 'compose:namespace',
        description: undefined,
        isScript: true,
        constraints: [],
        constraintChips: [],
        applies: true,
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
    ).toEqual([
      '/server-scripts/agent-sandbox/ContactActivate.js:default',
      '/server-scripts/agent-sandbox/RebuildSpace.js:default',
    ])
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

  it('offers a record-bound script on a record page', async () => {
    const wrapper = mountEditor([], {
      page: RECORD_PAGE,
      namespace: NAMESPACE,
      module: CONTACT_MODULE,
    })
    await flushPromises()

    expect(
      scriptsTab(wrapper)
        .props('items')
        .map(i => i.script),
    ).toEqual([
      '/client-scripts/compose/agent-sandbox/ContactGreet.js:default',
      '/server-scripts/agent-sandbox/ContactActivate.js:default',
      '/server-scripts/agent-sandbox/RebuildSpace.js:default',
    ])
  })

  it('leaves a record-bound script out on a page showing no record', async () => {
    const wrapper = mountEditor([], { page: DASHBOARD_PAGE, namespace: NAMESPACE })
    await flushPromises()

    expect(
      scriptsTab(wrapper)
        .props('items')
        .map(i => i.script),
    ).toEqual(['/server-scripts/agent-sandbox/RebuildSpace.js:default'])
  })

  it('leaves a record-bound script out for a block that hands no record', async () => {
    const wrapper = mountEditor([], {
      page: RECORD_PAGE,
      namespace: NAMESPACE,
      canSupplyRecord: false,
    })
    await flushPromises()

    expect(
      scriptsTab(wrapper)
        .props('items')
        .map(i => i.script),
    ).toEqual(['/server-scripts/agent-sandbox/RebuildSpace.js:default'])
  })

  it('lists a script the page contradicts, and refuses to add it', async () => {
    const wrapper = mountEditor([], {
      page: RECORD_PAGE,
      namespace: NAMESPACE,
      module: { handle: 'agent-deal', name: 'Deal' },
    })
    await flushPromises()

    const items = scriptsTab(wrapper).props('items')

    expect(items.map(i => i.applies)).toEqual([false, false, true])

    scriptsTab(wrapper).vm.$emit('select', items[0])

    expect(wrapper.emitted('update:buttons')).toBeUndefined()
  })

  it('adds a script whose constraints the page satisfies', async () => {
    const wrapper = mountEditor([], {
      page: RECORD_PAGE,
      namespace: NAMESPACE,
      module: CONTACT_MODULE,
    })
    await flushPromises()

    scriptsTab(wrapper).vm.$emit('select', scriptsTab(wrapper).props('items')[0])

    expect(wrapper.emitted('update:buttons')[0][0][0].script).toBe(
      '/client-scripts/compose/agent-sandbox/ContactGreet.js:default',
    )
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

describe('AutomationButtonsEditor visibility condition', () => {
  it('stores the condition on the button and drops it once cleared', async () => {
    const wrapper = mountEditor([{ label: 'Go', workflowID: 'W1', scriptType: 'workflow' }])
    await flushPromises()

    wrapper.vm.updateVisibility(0, 'record.values.status == "draft"')
    expect(wrapper.emitted('update:buttons')[0][0][0].visibility).toEqual({
      expression: 'record.values.status == "draft"',
    })

    await wrapper.setProps({ buttons: wrapper.emitted('update:buttons')[0][0] })
    wrapper.vm.updateVisibility(0, '')
    expect(wrapper.emitted('update:buttons')[1][0][0]).not.toHaveProperty('visibility')
  })
})

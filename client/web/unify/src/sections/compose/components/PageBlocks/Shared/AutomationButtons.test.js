import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

import { compose } from '@planetcrust/human-js'

import AutomationButtons from './AutomationButtons.vue'

// `input` is decoded into expr.Vars server-side, and expr.Vars reads only
// {"@type":…,"@value":…}. A bare value fails the whole request at parse time,
// before the automation is even looked up, so every entry is enveloped.

const NAMESPACE = { namespaceID: 'N1' }
const PAGE = { pageID: 'P1' }
const MODULE = { moduleID: 'M1' }
const RECORD = { recordID: 'R1', ownedBy: '3' }

let ngExec
let wfExec
let dispatch
let evaluate

function mountButtons(props, { conditions } = {}) {
  ngExec = vi.fn(() => Promise.resolve({}))
  wfExec = vi.fn(() => Promise.resolve({}))
  dispatch = vi.fn(() => Promise.resolve(null))
  evaluate = vi.fn(() => Promise.resolve(conditions || {}))

  return mount(AutomationButtons, {
    props: { namespace: NAMESPACE, page: PAGE, ...props },
    global: {
      stubs: { Button: { props: ['label'], template: '<button>{{ label }}</button>' } },
      directives: { tooltip: {} },
      provide: {
        $AutomationAPI: { ngAutomationExec: ngExec, workflowExec: wfExec },
        $ScriptBus: { Dispatch: dispatch },
        $SystemAPI: { expressionEvaluate: evaluate },
        $Auth: { user: { userID: '42' } },
        $toast: null,
      },
      mocks: { $t: k => k },
    },
  })
}

async function press(wrapper) {
  await wrapper.find('button').trigger('click')
  await flushPromises()
}

beforeEach(() => {
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

describe('AutomationButtons input', () => {
  it('envelopes every value it sends to a TAQ', async () => {
    const wrapper = mountButtons({
      buttons: [{ label: 'Go', automationID: 'A1' }],
      module: MODULE,
      record: RECORD,
      records: [RECORD],
      filter: 'name = "x"',
    })
    await press(wrapper)

    expect(ngExec).toHaveBeenCalledTimes(1)
    expect(ngExec.mock.calls[0][0].input).toEqual({
      namespace: { '@type': 'ComposeNamespace', '@value': NAMESPACE },
      page: { '@type': 'ComposePage', '@value': PAGE },
      module: { '@type': 'ComposeModule', '@value': MODULE },
      record: { '@type': 'ComposeRecord', '@value': RECORD },
      selected: { '@type': 'Array', '@value': [RECORD] },
      filter: { '@type': 'String', '@value': 'name = "x"' },
    })
  })

  it('envelopes the selection and filter it sends to a workflow', async () => {
    const wrapper = mountButtons({
      buttons: [{ label: 'Go', workflowID: 'W1' }],
      records: [RECORD],
      filter: 'q',
    })
    await press(wrapper)

    expect(wfExec.mock.calls[0][0].input.selected).toEqual({
      '@type': 'Array',
      '@value': [RECORD],
    })
    expect(wfExec.mock.calls[0][0].input.filter).toEqual({ '@type': 'String', '@value': 'q' })
  })

  it('omits what the page has not got', async () => {
    const wrapper = mountButtons({ buttons: [{ label: 'Go', automationID: 'A1' }] })
    await press(wrapper)

    expect(Object.keys(ngExec.mock.calls[0][0].input).sort()).toEqual(['namespace', 'page'])
  })
})

// A Corredor script is not executed here: it is handed to the script bus, which
// is what knows whether the script runs in the browser or inside Corredor.
describe('AutomationButtons script', () => {
  const namespace = new compose.Namespace({ namespaceID: '1001', slug: 'agent-sandbox' })
  const recordModule = new compose.Module(
    { moduleID: '1002', namespaceID: '1001', handle: 'agent-contact', fields: [{ name: 'name' }] },
    namespace,
  )
  const record = new compose.Record(recordModule, { recordID: '1003' })

  it('dispatches a record event for the script the button names', async () => {
    const wrapper = mountButtons({
      buttons: [
        {
          label: 'Greet',
          script: '/client-scripts/Greet.js:default',
          resourceType: 'compose:record',
        },
      ],
      namespace,
      module: recordModule,
      record,
    })
    await press(wrapper)

    expect(ngExec).not.toHaveBeenCalled()
    expect(wfExec).not.toHaveBeenCalled()
    expect(dispatch).toHaveBeenCalledTimes(1)

    const [event, script] = dispatch.mock.calls[0]
    expect(script).toBe('/client-scripts/Greet.js:default')
    expect(event.resourceType).toBe('compose:record')
    expect(event.eventType).toBe('onManual')
    expect(event.args.record).toEqual(record)
    expect(event.args.namespace).toEqual(namespace)
    expect(event.args.module).toEqual(recordModule)

    // The record goes by reference — a client script mutating it changes the
    // record the page is showing
    event.args.record.values.name = 'written by a script'
    expect(record.values.name).toBe('written by a script')

    // Constraints are matched against the namespace and module the page knows,
    // not against the module the record carries: a cached module does not know
    // the namespace it lives in
    expect(event.match({ Name: () => 'namespace', Match: v => v === 'agent-sandbox' })).toBe(true)
    expect(event.match({ Name: () => 'module', Match: v => v === 'agent-contact' })).toBe(true)
  })

  it('dispatches a namespace event when the trigger is bound to the namespace', async () => {
    const wrapper = mountButtons({
      buttons: [
        {
          label: 'Rebuild',
          script: '/server-scripts/Rebuild.js:default',
          resourceType: 'compose:namespace',
        },
      ],
      namespace,
    })
    await press(wrapper)

    expect(dispatch.mock.calls[0][0].resourceType).toBe('compose:namespace')
    expect(dispatch.mock.calls[0][0].args.namespace).toEqual(namespace)
  })

  it('refreshes the block once the script has run', async () => {
    const wrapper = mountButtons({
      buttons: [
        {
          label: 'Greet',
          script: '/client-scripts/Greet.js:default',
          resourceType: 'compose:record',
        },
      ],
      namespace,
      module: recordModule,
      record,
    })
    await press(wrapper)

    expect(wrapper.emitted('refresh')).toHaveLength(1)
  })

  it('does not refresh when the script aborts', async () => {
    const wrapper = mountButtons({
      buttons: [
        {
          label: 'Greet',
          script: '/client-scripts/Greet.js:default',
          resourceType: 'compose:record',
        },
      ],
      namespace,
      module: recordModule,
      record,
    })
    dispatch.mockRejectedValue(new Error('Aborted'))
    await press(wrapper)

    expect(wrapper.emitted('refresh')).toBeUndefined()
  })
})

// A button naming a script the server does not offer, or asking for a resource
// this page has not got, is marked rather than silently doing nothing.
describe('AutomationButtons unrunnable script', () => {
  const namespace = new compose.Namespace({ namespaceID: '1001', slug: 'agent-sandbox' })
  const recordModule = new compose.Module(
    { moduleID: '1002', namespaceID: '1001', handle: 'agent-contact', fields: [{ name: 'name' }] },
    namespace,
  )
  const record = new compose.Record(recordModule, { recordID: '1003' })

  const GREET = {
    label: 'Greet',
    script: '/client-scripts/Greet.js:default',
    resourceType: 'compose:record',
    variant: 'primary',
  }

  // The stub keeps the props the flag is made of, and the directive parks the
  // tooltip where a test can read it.
  const buttonStub = {
    props: ['label', 'severity', 'outlined', 'loading', 'disabled', 'size'],
    template: '<button :data-severity="severity" :data-outlined="String(!!outlined)" />',
  }

  const tooltip = {
    mounted(el, { value }) {
      if (value) el.setAttribute('data-tooltip', value)
    },
  }

  let toastWarning

  function mountFlagged(props, uiHooks) {
    toastWarning = vi.fn()
    dispatch = vi.fn(() => Promise.resolve(null))

    return mount(AutomationButtons, {
      props: { namespace, page: PAGE, ...props },
      global: {
        stubs: { Button: buttonStub },
        directives: { tooltip },
        provide: {
          $AutomationAPI: { ngAutomationExec: vi.fn(), workflowExec: vi.fn() },
          $ScriptBus: { Dispatch: dispatch },
          $Auth: { user: { userID: '42' } },
          $toast: { toastWarning },
          $UIHooks: uiHooks,
        },
        mocks: { $t: k => k },
      },
    })
  }

  const registry = scripts => ({ FindByScript: name => scripts.find(s => s === name) })

  it('flags a button naming a script the server does not offer', async () => {
    const wrapper = mountFlagged({ buttons: [GREET], module: recordModule, record }, registry([]))

    const button = wrapper.find('button')
    expect(button.attributes('data-severity')).toBe('danger')
    expect(button.attributes('data-outlined')).toBe('true')
    expect(button.attributes('data-tooltip')).toBe('block.automation.scriptNotLoaded')
  })

  it('explains itself when the flagged button is pressed, and dispatches nothing', async () => {
    const wrapper = mountFlagged({ buttons: [GREET], module: recordModule, record }, registry([]))
    await press(wrapper)

    expect(dispatch).not.toHaveBeenCalled()
    expect(toastWarning).toHaveBeenCalledWith('block.automation.scriptNotLoaded')
  })

  it('flags a record-bound button on a page carrying no record', async () => {
    const wrapper = mountFlagged({ buttons: [GREET] }, registry([GREET.script]))

    const button = wrapper.find('button')
    expect(button.attributes('data-severity')).toBe('danger')
    expect(button.attributes('data-tooltip')).toBe('block.automation.noRecord')
  })

  it('flags a record-bound selection button of a record list', async () => {
    const wrapper = mountFlagged(
      { buttons: [GREET], module: recordModule, records: [record] },
      registry([GREET.script]),
    )

    expect(wrapper.find('button').attributes('data-tooltip')).toBe('block.automation.noRecord')
  })

  it('leaves a button it can run alone', async () => {
    const wrapper = mountFlagged(
      { buttons: [GREET], module: recordModule, record },
      registry([GREET.script]),
    )

    const button = wrapper.find('button')
    expect(button.attributes('data-severity')).toBeUndefined()
    expect(button.attributes('data-outlined')).toBe('false')
    expect(button.attributes('data-tooltip')).toBeUndefined()

    await press(wrapper)
    expect(dispatch).toHaveBeenCalledTimes(1)
  })

  it('judges nothing when no script registry is installed', async () => {
    const wrapper = mountFlagged({ buttons: [GREET], module: recordModule, record }, null)

    expect(wrapper.find('button').attributes('data-severity')).toBeUndefined()
    expect(wrapper.find('button').attributes('data-tooltip')).toBeUndefined()
  })

  it('leaves a workflow button unjudged, registry or not', async () => {
    const wrapper = mountFlagged(
      { buttons: [{ label: 'Go', workflowID: 'W1', resourceType: 'compose:record' }] },
      registry([]),
    )

    expect(wrapper.find('button').attributes('data-severity')).toBeUndefined()
    expect(wrapper.find('button').attributes('data-tooltip')).toBeUndefined()
  })
})

// A button's condition is judged by the server with the same variables a
// block's visibility gets, and the button is kept out of the row until it holds.
describe('AutomationButtons visibility', () => {
  const labels = wrapper => wrapper.findAll('button').map(b => b.text())

  it('shows only the buttons whose condition holds', async () => {
    const wrapper = mountButtons(
      {
        buttons: [
          { label: 'Always', workflowID: 'W1' },
          { label: 'Hidden', workflowID: 'W2', visibility: { expression: 'false' } },
          { label: 'Shown', workflowID: 'W3', visibility: { expression: 'true' } },
        ],
        record: RECORD,
      },
      { conditions: { 1: false, 2: true } },
    )
    await flushPromises()

    expect(evaluate).toHaveBeenCalledTimes(1)
    expect(evaluate.mock.calls[0][0].expressions).toEqual({ 1: 'false', 2: 'true' })
    expect(evaluate.mock.calls[0][0].variables.user).toEqual({ userID: '42' })
    expect(labels(wrapper)).toEqual(['Always', 'Shown'])
  })

  it('keeps a conditioned button hidden until the server has answered', async () => {
    let answer
    const wrapper = mountButtons({
      buttons: [
        { label: 'Always', workflowID: 'W1' },
        { label: 'Later', workflowID: 'W2', visibility: { expression: 'true' } },
      ],
    })
    evaluate.mockImplementation(() => new Promise(resolve => (answer = resolve)))
    await wrapper.setProps({ buttons: [...wrapper.props('buttons')] })

    expect(labels(wrapper)).toEqual(['Always'])

    answer({ 1: true })
    await flushPromises()
    expect(labels(wrapper)).toEqual(['Always', 'Later'])
  })

  it('asks the server nothing when no button carries a condition', async () => {
    const wrapper = mountButtons({
      buttons: [
        { label: 'A', workflowID: 'W1' },
        { label: 'B', workflowID: 'W2' },
      ],
    })
    await flushPromises()

    expect(evaluate).not.toHaveBeenCalled()
    expect(labels(wrapper)).toEqual(['A', 'B'])
  })

  it('runs the button the pressed index names, not its place in the row', async () => {
    const wrapper = mountButtons(
      {
        buttons: [
          { label: 'Hidden', workflowID: 'W1', visibility: { expression: 'false' } },
          { label: 'Shown', workflowID: 'W2' },
        ],
      },
      { conditions: { 0: false } },
    )
    await flushPromises()
    await press(wrapper)

    expect(wfExec.mock.calls[0][0].workflowID).toBe('W2')
  })
})

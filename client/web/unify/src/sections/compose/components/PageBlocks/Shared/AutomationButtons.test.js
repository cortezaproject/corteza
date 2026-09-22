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

function mountButtons(props) {
  ngExec = vi.fn(() => Promise.resolve({}))
  wfExec = vi.fn(() => Promise.resolve({}))
  dispatch = vi.fn(() => Promise.resolve(null))

  return mount(AutomationButtons, {
    props: { namespace: NAMESPACE, page: PAGE, ...props },
    global: {
      stubs: { Button: { props: ['label'], template: '<button />' } },
      provide: {
        $AutomationAPI: { ngAutomationExec: ngExec, workflowExec: wfExec },
        $ScriptBus: { Dispatch: dispatch },
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

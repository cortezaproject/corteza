import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

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

function mountButtons(props) {
  ngExec = vi.fn(() => Promise.resolve({}))
  wfExec = vi.fn(() => Promise.resolve({}))

  return mount(AutomationButtons, {
    props: { namespace: NAMESPACE, page: PAGE, ...props },
    global: {
      stubs: { Button: { props: ['label'], template: '<button />' } },
      provide: {
        $AutomationAPI: { ngAutomationExec: ngExec, workflowExec: wfExec },
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

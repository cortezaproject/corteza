import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// A configured corredorExec step stores its script as the path inside the
// extension. The preview reads the label the script declares, and keeps the
// path underneath it, so a step never reads as a bare file name.

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (k: string) => k }) }))

import CViewCorredorScript from './CViewCorredorScript.vue'
import { resetCorredorScriptCache } from './corredor-scripts'

const SCRIPT = '/server-scripts/SystemPing.js:default'
const LABEL = 'System ping (server)'

function systemAPI(result: unknown = { set: [{ name: SCRIPT, label: LABEL }] }) {
  return { automationList: vi.fn().mockResolvedValue(result) }
}

let wrappers: ReturnType<typeof mount>[] = []

async function mountViewer(modelValue: string | null, $SystemAPI: unknown) {
  const wrapper = mount(CViewCorredorScript, {
    props: { modelValue },
    global: { provide: { $SystemAPI } },
  })
  wrappers.push(wrapper)
  await flushPromises()
  return wrapper
}

beforeEach(() => {
  resetCorredorScriptCache()
})

afterEach(() => {
  wrappers.forEach(w => w.unmount())
  wrappers = []
})

describe('CViewCorredorScript', () => {
  it('shows the script label with the stored name beneath it', async () => {
    const api = systemAPI()
    const wrapper = await mountViewer(SCRIPT, api)

    expect(wrapper.text()).toContain(LABEL)
    expect(wrapper.find('small').text()).toBe(SCRIPT)
    expect(wrapper.attributes('title')).toBe(SCRIPT)
  })

  it('asks the server only for manual system-level server scripts', async () => {
    const api = systemAPI()
    await mountViewer(SCRIPT, api)

    expect(api.automationList).toHaveBeenCalledWith({
      eventTypes: ['onManual'],
      resourceTypes: ['system'],
      excludeInvalid: true,
      excludeClientScripts: true,
    })
  })

  it('shows the raw name when the list carries no such script', async () => {
    const wrapper = await mountViewer('/server-scripts/Gone.js:default', systemAPI())

    expect(wrapper.text()).toContain('/server-scripts/Gone.js:default')
    // Nothing to show twice: the name is the whole answer.
    expect(wrapper.find('small').exists()).toBe(false)
  })

  it('shows the raw name when Corredor is unreachable', async () => {
    const api = { automationList: vi.fn().mockRejectedValue(new Error('no corredor')) }
    const wrapper = await mountViewer(SCRIPT, api)

    expect(wrapper.text()).toContain(SCRIPT)
    expect(wrapper.find('small').exists()).toBe(false)
  })

  it('asks again after a failed fetch rather than caching the failure', async () => {
    const api = {
      automationList: vi
        .fn()
        .mockRejectedValueOnce(new Error('no corredor'))
        .mockResolvedValue({ set: [{ name: SCRIPT, label: LABEL }] }),
    }

    const first = await mountViewer(SCRIPT, api)
    expect(first.text()).toContain(SCRIPT)

    const second = await mountViewer(SCRIPT, api)
    expect(second.text()).toContain(LABEL)
    expect(api.automationList).toHaveBeenCalledTimes(2)
  })

  it('fetches the list once for a preview carrying several steps', async () => {
    const api = systemAPI()
    await mountViewer(SCRIPT, api)
    await mountViewer(SCRIPT, api)
    await mountViewer(SCRIPT, api)

    expect(api.automationList).toHaveBeenCalledTimes(1)
  })

  it('reads as not set when no script is configured', async () => {
    const api = systemAPI()
    const wrapper = await mountViewer(null, api)

    expect(wrapper.text()).toContain('builder.preview.notSet')
    expect(api.automationList).not.toHaveBeenCalled()
  })

  it('resolves a script picked after the first render', async () => {
    const api = systemAPI()
    const wrapper = await mountViewer(null, api)

    await wrapper.setProps({ modelValue: SCRIPT })
    await flushPromises()

    expect(wrapper.text()).toContain(LABEL)
  })
})

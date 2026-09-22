import { describe, it, expect, beforeAll, afterEach, vi } from 'vitest'
import { mount, flushPromises, type VueWrapper } from '@vue/test-utils'
import CInputCorredorScript from './CInputCorredorScript.vue'

beforeAll(() => {
  // @ts-expect-error assigning a stub over a property jsdom does not define
  window.matchMedia = (q: string) => ({
    matches: false,
    media: q,
    onchange: null,
    addEventListener() {},
    removeEventListener() {},
    addListener() {},
    removeListener() {},
    dispatchEvent: () => false,
  })
})

const PING = '/server-scripts/SystemPing.js:default'
const UNLABELLED = '/server-scripts/agent-sandbox/Unlabelled.js:default'
const ARCHIVE = '/server-scripts/Archive.js:default'

const SCRIPTS = [
  {
    name: PING,
    label: 'System ping (server)',
    description: 'Answers a manual system-level call',
    triggers: [{ eventTypes: ['onManual'], resourceTypes: ['system'] }],
  },
  { name: UNLABELLED, triggers: [{ eventTypes: ['onManual'], resourceTypes: ['system'] }] },
  { name: ARCHIVE, label: 'Archive old records', description: '' },
]

let wrapper: VueWrapper | null = null

function apiReturning(result: () => Promise<unknown>) {
  return {
    automationListCancellable: vi.fn(() => ({ response: result, cancel: vi.fn() })),
  }
}

async function mountInput(
  props: Record<string, unknown> = {},
  api = apiReturning(async () => ({ set: SCRIPTS })),
) {
  wrapper = mount(CInputCorredorScript, { props, global: { provide: { $SystemAPI: api } } })
  await flushPromises()
  return { wrapper, api }
}

function select() {
  return wrapper!.findComponent({ name: 'Select' })
}

function offered() {
  return select().props('options') as Array<{ name: string; label: string; description?: string }>
}

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('CInputCorredorScript', () => {
  it('lists the manual server scripts on the system resource the caller may run', async () => {
    const { api } = await mountInput()

    expect(api.automationListCancellable).toHaveBeenCalledWith({
      eventTypes: ['onManual'],
      resourceTypes: ['system'],
      excludeInvalid: true,
      excludeClientScripts: true,
    })
  })

  it('offers each script by name, labelled by its label or else its name, sorted by label', async () => {
    await mountInput()

    expect(offered().map(o => [o.name, o.label])).toEqual([
      [UNLABELLED, UNLABELLED],
      [ARCHIVE, 'Archive old records'],
      [PING, 'System ping (server)'],
    ])
    expect(select().props('optionValue')).toBe('name')
    expect(select().props('optionLabel')).toBe('label')
  })

  it('keeps the description for the line under the label', async () => {
    await mountInput()

    expect(offered().find(o => o.name === PING)?.description).toBe(
      'Answers a manual system-level call',
    )
  })

  it('filters over label and name, and can be cleared', async () => {
    await mountInput()

    expect(select().props('filter')).toBe(true)
    expect(select().props('filterFields')).toEqual(['label', 'name'])
    expect(select().props('showClear')).toBe(true)
  })

  it('shows the label of a listed script it holds', async () => {
    await mountInput({ modelValue: PING })

    expect((select().find('input').element as HTMLInputElement).value).toBe('System ping (server)')
  })

  it('keeps a typed name that matches no listed script as typed', async () => {
    await mountInput()

    await select().find('input').setValue('/server-scripts/DeployedLater.js:default')

    expect(wrapper!.emitted('update:modelValue')!.at(-1)).toEqual([
      '/server-scripts/DeployedLater.js:default',
    ])
  })

  it('shows an unlisted name it holds as-is', async () => {
    await mountInput({ modelValue: '/server-scripts/Hidden.js:default' })

    expect(select().props('modelValue')).toBe('/server-scripts/Hidden.js:default')
    expect((select().find('input').element as HTMLInputElement).value).toBe(
      '/server-scripts/Hidden.js:default',
    )
  })

  it('emits the picked script name, and null once cleared', async () => {
    await mountInput()

    select().vm.$emit('update:modelValue', PING)
    select().vm.$emit('update:modelValue', null)
    select().vm.$emit('update:modelValue', '')

    expect(wrapper!.emitted('update:modelValue')).toEqual([[PING], [null], [null]])
  })

  it('passes placeholder and disabled through', async () => {
    await mountInput({ placeholder: 'Pick a script', disabled: true })

    expect(select().props('placeholder')).toBe('Pick a script')
    expect(select().props('disabled')).toBe(true)
  })

  it('shows loading while the list is on its way', async () => {
    let resolve: (_result: unknown) => void = () => {}
    const pending = new Promise(r => (resolve = r))
    await mountInput(
      {},
      apiReturning(() => pending),
    )

    expect(select().props('loading')).toBe(true)

    resolve({ set: SCRIPTS })
    await flushPromises()

    expect(select().props('loading')).toBe(false)
    expect(offered()).toHaveLength(3)
  })

  it('offers nothing when the list fails, and still takes a typed name', async () => {
    await mountInput(
      {},
      apiReturning(() => Promise.reject(new Error('corredor down'))),
    )

    expect(offered()).toEqual([])
    expect(select().props('loading')).toBe(false)

    await select().find('input').setValue('/server-scripts/SystemPing.js:default')
    expect(wrapper!.emitted('update:modelValue')!.at(-1)).toEqual([PING])
  })
})

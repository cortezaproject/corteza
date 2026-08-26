import { describe, it, expect, beforeAll, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import AgentToolDialog from './AgentToolDialog.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (k: string, fallback?: unknown) => (typeof fallback === 'string' ? fallback : k),
  }),
}))

beforeAll(() => {
  // PrimeVue overlays bind a matchMedia listener on mount; jsdom has none.
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

const tools = [
  { name: 'compose_record_lookup', title: 'Lookup record', groups: ['usage'], risk: 'read' },
  {
    name: 'compose_record_delete',
    title: 'Delete record',
    description: 'Removes a row. The delete is soft.',
    groups: ['usage'],
    risk: 'destructive',
  },
  { name: 'system_user_create', title: 'Create user', groups: ['configuring'], risk: 'write' },
]

function toolNames(vm: any): string[] {
  return vm.groups.flatMap((g: any) => g.areas.flatMap((a: any) => a.tools.map((t: any) => t.name)))
}

function mountDialog(grants: any[] = []) {
  return mount(AgentToolDialog, {
    props: { visible: true, tools, grants, disabled: false },
    global: { mocks: { $t: (k: string, p?: any) => (p ? `${k}:${JSON.stringify(p)}` : k) } },
  })
}

describe('AgentToolDialog', () => {
  // Two levels: a grant can name a group, and a person scans by area. One level
  // of either alone is a list of ninety.
  it('nests areas inside the group a grant can name', () => {
    const vm = mountDialog().vm as any
    expect(vm.groups.map((g: any) => g.key)).toEqual(['usage', 'configuring'])

    const usage = vm.groups.find((g: any) => g.key === 'usage')
    expect(usage.areas.map((a: any) => a.key)).toContain('compose_record')

    const configuring = vm.groups.find((g: any) => g.key === 'configuring')
    expect(configuring.areas.map((a: any) => a.key)).toContain('system_user')
  })

  // A group grant keeps meaning "every data tool" as tools are added; naming
  // them one at a time does not.
  it('stores a group grant with its risk ceiling', () => {
    const w = mountDialog()
    const vm = w.vm as any

    vm.toggleGroup('usage', true)
    vm.setGroupRisk('usage', 'write')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([{ group: 'usage', maxRisk: 'write', description: '', allow: [] }])
  })

  it('opens showing a group grant it was given', () => {
    const vm = mountDialog([{ group: 'configuring', maxRisk: 'destructive', allow: [] }]).vm as any
    expect(vm.groupGrant('configuring')).toEqual({ group: 'configuring', maxRisk: 'destructive' })
    expect(vm.groupGrant('usage')).toBeNull()
  })

  // The description says what a tool does; the handle says nothing a person
  // needs. Only the first sentence — the rest is written for the model.
  it('summarises a description to its first sentence', () => {
    const vm = mountDialog().vm as any
    expect(vm.summarise('Delete a page. The delete is soft: nothing is erased.')).toBe(
      'Delete a page.',
    )
    expect(vm.summarise('No full stop here')).toBe('No full stop here')
    expect(vm.summarise('')).toBe('')
  })

  it('opens showing what the agent already has', () => {
    const vm = mountDialog([{ name: 'compose_record_lookup', permission: 'ask' }]).vm as any
    expect([...vm.chosen]).toEqual(['compose_record_lookup'])
    expect(vm.modeOf({ name: 'compose_record_lookup' })).toBe('ask')
  })

  // A grant carries a namespace scope set elsewhere in the editor; the dialog
  // chooses tools and modes and must not quietly drop the rest of the entry.
  it('keeps a grant scope through an edit', () => {
    const scoped = {
      name: 'compose_record_lookup',
      description: 'Reads',
      allow: [{ namespaceID: '100', moduleIDs: [] }],
    }
    const w = mountDialog([scoped])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_lookup' }, 'always')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([{ ...scoped, permission: 'always' }])
  })

  // A group grant has no name; it belongs to the panel, not this dialog, and
  // applying must leave it alone rather than delete it.
  it('leaves group grants untouched', () => {
    const group = { group: 'usage', maxRisk: 'read', allow: [] }
    const w = mountDialog([group, { name: 'compose_record_lookup' }])
    const vm = w.vm as any

    vm.toggle({ name: 'compose_record_lookup' }, false)
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([group])
  })

  it('narrows by search, over descriptions too', async () => {
    const w = mountDialog()
    const vm = w.vm as any

    vm.search = 'delete'
    await w.vm.$nextTick()
    expect(toolNames(vm)).toEqual(['compose_record_delete'])

    vm.search = 'removes a row'
    await w.vm.$nextTick()
    expect(toolNames(vm)).toEqual(['compose_record_delete'])
  })
})

describe('AgentToolDialog overrides', () => {
  // "All data tools, but ask before deleting" is a group grant plus a named
  // entry. Setting a mode on a tool the group covers has to write that entry.
  it('writes a named override beside a group grant', () => {
    const w = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_delete' }, 'deny')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([
      { group: 'usage', maxRisk: 'write', allow: [] },
      { name: 'compose_record_delete', description: '', allow: [], permission: 'deny' },
    ])
  })

  // Clearing an override the session added removes the entry rather than
  // leaving a grant that says nothing.
  it('drops an override cleared back to default', () => {
    const w = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_delete' }, 'deny')
    vm.setMode({ name: 'compose_record_delete' }, '')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([{ group: 'usage', maxRisk: 'write', allow: [] }])
  })
})

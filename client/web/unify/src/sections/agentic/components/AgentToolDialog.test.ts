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
  { name: 'compose_record_delete', title: 'Delete record', groups: ['usage'], risk: 'destructive' },
  { name: 'system_user_create', title: 'Create user', groups: ['configuring'], risk: 'write' },
]

function mountDialog(grants: any[] = []) {
  return mount(AgentToolDialog, {
    props: { visible: true, tools, grants, disabled: false },
    global: { mocks: { $t: (k: string, p?: any) => (p ? `${k}:${JSON.stringify(p)}` : k) } },
  })
}

describe('AgentToolDialog', () => {
  // Tools are grouped by what they act on. Grouping by the two coarse groups
  // instead puts ninety tools under one heading, which is the dropdown again.
  it('sections tools by the area their name carries', () => {
    const vm = mountDialog().vm as any
    const keys = vm.sections.map((s: any) => s.key)
    expect(keys).toContain('compose_record')
    expect(keys).toContain('system_user')
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

  it('narrows by search', async () => {
    const w = mountDialog()
    const vm = w.vm as any
    vm.search = 'delete'
    await w.vm.$nextTick()
    const names = vm.sections.flatMap((s: any) => s.tools.map((t: any) => t.name))
    expect(names).toEqual(['compose_record_delete'])
  })
})

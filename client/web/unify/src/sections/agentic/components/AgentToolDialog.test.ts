import { describe, it, expect, beforeAll, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
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
  return vm.domains.flatMap((d: any) =>
    d.areas.flatMap((a: any) => a.tools.map((t: any) => t.name)),
  )
}

function domain(vm: any, key: string) {
  return vm.domains.find((d: any) => d.key === key)
}

function mountDialog(grants: any[] = [], list = tools) {
  return mount(AgentToolDialog, {
    props: { visible: true, tools: list, grants, disabled: false },
    global: { mocks: { $t: (k: string, p?: any) => (p ? `${k}:${JSON.stringify(p)}` : k) } },
  })
}

describe('AgentToolDialog grouping', () => {
  // Grouping by `usage` and `configuring` split five subjects across both
  // sections — every TAQ tool but `exec` in one, `exec` in the other. Asking
  // "what can this agent do to TAQs?" then meant looking in two places.
  it('keeps a subject whole, whichever group each tool belongs to', () => {
    const split = [
      { name: 'automation_taq_exec', title: 'Execute TAQ', groups: ['usage'], risk: 'write' },
      {
        name: 'automation_taq_create',
        title: 'Create TAQ',
        groups: ['configuring'],
        risk: 'write',
      },
    ]
    const vm = mountDialog([], split).vm as any

    expect(vm.domains.map((d: any) => d.key)).toEqual(['automation'])
    expect(toolNames(vm)).toEqual(['automation_taq_create', 'automation_taq_exec'])
  })

  it('counts what a section holds', () => {
    const vm = mountDialog().vm as any
    expect(domain(vm, 'data').count).toBe(2)
    expect(domain(vm, 'people').count).toBe(1)
  })

  // A tool family nobody has placed yet gets a section of its own rather than
  // disappearing from a dialog that claims to list everything.
  it('gives an unplaced family a section rather than dropping it', () => {
    const vm = mountDialog(
      [],
      [...tools, { name: 'discovery_search', title: 'Search', groups: ['usage'], risk: 'read' }],
    ).vm as any

    expect(domain(vm, 'other').areas.map((a: any) => a.key)).toEqual(['discovery'])
    expect(vm.domains.at(-1).key).toBe('other')
  })

  // Alphabetical put Charts first, which is nobody's starting point.
  it('leads with what an agent works with, not with C', () => {
    const many = [
      { name: 'compose_chart_lookup', title: 'Lookup chart', groups: ['usage'], risk: 'read' },
      { name: 'compose_record_lookup', title: 'Lookup record', groups: ['usage'], risk: 'read' },
      { name: 'system_user_create', title: 'Create user', groups: ['configuring'], risk: 'write' },
    ]
    const vm = mountDialog([], many).vm as any
    expect(vm.domains.map((d: any) => d.key)).toEqual(['data', 'structure', 'people'])
  })

  // Reads first, then writes, then the one that cannot be taken back.
  it('orders an area by risk, then by name', () => {
    const records = [
      {
        name: 'compose_record_delete',
        title: 'Delete record',
        groups: ['usage'],
        risk: 'destructive',
      },
      { name: 'compose_record_create', title: 'Create record', groups: ['usage'], risk: 'write' },
      { name: 'compose_record_lookup', title: 'Lookup record', groups: ['usage'], risk: 'read' },
    ]
    const vm = mountDialog([], records).vm as any
    expect(toolNames(vm)).toEqual([
      'compose_record_lookup',
      'compose_record_create',
      'compose_record_delete',
    ])
  })

  // A skill is injected when the tool that triggers it is used; offering it
  // here invites a choice that changes nothing.
  it('does not offer skills', () => {
    const vm = mountDialog(
      [],
      [
        ...tools,
        {
          name: 'system_skill_lookup',
          title: 'Lookup skill',
          groups: ['configuring'],
          risk: 'read',
        },
      ],
    ).vm as any
    expect(toolNames(vm)).not.toContain('system_skill_lookup')
  })
})

describe('AgentToolDialog sections', () => {
  // Six shut sections give a fresh agent nothing to act on; six open ones are
  // a list of ninety.
  it('opens the first section when the agent has no named tool', () => {
    const vm = mountDialog().vm as any
    expect(vm.isOpen('data')).toBe(true)
    expect(vm.isOpen('people')).toBe(false)
  })

  it('opens on the sections holding what the agent already has', () => {
    const vm = mountDialog([{ name: 'system_user_create' }]).vm as any
    expect(vm.isOpen('people')).toBe(true)
    expect(vm.isOpen('data')).toBe(false)
  })

  it('folds a section shut and back open', () => {
    const vm = mountDialog().vm as any
    expect(vm.isOpen('data')).toBe(true)
    vm.toggleCollapsed('data')
    expect(vm.isOpen('data')).toBe(false)
    vm.toggleCollapsed('data')
    expect(vm.isOpen('data')).toBe(true)
  })

  // A search that only looked inside open sections would report nothing while
  // showing a shut section holding the match.
  it('shows a match inside a section that is shut', async () => {
    const w = mountDialog()
    const vm = w.vm as any

    expect(vm.isOpen('people')).toBe(false)
    vm.search = 'Create user'
    await w.vm.$nextTick()

    expect(toolNames(vm)).toEqual(['system_user_create'])
    expect(vm.isOpen('people')).toBe(true)
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

describe('AgentToolDialog family grants', () => {
  // A family grant keeps meaning "every data tool" as tools are added; naming
  // them one at a time does not.
  it('stores a family grant with its risk ceiling', () => {
    const w = mountDialog()
    const vm = w.vm as any

    vm.toggleGroup('usage', true)
    vm.setGroupRisk('usage', 'write')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([{ group: 'usage', maxRisk: 'write', description: '', allow: [] }])
  })

  it('opens showing a family grant it was given', () => {
    const vm = mountDialog([{ group: 'configuring', maxRisk: 'destructive', allow: [] }]).vm as any
    expect(vm.groupGrant('configuring')).toEqual({ group: 'configuring', maxRisk: 'destructive' })
    expect(vm.groupGrant('usage')).toBeNull()
  })

  // A ceiling admits its own level and everything below it. Reading coverage
  // off the section instead said a `write` family covered the delete tool
  // sitting in it, and left every covered row's checkbox unticked.
  it('covers a tool at or below the ceiling, and no higher', () => {
    const vm = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }]).vm as any

    expect(vm.coveredByGroup({ groups: ['usage'], risk: 'read' })).toBe(true)
    expect(vm.coveredByGroup({ groups: ['usage'], risk: 'write' })).toBe(true)
    expect(vm.coveredByGroup({ groups: ['usage'], risk: 'destructive' })).toBe(false)
    expect(vm.coveredByGroup({ groups: ['configuring'], risk: 'read' })).toBe(false)
  })

  // A family grant has no name; applying must leave it alone rather than
  // delete it.
  it('leaves family grants untouched', () => {
    const group = { group: 'usage', maxRisk: 'read', allow: [] }
    const w = mountDialog([group, { name: 'compose_record_lookup' }])
    const vm = w.vm as any

    vm.toggle({ name: 'compose_record_lookup' }, false)
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([group])
  })
})

describe('AgentToolDialog modes', () => {
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
    expect(vm.effectiveMode({ name: 'compose_record_lookup', risk: 'read' })).toBe('ask')
  })

  // There is no fourth "default" segment: the default IS one of the three, and
  // a row with nothing set shows the one its risk decides.
  it('shows the risk default as the chosen segment', () => {
    const vm = mountDialog().vm as any
    expect(vm.effectiveMode({ name: 'compose_record_lookup', risk: 'read' })).toBe('always')
    expect(vm.effectiveMode({ name: 'compose_record_create', risk: 'write' })).toBe('ask')
    expect(vm.effectiveMode({ name: 'compose_record_delete', risk: 'destructive' })).toBe('ask')
  })

  // The rail marks the one state the risk rule would not have produced: an
  // allow set by hand on a tool that writes. Tightening never needs a warning.
  it('rails only a hand-set allow on a tool that writes', () => {
    const w = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }])
    const vm = w.vm as any
    const create = { name: 'compose_record_create', risk: 'write' }

    expect(vm.loosened(create)).toBe(false)
    vm.setMode(create, 'always')
    expect(vm.loosened(create)).toBe(true)

    vm.setMode({ name: 'compose_record_lookup', risk: 'read' }, 'always')
    expect(vm.loosened({ name: 'compose_record_lookup', risk: 'read' })).toBe(false)
  })

  // Choosing the mode the risk would have given anyway is not a decision to
  // store: pinning it would freeze the tool if the rule ever changed.
  it('stores an override only when it differs from the risk default', () => {
    const w = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_create', risk: 'write' }, 'ask')
    expect(vm.chosen.has('compose_record_create')).toBe(false)

    vm.setMode({ name: 'compose_record_create', risk: 'write' }, 'always')
    expect(vm.chosen.has('compose_record_create')).toBe(true)
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

    vm.setMode({ name: 'compose_record_lookup', risk: 'read' }, 'ask')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([{ ...scoped, permission: 'ask' }])
  })

  // "All data tools, but ask before deleting" is a family grant plus a named
  // entry. Setting a mode on a tool the family covers has to write that entry.
  it('writes a named override beside a family grant', () => {
    const w = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_delete', risk: 'destructive' }, 'deny')
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

    vm.setMode({ name: 'compose_record_delete', risk: 'destructive' }, 'deny')
    vm.setMode({ name: 'compose_record_delete', risk: 'destructive' }, 'ask')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([{ group: 'usage', maxRisk: 'write', allow: [] }])
  })
})

describe('AgentToolDialog footer', () => {
  // Every other dialog in the app footers at `size="small"`; this one shipped
  // at the default and read a size larger than all of them.
  it('sizes its footer buttons like every other dialog', async () => {
    const PrimeVue = (await import('primevue/config')).default
    const Dialog = (await import('primevue/dialog')).default
    const Button = (await import('primevue/button')).default

    const w = mount(AgentToolDialog, {
      props: { visible: true, tools, grants: [], disabled: false },
      global: {
        plugins: [PrimeVue],
        components: { Dialog, Button },
        mocks: { $t: (k: string) => k },
        stubs: { Checkbox: true, Select: true, SelectButton: true },
      },
    })

    await flushPromises()

    // The dialog teleports to body, so the footer is not inside the wrapper.
    const buttons = [...document.querySelectorAll('.p-dialog-footer button')]
    expect(buttons.length).toBe(2)
    for (const b of buttons) expect([...b.classList]).toContain('p-button-sm')

    w.unmount()
  })
})

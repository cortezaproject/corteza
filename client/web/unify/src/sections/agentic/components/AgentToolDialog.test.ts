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

  // Reminders are a personal surface — create, snooze, dismiss — and share
  // nothing with records but the group they are registered in.
  it('keeps a subject out of a section it only shares a group with', () => {
    const vm = mountDialog(
      [],
      [
        { name: 'compose_record_lookup', title: 'Lookup record', groups: ['usage'], risk: 'read' },
        {
          name: 'system_reminder_snooze',
          title: 'Snooze reminder',
          groups: ['usage'],
          risk: 'write',
        },
      ],
    ).vm as any

    expect(domain(vm, 'records').areas.map((a: any) => a.key)).toEqual(['compose_record'])
    expect(domain(vm, 'reminders').areas.map((a: any) => a.key)).toEqual(['system_reminder'])
  })

  // The schema a namespace defines is a different job from the pages that
  // display it, and one 26-tool section buried the twelve page tools.
  it('separates what defines the data from what displays it', () => {
    const vm = mountDialog(
      [],
      [
        {
          name: 'compose_module_create',
          title: 'Create module',
          groups: ['configuring'],
          risk: 'write',
        },
        {
          name: 'compose_page_create',
          title: 'Create page',
          groups: ['configuring'],
          risk: 'write',
        },
      ],
    ).vm as any

    expect(domain(vm, 'datamodel').areas.map((a: any) => a.key)).toEqual(['compose_module'])
    expect(domain(vm, 'interface').areas.map((a: any) => a.key)).toEqual(['compose_page'])
  })

  // A subject listed under two sections renders its tools twice, and a mode set
  // on one copy reads as unset on the other.
  it('places every subject in exactly one section', () => {
    const areas = [
      'compose_record',
      'compose_namespace',
      'compose_module',
      'compose_page',
      'compose_chart',
      'automation_taq',
      'automation_workflow',
      'automation_trigger',
      'automation_event',
      'system_user',
      'system_role',
      'system_auth',
      'system_agent',
      'system_chatbot',
      'system_reminder',
      'system_application',
      'system_theme',
    ]
    const vm = mountDialog(
      [],
      areas.map(a => ({ name: `${a}_lookup`, title: a, groups: ['usage'], risk: 'read' })),
    ).vm as any

    const names = toolNames(vm)
    expect(names).toHaveLength(areas.length)
    expect(new Set(names).size).toBe(areas.length)
  })

  // Alphabetical put Charts first, which is nobody's starting point.
  it('leads with what an agent works with, not with C', () => {
    const many = [
      { name: 'compose_chart_lookup', title: 'Lookup chart', groups: ['usage'], risk: 'read' },
      { name: 'compose_record_lookup', title: 'Lookup record', groups: ['usage'], risk: 'read' },
      { name: 'system_user_create', title: 'Create user', groups: ['configuring'], risk: 'write' },
    ]
    const vm = mountDialog([], many).vm as any
    expect(vm.domains.map((d: any) => d.key)).toEqual(['records', 'interface', 'people'])
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
    expect(vm.isOpen('records')).toBe(true)
    expect(vm.isOpen('people')).toBe(false)
  })

  it('opens on the sections holding what the agent already has', () => {
    const vm = mountDialog([{ name: 'system_user_create' }]).vm as any
    expect(vm.isOpen('people')).toBe(true)
    expect(vm.isOpen('records')).toBe(false)
  })

  it('folds a section shut and back open', () => {
    const vm = mountDialog().vm as any
    expect(vm.isOpen('records')).toBe(true)
    vm.toggleCollapsed('records')
    expect(vm.isOpen('records')).toBe(false)
    vm.toggleCollapsed('records')
    expect(vm.isOpen('records')).toBe(true)
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
  // The dialog chooses tools; a grant naming a whole family is set elsewhere in
  // the editor. It still has to know what a family already covers, and hand
  // every family back exactly as it found it.
  it('shows a covered tool as one the agent has', () => {
    const vm = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }]).vm as any
    const lookup = { name: 'compose_record_lookup', groups: ['usage'], risk: 'read' }

    expect(vm.chosen.has(lookup.name)).toBe(false)
    expect(vm.toolOn(lookup)).toBe(true)
  })

  it('leaves family grants untouched', () => {
    const group = { group: 'usage', maxRisk: 'read', allow: [] }
    const w = mountDialog([group, { name: 'compose_record_lookup' }])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_lookup' }, 'deny')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([group])
  })
})

describe('AgentToolDialog section control', () => {
  function recordsSection(vm: any) {
    return domain(vm, 'records')
  }

  // A subject is the unit an agent is actually given: every record tool, or
  // none of them. Ticking six boxes to say that is six chances to miss one.
  it('turns a whole section on at one mode', () => {
    const w = mountDialog()
    const vm = w.vm as any

    vm.setSectionMode(recordsSection(vm), 'always')
    expect(vm.sectionMode(recordsSection(vm))).toBe('always')
    expect([...vm.chosen].sort()).toEqual(['compose_record_delete', 'compose_record_lookup'])
  })

  it('turns a whole section off', () => {
    const w = mountDialog()
    const vm = w.vm as any

    vm.setSectionMode(recordsSection(vm), 'ask')
    expect(vm.sectionMode(recordsSection(vm))).toBe('ask')

    vm.setSectionMode(recordsSection(vm), 'deny')
    expect(vm.sectionMode(recordsSection(vm))).toBe('deny')
    expect([...vm.chosen]).toEqual([])
  })

  // Custom is what the control reads as once the tools disagree — the state it
  // shows back, never a state it is asked to produce.
  it('reads as custom once the tools disagree, and setting it changes nothing', () => {
    const w = mountDialog()
    const vm = w.vm as any

    vm.setSectionMode(recordsSection(vm), 'always')
    vm.setMode({ name: 'compose_record_delete', risk: 'destructive' }, 'ask')
    expect(vm.sectionMode(recordsSection(vm))).toBe('custom')

    const before = [...vm.chosen].sort()
    vm.setSectionMode(recordsSection(vm), 'custom')
    expect([...vm.chosen].sort()).toEqual(before)
  })

  // A total says how big the section is; the split says what it lets the agent
  // do, which is the question the header is there to answer.
  it('shows a section as a split by mode rather than a total', () => {
    const w = mountDialog([{ name: 'compose_record_lookup' }])
    const vm = w.vm as any

    expect(vm.sectionCounts(domain(vm, 'records'))).toEqual([
      { mode: 'always', n: 1 },
      { mode: 'deny', n: 1 },
    ])
  })

  it('reads a section nothing has been chosen in as blocked', () => {
    const vm = mountDialog().vm as any
    expect(vm.sectionMode(domain(vm, 'people'))).toBe('deny')
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
    expect(vm.rowMode({ name: 'compose_record_lookup', risk: 'read' })).toBe('ask')
  })

  // There is no fourth "default" segment: the default IS one of the three, and
  // a row with nothing set shows the one its risk decides.
  it('shows the risk default as the chosen segment', () => {
    const vm = mountDialog([
      { name: 'compose_record_lookup' },
      { name: 'compose_record_create' },
      { name: 'compose_record_delete' },
    ]).vm as any

    expect(vm.rowMode({ name: 'compose_record_lookup', risk: 'read' })).toBe('always')
    expect(vm.rowMode({ name: 'compose_record_create', risk: 'write' })).toBe('ask')
    expect(vm.rowMode({ name: 'compose_record_delete', risk: 'destructive' })).toBe('ask')
  })

  // Choosing the mode the risk would have given anyway is not a decision to
  // store: pinning it would freeze the tool if the rule ever changed.
  it('stores an override only when it differs from the risk default', () => {
    const w = mountDialog([{ group: 'usage', maxRisk: 'write', allow: [] }])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_create', groups: ['usage'], risk: 'write' }, 'ask')
    expect(vm.chosen.has('compose_record_create')).toBe(false)

    vm.setMode({ name: 'compose_record_create', groups: ['usage'], risk: 'write' }, 'always')
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
    const family = { group: 'usage', maxRisk: 'destructive', allow: [] }
    const w = mountDialog([family])
    const vm = w.vm as any

    vm.setMode({ name: 'compose_record_delete', groups: ['usage'], risk: 'destructive' }, 'deny')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([
      family,
      { name: 'compose_record_delete', description: '', allow: [], permission: 'deny' },
    ])
  })

  // Blocking a tool nothing else grants is saying nothing about it: the entry
  // goes, rather than staying behind as a grant that grants nothing.
  it('drops the entry when blocking a tool no family covers', () => {
    const w = mountDialog()
    const vm = w.vm as any
    const create = { name: 'compose_record_create', risk: 'write' }

    vm.setMode(create, 'ask')
    expect(vm.toolOn(create)).toBe(true)

    vm.setMode(create, 'deny')
    expect(vm.toolOn(create)).toBe(false)

    vm.apply()
    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([])
  })

  // Clearing an override the session added removes the entry rather than
  // leaving a grant that says nothing.
  it('drops an override cleared back to default', () => {
    const family = { group: 'usage', maxRisk: 'destructive', allow: [] }
    const w = mountDialog([family])
    const vm = w.vm as any
    const del = { name: 'compose_record_delete', groups: ['usage'], risk: 'destructive' }

    vm.setMode(del, 'deny')
    vm.setMode(del, 'ask')
    vm.apply()

    const [applied] = w.emitted('apply') as any[]
    expect(applied[0]).toEqual([family])
  })
})

// The dialog teleports to body, so what it renders is not inside the wrapper.
async function mountRendered(grants: any[] = []) {
  const PrimeVue = (await import('primevue/config')).default
  const Dialog = (await import('primevue/dialog')).default
  const Button = (await import('primevue/button')).default

  const w = mount(AgentToolDialog, {
    props: { visible: true, tools, grants, disabled: false },
    global: {
      plugins: [PrimeVue],
      components: { Dialog, Button },
      mocks: { $t: (k: string) => k },
      stubs: {
        Checkbox: true,
        Select: true,
        Message: true,
        SelectButton: { template: '<div data-testid="mode-toggle" />' },
      },
    },
  })

  await flushPromises()
  return w
}

describe('AgentToolDialog footer', () => {
  // Every other dialog in the app footers at `size="small"`; this one shipped
  // at the default and read a size larger than all of them.
  it('sizes its footer buttons like every other dialog', async () => {
    const w = await mountRendered()

    const buttons = [...document.querySelectorAll('.p-dialog-footer button')]
    expect(buttons.length).toBe(2)
    for (const b of buttons) expect([...b.classList]).toContain('p-button-sm')

    w.unmount()
  })
})

describe('AgentToolDialog section headers', () => {
  // Thirty rows go past in People and access; losing the subject and its
  // permission control off the top means scrolling back to use either.
  it('pins a section header while its rows scroll', async () => {
    const w = await mountRendered()
    const heads = [...document.querySelectorAll('.tool-section-head')]

    expect(heads.length).toBeGreaterThan(0)
    for (const h of heads) {
      expect([...h.classList]).toContain('sticky')
      expect([...h.classList]).toContain('top-0')
    }

    w.unmount()
  })

  // One chevron that turns, rather than two icons swapped for each other.
  it('turns the chevron with the section rather than swapping icons', async () => {
    const w = await mountRendered()
    const chevrons = [...document.querySelectorAll('.tool-section-chevron')]

    expect(chevrons.length).toBeGreaterThan(1)
    expect(chevrons.every(c => c.classList.contains('pi-chevron-right'))).toBe(true)
    expect(chevrons.map(c => c.getAttribute('data-open'))).toContain('true')
    expect(chevrons.map(c => c.getAttribute('data-open'))).toContain('false')

    w.unmount()
  })
})

describe('AgentToolDialog blocked rows', () => {
  // Blocked is the off state, so a row the agent does not have reads as blocked
  // rather than showing a default it is not following.
  it('reads a tool the agent does not have as blocked', () => {
    const vm = mountDialog().vm as any
    const lookup = { name: 'compose_record_lookup', risk: 'read' }

    expect(vm.rowMode(lookup)).toBe('deny')
    vm.setMode(lookup, 'always')
    expect(vm.rowMode(lookup)).toBe('always')
  })

  // A blocked row's control recedes rather than vanishes: at full strength on
  // every one of ninety tools the agent does not have it reads as ninety
  // decisions waiting to be made, and gone it leaves no way back.
  it('leaves the control in place on a blocked row, receded', async () => {
    const w = await mountRendered()
    const rows = [...document.querySelectorAll('[data-blocked]')]

    expect(rows.length).toBe(2)
    expect(rows.every(r => r.getAttribute('data-blocked') === 'true')).toBe(true)
    expect(rows.every(r => r.classList.contains('tool-mode'))).toBe(true)
    expect(rows.every(r => r.querySelector('[data-testid="mode-toggle"]'))).toBe(true)

    w.unmount()
  })

  it('leaves out the sub-label when a section holds one subject', async () => {
    const w = await mountRendered()
    const labels = [...document.querySelectorAll('.p-dialog .uppercase')].map(e =>
      e.textContent.trim(),
    )

    // The open section is Records, whose only area repeats its own name.
    expect(labels).toEqual([])
    w.unmount()
  })

  it('shows the control on a tool the agent has', async () => {
    const w = await mountRendered([{ name: 'compose_record_lookup' }])
    const shown = [...document.querySelectorAll('[data-blocked="false"]')]

    expect(shown.length).toBe(1)
    expect(shown[0].getAttribute('data-blocked')).toBe('false')

    w.unmount()
  })
})

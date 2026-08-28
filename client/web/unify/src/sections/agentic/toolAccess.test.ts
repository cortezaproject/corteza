import { describe, it, expect } from 'vitest'
import {
  canScope,
  confineTo,
  coveredByFamily,
  defaultModeFor,
  hasSettings,
  modeOf,
  scopeBlocked,
  scopesModules,
  summaryOf,
  sectionsOf,
  splitGrants,
  tally,
} from './toolAccess'

const read = { name: 'compose_record_lookup', groups: ['usage'], risk: 'read' }
const write = { name: 'compose_record_create', groups: ['usage'], risk: 'write' }
const destroy = { name: 'compose_record_delete', groups: ['usage'], risk: 'destructive' }

function grants(...gg: any[]) {
  return splitGrants(gg)
}

describe('toolAccess', () => {
  // The runtime resolves a grant that names no mode when it runs the agent, so
  // the editor has to resolve it the same way to show the agent back.
  it('resolves a grant that names no mode the way the runtime does', () => {
    expect(defaultModeFor('read')).toBe('always')
    expect(defaultModeFor('write')).toBe('ask')
    expect(defaultModeFor('destructive')).toBe('ask')
    expect(defaultModeFor(undefined)).toBe('always')
  })

  it('reads a tool nothing grants as blocked', () => {
    const { named, families } = grants()
    expect(modeOf(read, named, families)).toBe('deny')
  })

  it('reads a named grant at the mode it states, or its risk default', () => {
    const { named, families } = grants(
      { name: read.name },
      { name: write.name, permission: 'always' },
    )
    expect(modeOf(read, named, families)).toBe('always')
    expect(modeOf(write, named, families)).toBe('always')
  })

  // A ceiling admits its own level and everything below it.
  it('covers a tool at or below a family ceiling, and no higher', () => {
    const families = [{ group: 'usage', maxRisk: 'write' }]
    expect(coveredByFamily(read, families)).toBe(true)
    expect(coveredByFamily(write, families)).toBe(true)
    expect(coveredByFamily(destroy, families)).toBe(false)
    expect(coveredByFamily({ groups: ['configuring'], risk: 'read' }, families)).toBe(false)
  })

  it('gives a covered tool its risk default without an entry of its own', () => {
    const { named, families } = grants({ group: 'usage', maxRisk: 'write' })
    expect(modeOf(read, named, families)).toBe('always')
    expect(modeOf(write, named, families)).toBe('ask')
    expect(modeOf(destroy, named, families)).toBe('deny')
  })

  // "All data tools, but never delete" is a family grant plus a named entry.
  it('lets a named block override the family that covers it', () => {
    const { named, families } = grants(
      { group: 'usage', maxRisk: 'destructive' },
      { name: destroy.name, permission: 'deny' },
    )
    expect(modeOf(destroy, named, families)).toBe('deny')
    expect(modeOf(read, named, families)).toBe('always')
  })

  it('counts each mode, leaving out the ones nothing sits in', () => {
    const { named, families } = grants({ name: read.name }, { name: write.name })
    expect(tally([read, write, destroy], named, families)).toEqual([
      { mode: 'always', n: 1 },
      { mode: 'ask', n: 1 },
      { mode: 'deny', n: 1 },
    ])
  })

  it('counts nothing granted as all blocked', () => {
    const { named, families } = grants()
    expect(tally([read, write], named, families)).toEqual([{ mode: 'deny', n: 2 }])
  })
})

describe('scoping applies to compose resources only', () => {
  // An allow entry describes a compose namespace and its modules. On 94 of the
  // 111 tools the policy check ignores it, and on the two exec tools it makes
  // the runtime refuse the call — so the control belongs on neither.
  it('scopes the compose tools a namespace can narrow', () => {
    expect(canScope('compose_record_lookup')).toBe(true)
    expect(canScope('compose_module_create')).toBe(true)
    expect(canScope('compose_namespace_lookup')).toBe(true)
  })

  it('does not scope a tool the check cannot narrow', () => {
    for (const name of [
      'system_user_create',
      'compose_page_create',
      'compose_chart_lookup',
      'automation_taq_create',
      'system_reminder_snooze',
    ]) {
      expect(canScope(name)).toBe(false)
    }
  })

  // Not merely ignored — an allow entry here makes the call fail.
  it('names the two an allow entry would block', () => {
    expect(scopeBlocked('automation_taq_exec')).toBe(true)
    expect(scopeBlocked('automation_workflow_exec')).toBe(true)
    expect(scopeBlocked('compose_record_lookup')).toBe(false)
    expect(scopeBlocked('system_user_create')).toBe(false)
  })

  // A namespace resource has no module segment, so modules cannot narrow it.
  it('offers modules only where the resource has them', () => {
    expect(scopesModules('compose_record_lookup')).toBe(true)
    expect(scopesModules('compose_module_lookup')).toBe(true)
    expect(scopesModules('compose_namespace_lookup')).toBe(false)
    expect(scopesModules('system_user_create')).toBe(false)
  })
})

describe('sectionsOf', () => {
  const label = (k: string) => k
  const tools = [
    { name: 'compose_record_lookup', groups: ['usage'], risk: 'read' },
    { name: 'compose_record_delete', groups: ['usage'], risk: 'destructive' },
    { name: 'system_user_create', groups: ['configuring'], risk: 'write' },
  ]

  // The panel summarises by section and the dialog groups by it; reading the
  // taxonomy from one place is what stops them disagreeing.
  it('reports each section by what it lets the agent do', () => {
    const { named, families } = splitGrants([
      { name: 'compose_record_lookup' },
      { name: 'system_user_create' },
    ])

    expect(sectionsOf(tools, named, families, label)).toEqual([
      {
        key: 'records',
        label: 'records',
        counts: [
          { mode: 'always', n: 1 },
          { mode: 'deny', n: 1 },
        ],
      },
      { key: 'people', label: 'people', counts: [{ mode: 'ask', n: 1 }] },
    ])
  })

  // A section the agent has nothing in is not worth a line.
  it('leaves out a section holding nothing the agent has', () => {
    const { named, families } = splitGrants([{ name: 'system_user_create' }])
    expect(sectionsOf(tools, named, families, label).map(s => s.key)).toEqual(['people'])
  })
})

describe('hasSettings', () => {
  it('is false for a grant that says nothing but its mode', () => {
    expect(hasSettings({ name: 'compose_record_lookup', permission: 'always' })).toBe(false)
    expect(hasSettings({ name: 'x', description: '', allow: [] })).toBe(false)
    expect(hasSettings(undefined)).toBe(false)
  })

  it('is true for a note or a narrowing', () => {
    expect(hasSettings({ name: 'x', description: 'only open leads' })).toBe(true)
    expect(hasSettings({ name: 'x', allow: [{ namespaceID: '1', moduleIDs: ['7'] }] })).toBe(true)
  })
})

// The policy check reads the agent's scope and the tool's and requires both, so
// a narrowing left naming the previous namespace denies that tool everything it
// can now reach — silently, since the dialog only renders narrowings for
// namespaces the agent still holds.
describe('confineTo', () => {
  const grants = [
    { name: 'compose_record_lookup', allow: [{ namespaceID: '100', moduleIDs: ['7'] }] },
    { name: 'compose_record_create', allow: [{ namespaceID: '200', moduleIDs: ['9'] }] },
    { name: 'system_user_lookup' },
  ]

  it('drops a narrowing against a namespace the agent no longer works in', () => {
    const { tools, cleared } = confineTo(grants, '200')

    expect(cleared).toBe(1)
    expect(tools[0].allow).toEqual([])
    expect(tools[1].allow).toEqual([{ namespaceID: '200', moduleIDs: ['9'] }])
  })

  it('drops every narrowing when the agent is confined to nothing', () => {
    const { tools, cleared } = confineTo(grants, null)

    expect(cleared).toBe(2)
    expect(tools.every(t => t.allow.length === 0)).toBe(true)
  })

  // The count reads as "N tools lost a limit", so a tool that loses three
  // entries is still one tool.
  it('counts tools, not entries', () => {
    const { tools, cleared } = confineTo(
      [
        {
          name: 'compose_record_lookup',
          allow: [
            { namespaceID: '100', moduleIDs: ['7'] },
            { namespaceID: '300', moduleIDs: ['8'] },
            { namespaceID: '400', moduleIDs: ['9'] },
          ],
        },
      ],
      '100',
    )

    expect(cleared).toBe(1)
    expect(tools[0].allow).toEqual([{ namespaceID: '100', moduleIDs: ['7'] }])
  })

  it('reports nothing cleared when every narrowing already names the namespace', () => {
    expect(confineTo([grants[1], grants[2]], '200').cleared).toBe(0)
  })

  it('compares IDs by value, so a number and its string are one namespace', () => {
    const { cleared } = confineTo(
      [{ name: 'x', allow: [{ namespaceID: 200, moduleIDs: ['9'] }] }],
      '200',
    )
    expect(cleared).toBe(0)
  })
})

describe('summaryOf', () => {
  const label = (key: string) => key
  const tools = [
    { name: 'compose_record_lookup', groups: ['usage'], risk: 'read' },
    { name: 'compose_record_create', groups: ['usage'], risk: 'write' },
    { name: 'compose_record_delete', groups: ['usage'], risk: 'destructive' },
    { name: 'system_user_create', groups: ['configuring'], risk: 'write' },
    // A skill is attached to the tool that triggers it, not chosen, so no
    // section holds it.
    { name: 'system_skill_lookup', groups: ['usage'], risk: 'read' },
  ]

  it('gives each subject one row carrying both counts', () => {
    const { named, families } = splitGrants([
      { name: 'compose_record_lookup' },
      { name: 'compose_record_create' },
      { name: 'system_user_create' },
    ])

    expect(summaryOf(tools, named, families, label)).toEqual({
      total: 3,
      totals: { always: 1, ask: 2 },
      subjects: [
        { key: 'records', label: 'records', always: 1, ask: 1 },
        { key: 'people', label: 'people', always: 0, ask: 1 },
      ],
    })
  })

  // A subject appears once however many permissions its tools sit in: reading
  // the split is reading across the row, not diffing two lists.
  it('names a subject once even when its tools sit in both columns', () => {
    const { named, families } = splitGrants([
      { name: 'compose_record_lookup' },
      { name: 'compose_record_create' },
    ])

    const { subjects } = summaryOf(tools, named, families, label)
    expect(subjects.map(s => s.key)).toEqual(['records'])
    expect(subjects[0]).toEqual({ key: 'records', label: 'records', always: 1, ask: 1 })
  })

  // The foot of the table is the sum of the rows above it. Counted over the
  // whole tool list instead, a granted skill — which no subject holds — put the
  // headline one above what it was heading.
  it('totals what the rows show and nothing else', () => {
    const { named, families } = splitGrants([
      { name: 'compose_record_lookup' },
      { name: 'system_skill_lookup' },
    ])

    const summary = summaryOf(tools, named, families, label)
    const shown = summary.subjects.reduce((n, s) => n + s.always + s.ask, 0)

    expect(summary.total).toBe(shown)
    expect(summary.total).toBe(1)
    expect(summary.totals).toEqual({ always: 1, ask: 0 })
    expect(summary.subjects.map(s => s.key)).toEqual(['records'])
  })

  // A column a subject has nothing in is a zero the row still carries, so the
  // readout can draw a dash there rather than leave the cell out of line.
  it('carries a zero for a column the subject has nothing in', () => {
    const { named, families } = splitGrants([{ name: 'system_user_create' }])
    expect(summaryOf(tools, named, families, label).subjects).toEqual([
      { key: 'people', label: 'people', always: 0, ask: 1 },
    ])
  })

  it('is empty for an agent holding nothing', () => {
    const { named, families } = splitGrants([])
    expect(summaryOf(tools, named, families, label)).toEqual({
      total: 0,
      totals: { always: 0, ask: 0 },
      subjects: [],
    })
  })
})

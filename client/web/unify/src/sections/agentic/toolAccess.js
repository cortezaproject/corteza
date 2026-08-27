// What an agent's tool grants add up to.
//
// The runtime resolves a grant that names no mode when it runs the agent
// (`withResolvedPermissions`), so the editor has to resolve it the same way to
// show the agent back to its author. One definition, used by the tool dialog
// and by the panel that summarises it.

// Reads first, then writes, then deletes: the harmless surface leads, and the
// one that cannot be taken back is last.
export const RISK_ORDER = { read: 0, write: 1, destructive: 2 }

// Custom has none. It is not a fourth state a tool can be in, it is three
// sections' worth of disagreement, and a glyph beside the word claimed it was
// a choice like the others.
export const MODE_ICONS = {
  always: 'pi pi-check-circle',
  ask: 'pi pi-question-circle',
  deny: 'pi pi-ban',
  custom: '',
}

// Ask keeps the body colour: it is the mode most tools sit in until someone
// decides otherwise, and colouring it too would leave nothing quiet to read the
// other two against.
export const MODE_COLOURS = {
  always: 'text-green-500',
  ask: '',
  deny: 'text-red-500',
  custom: '',
}

// The order they are counted and shown in: what runs freely, what stops to ask,
// what the agent does not have.
export const MODES = ['always', 'ask', 'deny']

// Reading changes nothing and runs unannounced; anything that writes is put to
// the user first.
export function defaultModeFor(risk) {
  return risk && risk !== 'read' ? 'ask' : 'always'
}

// A ceiling admits its own level and everything below it. A destructive tool
// under a write ceiling is not covered and stays available to grant by name.
export function coveredByFamily(tool, families) {
  return (families || []).some(
    f =>
      (tool.groups || []).includes(f.group) &&
      (RISK_ORDER[tool.risk] ?? 0) <= (RISK_ORDER[f.maxRisk || 'read'] ?? 0),
  )
}

// What one tool ends up doing. Blocked and absent are one state: an agent that
// may not use a tool and an agent that was never given it come to the same
// thing.
export function modeOf(tool, named, families) {
  const permission = named.get(tool.name)

  if (coveredByFamily(tool, families)) {
    return permission === 'deny' ? 'deny' : permission || defaultModeFor(tool.risk)
  }

  if (permission === undefined) return 'deny'
  return permission || defaultModeFor(tool.risk)
}

// Grants as stored, split into the two shapes the rest of this reads.
export function splitGrants(grants) {
  return {
    named: new Map((grants || []).filter(g => g.name).map(g => [g.name, g.permission || ''])),
    families: (grants || []).filter(g => g.group),
  }
}

// How many tools sit in each mode, in MODES order, leaving out the ones nothing
// sits in.
export function tally(tools, named, families) {
  const counts = { always: 0, ask: 0, deny: 0 }
  for (const tool of tools) counts[modeOf(tool, named, families)]++
  return MODES.filter(mode => counts[mode]).map(mode => ({ mode, n: counts[mode] }))
}

// Skills are attached to a tool, not chosen: the runtime injects one when the
// tool that triggers it is used. Offering them here invites a choice that
// changes nothing.
export const HIDDEN_AREAS = new Set(['system_skill'])

// Sections are subjects, in the order an agent meets them: what it works with,
// then what it runs, then what it builds on, then who it touches.
//
// Grouping by `usage` and `configuring` instead put 17 tools in one section and
// 94 in the other, and split five subjects across both — every TAQ tool but
// `exec` in one section, `exec` in the other. A subject now appears once, with
// all of its tools, whichever group each one belongs to.
//
// A section holds one subject and not two: reminders are a personal surface and
// share nothing with records but their group, and the schema a namespace
// defines is a different job from the pages that display it.
export const DOMAINS = [
  { key: 'records', areas: ['compose_record'] },
  { key: 'datamodel', areas: ['compose_namespace', 'compose_module'] },
  { key: 'interface', areas: ['compose_page', 'compose_chart'] },
  {
    key: 'automation',
    areas: ['automation_taq', 'automation_workflow', 'automation_trigger', 'automation_event'],
  },
  { key: 'people', areas: ['system_user', 'system_role', 'system_auth'] },
  { key: 'ai', areas: ['system_agent', 'system_chatbot'] },
  { key: 'reminders', areas: ['system_reminder'] },
  { key: 'workspace', areas: ['system_application', 'system_theme'] },
]

export const PLACED_AREAS = new Set(DOMAINS.flatMap(d => d.areas))

// compose_record_lookup -> compose_record; discovery_search -> discovery.
export function areaOf(name) {
  const parts = String(name).split('_')
  return parts.length > 2 ? `${parts[0]}_${parts[1]}` : parts[0]
}

// Which tools a namespace/module scope actually narrows.
//
// An allow entry describes a compose namespace and its modules, and the policy
// check can say nothing about a resource of any other kind: on 94 of the tools
// it is ignored, and on the two exec tools it makes the runtime refuse the call
// outright. Offering the control on all of them is how a setting that does
// nothing came to look like one that does.
const SCOPABLE = ['compose_record_', 'compose_module_', 'compose_namespace_']

// Scoped by naming the TAQ or workflow instead, not by namespace.
const SCOPE_BLOCKED = ['automation_taq_exec', 'automation_workflow_exec']

export function canScope(name) {
  return SCOPABLE.some(prefix => String(name).startsWith(prefix))
}

export function scopeBlocked(name) {
  return SCOPE_BLOCKED.includes(name)
}

// A namespace scope reaches modules everywhere except on the namespace tools,
// whose resource has no module segment to narrow.
export function scopesModules(name) {
  return canScope(name) && !String(name).startsWith('compose_namespace_')
}

// The sections a set of tools falls into, each with what it lets the agent do.
// The dialog groups its rows by this and the panel summarises by it, so they
// cannot disagree about which subject a tool belongs to.
export function sectionsOf(tools, named, families, label) {
  const byArea = new Map()

  for (const tool of tools) {
    const area = areaOf(tool.name)
    if (HIDDEN_AREAS.has(area)) continue
    if (!byArea.has(area)) byArea.set(area, [])
    byArea.get(area).push(tool)
  }

  const spec = [
    ...DOMAINS,
    { key: 'other', areas: [...byArea.keys()].filter(k => !PLACED_AREAS.has(k)).sort() },
  ]

  return spec
    .map(d => {
      const held = d.areas.filter(k => byArea.has(k)).flatMap(k => byArea.get(k))
      return { key: d.key, label: label(d.key), counts: tally(held, named, families) }
    })
    .filter(d => d.counts.some(c => c.mode !== 'deny'))
}

// Whether a grant carries anything beyond its mode — the note the model reads
// before calling the tool, or a narrowing. A row shows this back so a
// configured tool can be told from an untouched one without opening it.
export function hasSettings(entry) {
  return Boolean(entry?.description || entry?.allow?.length)
}

// The grants an agent keeps once it is confined to one namespace.
//
// A tool's own narrowing names a namespace, and the modules within it. Left
// naming one the agent no longer works in, it names nothing the agent can
// reach: the policy check reads the agent's scope and the tool's and requires
// both, so such a tool is denied everything. Dropping the narrowing is the only
// reading that leaves it usable.
export function confineTo(grants, namespaceID) {
  let cleared = 0

  const tools = (grants || []).map(grant => {
    const allow = grant.allow || []
    const kept = allow.filter(a => String(a.namespaceID) === String(namespaceID))
    if (kept.length !== allow.length) cleared++
    return { ...grant, allow: kept }
  })

  return { tools, cleared }
}

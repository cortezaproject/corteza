// What an agent's tool grants add up to.
//
// The runtime resolves a grant that names no mode when it runs the agent
// (`withResolvedPermissions`), so the editor has to resolve it the same way to
// show the agent back to its author. One definition, used by the tool dialog
// and by the panel that summarises it.

// Reads first, then writes, then deletes: the harmless surface leads, and the
// one that cannot be taken back is last.
export const RISK_ORDER = { read: 0, write: 1, destructive: 2 }

export const MODE_ICONS = {
  always: 'pi pi-check-circle',
  ask: 'pi pi-question-circle',
  deny: 'pi pi-ban',
  custom: 'pi pi-ellipsis-h',
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

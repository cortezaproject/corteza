// The inventory the dashboard shows: one row per resource the stats endpoint
// counts, in display order, with where a click lands and how its statuses
// stack. Status keys not listed here are still shown, after the listed ones.
export const RESOURCES = [
  { key: 'users', route: { name: 'system.users' }, statuses: ['active', 'suspended', 'deleted'] },
  { key: 'roles', route: { name: 'system.roles' }, statuses: ['active', 'archived', 'deleted'] },
  {
    key: 'projects',
    route: { name: 'project.list' },
    statuses: ['active', 'published', 'draft', 'suspended', 'archived', 'deprecated', 'deleted'],
  },
  {
    key: 'namespaces',
    route: { name: 'namespace.list' },
    statuses: ['active', 'disabled', 'deleted'],
  },
  { key: 'modules', route: { name: 'namespace.list' }, statuses: ['active', 'deleted'] },
  {
    key: 'workflows',
    route: { name: 'automation.workflows' },
    statuses: ['enabled', 'disabled', 'deleted'],
  },
  { key: 'taqs', route: { name: 'automation.taq' }, statuses: ['enabled', 'disabled', 'deleted'] },
  { key: 'agents', route: { name: 'agentic' }, statuses: ['active', 'deleted'] },
  { key: 'chatbots', route: { name: 'chatbot' }, statuses: ['enabled', 'disabled', 'deleted'] },
  {
    key: 'applications',
    route: { name: 'system.applications' },
    statuses: ['enabled', 'disabled', 'deleted'],
  },
  {
    key: 'authClients',
    route: { name: 'system.authClients' },
    statuses: ['enabled', 'disabled', 'deleted'],
  },
]

// Statuses that mean "in use": the number a tile leads with.
export const LIVE_STATUSES = ['active', 'enabled', 'published']

// Workflow session outcomes, worst first for the stacked chart; the three
// non-terminal states fold into one "running" series.
export const WORKFLOW_OUTCOMES = [
  { key: 'failed', from: ['failed'] },
  { key: 'running', from: ['started', 'prompted', 'suspended'] },
  { key: 'canceled', from: ['canceled'] },
  { key: 'completed', from: ['completed'] },
]

export const TAQ_OUTCOMES = [
  { key: 'failed', from: ['failed'] },
  { key: 'canceled', from: ['cancelled'] },
  { key: 'completed', from: ['completed'] },
]

// Statuses a resource may report, ordered so a stacked bar reads live → gone.
export function orderedStatuses(resource, status) {
  const listed = resource.statuses.filter(k => status[k] !== undefined)
  const extra = Object.keys(status)
    .filter(k => !resource.statuses.includes(k))
    .sort()
  return [...listed, ...extra]
}

export function liveCount(status) {
  return LIVE_STATUSES.reduce((sum, k) => sum + (status[k] || 0), 0)
}

// Sums the raw per-status series of a runs section into the display outcomes.
export function foldOutcomes(runs, outcomes, length) {
  const byStatus = {}
  const series = {}
  for (const o of outcomes) {
    byStatus[o.key] = o.from.reduce((sum, k) => sum + (runs?.byStatus?.[k] || 0), 0)
    const s = new Array(length).fill(0)
    for (const k of o.from) {
      const src = runs?.series?.[k] || []
      for (let i = 0; i < length; i++) s[i] += src[i] || 0
    }
    series[o.key] = s
  }
  return { byStatus, series }
}

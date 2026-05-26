// The project build pipeline — single source of truth for steps, gates and
// mode visibility. Reused by the wizard later; used here for the list's
// "Build progress" column. See ../../IDEA.md for the conceptual model.

export const STEPS = [
  { key: 'summary', label: 'Project Summary', type: 'form', gatedOnly: true, gate: false },
  { key: 'members', label: 'Project Members', type: 'members', gatedOnly: true, gate: true },
  { key: 'data-model', label: 'Data Model', type: 'resource', kind: 'module', gatedOnly: false, gate: false },
  { key: 'connections', label: 'Connections', type: 'resource', kind: 'connection', gatedOnly: false, gate: false },
  { key: 'data-sensitivity', label: 'Data Sensitivity', type: 'form', gatedOnly: true, gate: true },
  { key: 'automations', label: 'Automations', type: 'resource', kind: 'automation', gatedOnly: false, gate: false },
  { key: 'agents', label: 'Agents', type: 'resource', kind: 'agent', gatedOnly: false, gate: false },
  { key: 'chatbots', label: 'Chatbots', type: 'resource', kind: 'chatbot', gatedOnly: false, gate: true },
  { key: 'pages', label: 'Pages', type: 'resource', kind: 'page', gatedOnly: false, gate: true },
  { key: 'rbac', label: 'RBAC', type: 'resource', kind: 'role', gatedOnly: false, gate: false },
  { key: 'users', label: 'Users', type: 'resource', kind: 'user', gatedOnly: false, gate: false },
  { key: 'security', label: 'Security', type: 'form', gatedOnly: true, gate: true },
  { key: 'architecture', label: 'Technical Architecture', type: 'groups', gatedOnly: false, gate: false },
  { key: 'risk-classification', label: 'Risk Classification', type: 'form', gatedOnly: true, gate: true },
  { key: 'risk', label: 'Risk Management', type: 'form', gatedOnly: true, gate: false },
  { key: 'monitoring', label: 'Monitoring', type: 'form', gatedOnly: true, gate: true },
  {
    key: 'preview',
    label: 'Preview',
    type: 'milestone',
    gatedOnly: false,
    gate: false,
    description:
      'Review the whole project and provision the real resources with their RBAC. Sets the project to preview status — until now only the developer/approver can view the resources (e.g. a module in development is only visible to them).',
  },
  {
    key: 'publish',
    label: 'Publish',
    type: 'milestone',
    gatedOnly: false,
    gate: false,
    description:
      'Lock the configuration as a version. After publishing, any further change goes through the change-request procedure.',
  },
]

// Steps visible for a given build mode (Free hides the governance-only steps).
export const stepsForMode = mode => STEPS.filter(s => mode === 'gated' || !s.gatedOnly)

// Steps shown in a wizard tab. Build = the Free-mode subset; Governance = everything.
export const stepsForTab = tab => (tab === 'build' ? STEPS.filter(s => !s.gatedOnly) : STEPS)

// Total gates a mode must clear (gates only apply in gated mode).
export const gateCount = mode => (mode === 'gated' ? STEPS.filter(s => s.gate).length : 0)

// Group the pipeline into gated sections; each section ends at a gate step
// (inclusive). A trailing group with gateKey=null holds any post-final-gate steps.
export function sections(stepList = STEPS) {
  const out = []
  let current = []
  for (const s of stepList) {
    current.push(s)
    if (s.gate) {
      out.push({ gateKey: s.key, steps: current })
      current = []
    }
  }
  if (current.length) out.push({ gateKey: null, steps: current })
  return out
}

// The project build pipeline — single source of truth for steps, gates and
// mode visibility. We grow this one verified step at a time; only steps whose
// persistence and approval flow are fully backend-backed belong here. The full
// conceptual pipeline lives in the PoC copy (../project-poc).
//
// Project Members carries Gate 1, closing the section: Summary + Resource
// Management + Members. Data Sensitivity carries Gate 2, closing the Data Model
// + Data Sensitivity section. Free mode skips the governance steps entirely and
// starts at the Data Model.
// `labelKey` is an i18n key; components resolve it with $t for display.
// `icon` (PrimeIcons class) is shown in the step nav for steps without a resource
// `kind`; resource steps take their icon from config/kinds so the sidebar matches
// the metrics strip and resource graph exactly.
export const STEPS = [
  { key: 'summary', labelKey: 'project.steps.summary.label', type: 'form', icon: 'pi-file', gatedOnly: true, gate: false },
  { key: 'resource-management', labelKey: 'project.steps.resource-management.label', type: 'form', icon: 'pi-sliders-h', gatedOnly: true, gate: false },
  { key: 'members', labelKey: 'project.steps.members.label', type: 'members', icon: 'pi-users', gatedOnly: true, gate: true },
  { key: 'data-model', labelKey: 'project.steps.data-model.label', type: 'resource', kind: 'module', gatedOnly: false, gate: false },
  { key: 'data-sensitivity', labelKey: 'project.steps.data-sensitivity.label', type: 'sensitivity', icon: 'pi-eye-slash', gatedOnly: true, gate: true },
  { key: 'connections', labelKey: 'project.steps.connections.label', type: 'resource', kind: 'connection', gatedOnly: false, gate: false },
  { key: 'automations', labelKey: 'project.steps.automations.label', type: 'resource', kind: 'automation', gatedOnly: false, gate: false },
  { key: 'agents', labelKey: 'project.steps.agents.label', type: 'resource', kind: 'agent', gatedOnly: false, gate: false },
  { key: 'chatbots', labelKey: 'project.steps.chatbots.label', type: 'resource', kind: 'chatbot', gatedOnly: false, gate: false },
  { key: 'pages', labelKey: 'project.steps.pages.label', type: 'resource', kind: 'page', gatedOnly: false, gate: false },
  { key: 'roles', labelKey: 'project.steps.roles.label', type: 'resource', kind: 'role', gatedOnly: false, gate: false },
  { key: 'permissions', labelKey: 'project.steps.permissions.label', type: 'permissions', icon: 'pi-lock', gatedOnly: false, gate: false },
  { key: 'users', labelKey: 'project.steps.users.label', type: 'resource', kind: 'user', gatedOnly: false, gate: false },
  { key: 'publish', labelKey: 'project.steps.publish.label', type: 'publish', icon: 'pi-cloud-upload', gatedOnly: false, gate: false },
]

// Resolve the pipeline for a project context. Kept as a hook for future
// conditional steps (e.g. FRIA); currently the pipeline is static.
export const resolveSteps = () => STEPS.map(s => ({ ...s }))

// Steps visible for a given build mode (Free hides the governance-only steps).
export const stepsForMode = mode =>
  resolveSteps().filter(s => mode === 'gated' || !s.gatedOnly)

// Steps shown in a wizard tab. Build = the Free-mode subset; Governance = everything.
export const stepsForTab = tab =>
  resolveSteps().filter(s => (tab === 'build' ? !s.gatedOnly : true))

// Total gates a mode must clear (gates only apply in gated mode).
export const gateCount = mode => (mode === 'gated' ? STEPS.filter(s => s.gate).length : 0)

// The build is a process: a resource kind only exists once its step is reached.
// Returns the resource kinds introduced by every step up to and including
// `stepKey` (inclusive) — drives the resource graph so it shows only what's been
// built so far. Order is owned here by STEPS; add a resource step and it joins
// automatically. Steps before any resource step contribute nothing.
export function kindsThroughStep(stepKey) {
  const out = new Set()
  for (const s of STEPS) {
    if (s.kind) out.add(s.kind)
    if (s.key === stepKey) break
  }
  return out
}

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

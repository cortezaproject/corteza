// The project build pipeline — single source of truth for steps and mode
// visibility. We grow this one verified step at a time; only steps whose
// persistence flow is fully backend-backed belong here.
//
// Free mode skips the governance-only steps entirely and starts at the Data
// Model. Gated mode shows every step; approval only gates publishing itself
// (see the Publish step), not any individual step here.
// `labelKey` is an i18n key; components resolve it with $t for display.
// `icon` (PrimeIcons class) is shown in the step nav for steps without a resource
// `kind`; resource steps take their icon from config/kinds so the sidebar matches
// the metrics strip and resource graph exactly.
// The well-known governance step key that drives the publish-time submit →
// approve/request-changes cycle (gated mode only) — mirrors the backend's
// types.ProjectGovernanceStepPublish. Every other step key only ever persists
// form values via SaveGovernanceStep and never locks (see stores/projects.js).
export const PUBLISH_GOVERNANCE_STEP_KEY = 'publish'

export const STEPS = [
  {
    key: 'summary',
    labelKey: 'project.steps.summary.label',
    type: 'form',
    icon: 'pi-file',
    gatedOnly: true,
  },
  {
    key: 'resource-management',
    labelKey: 'project.steps.resource-management.label',
    type: 'form',
    icon: 'pi-sliders-h',
    gatedOnly: true,
  },
  {
    key: 'members',
    labelKey: 'project.steps.members.label',
    type: 'members',
    icon: 'pi-users',
    gatedOnly: false,
  },
  {
    key: 'data-model',
    labelKey: 'project.steps.data-model.label',
    type: 'resource',
    kind: 'module',
    gatedOnly: false,
  },
  {
    key: 'data-sensitivity',
    labelKey: 'project.steps.data-sensitivity.label',
    type: 'sensitivity',
    icon: 'pi-eye-slash',
    gatedOnly: true,
  },
  {
    key: 'connections',
    labelKey: 'project.steps.connections.label',
    type: 'resource',
    kind: 'connection',
    gatedOnly: false,
  },
  {
    key: 'automations',
    labelKey: 'project.steps.automations.label',
    type: 'resource',
    kind: 'automation',
    gatedOnly: false,
  },
  {
    key: 'agents',
    labelKey: 'project.steps.agents.label',
    type: 'resource',
    kind: 'agent',
    gatedOnly: false,
  },
  {
    key: 'chatbots',
    labelKey: 'project.steps.chatbots.label',
    type: 'resource',
    kind: 'chatbot',
    gatedOnly: false,
  },
  {
    key: 'pages',
    labelKey: 'project.steps.pages.label',
    type: 'resource',
    kind: 'page',
    gatedOnly: false,
  },
  {
    key: 'roles',
    labelKey: 'project.steps.roles.label',
    type: 'resource',
    kind: 'role',
    gatedOnly: false,
  },
  {
    key: 'permissions',
    labelKey: 'project.steps.permissions.label',
    type: 'permissions',
    icon: 'pi-lock',
    gatedOnly: false,
  },
  {
    key: 'users',
    labelKey: 'project.steps.users.label',
    type: 'resource',
    kind: 'user',
    gatedOnly: false,
  },
  {
    key: 'publish',
    labelKey: 'project.steps.publish.label',
    type: 'publish',
    icon: 'pi-cloud-upload',
    gatedOnly: false,
  },
]

// Resolve the pipeline for a project context. Kept as a hook for future
// conditional steps (e.g. FRIA); currently the pipeline is static.
export const resolveSteps = () => STEPS.map(s => ({ ...s }))

// Steps visible for a given build mode (Free hides the governance-only steps).
export const stepsForMode = mode => resolveSteps().filter(s => mode === 'gated' || !s.gatedOnly)

// Steps shown in a wizard tab. Build = the Free-mode subset; Governance = everything.
export const stepsForTab = tab =>
  resolveSteps().filter(s => (tab === 'build' ? !s.gatedOnly : true))

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

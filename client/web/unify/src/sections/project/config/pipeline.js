// The project build pipeline — single source of truth for steps and which
// wizard tab each belongs to. We grow this one verified step at a time; only
// steps whose persistence flow is fully backend-backed belong here.
//
// Every project behaves identically now (there is no build mode): the wizard
// always shows all three tabs — Build, Govern, Manage & Monitor — to every
// member. `tab` below assigns each step to Build or Govern; Manage & Monitor
// has no steps of its own (see Wizard.vue).
// `labelKey` is an i18n key; components resolve it with $t for display.
// `icon` (PrimeIcons class) is shown in the step nav for steps without a resource
// `kind`; resource steps take their icon from config/kinds so the sidebar matches
// the metrics strip and resource graph exactly.
// The well-known governance step key that drives the publish-time submit →
// approve/request-changes cycle. It is not a wizard step any more (Publish
// lives in the topbar toolbar cluster instead — see Wizard.vue); every other
// step key only ever persists form values, and its governance status via the
// direct approve/request-changes actions (see stores/projects.js
// transitionStep — the whole governance workflow is session-local scaffolding
// there, pending a redesign).
export const PUBLISH_GOVERNANCE_STEP_KEY = 'publish'

export const STEPS = [
  {
    key: 'summary',
    labelKey: 'project.steps.summary.label',
    type: 'form',
    icon: 'pi-file',
    tab: 'govern',
  },
  {
    key: 'resource-management',
    labelKey: 'project.steps.resource-management.label',
    type: 'form',
    icon: 'pi-sliders-h',
    tab: 'govern',
  },
  {
    key: 'data-model',
    labelKey: 'project.steps.data-model.label',
    type: 'resource',
    kind: 'module',
    tab: 'build',
  },
  {
    // Positioned right after data-model (not with the other Govern steps)
    // so kindsThroughStep('data-sensitivity') already includes 'module' —
    // this step classifies the fields data-model just built, so the graph
    // should show those modules while you're on it. Tab membership (Govern)
    // is independent of this array position; see stepsForTab.
    key: 'data-sensitivity',
    labelKey: 'project.steps.data-sensitivity.label',
    type: 'sensitivity',
    icon: 'pi-eye-slash',
    tab: 'govern',
  },
  // --- FRIA (Fundamental Rights Impact Assessment, EU AI Act Art. 27) -----
  // Several Govern steps, one per section of the design mockup (see
  // client/web/unify/src/sections/project/components/wizard/steps/Fria*.vue):
  // a determination step (the AI Act deployer-category questions, moved here
  // from the create flow per the 2026-07-28 ruling), a scenario list/summary
  // step, and five section-editor steps that each edit one slice of whichever
  // scenario is currently active (see composables/useFriaActiveScenario.js —
  // the `scenario` route query param threads the active scenario id through
  // all five, since there can be many scenarios per project but the mockup's
  // five sections are one continuous per-scenario form). All FRIA copy lives
  // in locale/en/human-webapp/fria.yaml, not project.yaml, on purpose.
  {
    key: 'fria-determination',
    labelKey: 'fria.steps.determination.label',
    type: 'form',
    icon: 'pi-question-circle',
    tab: 'govern',
  },
  {
    key: 'fria-scenarios',
    labelKey: 'fria.steps.scenarios.label',
    type: 'fria-scenarios',
    icon: 'pi-table',
    tab: 'govern',
  },
  {
    key: 'fria-harm',
    labelKey: 'fria.steps.harm.label',
    type: 'fria-scenario',
    icon: 'pi-exclamation-triangle',
    tab: 'govern',
  },
  {
    key: 'fria-trigger',
    labelKey: 'fria.steps.trigger.label',
    type: 'fria-scenario',
    icon: 'pi-bolt',
    tab: 'govern',
  },
  {
    key: 'fria-parties',
    labelKey: 'fria.steps.parties.label',
    type: 'fria-scenario',
    icon: 'pi-users',
    tab: 'govern',
  },
  {
    key: 'fria-rights',
    labelKey: 'fria.steps.rights.label',
    type: 'fria-scenario',
    icon: 'pi-verified',
    tab: 'govern',
  },
  {
    key: 'fria-vectors',
    labelKey: 'fria.steps.vectors.label',
    type: 'fria-scenario',
    icon: 'pi-sitemap',
    tab: 'govern',
  },
  {
    key: 'connections',
    labelKey: 'project.steps.connections.label',
    type: 'resource',
    kind: 'connection',
    tab: 'build',
  },
  {
    key: 'automations',
    labelKey: 'project.steps.automations.label',
    type: 'resource',
    kind: 'automation',
    tab: 'build',
  },
  {
    key: 'agents',
    labelKey: 'project.steps.agents.label',
    type: 'resource',
    kind: 'agent',
    tab: 'build',
  },
  {
    key: 'chatbots',
    labelKey: 'project.steps.chatbots.label',
    type: 'resource',
    kind: 'chatbot',
    tab: 'build',
  },
  {
    key: 'pages',
    labelKey: 'project.steps.pages.label',
    type: 'resource',
    kind: 'page',
    tab: 'build',
  },
  {
    key: 'roles',
    labelKey: 'project.steps.roles.label',
    type: 'resource',
    kind: 'role',
    tab: 'build',
  },
  {
    key: 'permissions',
    labelKey: 'project.steps.permissions.label',
    type: 'permissions',
    icon: 'pi-lock',
    tab: 'build',
  },
  {
    key: 'users',
    labelKey: 'project.steps.users.label',
    type: 'resource',
    kind: 'user',
    tab: 'build',
  },
]

// Resolve the pipeline for a project context. Kept as a hook for future
// conditional steps; currently the pipeline is static (the FRIA steps above
// are unconditional too — every project gets them, same as every other step).
export const resolveSteps = () => STEPS.map(s => ({ ...s }))

// Steps shown in a wizard tab ('build' | 'govern'). Manage & Monitor has no
// steps of its own — callers should treat it as an empty list.
export const stepsForTab = tab => resolveSteps().filter(s => s.tab === tab)

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

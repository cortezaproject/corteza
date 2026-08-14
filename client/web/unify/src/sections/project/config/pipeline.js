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
// The well-known governance step key the publish-time review hangs off. It is
// not a wizard step any more (Publish lives in its own tab instead — see
// Wizard.vue), but it runs exactly the cycle every step below runs: submit →
// approve/request-changes, with its own extra gate that no other step may sit
// flagged when the revision is approved (see stores/projects.js transitionStep
// — the whole governance workflow is session-local scaffolding there, pending
// a redesign).
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
    // because it classifies the fields data-model just built. Tab membership
    // (Govern) is independent of this array position; see stepsForTab.
    key: 'data-sensitivity',
    labelKey: 'project.steps.data-sensitivity.label',
    type: 'sensitivity',
    icon: 'pi-eye-slash',
    tab: 'govern',
  },
  // --- FRIA (Fundamental Rights Impact Assessment, EU AI Act Art. 27) -----
  // THREE Govern steps, in the order the assessment actually reasons:
  //
  //   1. determination — are WE, as deployer, in scope at all? Organisation-
  //      level, so it stays a revision-wide answer.
  //   2. ai-systems — WHAT are we assessing? A project can hold several AI
  //      systems (Art. 3(1) makes the AI system the regulated unit, not the
  //      project), each a named set of project resources carrying its own
  //      intended purpose and Art. 6 risk class. Nothing downstream can be
  //      classified before this exists, which is why it precedes scenarios.
  //   3. scenarios — HOW could each system cause harm? Every scenario names
  //      exactly one AI system (config/friaScenario.js).
  //
  // Determination is deliberately split from per-system classification:
  // "are we a public authority / essential-services / banking deployer" is one
  // answer for the whole revision, while "is THIS system
  // Annex III high-risk" differs per system — a project with one high-risk and
  // one minimal-risk system is the common case, and one conflated step cannot
  // express it.
  //
  // The scenario step's body switches between the scenario list and a
  // single scrolling per-scenario editor (see
  // components/wizard/steps/FriaScenariosStep.vue and
  // components/wizard/fria/FriaScenarioEditor.vue, which stacks the design
  // mockup's five sections as headed regions on one page) — which of the two
  // it shows is driven by composables/useFriaActiveScenario.js's `scenario`
  // route query param, so an open scenario stays deep-linkable and back/
  // forward still work, without needing a nav entry per section. All FRIA
  // copy lives in locale/en/human-webapp/fria.yaml, not project.yaml, on
  // purpose.
  {
    key: 'fria-determination',
    labelKey: 'fria.steps.determination.label',
    type: 'form',
    icon: 'pi-question-circle',
    tab: 'govern',
  },
  {
    // NOT a resource step and NOT a kind: an AI system is a boundary drawn over
    // resources the Build steps already created, so it stays out of
    // config/kinds.js, out of the graph as a node, and out of the permission
    // matrix. Its own type, like the FRIA steps.
    key: 'ai-systems',
    labelKey: 'fria.steps.aiSystems.label',
    type: 'ai-systems',
    icon: 'pi-sitemap',
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

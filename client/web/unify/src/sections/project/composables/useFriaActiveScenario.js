import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// Shared "which FRIA risk scenario is being edited" state for the Govern
// tab's fria-scenarios step (see config/pipeline.js) — that ONE step's body
// switches between the scenario list and a single scrolling scenario editor
// (see components/wizard/steps/FriaScenariosStep.vue and
// components/wizard/fria/FriaScenarioEditor.vue): this composable is what
// lets those two agree on whether a scenario is open, and which one,
// without Wizard.vue owning any FRIA-specific state itself (out of scope for
// this pass — see project.intent.md, Wizard.vue is dispatch only).
//
// Backed by the `scenario` route query param — same idiom Wizard.vue itself
// uses for `step`/`section`/`tab` — so the list <-> editor transition is
// deep-linkable and survives back/forward navigation.
//
// The sentinel value 'new' marks "the editor is open in create mode" rather
// than naming a real scenario. This matters for the draft lifecycle (see
// components/wizard/fria/FriaScenarioEditor.vue and stores/projects.js'
// createFriaScenario): a brand-new scenario is held as a local draft and
// never written to the store until the editor's explicit Save, so opening
// with this sentinel — then cancelling — leaves the store untouched; there
// is nothing to clean up.
const NEW_SCENARIO_SENTINEL = 'new'

export function useFriaActiveScenario() {
  const route = useRoute()
  const router = useRouter()

  const activeScenarioId = computed(() => route.query.scenario || null)
  const isCreatingScenario = computed(() => activeScenarioId.value === NEW_SCENARIO_SENTINEL)

  // Open an existing scenario's editor — used by the list's Edit actions.
  function openScenario(scenarioId) {
    router.replace({ query: { ...route.query, scenario: scenarioId } })
  }

  // Open the editor in create mode (see the sentinel note above) — used by
  // the list's New/"Document new" actions. Deliberately does not touch the
  // store: the draft only becomes real on Save.
  function openNewScenario() {
    router.replace({ query: { ...route.query, scenario: NEW_SCENARIO_SENTINEL } })
  }

  // Back to the scenario list. Used by both the editor's Cancel (the local
  // draft is simply dropped — nothing to undo in the store) and its Save
  // (the draft was already committed by then).
  function closeScenario() {
    router.replace({ query: { ...route.query, scenario: undefined } })
  }

  return { activeScenarioId, isCreatingScenario, openScenario, openNewScenario, closeScenario }
}

import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// Shared "which FRIA risk scenario is being edited" state for the Govern
// tab's fria-scenarios (list) step and its five section-editor steps
// (fria-harm/trigger/parties/rights/vectors — see config/pipeline.js). There
// can be many scenarios per project, but the design mockup's five sections
// are one continuous per-scenario form, split here into five separate wizard
// steps (per the ruling recorded in config/pipeline.js) — this composable is
// what lets all six of those step components agree on which scenario they're
// looking at without Wizard.vue owning any FRIA-specific state itself (out of
// scope for this pass — see project.intent.md, Wizard.vue is dispatch only).
//
// Backed by the `scenario` route query param — same idiom Wizard.vue itself
// uses for `step`/`section`/`tab` — so every mounted step component reads and
// writes the exact same value reactively via the shared route/router
// singleton, and the active scenario survives switching between the five
// section-editor steps (though not a reload — scenario data itself is
// session-local scaffolding per stores/projects.js, so a reload has nothing
// to resume to anyway).
export function useFriaActiveScenario() {
  const route = useRoute()
  const router = useRouter()

  const activeScenarioId = computed(() => route.query.scenario || null)

  // Jump straight into a scenario's editor (defaults to the first section —
  // Harm Scenario Description) — used by the list step's New/Edit actions.
  function openScenario(scenarioId, stepKey = 'fria-harm') {
    router.replace({ query: { ...route.query, scenario: scenarioId, step: stepKey } })
  }

  // Switch the active scenario without leaving the current section-editor
  // step — used by the editor shell's scenario switcher.
  function setActiveScenarioId(scenarioId) {
    router.replace({ query: { ...route.query, scenario: scenarioId || undefined } })
  }

  // Back to the scenario list step, keeping whichever scenario was active (so
  // returning to an editor step resumes on the same one).
  function goToScenarios() {
    router.replace({ query: { ...route.query, step: 'fria-scenarios' } })
  }

  return { activeScenarioId, openScenario, setActiveScenarioId, goToScenarios }
}

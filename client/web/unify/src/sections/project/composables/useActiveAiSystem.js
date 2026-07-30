import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// Shared "which AI system is open" state for the Govern tab's ai-systems step
// (see config/pipeline.js) — that ONE step's body switches between the system
// list and a single system's editor, exactly like the fria-scenarios step does
// for scenarios. Twin of composables/useFriaActiveScenario.js; kept as its own
// composable rather than generalised because the two differ where it matters
// (see the create-mode note below).
//
// Backed by the `aiSystem` route query param — same idiom Wizard.vue uses for
// `step`/`section`/`tab` — so the list <-> editor transition is deep-linkable
// and survives back/forward navigation.
//
// NO 'new' SENTINEL, deliberately — this is the one place it diverges from
// useFriaActiveScenario. A FRIA scenario is held as a local draft and written
// only on an explicit Save, so its editor can open against nothing. An AI
// system is REAL PERSISTED backend state (ProjectAiSystem) whose membership
// entries are child rows keyed by its ID: there is no ID to hang entries off
// until the row exists. So creation is a small inline form on the list that
// creates the row first, then opens the editor against a real ID.
export function useActiveAiSystem() {
  const route = useRoute()
  const router = useRouter()

  const activeAiSystemId = computed(() => route.query.aiSystem || null)

  function openAiSystem(aiSystemId) {
    router.replace({ query: { ...route.query, aiSystem: String(aiSystemId) } })
  }

  function closeAiSystem() {
    router.replace({ query: { ...route.query, aiSystem: undefined } })
  }

  return { activeAiSystemId, openAiSystem, closeAiSystem }
}

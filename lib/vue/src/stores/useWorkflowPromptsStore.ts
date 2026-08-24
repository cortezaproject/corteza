import { automation } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'
import { promptDefinitions } from '../components/prompts'

function onlyFresh(
  existing: Array<automation.Prompt>,
  fresh: Array<automation.Prompt>,
): Array<automation.Prompt> {
  const index = existing.map(({ stateID }) => stateID)
  return fresh.filter(({ stateID = undefined }) => stateID && !index.includes(stateID))
}

export const useWorkflowPromptsStore = defineStore('wfPrompts', () => {
  const $AutomationAPI = inject<any>('$AutomationAPI')

  const loading = ref(false)
  const prompts = ref<Array<automation.Prompt>>([])
  const active = ref<automation.Prompt | boolean>(false)

  // Sessions this screen started, the prompts held back for a webapp that can
  // render them, and the ones stepped over. The prompt list is the whole
  // user's, so a kind this webapp cannot render may still be meant for another
  // tab; only a session started here is ours to resume past.
  const ownedSessions = ref<Set<string>>(new Set())
  const foreign = ref<Array<automation.Prompt>>([])
  const skipped = ref<Array<automation.Prompt>>([])

  const all = computed(() => prompts.value)
  const isLoading = computed(() => loading.value)
  const isActive = computed(() => active.value !== false)
  const current = computed(() => (typeof active.value === 'boolean' ? undefined : active.value))

  function activate(prompt?: true | automation.Prompt) {
    active.value = prompt ?? true
  }

  function deactivate() {
    active.value = false
  }

  async function update(webapp: string) {
    const { set = [] } = await $AutomationAPI.sessionListPrompts()
    if (!Array.isArray(set) || set.length === 0) {
      clearAll()
      return
    }

    const mapped = set.map((prompt: automation.Prompt) => new automation.Prompt(prompt))
    const fresh = onlyFresh(prompts.value, mapped)
    if (fresh.length > 0) {
      appendPrompts(fresh, webapp)
    }
  }

  function newPrompt(prompt: automation.Prompt, webapp: string) {
    appendPrompts([new automation.Prompt(prompt)], webapp)
  }

  async function resume(prompt: automation.Prompt, input: automation.Vars) {
    loading.value = true
    try {
      await $AutomationAPI.sessionResumeState({
        sessionID: prompt.sessionID,
        stateID: prompt.stateID,
        input,
      })
    } finally {
      remove(prompt)
      loading.value = false
    }
  }

  async function cancel(prompt: automation.Prompt) {
    loading.value = true
    try {
      await $AutomationAPI.sessionCancel({
        sessionID: prompt.sessionID,
        stateID: prompt.stateID,
      })
    } finally {
      remove(prompt)
      loading.value = false
    }
  }

  function clear(prompt: automation.Prompt) {
    remove(prompt)
  }

  function clearAll() {
    prompts.value = []
    active.value = false
  }

  function appendPrompts(next: Array<automation.Prompt>, webapp: string) {
    const allowed: Array<automation.Prompt> = []
    const elsewhere: Array<automation.Prompt> = []

    for (const prompt of next) {
      const definition = promptDefinitions.find(({ ref }) => ref === prompt.ref)
      if (!definition) continue

      if (!definition.meta.webapps || definition.meta.webapps.includes(webapp)) {
        allowed.push(prompt)
      } else {
        elsewhere.push(prompt)
      }
    }

    prompts.value.push(...allowed)

    // Held rather than dropped: the session's own screen may only claim it a
    // moment later. The prompt push beats the exec response that names the
    // session, so deciding this once, on arrival, would always decide it too
    // early.
    foreign.value.push(...elsewhere)
    drainForeign()
  }

  function drainForeign() {
    const mine = foreign.value.filter(p => ownedSessions.value.has(p.sessionID))
    if (!mine.length) return

    foreign.value = foreign.value.filter(p => !ownedSessions.value.has(p.sessionID))
    mine.forEach(stepOver)
  }

  // A prompt whose kind belongs to another webapp still holds its session open.
  // Resuming without running the handler lets the workflow carry on; the step's
  // effect is simply not something this screen can perform.
  async function stepOver(prompt: automation.Prompt) {
    skipped.value.push(prompt)
    try {
      await $AutomationAPI.sessionResumeState({
        sessionID: prompt.sessionID,
        stateID: prompt.stateID,
        input: {},
      })
    } catch {
      // The session may already be gone; the notice stands either way.
    }
  }

  function ownSession(sessionID: string) {
    ownedSessions.value.add(sessionID)
    drainForeign()
  }

  function disownSession(sessionID: string) {
    ownedSessions.value.delete(sessionID)
    foreign.value = foreign.value.filter(p => p.sessionID !== sessionID)
    skipped.value = skipped.value.filter(p => p.sessionID !== sessionID)
  }

  function remove(prompt: automation.Prompt) {
    prompts.value = prompts.value.filter(({ stateID }) => stateID !== prompt.stateID)
    if (typeof active.value === 'object' && active.value.stateID === prompt.stateID) {
      active.value = prompts.value.length > 0
    }
  }

  return {
    loading,
    prompts,
    active,
    all,
    isLoading,
    isActive,
    current,
    activate,
    deactivate,
    update,
    newPrompt,
    resume,
    cancel,
    clear,
    clearAll,
    remove,
    skipped,
    ownSession,
    disownSession,
  }
})

import { automation } from '@cortezaproject/corteza-js-next'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { promptDefinitions } from '../components/prompts'

type AutomationAPI = {
  sessionListPrompts: () => Promise<any>;
  sessionResumeState: (args: Record<string, unknown>) => Promise<unknown>;
  sessionCancel: (args: Record<string, unknown>) => Promise<unknown>;
}

function onlyFresh(existing: Array<automation.Prompt>, fresh: Array<automation.Prompt>): Array<automation.Prompt> {
  const index = existing.map(({ stateID }) => stateID)
  return fresh.filter(({ stateID = undefined }) => stateID && !index.includes(stateID))
}

export const useWorkflowPromptsStore = defineStore('wfPrompts', () => {
  const loading = ref(false)
  const prompts = ref<Array<automation.Prompt>>([])
  const active = ref<automation.Prompt | boolean>(false)

  const all = computed(() => prompts.value)
  const isLoading = computed(() => loading.value)
  const isActive = computed(() => active.value !== false)
  const current = computed(() => typeof active.value === 'boolean' ? undefined : active.value)

  function activate(prompt?: true | automation.Prompt) {
    active.value = prompt ?? true
  }

  function deactivate() {
    active.value = false
  }

  async function update(api: AutomationAPI, webapp: string) {
    const { set = [] } = await api.sessionListPrompts()
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

  async function resume(api: AutomationAPI, prompt: automation.Prompt, input: automation.Vars) {
    loading.value = true
    try {
      await api.sessionResumeState({
        sessionID: prompt.sessionID,
        stateID: prompt.stateID,
        input,
      })
    } finally {
      remove(prompt)
      loading.value = false
    }
  }

  async function cancel(api: AutomationAPI, prompt: automation.Prompt) {
    loading.value = true
    try {
      await api.sessionCancel({
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
    const allowed = next.filter(({ ref }) => {
      return promptDefinitions.some(definition => {
        return definition.ref === ref && (!definition.meta.webapps || definition.meta.webapps.includes(webapp))
      })
    })

    prompts.value.push(...allowed)
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
  }
})

import { apiClients, automation } from '@cortezaproject/corteza-js'
import { ActionContext, StoreOptions } from 'vuex'
import { promptDefinitions } from '../../components/prompts'

interface Options {
  api: apiClients.Automation;
  ws: WebSocket;
  watchInterval: number;
  webapp: string;
}

interface State {
  loading: boolean;
  prompts: Array<automation.Prompt>;

  /**
   * Is prompt component active (modal open)?
   *   prompt = modal is open, show this/current prompt
   *   true   = modal is open, show list of pending prompts
   *   false  = modal is closed
   */
  active: automation.Prompt | boolean;
}

interface ResumePayload {
  prompt: automation.Prompt;
  input: automation.Vars;
}

async function loadPrompts (api: apiClients.Automation): Promise<Array<automation.Prompt>> {
  return api.sessionListPrompts().then(({ set } = {}) => {
    if (!Array.isArray(set)) {
      return []
    }

    return set.map((p: automation.Prompt) => new automation.Prompt(p))
  })
}

async function resumeState (api: apiClients.Automation, { sessionID, stateID }: automation.Prompt, input: automation.Vars): Promise<unknown> {
  return api.sessionResumeState({ sessionID, stateID, input })
}

async function cancelState (api: apiClients.Automation, { sessionID, stateID }: automation.Prompt): Promise<unknown> {
  return api.sessionCancel({ sessionID, stateID })
}

function onlyFresh (existing: Array<automation.Prompt>, fresh: Array<automation.Prompt>): Array<automation.Prompt> {
  const index = existing.map(({ stateID }) => stateID)
  return fresh.filter(({ stateID = undefined }) => stateID && !index.includes(stateID))
}

export default function ({ api, webapp }: Options): StoreOptions<State> {
  return {
    strict: true,

    state: {
      loading: false,
      active: false,
      prompts: [],
    },

    getters: {
      all ({ prompts }: State): Array<automation.Prompt> {
        return prompts
      },

      isLoading ({ loading }: State): boolean {
        return loading
      },

      isActive ({ active }: State): boolean {
        return active !== false
      },

      current ({ active }: State): automation.Prompt | undefined {
        if (typeof active === 'boolean') {
          return undefined
        } else {
          return active
        }
      },
    },

    actions: {
      async activate ({ commit }: ActionContext<State, State>, m?: true | automation.Prompt): Promise<void> {
        commit('active', m ?? true)
      },

      async deactivate ({ commit }: ActionContext<State, State>): Promise<void> {
        commit('active', false)
      },

      // fetch calls automation API endpoint and collects all states
      async update ({ commit, state }: ActionContext<State, State>): Promise<void> {
        return loadPrompts(api)
          .then(pp => {
            if (pp.length === 0) {
              // purge all
              commit('clear')
              return
            }

            // filter fresh prompts before committing
            // to minimize store traffic & state changes
            const fresh = onlyFresh(state.prompts, pp)
            if (fresh.length > 0) {
              commit('update', fresh)
            }
          })
      },

      new ({ commit }: ActionContext<State, State>, prompt: automation.Prompt): void {
        commit('update', [prompt])
      },

      async resume ({ commit }: ActionContext<State, State>, { prompt, input }: ResumePayload): Promise<void> {
        commit('loading')
        return resumeState(api, prompt, input)
          .then(() => commit('remove', prompt))
          .catch(() => commit('remove', prompt))
          .finally(() => commit('loading', false))
      },

      /**
       * Used when prompt is handled in this session and we need to
       * send the cancellation state back to the API
       */
      async cancel ({ commit }: ActionContext<State, State>, prompt: automation.Prompt): Promise<void> {
        commit('loading')
        return cancelState(api, prompt)
          .then(() => commit('remove', prompt))
          .finally(() => commit('loading', false))
      },

      /**
       * Used when prompt is handled in another session
       */
      clear ({ commit }: ActionContext<State, State>, prompt: automation.Prompt): void {
        commit('remove', prompt)
      },
    },

    mutations: {
      loading (state: State, flag = true): void {
        state.loading = flag
      },

      active (state: State, m: boolean | automation.Prompt): void {
        state.active = m
      },

      clear (state: State): void {
        state.prompts = []
      },

      update (state: State, pp: Array<automation.Prompt>): void {
        // // Check if prompt should be run only in certain webapps
        state.prompts.push(...pp.filter(({ ref }) => {
          return promptDefinitions.some(p => {
            return p.ref === ref && (!p.meta.webapps || p.meta.webapps.includes(webapp))
          })
        }))
      },

      remove (state: State, prompt: automation.Prompt): void {
        state.prompts = state.prompts.filter(({ stateID }) => stateID !== prompt.stateID)

        if (typeof state.active === 'object' && state.active.stateID === prompt.stateID) {
          // When removed prompt is also active,
          // reset active value.
          //
          // Set it to true but only if there are more prompts pending
          state.active = state.prompts.length > 0
        }
      },
    },
  }
}

import { ActionContext, StoreOptions } from 'vuex'

// all resources are namespaced/prefixed
// we're interested only in corteza resources
const resourcePrefix = 'corteza::'

const types = {
  UPDATE: 'UPDATE',
  UNLOAD: 'UNLOAD',
}

interface Rule {
  resource: string;
  operation: string;
  allow: boolean;
}

interface State {
  loaded: boolean;
  rules: Array<Rule>;
}

interface Fetcher {
  permissionsEffective: () => Promise<Array<Rule>>;
}

// This store serves as effective permission loader (read-only)
export default function (...apis: Array<Fetcher>): StoreOptions<State> {
  return {
    strict: true,

    state: {
      loaded: false,
      rules: [],
    },

    getters: {
      can (state: State) {
        return (res: string, op: string): boolean => {
          return (state.rules.find(({ resource, operation }) => resource === resourcePrefix + res && op === operation) || { allow: false }).allow
        }
      },

      isLoaded (state: State): boolean {
        return state.loaded
      },
    },

    actions: {
      // The apis must be valid and enabled, this is the responsibility of the callee
      load ({ commit }: ActionContext<State, State>, enabledAPIs: Array<Fetcher> = apis || []): void {
        commit(types.UNLOAD)
        Promise.all(enabledAPIs.map(api => {
          return api.permissionsEffective().catch(() => [])
        })).then((rr: Array<Array<Rule>>) => {
          commit(types.UPDATE, ([] as Array<Rule>)
            .concat(...rr)
            .filter(({ resource }) => resource.startsWith(resourcePrefix)))
        })
      },
    },

    mutations: {
      [types.UNLOAD] (state: State): void {
        state.loaded = false
      },

      [types.UPDATE] (state: State, rules: Array<Rule>): void {
        state.rules = rules
        state.loaded = true
      },
    },
  }
}

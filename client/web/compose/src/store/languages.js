const types = {
  pending: 'pending',
  completed: 'completed',
  updateSet: 'updateSet',
}

export default function (SystemAPI) {
  return {
    namespaced: true,

    state: {
      pending: false,
      set: [],
    },

    getters: {
      pending: (state) => state.pending,
      default: (state) => state.set.length > 0 ? state.set[0] : undefined,
      set: (state) => state.set,
    },

    actions: {
      async load ({ commit }) {
        commit(types.pending)
        return SystemAPI.localeList().then(({ set }) => {
          commit(types.updateSet, set)
        }).finally(() => {
          commit(types.completed)
        })
      },
    },

    mutations: {
      [types.pending] (state) {
        state.pending = true
      },

      [types.completed] (state) {
        state.pending = false
      },

      [types.updateSet] (state, set) {
        state.set = set
      },
    },
  }
}

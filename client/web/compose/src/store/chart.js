import { compose } from '@cortezaproject/corteza-js'

const types = {
  loading: 'loading',
  loaded: 'loaded',
  pending: 'pending',
  completed: 'completed',
  updateSet: 'updateSet',
  removeFromSet: 'removeFromSet',
  clearSet: 'clearSet',
}

export default function (ComposeAPI) {
  return {
    namespaced: true,

    state: {
      loading: false,
      pending: false,
      set: [],
    },

    getters: {
      loading: (state) => state.loading,

      pending: (state) => state.pending,

      getByID (state) {
        return (ID) => state.set.find(({ chartID }) => ID === chartID)
      },

      set (state) {
        return state.set
      },
    },

    actions: {
      async load ({ commit, getters }, { namespaceID, clear = false, force = false } = {}) {
        if (clear) {
          commit(types.clearSet)
        }

        if (!force && getters.set.length > 1) {
          // When there's forced load, make sure we have more than 1 item in the set
          // in the scenario when user came to detail page first and has one item loaded
          // > 0 would not be sufficient.
          return new Promise((resolve) => resolve(getters.set))
        }

        commit(types.loading)
        commit(types.pending)
        return ComposeAPI.chartList({ namespaceID, sort: 'name ASC' }).then(({ set, filter }) => {
          if (set && set.length > 0) {
            commit(types.updateSet, set.map(c => new compose.Chart(c)))
          }

          return getters.set
        }).finally(() => {
          commit(types.loaded)
          commit(types.completed)
        })
      },

      async findByID ({ commit, getters }, { namespaceID, chartID, force = false } = {}) {
        if (!force) {
          const oldItem = getters.getByID(chartID)
          if (oldItem) {
            return new Promise((resolve) => resolve(oldItem))
          }
        }

        commit(types.pending)
        return ComposeAPI.chartRead({ namespaceID, chartID }).then(raw => {
          const chart = new compose.Chart(raw)
          commit(types.updateSet, [chart])
          return chart
        }).finally(() => {
          commit(types.completed)
        })
      },

      async create ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.chartCreate(item).then(raw => {
          const chart = new compose.Chart(raw)
          commit(types.updateSet, [chart])
          return chart
        }).finally(() => {
          commit(types.completed)
        })
      },

      async update ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.chartUpdate(item).then(raw => {
          const chart = new compose.Chart(raw)
          commit(types.updateSet, [chart])
          return chart
        }).finally(() => {
          commit(types.completed)
        })
      },

      async delete ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.chartDelete(item).then(() => {
          commit(types.removeFromSet, [item])
          return true
        }).finally(() => {
          commit(types.completed)
        })
      },

      clearSet ({ commit }) {
        commit(types.clearSet)
      },
    },

    mutations: {
      [types.loading] (state) {
        state.loading = true
      },

      [types.loaded] (state) {
        state.loading = false
      },

      [types.pending] (state) {
        state.pending = true
      },

      [types.completed] (state) {
        state.pending = false
      },

      [types.updateSet] (state, set) {
        set = set.map(i => Object.freeze(i))

        if (state.set.length === 0) {
          state.set = set
          return
        }

        set.forEach(newItem => {
          const oldIndex = state.set.findIndex(({ chartID }) => chartID === newItem.chartID)
          if (oldIndex > -1) {
            state.set.splice(oldIndex, 1, newItem)
          } else {
            state.set.push(newItem)
          }
        })
      },

      [types.removeFromSet] (state, removedSet) {
        (removedSet || []).forEach(removedItem => {
          const i = state.set.findIndex(({ chartID }) => chartID === removedItem.chartID)
          if (i > -1) {
            state.set.splice(i, 1)
          }
        })
      },

      [types.clearSet] (state) {
        state.pending = false
        state.set.splice(0)
      },
    },
  }
}

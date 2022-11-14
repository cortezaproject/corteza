const state = {
  set: [],
}

const getters = {
  set (state) {
    return state.set
  },

  unifyOnly (state) {
    return state.set.filter(({ unify: { listed } = { listed: false } }) => listed)
  },
}

const mutations = {
  updateSet (state, set) {
    state.set = set
  },
}

/**
 * @param localStorage
 * @param api
 */
export default ({ api }) => {
  return {
    namespaced: true,

    state,
    getters,
    mutations,

    actions: {
      async load ({ commit }) {
        if (api && api.applicationList) {
          return api.applicationList({ sort: 'weight', incFlags: 0 })
            .then(({ set }) => commit('updateSet', set))
        }
      },

      async reorder ({ dispatch }, applicationIDs,) {
        return api.applicationReorder({ applicationIDs }).then(() => {
          return dispatch('load')
        })
      },

      async pin ({ dispatch }, { applicationID, ownedBy }) {
        return api.applicationFlagCreate({ applicationID, flag: 'pinned', ownedBy }).then(() => {
          return dispatch('load')
        }).catch(() => {})
      },

      async unpin ({ dispatch }, { applicationID, ownedBy }) {
        return api.applicationFlagDelete({ applicationID, flag: 'pinned', ownedBy }).then(() => {
          return dispatch('load')
        }).catch(() => {})
      },
    },
  }
}

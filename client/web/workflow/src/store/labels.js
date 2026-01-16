import Vue from 'vue'

export default {
  namespaced: true,

  state: {
    namespaces: {},
    modules: {},
  },

  getters: {
    getNamespace: (state) => (namespaceID) => {
      return state.namespaces[namespaceID]
    },

    getModule: (state) => (moduleID) => {
      return state.modules[moduleID]
    },
  },

  mutations: {
    setNamespace (state, { namespaceID, name }) {
      Vue.set(state.namespaces, namespaceID, name)
    },

    setModule (state, { moduleID, name }) {
      Vue.set(state.modules, moduleID, name)
    },
  },

  actions: {
    async resolveNamespace ({ commit, state }, { namespaceID, api }) {
      // Return cached value if available
      if (state.namespaces[namespaceID] !== undefined) {
        return state.namespaces[namespaceID]
      }

      // Mark as loading to prevent duplicate requests
      commit('setNamespace', { namespaceID, name: null })

      try {
        const namespace = await api.namespaceRead({ namespaceID })
        const name = namespace.name || namespace.slug || namespaceID
        commit('setNamespace', { namespaceID, name })
        return name
      } catch (error) {
        commit('setNamespace', { namespaceID, name: namespaceID })
        return namespaceID
      }
    },

    async resolveModule ({ commit, state }, { moduleID, namespaceID, api }) {
      // Return cached value if available
      if (state.modules[moduleID] !== undefined) {
        return state.modules[moduleID]
      }

      // Mark as loading to prevent duplicate requests
      commit('setModule', { moduleID, name: null })

      try {
        const module = await api.moduleRead({ namespaceID, moduleID })
        const name = module.name || module.handle || moduleID
        commit('setModule', { moduleID, name })
        return name
      } catch (error) {
        commit('setModule', { moduleID, name: moduleID })
        return moduleID
      }
    },

    async resolveMultipleNamespaces ({ dispatch }, { namespaceIDs, api }) {
      return Promise.all(
        namespaceIDs.map(namespaceID =>
          dispatch('resolveNamespace', { namespaceID, api }),
        ),
      )
    },

    async resolveMultipleModules ({ dispatch }, { modules, api }) {
      return Promise.all(
        modules.map(({ moduleID, namespaceID }) =>
          dispatch('resolveModule', { moduleID, namespaceID, api }),
        ),
      )
    },
  },
}

import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    processing: false,
    types: [],
    aggregations: [],
    modules: [],
    namespaces: [],
  },
  mutations: {
    updateProcessing (state, value = false) {
      state.processing = value
    },

    updateTypes (state, types = []) {
      state.types = types
    },

    updateAggregations (state, aggs = []) {
      state.aggregations = aggs
    },

    updateModules (state, modules = []) {
      state.modules = modules
    },

    updateNamespaces (state, namespaces = []) {
      state.namespaces = namespaces
    },
  },
  actions: {
  },
  modules: {
  },
})

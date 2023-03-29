import { compose } from '@cortezaproject/corteza-js'
import * as request from '../lib/request'

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
        return (ID) => state.set.find(({ pageLayoutID }) => ID === pageLayoutID)
      },

      getByHandle (state) {
        return (handle) => state.set.find((pl) => handle === pl.handle)
      },

      getByPageID (state) {
        return (ID) => state.set.filter(({ pageID }) => ID === pageID).sort((a, b) => a.weight - b.weight)
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
          return new Promise((resolve) => resolve(getters.set))
        }

        commit(types.loading)
        commit(types.pending)
        return ComposeAPI.pageLayoutListNamespace({ namespaceID, sort: 'weight ASC' }).then(({ set, filter }) => {
          if (set && set.length > 0) {
            commit(types.updateSet, set.map(pl => new compose.PageLayout(pl)))
          }

          return getters.set
        }).finally(() => {
          commit(types.loaded)
          commit(types.completed)
        })
      },

      async findByID ({ commit, getters }, { namespaceID, pageID, pageLayoutID, force = false } = {}) {
        if (!force) {
          const oldItem = getters.getByID(pageLayoutID)
          return new Promise((resolve) => resolve(oldItem))
        }

        commit(types.pending)
        return ComposeAPI.pageLayoutRead({ namespaceID, pageID, pageLayoutID }).then(pl => {
          const pageLayout = new compose.PageLayout(pl)

          commit(types.updateSet, [pageLayout])
          return pageLayout
        }).finally(() => {
          commit(types.completed)
        })
      },

      async findByPageID ({ commit, getters }, { namespaceID, pageID, force = false } = {}) {
        if (!force) {
          const oldItems = getters.getByPageID(pageID)
          return new Promise((resolve) => resolve(oldItems))
        }

        commit(types.pending)
        return ComposeAPI.pageLayoutList({ namespaceID, pageID, sort: 'weight ASC' }).then(({ set }) => {
          commit(types.updateSet, set.map(pl => new compose.PageLayout(pl)))
          return set
        }).finally(() => {
          commit(types.completed)
        })
      },

      async create ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.pageLayoutCreate(item, request.config(item)).then(pl => {
          const pageLayout = new compose.PageLayout(pl)
          commit(types.updateSet, [pageLayout])
          return pageLayout
        }).finally(() => {
          commit(types.completed)
        })
      },

      async update ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.pageLayoutUpdate(item, request.config(item)).then(pl => {
          const pageLayout = new compose.PageLayout(pl)
          commit(types.updateSet, [pageLayout])
          return pageLayout
        }).finally(() => {
          commit(types.completed)
        })
      },

      async delete ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.pageLayoutDelete(item).then(() => {
          commit(types.removeFromSet, [item])
          return true
        }).finally(() => {
          commit(types.completed)
        })
      },

      updateSet ({ commit }, pageLayout) {
        commit(types.updateSet, [pageLayout])
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
          const oldIndex = state.set.findIndex(({ pageLayoutID }) => pageLayoutID === newItem.pageLayoutID)
          if (oldIndex > -1) {
            state.set.splice(oldIndex, 1, newItem)
          } else {
            state.set.push(newItem)
          }
        })
      },

      [types.removeFromSet] (state, removedSet) {
        (removedSet || []).forEach(removedItem => {
          const i = state.set.findIndex(({ pageLayoutID }) => pageLayoutID === removedItem.pageLayoutID)
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

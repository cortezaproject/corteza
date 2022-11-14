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
        return (ID) => state.set.find(({ pageID }) => ID === pageID)
      },

      getByHandle (state) {
        return (handle) => state.set.find((p) => handle === p.handle)
      },

      homePage: (state) => state.set.find(p => p.visible && p.firstLevel && !p.isRecordPage),

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
        return ComposeAPI.pageList({ namespaceID, sort: 'weight ASC' }).then(({ set, filter }) => {
          if (set && set.length > 0) {
            commit(types.updateSet, set.map(p => new compose.Page(p)))
          }

          return getters.set
        }).finally(() => {
          commit(types.loaded)
          commit(types.completed)
        })
      },

      async findByID ({ commit, getters }, { namespaceID, pageID, force = false } = {}) {
        if (!force) {
          const oldItem = getters.getByID(pageID)
          if (oldItem) {
            return new Promise((resolve) => resolve(oldItem))
          }
        }

        commit(types.pending)
        return ComposeAPI.pageRead({ namespaceID, pageID }).then(raw => {
          const page = new compose.Page(raw)
          commit(types.updateSet, [page])
          return page
        }).finally(() => {
          commit(types.completed)
        })
      },

      async create ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.pageCreate(item, request.config(item)).then(raw => {
          const page = new compose.Page(raw)
          commit(types.updateSet, [page])
          return page
        }).finally(() => {
          commit(types.completed)
        })
      },

      async update ({ commit }, item) {
        commit(types.pending)
        return ComposeAPI.pageUpdate(item, request.config(item)).then(raw => {
          const page = new compose.Page(raw)
          commit(types.updateSet, [page])
          return page
        }).finally(() => {
          commit(types.completed)
        })
      },

      async delete ({ commit, dispatch }, item) {
        commit(types.pending)
        return ComposeAPI.pageDelete(item).then(() => {
          commit(types.removeFromSet, [item])
          const { namespaceID } = item || {}
          if (namespaceID) {
            dispatch('load', { namespaceID: item.namespaceID, clear: true })
          }
          return true
        }).finally(() => {
          commit(types.completed)
        })
      },

      updateSet ({ commit }, page) {
        commit(types.updateSet, [page])
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
          const oldIndex = state.set.findIndex(({ pageID }) => pageID === newItem.pageID)
          if (oldIndex > -1) {
            state.set.splice(oldIndex, 1, newItem)
          } else {
            state.set.push(newItem)
          }
        })
      },

      [types.removeFromSet] (state, removedSet) {
        (removedSet || []).forEach(removedItem => {
          const i = state.set.findIndex(({ pageID }) => pageID === removedItem.pageID)
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

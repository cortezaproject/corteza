import { system } from '@cortezaproject/corteza-js'

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

      findByID (state) {
        return (ID) => state.set.find(({ userID }) => ID === userID)
      },

      findByUsername: (state) => (username) => {
        return state.set.filter(user => user.username === username)[0] || undefined
      },

      set (state) {
        return state.set
      },
    },

    actions: {
      async load ({ commit }, filter) {
        commit(types.pending)
        return SystemAPI.userList(filter).then(({ set }) => {
          commit(types.updateSet, set)
        }).finally(() => {
          commit(types.completed)
        })
      },

      push ({ commit }, user) {
        commit(types.updateSet, user)
      },

      async fetchUsers ({ commit }, userID) {
        commit(types.pending)

        if (userID.length === 0) {
          return null
        }

        return SystemAPI.userList({ userID }).then(({ set }) => {
          commit(types.updateSet, set)
        }).finally(() => {
          commit(types.completed)
        })
      },

      /**
       * Similar to fetchUsers but it only fetches unknown (not in set) ids
       */
      async resolveUsers ({ commit, getters }, list) {
        if (list.length === 0) {
          // save ourselves some work
          return
        }

        // exclude existing & make unique
        const existing = new Set(getters.set.map(({ userID }) => userID))
        list = [...new Set(list.filter(userID => userID && !existing.has(userID)))]

        if (list.length === 0) {
          // Check for values again
          return
        }

        return SystemAPI.userList({ userID: list }).then(({ set }) => {
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
        set = (Array.isArray(set) ? set : [set]).filter(u => !!u).map(i => new system.User(i))

        if (state.set.length === 0) {
          state.set = set
          return
        }

        set.forEach(newItem => {
          const oldIndex = state.set.findIndex(({ userID }) => userID === newItem.userID)
          if (oldIndex > -1) {
            state.set.splice(oldIndex, 1, newItem)
          } else {
            state.set.push(newItem)
          }
        })
      },
    },
  }
}

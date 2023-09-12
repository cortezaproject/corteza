// Records in this store not classes, but raw objects instead
const types = {
  pending: 'pending',
  completed: 'completed',
  updateSet: 'updateSet',
  clearSet: 'clearSet',
}

export default function (ComposeAPI) {
  return {
    namespaced: true,

    state: {
      pending: false,
      set: [],
    },

    getters: {
      pending: (state) => state.pending,

      findByID (state) {
        return (ID) => state.set.find(({ recordID }) => ID === recordID)
      },

      findByIDs (state) {
        return (IDs) => {
          IDs = IDs.flat()
          return state.set.filter(({ recordID }) => IDs.includes(recordID))
        }
      },

      set (state) {
        return state.set
      },
    },

    actions: {
      /**
       * Similar to fetchRecords but it only fetches unknown (not in set) ids
       */
      async resolveRecords ({ commit }, { namespaceID, moduleID, recordIDs }) {
        if (recordIDs.length === 0) {
          // save ourselves some work
          return
        }

        if (recordIDs.length === 0) {
          // Check for values again
          return
        }

        const query = recordIDs.map(recordID => `recordID = ${recordID}`).join(' OR ')

        return ComposeAPI.recordList({ namespaceID, moduleID, query, deleted: 1 }).then(({ set }) => {
          commit(types.updateSet, set)
        }).finally(() => {
          commit(types.completed)
          recordIDs = []
        })
      },

      updateRecords ({ commit }, records) {
        commit(types.updateSet, records)
      },

      push ({ commit }, record) {
        commit(types.updateSet, record)
      },

      clearSet ({ commit }) {
        commit(types.clearSet)
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
        set = (Array.isArray(set) ? set : [set]).filter(r => !!r)

        if (state.set.length === 0) {
          state.set = set
          return
        }

        const existing = new Set(state.set.map(({ recordID }) => recordID))

        set.forEach(newItem => {
          const oldIndex = existing.has(newItem.recordID) ? state.set.findIndex(({ recordID }) => recordID === newItem.recordID) : -1
          if (oldIndex > -1) {
            state.set.splice(oldIndex, 1, newItem)
          } else {
            state.set.push(newItem)
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

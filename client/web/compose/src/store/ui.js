const types = {
  loading: 'loading',
  loaded: 'loaded',
  pending: 'pending',
  completed: 'completed',
  setRecordPagination: 'setRecordPagination',
  clearRecordPagination: 'clearRecordPagination',
  recordPaginationUsable: 'recordPaginationUsable',
  setRecordPaginationUsable: 'setRecordPaginationUsable',
}

export default function (ComposeAPI) {
  return {
    namespaced: true,

    state: {
      loading: false,
      pending: false,
      recordPaginationIDs: [],
      recordPaginationUsable: false,
    },

    getters: {
      loading: (state) => state.loading,

      pending: (state) => state.pending,

      recordPaginationUsable: (state) => state.recordPaginationUsable,

      getNextAndPrevRecord: ({ recordPaginationIDs }) => (recordID) => {
        const recordIndex = recordPaginationIDs.indexOf(recordID)
        const prev = recordIndex >= 0 ? recordPaginationIDs[recordIndex - 1] : undefined
        const next = recordIndex >= 0 ? recordPaginationIDs[recordIndex + 1] : undefined

        return { next, prev }
      },
    },

    actions: {
      async loadPaginationRecords ({ commit }, { filter } = {}) {
        commit(types.pending)
        commit(types.recordPaginationUsable, true)

        const { pageCursor, prevPage } = filter

        return Promise.all([prevPage, pageCursor].map(pageCursor => {
          return ComposeAPI.recordList({ ...filter, pageCursor })
            .then(({ set }) => {
              return set.map(({ recordID }) => recordID)
            })
        })).then(([prevRecords, nextRecords]) => {
          commit(types.setRecordPagination, [...new Set([...prevRecords, ...nextRecords])])
        }).finally(() => {
          commit(types.completed)
        })
      },

      clearRecordIDs ({ commit }) {
        commit(types.clearRecordPagination)
      },

      setRecordPaginationUsable ({ commit }, value) {
        commit(types.recordPaginationUsable, value)
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

      [types.setRecordPagination] (state, recordIDs) {
        state.recordPaginationIDs = recordIDs
      },

      [types.clearRecordPagination] (state) {
        state.recordPaginationIDs = []
      },

      [types.recordPaginationUsable] (state, value) {
        state.recordPaginationUsable = value
      },
    },
  }
}

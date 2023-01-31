const types = {
  loading: 'loading',
  loaded: 'loaded',
  pending: 'pending',
  completed: 'completed',
  setRecordPagination: 'setRecordPagination',
  clearRecordPagination: 'clearRecordPagination',
  clearRecordPageNavigation: 'clearRecordPageNavigation',
  setClearRecordPageNavigation: 'setClearRecordPageNavigation',
}

export default function (ComposeAPI) {
  return {
    namespaced: true,

    state: {
      loading: false,
      pending: false,
      recordPaginationIds: [],
      recordPageVisited: null,
      clearRecordPageNavigation: true,
    },

    getters: {
      loading: (state) => state.loading,

      pending: (state) => state.pending,

      clearRecordPageNavigation: (state) => state.clearRecordPageNavigation,

      getRecordNavigationIndex: (state) => (recordID) => {
        return state.recordPaginationIds.indexOf(recordID)
      },

      nextRecordNavigation: ({ recordPaginationIds }, { getRecordNavigationIndex }) => (recordID) => {
        const recordIndex = getRecordNavigationIndex(recordID)
        const index = recordIndex !== undefined ? recordIndex : 1

        return recordPaginationIds[index - 1]
      },

      prevRecordNavigation: ({ recordPaginationIds }, { getRecordNavigationIndex }) => (recordID) => {
        const recordIndex = getRecordNavigationIndex(recordID)
        const index = recordIndex !== undefined ? recordIndex : 1

        return recordPaginationIds[index + 1]
      },

      getNextAndPrevRecord: (_, { nextRecordNavigation, prevRecordNavigation }) => (recordID) => {
        return {
          next: prevRecordNavigation(recordID),
          prev: nextRecordNavigation(recordID),
        }
      },
    },

    actions: {
      loadPaginationRecords ({ commit }, { filter, moduleID, namespaceID, filterCursors, incTotal, incPageNavigation, options } = {}) {
        commit(types.pending)
        commit(types.setClearRecordPageNavigation, true)

        return Promise.all(filterCursors.map((cursor) => {
          filter.pageCursor = cursor

          return ComposeAPI.recordList({ ...filter, moduleID, namespaceID, incTotal, incPageNavigation })
            .then(({ set }) => {
              return set.map(({ recordID }) => recordID)
            })
        })).finally(() => {
          commit(types.completed)
        }).then(([prevRecords, nextRecords]) => {
          commit(types.setRecordPagination, [...new Set([...prevRecords, ...nextRecords])])
        })
      },

      clearRecordIds ({ commit }) {
        commit(types.clearRecordPagination)
      },

      setClearRecordPageNavigation ({ commit }, value) {
        commit(types.setClearRecordPageNavigation, value)
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

      [types.setRecordPagination] (state, recordIds) {
        state.recordPaginationIds = recordIds
      },

      [types.clearRecordPagination] (state) {
        state.recordPaginationIds = []
      },

      [types.setClearRecordPageNavigation] (state, value) {
        state.clearRecordPageNavigation = value
      },
    },
  }
}

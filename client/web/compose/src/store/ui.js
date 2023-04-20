const types = {
  loading: 'loading',
  loaded: 'loaded',
  pending: 'pending',
  completed: 'completed',
  setRecordPagination: 'setRecordPagination',
  clearRecordPagination: 'clearRecordPagination',
  recordPaginationUsable: 'recordPaginationUsable',
  setRecordPaginationUsable: 'setRecordPaginationUsable',
  previousPages: 'previousPages',
  setPreviousPages: 'setPreviousPages',
  pushPreviousPages: 'pushPreviousPages',
  popPreviousPages: 'popPreviousPages',
}

export default function (ComposeAPI) {
  return {
    namespaced: true,

    state: {
      loading: false,
      pending: false,
      recordPaginationIDs: [],
      recordPaginationUsable: false,

      previousPages: [],
    },

    getters: {
      loading: (state) => state.loading,

      pending: (state) => state.pending,

      recordPaginationUsable: (state) => state.recordPaginationUsable,

      previousPages: (state) => state.previousPages,

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

      setPreviousPages ({ commit }, value) {
        commit(types.setPreviousPages, value)
      },

      pushPreviousPages ({ commit }, value) {
        commit(types.pushPreviousPages, value)
      },

      popPreviousPages ({ commit, state }) {
        const previousPage = state.previousPages.slice(-1)[0]
        commit(types.popPreviousPages)
        return new Promise((resolve) => resolve(previousPage))
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

      [types.setPreviousPages] (state, value) {
        state.previousPages = value
      },

      [types.pushPreviousPages] (state, value) {
        state.previousPages.push(value)
      },

      [types.popPreviousPages] (state) {
        return state.previousPages.pop()
      },
    },
  }
}

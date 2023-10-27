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
  previousPage: 'previousPage',
  setPreviousPage: 'setPreviousPage',
  modalPreviousPages: 'modalPreviousPages',
  pushModalPreviousPage: 'pushModalPreviousPage',
  clearModalPreviousPage: 'clearModalPreviousPage',
  popModalPreviousPage: 'popModalPreviousPage',
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
      previousPage: null,
      modalPreviousPages: [],
    },

    getters: {
      loading: (state) => state.loading,

      pending: (state) => state.pending,

      recordPaginationUsable: (state) => state.recordPaginationUsable,

      previousPages: (state) => state.previousPages,

      previousPage: (state) => state.previousPage,

      getNextAndPrevRecord: ({ recordPaginationIDs }) => (recordID) => {
        const recordIndex = recordPaginationIDs.indexOf(recordID)
        const prev = recordIndex >= 0 ? recordPaginationIDs[recordIndex - 1] : undefined
        const next = recordIndex >= 0 ? recordPaginationIDs[recordIndex + 1] : undefined

        return { prev, next }
      },

      modalPreviousPages: (state) => state.modalPreviousPages,
    },

    actions: {
      async loadPaginationRecords ({ commit }, { filter } = {}) {
        commit(types.pending)
        commit(types.recordPaginationUsable, true)

        const { prevPage, pageCursor, nextPage } = filter

        const cursors = new Set([prevPage, pageCursor, nextPage])

        return Promise.all([...cursors].map(pageCursor => {
          return ComposeAPI.recordList({ ...filter, pageCursor })
            .then(({ set }) => {
              return set.map(({ recordID }) => recordID)
            })
        })).then(([...records]) => {
          commit(types.setRecordPagination, [...new Set(records.flatMap(r => r))])
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

      setPreviousPage ({ commit }, value) {
        // This prevents saving previous page for pages that causes incorrect redirection even though they are previous pages
        const shouldNotSavePage = value.name !== 'admin.pages.builder' &&
              !value.query.layoutID && value.name !== 'admin.modules.create' &&
                value.name !== 'admin.charts.create' &&
                  value.name !== 'namespace.create'

        if (value && value.name && shouldNotSavePage) {
          commit(types.setPreviousPage, value)
        }
      },

      pushModalPreviousPage ({ commit, state }, value) {
        const previousPage = state.modalPreviousPages[state.modalPreviousPages.length - 1] || {}
        if (previousPage.recordID === value.recordID && previousPage.recordPageID === value.recordPageID) {
          return
        }

        commit(types.pushModalPreviousPage, value)
      },

      clearModalPreviousPage ({ commit }) {
        commit(types.clearModalPreviousPage)
      },

      popModalPreviousPage ({ commit, state }) {
        const previousPage = state.modalPreviousPages[state.modalPreviousPages.length - 2] || {}
        commit(types.popModalPreviousPage)
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

      [types.setPreviousPage] (state, value) {
        state.previousPage = value
      },

      [types.pushModalPreviousPage] (state, value) {
        state.modalPreviousPages.push(value)
      },

      [types.clearModalPreviousPage] (state) {
        state.modalPreviousPages = []
      },

      [types.popModalPreviousPage] (state) {
        state.modalPreviousPages.pop()
      },
    },
  }
}

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
  setNamespaceSlug: 'setNamespaceSlug',
  setPageHandle: 'setPageHandle',
  setLayoutHandle: 'setLayoutHandle',
  setModalPageHandle: 'setModalPageHandle',
  setModalLayoutHandle: 'setModalLayoutHandle',
  setLayoutRequiredFields: 'setLayoutRequiredFields',
  clearLayoutRequiredFields: 'clearLayoutRequiredFields',
  setBlockUnsaved: 'setBlockUnsaved',
  clearUnsavedBlocks: 'clearUnsavedBlocks',
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

      namespaceSlug: '',
      pageHandle: '',
      layoutHandle: '',

      modalPageHandle: '',
      modalLayoutHandle: '',

      layoutRequiredFields: [],

      // IDs of page blocks holding changes that are not saved yet
      // (record lists in inline-edit mode); used to warn before leaving
      unsavedBlocks: [],
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

      namespaceSlug: (state) => state.namespaceSlug,
      pageHandle: (state) => state.pageHandle,
      layoutHandle: (state) => state.layoutHandle,

      modalPageHandle: (state) => state.modalPageHandle,
      modalLayoutHandle: (state) => state.modalLayoutHandle,

      layoutRequiredFields: (state) => state.layoutRequiredFields,
      isFieldRequiredByLayout: (state) => (fieldID) => state.layoutRequiredFields.includes(fieldID),
      hasUnsavedBlocks: (state) => state.unsavedBlocks.length > 0,
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

      clearRecordPagination ({ commit }) {
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

      setNamespaceSlug ({ commit }, value) {
        commit(types.setNamespaceSlug, value)
      },

      setPageHandle ({ commit }, value) {
        commit(types.setPageHandle, value)
      },

      setLayoutHandle ({ commit }, value) {
        commit(types.setLayoutHandle, value)
      },

      setModalPageHandle ({ commit }, value) {
        commit(types.setModalPageHandle, value)
      },

      setModalLayoutHandle ({ commit }, value) {
        commit(types.setModalLayoutHandle, value)
      },

      setLayoutRequiredFields ({ commit }, fields) {
        commit(types.setLayoutRequiredFields, fields)
      },

      clearLayoutRequiredFields ({ commit }) {
        commit(types.clearLayoutRequiredFields)
      },

      setBlockUnsaved ({ commit }, { id, unsaved }) {
        commit(types.setBlockUnsaved, { id, unsaved })
      },

      clearUnsavedBlocks ({ commit }) {
        commit(types.clearUnsavedBlocks)
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

      [types.setNamespaceSlug] (state, value) {
        state.namespaceSlug = value || ''
      },

      [types.setPageHandle] (state, value) {
        state.pageHandle = value || ''
      },

      [types.setLayoutHandle] (state, value) {
        state.layoutHandle = value || ''
      },

      [types.setModalPageHandle] (state, value) {
        state.modalPageHandle = value || ''
      },

      [types.setModalLayoutHandle] (state, value) {
        state.modalLayoutHandle = value || ''
      },

      [types.setLayoutRequiredFields] (state, fields) {
        state.layoutRequiredFields = fields || []
      },

      [types.clearLayoutRequiredFields] (state) {
        state.layoutRequiredFields = []
      },

      [types.setBlockUnsaved] (state, { id, unsaved }) {
        if (!id) return

        const others = state.unsavedBlocks.filter(b => b !== id)
        state.unsavedBlocks = unsaved ? [...others, id] : others
      },

      [types.clearUnsavedBlocks] (state) {
        state.unsavedBlocks = []
      },
    },
  }
}

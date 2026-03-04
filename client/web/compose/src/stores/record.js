import { compose } from '@cortezaproject/corteza-js-next'
import { defineStore } from 'pinia'
import { inject, reactive, toRef } from 'vue'
import { useModuleStore } from './module'

export const useRecordStore = defineStore('record', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    // Map of records keyed by recordID
    records: new Map(),
    // Lightweight cache for raw API records used by field viewers (no module required)
    labelCache: new Map(),
  })

  /**
   * Fetch a list of records for a given module.
   *
   * @param {Object} params
   * @param {string} params.namespaceID
   * @param {string} params.moduleID
   * @param {string} [params.query] - Search query
   * @param {string} [params.sort] - Sort expression
   * @param {number} [params.limit] - Page size
   * @param {string} [params.pageCursor] - Cursor for pagination
   * @param {boolean} [params.incTotal] - Include total count
   * @returns {Promise<{ set: compose.Record[], filter: Object }>}
   */
  async function list({
    namespaceID,
    moduleID,
    query,
    sort,
    limit,
    pageCursor,
    incTotal = true,
    deleted,
  } = {}) {
    if (!$ComposeAPI) {
      throw new Error('ComposeAPI not available via inject')
    }

    if (!namespaceID || !moduleID) {
      throw new Error('namespaceID and moduleID are required')
    }

    const moduleStore = useModuleStore()
    const mod = moduleStore.getByID(moduleID)
    if (!mod) {
      throw new Error(`Module ${moduleID} not found in store`)
    }

    state.loading = true

    try {
      const response = await $ComposeAPI.recordList({
        namespaceID,
        moduleID,
        query,
        sort,
        limit,
        pageCursor,
        incTotal,
        deleted,
      })

      const { set = [], filter = {} } = response

      const records = set.map(r => {
        const record = new compose.Record(mod, r)
        // Cache each record
        state.records.set(record.recordID, record)
        return record
      })

      return { set: records, filter }
    } catch (error) {
      console.error('Failed to list records:', error)
      throw error
    } finally {
      state.loading = false
    }
  }

  /**
   * Fetch a single record by ID.
   *
   * @param {Object} params
   * @param {string} params.namespaceID
   * @param {string} params.moduleID
   * @param {string} params.recordID
   * @param {boolean} [params.force=false] - Force fetch even if cached
   * @returns {Promise<compose.Record>}
   */
  async function findByID({ namespaceID, moduleID, recordID, force = false } = {}) {
    if (!force) {
      const cached = state.records.get(recordID)
      if (cached) return cached
    }

    if (!$ComposeAPI) {
      throw new Error('ComposeAPI not available via inject')
    }

    if (!namespaceID || !moduleID || !recordID) {
      throw new Error('namespaceID, moduleID, and recordID are required')
    }

    const moduleStore = useModuleStore()
    const mod = moduleStore.getByID(moduleID)
    if (!mod) {
      throw new Error(`Module ${moduleID} not found in store`)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.recordRead({ namespaceID, moduleID, recordID })
      const record = new compose.Record(mod, raw)
      state.records.set(record.recordID, record)
      return record
    } catch (error) {
      console.error('Failed to read record:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  /**
   * Create a new record.
   *
   * @param {compose.Record} record - Record instance to create
   * @returns {Promise<compose.Record>}
   */
  async function create(record) {
    if (!$ComposeAPI) {
      throw new Error('ComposeAPI not available via inject')
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.recordCreate({
        namespaceID: record.namespaceID,
        moduleID: record.moduleID,
        values: record.serializeValues(),
        ownedBy: record.ownedBy,
        meta: record.meta,
      })

      const mod = useModuleStore().getByID(record.moduleID)
      const created = new compose.Record(mod, raw)
      state.records.set(created.recordID, created)
      return created
    } catch (error) {
      console.error('Failed to create record:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  /**
   * Update an existing record.
   *
   * @param {compose.Record} record - Record instance to update
   * @returns {Promise<compose.Record>}
   */
  async function update(record) {
    if (!$ComposeAPI) {
      throw new Error('ComposeAPI not available via inject')
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.recordUpdate({
        namespaceID: record.namespaceID,
        moduleID: record.moduleID,
        recordID: record.recordID,
        values: record.serializeValues(),
        ownedBy: record.ownedBy,
        meta: record.meta,
      })

      const mod = useModuleStore().getByID(record.moduleID)
      const updated = new compose.Record(mod, raw)
      state.records.set(updated.recordID, updated)
      return updated
    } catch (error) {
      console.error('Failed to update record:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  /**
   * Delete a record.
   *
   * @param {Object} params
   * @param {string} params.namespaceID
   * @param {string} params.moduleID
   * @param {string} params.recordID
   * @returns {Promise<boolean>}
   */
  async function deleteRecord({ namespaceID, moduleID, recordID } = {}) {
    if (!$ComposeAPI) {
      throw new Error('ComposeAPI not available via inject')
    }

    state.pending = true

    try {
      await $ComposeAPI.recordDelete({ namespaceID, moduleID, recordID })
      state.records.delete(recordID)
      return true
    } catch (error) {
      console.error('Failed to delete record:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  /**
   * Get a cached record by ID (synchronous, no API call).
   * Checks the full records map first, then the label cache.
   */
  function getByID(recordID) {
    return state.records.get(recordID) || state.labelCache.get(recordID) || null
  }

  /**
   * Resolve record labels for viewer components without requiring the module in moduleStore.
   * Fetches missing records directly via API and caches them in labelCache.
   *
   * @param {Object} params
   * @param {string} params.namespaceID
   * @param {string} params.moduleID
   * @param {string[]} params.recordIDs
   */
  async function resolveRecordLabels({ namespaceID, moduleID, recordIDs } = {}) {
    if (!$ComposeAPI || !namespaceID || !moduleID || !recordIDs?.length) return

    const missing = recordIDs.filter(id => id && !state.records.has(id) && !state.labelCache.has(id))
    if (!missing.length) return

    try {
      const results = await Promise.all(
        missing.map(recordID =>
          $ComposeAPI.recordRead({ namespaceID, moduleID, recordID }).catch(() => null),
        ),
      )
      results.forEach(r => {
        if (r?.recordID) state.labelCache.set(r.recordID, r)
      })
    } catch (e) {
      console.error('Failed to resolve record labels:', e)
    }
  }

  /**
   * Clear all cached records.
   */
  function clearAll() {
    state.records.clear()
    state.labelCache.clear()
    state.loading = false
    state.pending = false
  }

  return {
    // state
    loading: toRef(state, 'loading'),
    pending: toRef(state, 'pending'),

    // getters
    getByID,

    // actions
    list,
    findByID,
    create,
    update,
    delete: deleteRecord,
    resolveRecordLabels,
    clearAll,
  }
})

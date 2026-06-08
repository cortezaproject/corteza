import { compose } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { inject, reactive, ref, toRef } from 'vue'
import { useModuleStore } from './useModuleStore'

export const useRecordStore = defineStore('record', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    records: new Map(),
    labelCache: new Map(),
  })

  const paginationRecordIDs = ref([])

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

  async function findByID({ namespaceID, moduleID, recordID } = {}) {
    const cached = state.records.get(recordID)
    if (cached) return cached

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

  async function create(record) {
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

  async function update(record) {
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

  async function deleteRecord({ namespaceID, moduleID, recordID } = {}) {
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

  function getByID(recordID) {
    return state.records.get(recordID) || state.labelCache.get(recordID) || null
  }

  async function resolveRecordLabels({ namespaceID, moduleID, recordIDs } = {}) {
    if (!namespaceID || !moduleID || !recordIDs?.length) return

    const missing = recordIDs.filter(id => id && !state.records.has(id) && !state.labelCache.has(id))
    if (!missing.length) return

    const mod = useModuleStore().getByID(moduleID)

    try {
      const results = await Promise.all(
        missing.map(recordID =>
          $ComposeAPI.recordRead({ namespaceID, moduleID, recordID }).catch(() => null),
        ),
      )
      results.forEach(r => {
        if (r?.recordID) {
          state.labelCache.set(r.recordID, mod ? new compose.Record(mod, r) : r)
        }
      })
    } catch (e) {
      console.error('Failed to resolve record labels:', e)
    }
  }

  function setNavigationIDs(ids) {
    paginationRecordIDs.value = ids || []
  }

  function getNextAndPrev(recordID) {
    const ids = paginationRecordIDs.value
    const idx = ids.indexOf(recordID)
    if (idx === -1) return { prev: undefined, next: undefined }
    return {
      prev: idx > 0 ? ids[idx - 1] : undefined,
      next: idx < ids.length - 1 ? ids[idx + 1] : undefined,
    }
  }

  function clearAll() {
    state.records.clear()
    state.labelCache.clear()
    state.loading = false
    state.pending = false
    paginationRecordIDs.value = []
  }

  return {
    loading: toRef(state, 'loading'),
    pending: toRef(state, 'pending'),
    paginationRecordIDs,

    getByID,
    getNextAndPrev,

    list,
    findByID,
    create,
    update,
    delete: deleteRecord,
    resolveRecordLabels,
    setNavigationIDs,
    clearAll,
  }
})

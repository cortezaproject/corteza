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

  // --- Label batching ---
  // Coalesce all label requests made in the same render pass into one
  // `recordList` by ID per (namespace, module). pendingLabels
  // holds IDs awaiting the next flush; inflightLabels dedups reads across ticks
  // and lets awaiting callers (e.g. tests) block until their IDs settle.
  const LABEL_CHUNK = 100
  const pendingLabels = new Map() // `${namespaceID}:${moduleID}` -> Set<recordID>
  const inflightLabels = new Map() // recordID -> { promise, resolve }
  let flushScheduled = false

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

  async function findByID({ namespaceID, moduleID, recordID, force = false, signal } = {}) {
    // force bypasses the cache to guarantee fresh data — used when navigating to
    // a record page / opening a record for edit, where stale values are unsafe.
    if (!force) {
      const cached = state.records.get(recordID)
      if (cached) return cached
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
      // Pass an optional AbortSignal so callers can cancel an in-flight load
      // (e.g. navigating away from a record page) — axios honors `signal`.
      const raw = await $ComposeAPI.recordRead({ namespaceID, moduleID, recordID }, { signal })
      const record = new compose.Record(mod, raw)
      state.records.set(record.recordID, record)
      return record
    } catch (error) {
      // Don't log/alarm on intentional cancellation.
      if (!signal?.aborted) console.error('Failed to read record:', error)
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

  // Lifts the soft delete. The cached copy still carries deletedAt, so it is
  // evicted rather than patched — a caller that keeps showing the record reads
  // it back and gets the whole server-side state, permissions included.
  async function undeleteRecord({ namespaceID, moduleID, recordID } = {}) {
    state.pending = true
    try {
      await $ComposeAPI.recordUndelete({ namespaceID, moduleID, recordID })
      state.records.delete(recordID)
      state.labelCache.delete(recordID)
      return true
    } catch (error) {
      console.error('Failed to restore record:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  function getByID(recordID) {
    return state.records.get(recordID) || state.labelCache.get(recordID) || null
  }

  function settleLabel(recordID) {
    const entry = inflightLabels.get(recordID)
    if (entry) {
      inflightLabels.delete(recordID)
      entry.resolve()
    }
  }

  // Fetch one (namespace, module) group of IDs: one recordList per chunk, deleted
  // records included. An ID the batch does not return cannot be read or does not
  // exist, so it stays unresolved; only a batch that fails outright falls back to
  // reading its IDs one by one.
  async function fetchLabelGroup(namespaceID, moduleID, ids) {
    const mod = useModuleStore().getByID(moduleID)
    const wrap = r => (mod ? new compose.Record(mod, r) : r)

    for (let i = 0; i < ids.length; i += LABEL_CHUNK) {
      const chunk = ids.slice(i, i + LABEL_CHUNK)
      try {
        const { set = [] } = await $ComposeAPI.recordList({
          namespaceID,
          moduleID,
          recordID: chunk,
          deleted: 1,
          limit: chunk.length,
          incTotal: false,
        })
        set.forEach(r => {
          if (r?.recordID) state.labelCache.set(r.recordID, wrap(r))
        })
      } catch (e) {
        console.error('Failed to batch-resolve record labels:', e)
        await Promise.all(
          chunk.map(recordID =>
            $ComposeAPI
              .recordRead({ namespaceID, moduleID, recordID })
              .then(r => {
                if (r?.recordID) state.labelCache.set(r.recordID, wrap(r))
              })
              .catch(() => null),
          ),
        )
      }

      chunk.forEach(settleLabel)
    }
  }

  async function flushLabels() {
    flushScheduled = false
    const groups = [...pendingLabels.entries()]
    pendingLabels.clear()

    await Promise.all(
      groups.map(([key, idSet]) => {
        const [namespaceID, moduleID] = key.split(':')
        return fetchLabelGroup(namespaceID, moduleID, [...idSet])
      }),
    )
  }

  async function resolveRecordLabels({ namespaceID, moduleID, recordIDs } = {}) {
    if (!namespaceID || !moduleID || !recordIDs?.length) return

    const missing = recordIDs.filter(
      id =>
        id &&
        /^\d+$/.test(String(id)) &&
        !state.records.has(id) &&
        !state.labelCache.has(id) &&
        !inflightLabels.has(id),
    )

    if (missing.length) {
      const key = `${namespaceID}:${moduleID}`
      const group = pendingLabels.get(key) || new Set()
      missing.forEach(id => {
        group.add(id)
        let resolve
        const promise = new Promise(res => {
          resolve = res
        })
        inflightLabels.set(id, { promise, resolve })
      })
      pendingLabels.set(key, group)

      if (!flushScheduled) {
        flushScheduled = true
        queueMicrotask(flushLabels)
      }
    }

    // Block awaiting callers until every requested ID has settled (resolved or
    // already cached). The viewer's watch is fire-and-forget, so it's unaffected.
    await Promise.all(recordIDs.map(id => inflightLabels.get(id)?.promise ?? Promise.resolve()))
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
    pendingLabels.clear()
    inflightLabels.forEach(entry => entry.resolve())
    inflightLabels.clear()
    flushScheduled = false
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
    undelete: undeleteRecord,
    resolveRecordLabels,
    setNavigationIDs,
    clearAll,
  }
})

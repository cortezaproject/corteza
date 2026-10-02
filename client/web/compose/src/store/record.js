// Records in this store not classes, but raw objects instead
const types = {
  pending: 'pending',
  completed: 'completed',
  updateSet: 'updateSet',
  clearSet: 'clearSet',
}

export default function (ComposeAPI) {
  // Records are keyed by moduleID + recordID, not recordID alone: modules backed
  // by external databases can reuse the same (int) recordID across different
  // modules, so a bare recordID is not a unique key across the shared set.
  const recordKey = (moduleID, recordID) => `${moduleID}/${recordID}`

  // Batching state for resolveRecords (shared across all dispatches)
  const pendingBatches = new Map() // key: `${namespaceID}/${moduleID}` -> { ids, resolvers, namespaceID, moduleID }
  const inflightKeys = new Set() // moduleID/recordID keys currently being fetched
  let flushTimer = null

  function flushResolves (commit) {
    const batches = new Map(pendingBatches)
    pendingBatches.clear()

    for (const [, { ids, resolvers, namespaceID, moduleID }] of batches) {
      const recordIDs = [...ids]

      if (recordIDs.length === 0) {
        resolvers.forEach(r => r())
        continue
      }

      recordIDs.forEach(id => inflightKeys.add(recordKey(moduleID, id)))
      commit(types.pending)

      const query = recordIDs.map(id => `recordID = ${id}`).join(' OR ')

      ComposeAPI.recordList({ namespaceID, moduleID, query, deleted: 1 })
        .then(({ set }) => {
          commit(types.updateSet, set)
        })
        .finally(() => {
          recordIDs.forEach(id => inflightKeys.delete(recordKey(moduleID, id)))
          commit(types.completed)
          resolvers.forEach(r => r())
        })
    }
  }

  return {
    namespaced: true,

    state: {
      pending: false,
      set: [],
    },

    getters: {
      pending: (state) => state.pending,

      findByID (state) {
        // moduleID is optional for backwards compatibility, but should be passed
        // whenever known to disambiguate records that share a recordID across modules
        return (ID, moduleID = undefined) => state.set.find(
          (r) => ID === r.recordID && (moduleID === undefined || r.moduleID === moduleID),
        )
      },

      findByIDs (state) {
        return (IDs, moduleID = undefined) => {
          const idSet = new Set(IDs.flat())
          return state.set.filter(
            (r) => idSet.has(r.recordID) && (moduleID === undefined || r.moduleID === moduleID),
          )
        }
      },

      set (state) {
        return state.set
      },
    },

    actions: {
      /**
       * Batched record resolver. Collects IDs from multiple callers within a
       * short window (50ms), deduplicates against the store and in-flight
       * requests, then fires a single API call per module.
       *
       * Returns a promise that resolves when the batch containing
       * the requested IDs has been fetched and committed to the store.
       */
      resolveRecords ({ commit, getters }, { namespaceID, moduleID, recordIDs }) {
        if (recordIDs.length === 0) {
          return Promise.resolve()
        }

        // Filter out records already in the store or currently being fetched.
        // Scope "known" to this module so a matching recordID in another module
        // doesn't mask a record we still need to fetch here.
        const knownIDs = new Set(
          getters.set.filter(({ moduleID: mID }) => mID === moduleID).map(({ recordID }) => recordID),
        )
        recordIDs = recordIDs.filter(id => !knownIDs.has(id) && !inflightKeys.has(recordKey(moduleID, id)))

        if (recordIDs.length === 0) {
          return Promise.resolve()
        }

        // Add to the pending batch for this module
        const key = `${namespaceID}/${moduleID}`

        if (!pendingBatches.has(key)) {
          pendingBatches.set(key, { ids: new Set(), resolvers: [], namespaceID, moduleID })
        }

        const batch = pendingBatches.get(key)
        recordIDs.forEach(id => batch.ids.add(id))

        // Each caller gets a promise that resolves when its batch completes
        const promise = new Promise(resolve => {
          batch.resolvers.push(resolve)
        })

        // Reset the debounce — flush after 50ms of no new requests
        clearTimeout(flushTimer)
        flushTimer = setTimeout(() => flushResolves(commit), 50)

        return promise
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
          state.set = set.map(r => JSON.parse(JSON.stringify(r)))
          return
        }

        // Build index map for O(1) lookups, keyed by moduleID+recordID so records
        // from different modules that share a recordID don't overwrite each other
        const indexByKey = new Map(state.set.map((r, i) => [recordKey(r.moduleID, r.recordID), i]))

        set.forEach(newItem => {
          newItem = JSON.parse(JSON.stringify(newItem))

          const key = recordKey(newItem.moduleID, newItem.recordID)
          const oldIndex = indexByKey.get(key)
          if (oldIndex !== undefined) {
            state.set.splice(oldIndex, 1, newItem)
          } else {
            indexByKey.set(key, state.set.length)
            state.set.push(newItem)
          }
        })
      },

      [types.clearSet] (state) {
        state.pending = false
        state.set.splice(0)

        // Clean up batching state
        inflightKeys.clear()
        clearTimeout(flushTimer)
        pendingBatches.clear()
      },
    },
  }
}

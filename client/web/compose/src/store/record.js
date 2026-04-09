// Records in this store not classes, but raw objects instead
const types = {
  pending: 'pending',
  completed: 'completed',
  updateSet: 'updateSet',
  clearSet: 'clearSet',
}

export default function (ComposeAPI) {
  // Batching state for resolveRecords (shared across all dispatches)
  const pendingBatches = new Map() // key: `${namespaceID}/${moduleID}` -> { ids, resolvers, namespaceID, moduleID }
  const inflightIDs = new Set()
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

      recordIDs.forEach(id => inflightIDs.add(id))
      commit(types.pending)

      const query = recordIDs.map(id => `recordID = ${id}`).join(' OR ')

      ComposeAPI.recordList({ namespaceID, moduleID, query, deleted: 1 })
        .then(({ set }) => {
          commit(types.updateSet, set)
        })
        .finally(() => {
          recordIDs.forEach(id => inflightIDs.delete(id))
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
        return (ID) => state.set.find(({ recordID }) => ID === recordID)
      },

      findByIDs (state) {
        return (IDs) => {
          const idSet = new Set(IDs.flat())
          return state.set.filter(({ recordID }) => idSet.has(recordID))
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

        // Filter out records already in the store or currently being fetched
        const knownIDs = new Set(getters.set.map(({ recordID }) => recordID))
        recordIDs = recordIDs.filter(id => !knownIDs.has(id) && !inflightIDs.has(id))

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

        // Build index map for O(1) lookups
        const indexByID = new Map(state.set.map(({ recordID }, i) => [recordID, i]))

        set.forEach(newItem => {
          newItem = JSON.parse(JSON.stringify(newItem))

          const oldIndex = indexByID.get(newItem.recordID)
          if (oldIndex !== undefined) {
            state.set.splice(oldIndex, 1, newItem)
          } else {
            indexByID.set(newItem.recordID, state.set.length)
            state.set.push(newItem)
          }
        })
      },

      [types.clearSet] (state) {
        state.pending = false
        state.set.splice(0)

        // Clean up batching state
        inflightIDs.clear()
        clearTimeout(flushTimer)
        pendingBatches.clear()
      },
    },
  }
}

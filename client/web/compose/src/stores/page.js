import { compose } from '@cortezaproject/corteza-js-next'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const usePageStore = defineStore('page', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    set: [],
    namespaceID: null,
  })

  // Getters
  const loading = computed(() => state.loading)
  const pending = computed(() => state.pending)
  const set = computed(() => state.set)

  const getByID = computed(() => {
    return ID => state.set.find(({ pageID }) => ID === pageID)
  })

  const getByHandle = computed(() => {
    return handle => state.set.find(p => p.handle === handle)
  })

  // Actions
  async function load({ namespace, namespaceID, clear = false, force = false } = {}) {
    const nsID = namespaceID || namespace?.namespaceID
    if (!nsID) {
      console.error('Cannot load pages without namespaceID')
      return Promise.reject(new Error('namespaceID required'))
    }

    if (clear) {
      clearSet()
    }

    // If same namespace and already loaded, skip
    if (!force && state.namespaceID === nsID && state.set.length > 1) {
      return Promise.resolve(state.set)
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to load pages:', error)
      return Promise.reject(error)
    }

    state.loading = true
    state.pending = true
    state.namespaceID = nsID

    try {
      const { set: pageSet } = await $ComposeAPI.pageList({
        namespaceID: nsID,
        sort: 'weight ASC',
      })

      if (pageSet && pageSet.length > 0) {
        updateSet(pageSet.map(p => new compose.Page(p)))
      }

      return state.set
    } catch (error) {
      console.error('Failed to load pages:', error)
      throw error
    } finally {
      state.loading = false
      state.pending = false
    }
  }

  async function findByID({ namespaceID, pageID, force = false } = {}) {
    if (!force) {
      const oldItem = getByID.value(pageID)
      if (oldItem) {
        return Promise.resolve(oldItem)
      }
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to find page:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.pageRead({ namespaceID, pageID })
      const page = new compose.Page(raw)
      updateSet([page])
      return page
    } catch (error) {
      console.error('Failed to find page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to create page:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.pageCreate(item)
      const page = new compose.Page(raw)
      updateSet([page])
      return page
    } catch (error) {
      console.error('Failed to create page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to update page:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.pageUpdate(item)
      const page = new compose.Page(raw)
      updateSet([page])
      return page
    } catch (error) {
      console.error('Failed to update page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deletePage(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to delete page:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      await $ComposeAPI.pageDelete(item)
      removeFromSet([item])
      return true
    } catch (error) {
      console.error('Failed to delete page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  function clearSet() {
    state.pending = false
    state.set.splice(0)
    state.namespaceID = null
  }

  // Helper functions
  function updateSet(newSet) {
    const frozenSet = newSet.map(i => Object.freeze(i))

    if (state.set.length === 0) {
      state.set = frozenSet
      return
    }

    frozenSet.forEach(newItem => {
      const oldIndex = state.set.findIndex(({ pageID }) => pageID === newItem.pageID)
      if (oldIndex > -1) {
        state.set.splice(oldIndex, 1, newItem)
      } else {
        state.set.push(newItem)
      }
    })
  }

  function removeFromSet(removedSet) {
    ;(removedSet || []).forEach(removedItem => {
      const i = state.set.findIndex(({ pageID }) => pageID === removedItem.pageID)
      if (i > -1) {
        state.set.splice(i, 1)
      }
    })
  }

  async function loadTree({ namespaceID } = {}) {
    const nsID = namespaceID || state.namespaceID
    if (!nsID) {
      return Promise.reject(new Error('namespaceID required'))
    }

    if (!$ComposeAPI) {
      return Promise.reject(new Error('ComposeAPI not available via inject'))
    }

    state.loading = true

    try {
      const pages = await $ComposeAPI.pageTree({ namespaceID: nsID })
      return pages || []
    } catch (error) {
      console.error('Failed to load page tree:', error)
      throw error
    } finally {
      state.loading = false
    }
  }

  async function reorder({ namespaceID, selfID, pageIDs } = {}) {
    if (!$ComposeAPI) {
      return Promise.reject(new Error('ComposeAPI not available via inject'))
    }

    try {
      await $ComposeAPI.pageReorder({ namespaceID, selfID, pageIDs })
    } catch (error) {
      console.error('Failed to reorder pages:', error)
      throw error
    }
  }

  async function reparent({ namespaceID, pageID, newParentID, siblingPageIDs } = {}) {
    if (!$ComposeAPI) {
      return Promise.reject(new Error('ComposeAPI not available via inject'))
    }

    try {
      // First, find the page to get current data for update
      const page = getByID.value(pageID)
      if (page) {
        await $ComposeAPI.pageUpdate({
          ...page,
          namespaceID,
          pageID,
          selfID: newParentID,
        })
      }

      // Then reorder children under the new parent
      if (siblingPageIDs && siblingPageIDs.length > 0) {
        await reorder({ namespaceID, selfID: newParentID, pageIDs: siblingPageIDs })
      }
    } catch (error) {
      console.error('Failed to reparent page:', error)
      throw error
    }
  }

  return {
    // state
    loading: toRef(state, 'loading'),
    pending: toRef(state, 'pending'),
    set: toRef(state, 'set'),
    namespaceID: toRef(state, 'namespaceID'),

    // getters
    getByID,
    getByHandle,

    // actions
    load,
    loadTree,
    findByID,
    create,
    update,
    delete: deletePage,
    reorder,
    reparent,
    clearSet,
    updateSet,
  }
})

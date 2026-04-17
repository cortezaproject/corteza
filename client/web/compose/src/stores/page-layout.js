import { compose } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const usePageLayoutStore = defineStore('pageLayout', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    set: [],
    namespaceID: null,
  })

  const getByID = computed(() => {
    return ID => state.set.find(({ pageLayoutID }) => ID === pageLayoutID)
  })

  const getByPageID = computed(() => {
    return pageID => state.set.filter(l => l.pageID === pageID)
  })

  // Actions
  async function load({ namespace, namespaceID, clear = false, force = false } = {}) {
    const nsID = namespaceID || namespace?.namespaceID
    if (!nsID) {
      console.error('Cannot load page layouts without namespaceID')
      return Promise.reject(new Error('namespaceID required'))
    }

    if (clear) {
      clearSet()
    }

    // If same namespace and already loaded, skip
    if (!force && state.namespaceID === nsID && state.set.length > 0) {
      return Promise.resolve(state.set)
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to load page layouts:', error)
      return Promise.reject(error)
    }

    state.loading = true
    state.pending = true
    state.namespaceID = nsID

    try {
      const { set: layoutSet } = await $ComposeAPI.pageLayoutListNamespace({
        namespaceID: nsID,
        sort: 'weight ASC',
      })

      if (layoutSet && layoutSet.length > 0) {
        updateSet(layoutSet.map(l => new compose.PageLayout(l)))
      }

      return state.set
    } catch (error) {
      console.error('Failed to load page layouts:', error)
      throw error
    } finally {
      state.loading = false
      state.pending = false
    }
  }

  async function findByID({ namespaceID, pageID, pageLayoutID, force = false } = {}) {
    if (!force) {
      const oldItem = getByID.value(pageLayoutID)
      if (oldItem) {
        return Promise.resolve(oldItem)
      }
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to find page layout:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.pageLayoutRead({ namespaceID, pageID, pageLayoutID })
      const layout = new compose.PageLayout(raw)
      updateSet([layout])
      return layout
    } catch (error) {
      console.error('Failed to find page layout:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function findByPageID({ namespaceID, pageID, force = false } = {}) {
    if (!force) {
      const cached = getByPageID.value(pageID)
      if (cached && cached.length > 0) {
        return cached
      }
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to find page layouts:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const { set: layoutSet } = await $ComposeAPI.pageLayoutList({
        namespaceID,
        pageID,
        sort: 'weight ASC',
      })

      const layouts = (layoutSet || []).map(l => new compose.PageLayout(l))
      updateSet(layouts)
      return layouts
    } catch (error) {
      console.error('Failed to find page layouts by page ID:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to create page layout:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.pageLayoutCreate(item)
      const layout = new compose.PageLayout(raw)
      updateSet([layout])
      return layout
    } catch (error) {
      console.error('Failed to create page layout:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to update page layout:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.pageLayoutUpdate(item)
      const layout = new compose.PageLayout(raw)
      updateSet([layout])
      return layout
    } catch (error) {
      console.error('Failed to update page layout:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deleteLayout(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to delete page layout:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      await $ComposeAPI.pageLayoutDelete(item)
      removeFromSet([item])
      return true
    } catch (error) {
      console.error('Failed to delete page layout:', error)
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
      const oldIndex = state.set.findIndex(
        ({ pageLayoutID }) => pageLayoutID === newItem.pageLayoutID,
      )
      if (oldIndex > -1) {
        state.set.splice(oldIndex, 1, newItem)
      } else {
        state.set.push(newItem)
      }
    })
  }

  function removeFromSet(removedSet) {
    ;(removedSet || []).forEach(removedItem => {
      const i = state.set.findIndex(
        ({ pageLayoutID }) => pageLayoutID === removedItem.pageLayoutID,
      )
      if (i > -1) {
        state.set.splice(i, 1)
      }
    })
  }

  return {
    // state
    loading: toRef(state, 'loading'),
    pending: toRef(state, 'pending'),
    set: toRef(state, 'set'),
    namespaceID: toRef(state, 'namespaceID'),

    // getters
    getByID,
    getByPageID,

    // actions
    load,
    findByID,
    findByPageID,
    create,
    update,
    delete: deleteLayout,
    clearSet,
    updateSet,
  }
})

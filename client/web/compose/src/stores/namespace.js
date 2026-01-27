import { compose } from '@cortezaproject/corteza-js-next'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const useNamespaceStore = defineStore('namespace', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    set: [],
  })

  // Getters
  const loading = computed(() => state.loading)
  const pending = computed(() => state.pending)
  const set = computed(() => state.set)

  const getByID = computed(() => {
    return ID => state.set.find(({ namespaceID }) => ID === namespaceID)
  })

  const getByUrlPart = computed(() => {
    return urlPart =>
      state.set.find(({ slug, namespaceID }) => urlPart === slug || urlPart === namespaceID)
  })

  // Actions
  async function load({ force = false } = {}) {
    if (!force && state.set.length > 1) {
      // When there's forced load, make sure we have more than 1 item in the set
      // in the scenario when user came to detail page first and has one item loaded
      // > 0 would not be sufficient.
      return Promise.resolve(state.set)
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to load namespaces:', error)
      return Promise.reject(error)
    }

    state.loading = true
    state.pending = true

    try {
      const { set: namespaceSet } = await $ComposeAPI.namespaceList({})

      if (namespaceSet && namespaceSet.length > 0) {
        updateSet(namespaceSet.map(n => new compose.Namespace(n)))
      }

      return state.set
    } catch (error) {
      console.error('Failed to load namespaces:', error)
      throw error
    } finally {
      state.loading = false
      state.pending = false
    }
  }

  async function findByID({ namespaceID, force = false } = {}) {
    if (!force) {
      const oldItem = getByID.value(namespaceID)
      if (oldItem) {
        return Promise.resolve(oldItem)
      }
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to find namespace:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceRead({ namespaceID })
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return namespace
    } catch (error) {
      console.error('Failed to find namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to create namespace:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceCreate(item)
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return namespace
    } catch (error) {
      console.error('Failed to create namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function clone(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to clone namespace:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceClone(item)
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return namespace
    } catch (error) {
      console.error('Failed to clone namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to update namespace:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceUpdate(item)
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return namespace
    } catch (error) {
      console.error('Failed to update namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deleteNamespace(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to delete namespace:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      await $ComposeAPI.namespaceDelete(item)
      removeFromSet([item])
      return true
    } catch (error) {
      console.error('Failed to delete namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  function clearSet() {
    state.pending = false
    state.set.splice(0)
  }

  // Helper functions for state management
  function updateSet(newSet) {
    const frozenSet = newSet.map(i => Object.freeze(i))

    if (state.set.length === 0) {
      state.set = frozenSet
      return
    }

    frozenSet.forEach(newItem => {
      const oldIndex = state.set.findIndex(({ namespaceID }) => namespaceID === newItem.namespaceID)
      if (oldIndex > -1) {
        state.set.splice(oldIndex, 1, newItem)
      } else {
        state.set.push(newItem)
      }
    })
  }

  function removeFromSet(removedSet) {
    ;(removedSet || []).forEach(removedItem => {
      const i = state.set.findIndex(({ namespaceID }) => namespaceID === removedItem.namespaceID)
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

    // getters
    getByID,
    getByUrlPart,

    // actions
    load,
    findByID,
    create,
    clone,
    update,
    delete: deleteNamespace,
    clearSet,
  }
})

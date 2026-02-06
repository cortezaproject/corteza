import { compose } from '@cortezaproject/corteza-js-next'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const useModuleStore = defineStore('module', () => {
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
    return ID => state.set.find(({ moduleID }) => ID === moduleID)
  })

  const getByHandle = computed(() => {
    return handle => state.set.find(m => m.handle === handle)
  })

  // Actions
  async function load({ namespace, namespaceID, clear = false, force = false } = {}) {
    const nsID = namespaceID || namespace?.namespaceID
    if (!nsID) {
      console.error('Cannot load modules without namespaceID')
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
      console.error('Failed to load modules:', error)
      return Promise.reject(error)
    }

    state.loading = true
    state.pending = true
    state.namespaceID = nsID

    try {
      const { set: moduleSet } = await $ComposeAPI.moduleList({
        namespaceID: nsID,
        sort: 'name ASC',
      })

      if (moduleSet && moduleSet.length > 0) {
        updateSet(moduleSet.map(m => new compose.Module(m)))
      }

      return state.set
    } catch (error) {
      console.error('Failed to load modules:', error)
      throw error
    } finally {
      state.loading = false
      state.pending = false
    }
  }

  async function findByID({ namespaceID, moduleID, force = false } = {}) {
    if (!force) {
      const oldItem = getByID.value(moduleID)
      if (oldItem) {
        return Promise.resolve(oldItem)
      }
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to find module:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.moduleRead({ namespaceID, moduleID })
      const module = new compose.Module(raw)
      updateSet([module])
      return module
    } catch (error) {
      console.error('Failed to find module:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to create module:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      // Wrap in Module class for proper serialization
      const moduleToCreate = new compose.Module(item)
      const raw = await $ComposeAPI.moduleCreate(moduleToCreate)
      const module = new compose.Module(raw)
      updateSet([module])
      return module
    } catch (error) {
      console.error('Failed to create module:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to update module:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      // Wrap in Module class for proper serialization
      const moduleToUpdate = new compose.Module(item)
      const raw = await $ComposeAPI.moduleUpdate(moduleToUpdate)
      const module = new compose.Module(raw)
      updateSet([module])
      return module
    } catch (error) {
      console.error('Failed to update module:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deleteModule(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to delete module:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      await $ComposeAPI.moduleDelete(item)
      removeFromSet([item])
      return true
    } catch (error) {
      console.error('Failed to delete module:', error)
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
      const oldIndex = state.set.findIndex(({ moduleID }) => moduleID === newItem.moduleID)
      if (oldIndex > -1) {
        state.set.splice(oldIndex, 1, newItem)
      } else {
        state.set.push(newItem)
      }
    })
  }

  function removeFromSet(removedSet) {
    ;(removedSet || []).forEach(removedItem => {
      const i = state.set.findIndex(({ moduleID }) => moduleID === removedItem.moduleID)
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
    getByHandle,

    // actions
    load,
    findByID,
    create,
    update,
    delete: deleteModule,
    clearSet,
    updateSet,
  }
})

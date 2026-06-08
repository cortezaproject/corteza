import { compose } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const useModuleStore = defineStore('module', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    // Multi-namespace cache: namespaceID -> Module[]. Lets the shared module
    // picker (CInputModule) hold modules for arbitrary namespaces without
    // disturbing the namespace the compose nav/views are looking at.
    byNamespace: {},
    // The "active" namespace whose modules `set` reflects (compose nav/views).
    namespaceID: null,
  })

  // Back-compat: existing consumers read `set` = the active namespace's modules.
  const set = computed(() => state.byNamespace[state.namespaceID] || [])

  // moduleID is globally unique, so resolve across every cached namespace —
  // a module loaded by any section/picker is findable everywhere.
  const getByID = computed(() => {
    return ID => {
      for (const list of Object.values(state.byNamespace)) {
        const found = list.find(({ moduleID }) => ID === moduleID)
        if (found) return found
      }
      return undefined
    }
  })

  const getByHandle = computed(() => {
    return handle => (state.byNamespace[state.namespaceID] || []).find(m => m.handle === handle)
  })

  // --- internal cache mutators ---
  function cacheModules(nsID, modules) {
    if (!nsID) return
    const existing = state.byNamespace[nsID] || []
    if (existing.length === 0) {
      state.byNamespace[nsID] = modules.map(m => Object.freeze(m))
      return
    }
    const next = existing.slice()
    modules.forEach(item => {
      const frozen = Object.freeze(item)
      const i = next.findIndex(({ moduleID }) => moduleID === frozen.moduleID)
      if (i > -1) next.splice(i, 1, frozen)
      else next.push(frozen)
    })
    state.byNamespace[nsID] = next
  }

  async function fetchInto(nsID) {
    const { set: moduleSet = [] } = await $ComposeAPI.moduleList({ namespaceID: nsID, sort: 'name ASC' })
    if (moduleSet && moduleSet.length > 0) {
      cacheModules(nsID, moduleSet.map(m => new compose.Module(m)))
    } else if (!state.byNamespace[nsID]) {
      state.byNamespace[nsID] = []
    }
    return state.byNamespace[nsID] || []
  }

  // --- namespace-scoped accessors (for the shared module picker) ---
  function modulesFor(namespaceID) {
    return state.byNamespace[namespaceID] || []
  }

  // Ensure a namespace's modules are cached, WITHOUT changing the active ns.
  // Cache-first: picker use case — don't refetch on every dropdown open.
  async function loadFor(namespaceID) {
    if (!namespaceID) return []
    if (state.byNamespace[namespaceID]?.length > 0) {
      return state.byNamespace[namespaceID]
    }
    state.loading = true
    state.pending = true
    try {
      return await fetchInto(namespaceID)
    } catch (error) {
      console.error('Failed to load modules:', error)
      throw error
    } finally {
      state.loading = false
      state.pending = false
    }
  }

  // --- actions ---
  async function load({ namespace, namespaceID, clear = false } = {}) {
    const nsID = namespaceID || namespace?.namespaceID
    if (!nsID) {
      console.error('Cannot load modules without namespaceID')
      return Promise.reject(new Error('namespaceID required'))
    }

    if (clear) clearSet()

    // Set as the active namespace (drives `set`/`getByHandle`).
    state.namespaceID = nsID

    state.loading = true
    state.pending = true
    try {
      return await fetchInto(nsID)
    } catch (error) {
      console.error('Failed to load modules:', error)
      throw error
    } finally {
      state.loading = false
      state.pending = false
    }
  }

  async function findByID({ namespaceID, moduleID } = {}) {
    const cached = getByID.value(moduleID)
    if (cached) return new compose.Module(cached)

    state.pending = true
    try {
      const raw = await $ComposeAPI.moduleRead({ namespaceID, moduleID })
      const module = new compose.Module(raw)
      cacheModules(namespaceID || module.namespaceID, [module])
      return new compose.Module(module)
    } catch (error) {
      console.error('Failed to find module:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    state.pending = true
    try {
      const moduleToCreate = new compose.Module(item)
      const raw = await $ComposeAPI.moduleCreate(moduleToCreate)
      const module = new compose.Module(raw)
      updateSet([module])
      return new compose.Module(module)
    } catch (error) {
      console.error('Failed to create module:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    state.pending = true
    try {
      const moduleToUpdate = new compose.Module(item)
      const raw = await $ComposeAPI.moduleUpdate(moduleToUpdate)
      const module = new compose.Module(raw)
      updateSet([module])
      return new compose.Module(module)
    } catch (error) {
      console.error('Failed to update module:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deleteModule(item) {
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
    if (state.namespaceID) delete state.byNamespace[state.namespaceID]
    state.namespaceID = null
  }

  function updateSet(newSet) {
    const groups = {}
    ;(newSet || []).forEach(m => {
      const ns = m.namespaceID || state.namespaceID
      if (!ns) return
      ;(groups[ns] = groups[ns] || []).push(m)
    })
    Object.entries(groups).forEach(([ns, mods]) => cacheModules(ns, mods))
  }

  function removeFromSet(removedSet) {
    ;(removedSet || []).forEach(removedItem => {
      const ns = removedItem.namespaceID || state.namespaceID
      const list = state.byNamespace[ns]
      if (!list) return
      const i = list.findIndex(({ moduleID }) => moduleID === removedItem.moduleID)
      if (i > -1) {
        const next = list.slice()
        next.splice(i, 1)
        state.byNamespace[ns] = next
      }
    })
  }

  return {
    // state
    loading: toRef(state, 'loading'),
    pending: toRef(state, 'pending'),
    set,
    namespaceID: toRef(state, 'namespaceID'),

    // getters
    getByID,
    getByHandle,

    // actions
    load,
    loadFor,
    modulesFor,
    findByID,
    create,
    update,
    delete: deleteModule,
    clearSet,
    updateSet,
    removeFromSet,
  }
})

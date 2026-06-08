import { compose } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const useNamespaceStore = defineStore('namespace', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    set: [],
  })

  const getByID = computed(() => {
    return ID => state.set.find(({ namespaceID }) => ID === namespaceID)
  })

  const getByUrlPart = computed(() => {
    return urlPart =>
      state.set.find(({ slug, namespaceID }) => urlPart === slug || urlPart === namespaceID)
  })

  async function load() {
    state.loading = true
    state.pending = true

    try {
      const { set: namespaceSet = [] } = await $ComposeAPI.namespaceList({})

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

  async function findByID({ namespaceID } = {}) {
    const cached = getByID.value(namespaceID)
    if (cached) {
      return new compose.Namespace(cached)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceRead({ namespaceID })
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return new compose.Namespace(namespace)
    } catch (error) {
      console.error('Failed to find namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceCreate(item)
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return new compose.Namespace(namespace)
    } catch (error) {
      console.error('Failed to create namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function clone(item) {
    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceClone(item)
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return new compose.Namespace(namespace)
    } catch (error) {
      console.error('Failed to clone namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    state.pending = true

    try {
      const raw = await $ComposeAPI.namespaceUpdate(item)
      const namespace = new compose.Namespace(raw)
      updateSet([namespace])
      return new compose.Namespace(namespace)
    } catch (error) {
      console.error('Failed to update namespace:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deleteNamespace(item) {
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

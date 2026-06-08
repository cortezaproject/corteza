import { compose } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const useChartStore = defineStore('chart', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    set: [],
    namespaceID: null,
  })

  const getByID = computed(() => {
    return ID => state.set.find(({ chartID }) => ID === chartID)
  })

  const getByHandle = computed(() => {
    return handle => state.set.find(c => c.handle === handle)
  })

  // Actions
  async function load({ namespace, namespaceID, clear = false, force = false } = {}) {
    const nsID = namespaceID || namespace?.namespaceID
    if (!nsID) {
      console.error('Cannot load charts without namespaceID')
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
      console.error('Failed to load charts:', error)
      return Promise.reject(error)
    }

    state.loading = true
    state.pending = true
    state.namespaceID = nsID

    try {
      const { set: chartSet = [] } = await $ComposeAPI.chartList({
        namespaceID: nsID,
        sort: 'name ASC',
      })

      if (chartSet && chartSet.length > 0) {
        updateSet(chartSet.map(c => new compose.Chart(c)))
      }

      return state.set
    } catch (error) {
      console.error('Failed to load charts:', error)
      throw error
    } finally {
      state.loading = false
      state.pending = false
    }
  }

  async function findByID({ namespaceID, chartID, force = false } = {}) {
    if (!force) {
      const oldItem = getByID.value(chartID)
      if (oldItem) {
        return Promise.resolve(oldItem)
      }
    }

    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to find chart:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.chartRead({ namespaceID, chartID })
      const chart = new compose.Chart(raw)
      updateSet([chart])
      return chart
    } catch (error) {
      console.error('Failed to find chart:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to create chart:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.chartCreate(item)
      const chart = new compose.Chart(raw)
      updateSet([chart])
      return chart
    } catch (error) {
      console.error('Failed to create chart:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to update chart:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const raw = await $ComposeAPI.chartUpdate(item)
      const chart = new compose.Chart(raw)
      updateSet([chart])
      return chart
    } catch (error) {
      console.error('Failed to update chart:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deleteChart(item) {
    if (!$ComposeAPI) {
      const error = new Error('ComposeAPI not available via inject')
      console.error('Failed to delete chart:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      await $ComposeAPI.chartDelete(item)
      removeFromSet([item])
      return true
    } catch (error) {
      console.error('Failed to delete chart:', error)
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
      const oldIndex = state.set.findIndex(({ chartID }) => chartID === newItem.chartID)
      if (oldIndex > -1) {
        state.set.splice(oldIndex, 1, newItem)
      } else {
        state.set.push(newItem)
      }
    })
  }

  function removeFromSet(removedSet) {
    ;(removedSet || []).forEach(removedItem => {
      const i = state.set.findIndex(({ chartID }) => chartID === removedItem.chartID)
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
    delete: deleteChart,
    clearSet,
    updateSet,
  }
})

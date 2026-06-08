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

  // Cache-first: picker use case — don't refetch on every dropdown open.
  async function loadFor(namespaceID) {
    if (!namespaceID) return []
    if (state.namespaceID === namespaceID && state.set.length > 0) {
      return state.set
    }
    return load({ namespaceID })
  }

  async function load({ namespace, namespaceID, clear = false } = {}) {
    const nsID = namespaceID || namespace?.namespaceID
    if (!nsID) {
      console.error('Cannot load charts without namespaceID')
      return Promise.reject(new Error('namespaceID required'))
    }

    if (clear) clearSet()

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

  async function findByID({ namespaceID, chartID } = {}) {
    const cached = getByID.value(chartID)
    if (cached) return new compose.Chart(cached)

    state.pending = true
    try {
      const raw = await $ComposeAPI.chartRead({ namespaceID, chartID })
      const chart = new compose.Chart(raw)
      updateSet([chart])
      return new compose.Chart(chart)
    } catch (error) {
      console.error('Failed to find chart:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    state.pending = true
    try {
      const raw = await $ComposeAPI.chartCreate(item)
      const chart = new compose.Chart(raw)
      updateSet([chart])
      return new compose.Chart(chart)
    } catch (error) {
      console.error('Failed to create chart:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    state.pending = true
    try {
      const raw = await $ComposeAPI.chartUpdate(item)
      const chart = new compose.Chart(raw)
      updateSet([chart])
      return new compose.Chart(chart)
    } catch (error) {
      console.error('Failed to update chart:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function deleteChart(item) {
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
      if (i > -1) state.set.splice(i, 1)
    })
  }

  return {
    loading: toRef(state, 'loading'),
    pending: toRef(state, 'pending'),
    set: toRef(state, 'set'),
    namespaceID: toRef(state, 'namespaceID'),

    getByID,
    getByHandle,

    load,
    loadFor,
    findByID,
    create,
    update,
    delete: deleteChart,
    clearSet,
    updateSet,
  }
})

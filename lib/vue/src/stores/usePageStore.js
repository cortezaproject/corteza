import { compose } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

import { usePageLayoutStore } from './usePageLayoutStore'

export const usePageStore = defineStore('page', () => {
  const $ComposeAPI = inject('$ComposeAPI')

  const state = reactive({
    loading: false,
    pending: false,
    set: [],
    namespaceID: null,
  })

  const getByID = computed(() => {
    return ID => state.set.find(({ pageID }) => ID === pageID)
  })

  const getByHandle = computed(() => {
    return handle => state.set.find(p => p.handle === handle)
  })

  async function load({ namespace, namespaceID, clear = false } = {}) {
    const nsID = namespaceID || namespace?.namespaceID
    if (!nsID) {
      console.error('Cannot load pages without namespaceID')
      return Promise.reject(new Error('namespaceID required'))
    }

    if (clear) clearSet()

    state.loading = true
    state.pending = true
    state.namespaceID = nsID

    try {
      const { set: pageSet = [] } = await $ComposeAPI.pageList({
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
    const cached = force ? null : getByID.value(pageID)
    if (cached) return new compose.Page(cached)

    state.pending = true
    try {
      const raw = await $ComposeAPI.pageRead({ namespaceID, pageID })
      const page = new compose.Page(raw)
      updateSet([page])
      return new compose.Page(page)
    } catch (error) {
      console.error('Failed to find page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function create(item) {
    state.pending = true
    try {
      const raw = await $ComposeAPI.pageCreate(item)
      const page = new compose.Page(raw)
      updateSet([page])
      await pullLayouts(page)
      return new compose.Page(page)
    } catch (error) {
      console.error('Failed to create page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function update(item) {
    state.pending = true
    try {
      const raw = await $ComposeAPI.pageUpdate(item)
      const page = new compose.Page(raw)
      updateSet([page])
      return new compose.Page(page)
    } catch (error) {
      console.error('Failed to update page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  // The primary layout is made server-side, inside the transaction that makes
  // the page, so a client that only records the page renders it blank: the
  // builder and the page view both position blocks from a layout they never
  // fetched. Pulling it here keeps the two stores as coupled as the server has
  // them, for every caller that creates a page.
  //
  // A layout that fails to arrive is not a failed create — the page exists.
  async function pullLayouts(page) {
    try {
      await usePageLayoutStore().findByPageID({
        namespaceID: page.namespaceID,
        pageID: page.pageID,
        force: true,
      })
    } catch (error) {
      console.error('Failed to load layouts for the new page:', error)
    }
  }

  async function deletePage(item) {
    state.pending = true
    try {
      await $ComposeAPI.pageDelete(item)
      removeFromSet([item])
      // cascade/rebase deletes/reparents children server-side — local tree is stale, refetch
      if (item.strategy === 'cascade' || item.strategy === 'rebase') {
        const nsID = item.namespaceID || state.namespaceID
        if (nsID) await load({ namespaceID: nsID })
      }
      return true
    } catch (error) {
      console.error('Failed to delete page:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function loadTree({ namespaceID } = {}) {
    const nsID = namespaceID || state.namespaceID
    if (!nsID) return Promise.reject(new Error('namespaceID required'))

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
    try {
      await $ComposeAPI.pageReorder({ namespaceID, selfID, pageIDs })
    } catch (error) {
      console.error('Failed to reorder pages:', error)
      throw error
    }
  }

  async function reparent({ namespaceID, pageID, newParentID, siblingPageIDs } = {}) {
    try {
      const page = getByID.value(pageID)
      if (page) {
        await $ComposeAPI.pageUpdate({ ...page, namespaceID, pageID, selfID: newParentID })
      }
      if (siblingPageIDs?.length) {
        await reorder({ namespaceID, selfID: newParentID, pageIDs: siblingPageIDs })
      }
    } catch (error) {
      console.error('Failed to reparent page:', error)
      throw error
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

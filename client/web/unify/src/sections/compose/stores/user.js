import { system } from '@planetcrust/human-js'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

export const useUserStore = defineStore('user', () => {
  const $SystemAPI = inject('$SystemAPI')

  const state = reactive({
    pending: false,
    set: [],
  })

  const findByID = computed(() => {
    return ID => state.set.find(({ userID }) => ID === userID)
  })

  const findByUsername = computed(() => {
    return username => state.set.filter(user => user.username === username)[0] || undefined
  })

  // Actions
  async function load(filter) {
    if (!$SystemAPI) {
      const error = new Error('SystemAPI not available via inject')
      console.error('Failed to load users:', error)
      return Promise.reject(error)
    }

    // Cache-guard: the bulk set is preloaded once at boot; skip refetching on
    // later calls unless explicitly forced (resolveUsers handles unknown IDs).
    if (!filter?.force && state.set.length > 0) {
      return Promise.resolve(state.set)
    }

    state.pending = true

    try {
      const { set: userSet } = await $SystemAPI.userList(filter)
      updateSet(userSet)
      return userSet
    } catch (error) {
      console.error('Failed to load users:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  async function fetchUsers(userID) {
    if (!$SystemAPI) {
      const error = new Error('SystemAPI not available via inject')
      console.error('Failed to fetch users:', error)
      return Promise.reject(error)
    }

    state.pending = true

    if (userID.length === 0) {
      state.pending = false
      return null
    }

    try {
      const { set: userSet } = await $SystemAPI.userList({ userID })
      updateSet(userSet)
      return userSet
    } catch (error) {
      console.error('Failed to fetch users:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  /**
   * Similar to fetchUsers but it only fetches unknown (not in set) ids
   */
  async function resolveUsers(list) {
    if (list.length === 0) {
      // save ourselves some work
      return
    }

    // exclude existing & make unique
    const existing = new Set(state.set.map(({ userID }) => userID))
    const filteredList = [...new Set(list.filter(userID => userID && !existing.has(userID)))]

    if (filteredList.length === 0) {
      // Check for values again
      return
    }

    if (!$SystemAPI) {
      const error = new Error('SystemAPI not available via inject')
      console.error('Failed to resolve users:', error)
      return Promise.reject(error)
    }

    state.pending = true

    try {
      const { set: userSet } = await $SystemAPI.userList({ userID: filteredList })
      updateSet(userSet)
      return userSet
    } catch (error) {
      console.error('Failed to resolve users:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  // Helper function for state management
  function updateSet(newSet) {
    const userSet = (Array.isArray(newSet) ? newSet : [newSet])
      .filter(u => !!u)
      .map(i => new system.User(i))

    if (state.set.length === 0) {
      state.set = userSet
      return
    }

    userSet.forEach(newItem => {
      const oldIndex = state.set.findIndex(({ userID }) => userID === newItem.userID)
      if (oldIndex > -1) {
        state.set.splice(oldIndex, 1, newItem)
      } else {
        state.set.push(newItem)
      }
    })
  }

  /**
   * Store pre-fetched users into the set (e.g. from search results), or upsert
   * after an edit elsewhere (e.g. admin) so the shared cache stays fresh.
   */
  function storeUsers(users) {
    updateSet(users)
  }

  /**
   * Drop users from the cache (e.g. after a delete in another section).
   */
  function removeUsers(ids) {
    const idList = (Array.isArray(ids) ? ids : [ids]).map(String).filter(Boolean)
    if (!idList.length) return
    state.set = state.set.filter(u => !idList.includes(String(u.userID)))
  }

  return {
    // state
    pending: toRef(state, 'pending'),
    set: toRef(state, 'set'),

    // getters
    findByID,
    findByUsername,

    // actions
    load,
    fetchUsers,
    resolveUsers,
    storeUsers,
    removeUsers,
  }
})

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

  async function load(filter) {
    state.pending = true
    try {
      const { set: userSet = [] } = await $SystemAPI.userList(filter)
      updateSet(userSet)
      return userSet
    } catch (error) {
      console.error('Failed to load users:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

  // Only fetches IDs not already in the cache — for batch resolution of encountered user refs.
  async function resolveUsers(list) {
    if (!list?.length) return

    const existing = new Set(state.set.map(({ userID }) => userID))
    const missing = [...new Set(list.filter(userID => userID && !existing.has(userID)))]

    if (missing.length === 0) return

    state.pending = true
    try {
      const { set: userSet = [] } = await $SystemAPI.userList({ userID: missing })
      updateSet(userSet)
      return userSet
    } catch (error) {
      console.error('Failed to resolve users:', error)
      throw error
    } finally {
      state.pending = false
    }
  }

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

  function storeUsers(users) {
    updateSet(users)
  }

  function removeUsers(ids) {
    const idList = (Array.isArray(ids) ? ids : [ids]).map(String).filter(Boolean)
    if (!idList.length) return
    state.set = state.set.filter(u => !idList.includes(String(u.userID)))
  }

  return {
    pending: toRef(state, 'pending'),
    set: toRef(state, 'set'),

    findByID,
    findByUsername,

    load,
    resolveUsers,
    storeUsers,
    removeUsers,
  }
})

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

  // Lookups not yet answered, by userID: a second caller waits for the first
  // rather than asking again, and IDs asked for in the same tick share one call.
  const inflightUsers = new Map() // userID -> Promise
  let pendingUsers = null // { ids, promise, resolve, reject } for the next flush

  async function flushUsers() {
    const batch = pendingUsers
    pendingUsers = null
    const ids = [...batch.ids]

    state.pending = true
    try {
      const { set: userSet = [] } = await $SystemAPI.userList({ userID: ids })
      updateSet(userSet)
      batch.resolve()
    } catch (error) {
      console.error('Failed to resolve users:', error)
      batch.reject(error)
    } finally {
      ids.forEach(id => inflightUsers.delete(id))
      state.pending = false
    }
  }

  // Fetches only IDs neither cached nor already being fetched — for batch
  // resolution of encountered user refs. Rejects when the lookup fails.
  async function resolveUsers(list) {
    if (!list?.length) return

    const existing = new Set(state.set.map(({ userID }) => userID))
    const wanted = [...new Set(list.filter(userID => userID && !existing.has(userID)))]

    if (wanted.length === 0) return

    const missing = wanted.filter(userID => !inflightUsers.has(userID))
    if (missing.length) {
      if (!pendingUsers) {
        const batch = { ids: new Set() }
        batch.promise = new Promise((resolve, reject) => {
          batch.resolve = resolve
          batch.reject = reject
        })
        pendingUsers = batch
        queueMicrotask(flushUsers)
      }

      missing.forEach(userID => {
        pendingUsers.ids.add(userID)
        inflightUsers.set(userID, pendingUsers.promise)
      })
    }

    await Promise.all(wanted.map(userID => inflightUsers.get(userID)))
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

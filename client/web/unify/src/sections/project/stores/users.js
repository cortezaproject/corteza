import { defineStore } from 'pinia'
import { computed, inject, ref } from 'vue'

// Real user directory for picking project teams (replaces the old mock list).
// Loaded once from the system API; `currentUserID` comes from the session.
export const useProjectUsersStore = defineStore('project-users', () => {
  const $SystemAPI = inject('$SystemAPI')
  const $Auth = inject('$Auth', inject('$auth', {}))

  const users = ref([])
  const loaded = ref(false)
  let loading = null

  const currentUserID = computed(() => String($Auth?.user?.userID || ''))

  async function load() {
    if (loaded.value) return
    if (loading) return loading
    loading = $SystemAPI
      .userList({ limit: 500, sort: 'name' })
      .then(({ set = [] } = {}) => {
        users.value = set.map(u => ({
          id: String(u.userID),
          name: u.name || u.username || u.email || u.handle,
          email: u.email || '',
        }))
        loaded.value = true
      })
      .catch(err => {
        console.error('Failed to load users', err)
      })
      .finally(() => {
        loading = null
      })
    return loading
  }

  // Force a re-fetch of the directory (e.g. after creating a new user) so newly
  // added users resolve to a name/email instead of a bare ID.
  async function reload() {
    loaded.value = false
    return load()
  }

  const findUser = id => users.value.find(u => u.id === String(id))

  const userName = id => findUser(id)?.name || String(id || '')

  const userInitials = id => {
    const name = findUser(id)?.name || ''
    return (
      name
        .split(/\s+/)
        .filter(Boolean)
        .slice(0, 2)
        .map(p => p[0]?.toUpperCase())
        .join('') || '?'
    )
  }

  return { users, loaded, load, reload, currentUserID, findUser, userName, userInitials }
})

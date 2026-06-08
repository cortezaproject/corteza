import { inject } from 'vue'
import { useUserStore } from '../stores/useUserStore'

interface UserLike {
  userID: string
  name?: string
  handle?: string
  email?: string
}

export function useUserResolver() {
  const $SystemAPI = inject<any>('$SystemAPI', null)
  const userStore = useUserStore()

  function formatUser(user: UserLike | null | undefined): string {
    if (!user) return ''
    return user.name || user.handle || user.email || user.userID || ''
  }

  function findCached(userID: string): UserLike | null {
    if (!userID) return null
    return userStore.findByID(userID) || null
  }

  async function resolveUser(userID: string): Promise<UserLike | null> {
    if (!userID) return null

    const cached = findCached(userID)
    if (cached) return cached

    if (!$SystemAPI) return null

    try {
      const user = await $SystemAPI.userRead({ userID })
      userStore.storeUsers([user])
      return user
    } catch {
      return null
    }
  }

  async function resolveUsers(userIDs: string[]): Promise<void> {
    if (!userIDs?.length) return
    await userStore.resolveUsers(userIDs)
  }

  function cacheUsers(users: UserLike[]): void {
    if (users.length) {
      userStore.storeUsers(users)
    }
  }

  return {
    formatUser,
    findCached,
    resolveUser,
    resolveUsers,
    cacheUsers,
  }
}

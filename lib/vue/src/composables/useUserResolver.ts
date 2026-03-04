import { inject } from 'vue'

interface UserLike {
  userID: string
  name?: string
  handle?: string
  email?: string
}

interface UserStore {
  findByID: { (id: string): UserLike | undefined }
  storeUsers: { (users: UserLike[]): void }
  resolveUsers: { (ids: string[]): Promise<void> }
}

export function useUserResolver() {
  const $SystemAPI = inject<any>('$SystemAPI', null)
  const $userStore = inject<UserStore | null>('$userStore', null)

  /**
   * Format a user object into a display string.
   */
  function formatUser(user: UserLike | null | undefined): string {
    if (!user) return ''
    return user.name || user.handle || user.email || user.userID || ''
  }

  /**
   * Synchronously look up a user from the store cache.
   * Returns null if not cached or store is unavailable.
   */
  function findCached(userID: string): UserLike | null {
    if (!userID || !$userStore) return null
    return $userStore.findByID(userID) || null
  }

  /**
   * Resolve a single user by ID — cache-first, then API fallback.
   */
  async function resolveUser(userID: string): Promise<UserLike | null> {
    if (!userID) return null

    const cached = findCached(userID)
    if (cached) return cached

    if (!$SystemAPI) return null

    try {
      const user = await $SystemAPI.userRead({ userID })
      if ($userStore) $userStore.storeUsers([user])
      return user
    } catch {
      return null
    }
  }

  /**
   * Resolve multiple user IDs via the store's batch resolution.
   * Only fetches IDs not already cached.
   */
  async function resolveUsers(userIDs: string[]): Promise<void> {
    if (!userIDs?.length || !$userStore) return
    await $userStore.resolveUsers(userIDs)
  }

  /**
   * Cache a list of already-fetched users into the store.
   */
  function cacheUsers(users: UserLike[]): void {
    if (users.length && $userStore) {
      $userStore.storeUsers(users)
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

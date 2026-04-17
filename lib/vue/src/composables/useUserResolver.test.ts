import { describe, it, expect, vi, beforeEach } from 'vitest'
import { defineComponent } from 'vue'
import { flushPromises } from '@vue/test-utils'
import { mountWithContext, makeUser, createMockSystemAPI } from '@planetcrust/human-test-utils'
import { useUserResolver } from './useUserResolver'

function makeUserStore(users: ReturnType<typeof makeUser>[] = []) {
  const cache = new Map(users.map(u => [u.userID, u]))
  return {
    findByID: (id: string) => cache.get(id),
    storeUsers: (us: typeof users) => us.forEach(u => cache.set(u.userID, u)),
    resolveUsers: vi.fn().mockResolvedValue(undefined),
  }
}

function mountResolver(ctx: Parameters<typeof mountWithContext>[1] = {}) {
  let resolver: ReturnType<typeof useUserResolver> | undefined

  const Comp = defineComponent({
    setup() {
      resolver = useUserResolver()
      return {}
    },
    template: '<div/>',
  })

  mountWithContext(Comp, ctx)
  return resolver!
}

describe('useUserResolver', () => {
  describe('formatUser', () => {
    it('returns empty string for null', () => {
      const { formatUser } = mountResolver()
      expect(formatUser(null)).toBe('')
    })

    it('prefers name over handle and email', () => {
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: '1', name: 'Alice', handle: 'alice', email: 'a@b.com' })).toBe('Alice')
    })

    it('falls back to handle when no name', () => {
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: '1', handle: 'alice' })).toBe('alice')
    })

    it('falls back to email when no name or handle', () => {
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: '1', email: 'a@b.com' })).toBe('a@b.com')
    })

    it('falls back to userID as last resort', () => {
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: 'uid-42' })).toBe('uid-42')
    })
  })

  describe('findCached', () => {
    it('returns null when no userStore', () => {
      const { findCached } = mountResolver({ userStore: null })
      expect(findCached('1')).toBeNull()
    })

    it('returns null for empty userID', () => {
      const { findCached } = mountResolver({ userStore: makeUserStore() })
      expect(findCached('')).toBeNull()
    })

    it('returns cached user when present', () => {
      const user = makeUser({ userID: 'u1', name: 'Alice' })
      const { findCached } = mountResolver({ userStore: makeUserStore([user]) })
      expect(findCached('u1')).toEqual(user)
    })

    it('returns null when user not in cache', () => {
      const { findCached } = mountResolver({ userStore: makeUserStore() })
      expect(findCached('missing')).toBeNull()
    })
  })

  describe('resolveUser', () => {
    it('returns null for empty userID', async () => {
      const { resolveUser } = mountResolver()
      expect(await resolveUser('')).toBeNull()
    })

    it('returns cached user without hitting API', async () => {
      const user = makeUser({ userID: 'u1' })
      const api = createMockSystemAPI()
      const { resolveUser } = mountResolver({ systemAPI: api, userStore: makeUserStore([user]) })
      const result = await resolveUser('u1')
      expect(result).toEqual(user)
      expect(api.userRead).not.toHaveBeenCalled()
    })

    it('fetches from API when not cached', async () => {
      const user = makeUser({ userID: 'u2', name: 'Bob' })
      const api = createMockSystemAPI({ userRead: vi.fn().mockResolvedValue(user) })
      const store = makeUserStore()
      const { resolveUser } = mountResolver({ systemAPI: api, userStore: store })

      const result = await resolveUser('u2')
      await flushPromises()

      expect(api.userRead).toHaveBeenCalledWith({ userID: 'u2' })
      expect(result).toEqual(user)
    })

    it('caches fetched user into store', async () => {
      const user = makeUser({ userID: 'u3' })
      const api = createMockSystemAPI({ userRead: vi.fn().mockResolvedValue(user) })
      const store = makeUserStore()
      const { resolveUser, findCached } = mountResolver({ systemAPI: api, userStore: store })

      await resolveUser('u3')
      expect(findCached('u3')).toEqual(user)
    })

    it('returns null when API not available', async () => {
      const { resolveUser } = mountResolver({ systemAPI: null })
      expect(await resolveUser('u1')).toBeNull()
    })

    it('returns null when API throws', async () => {
      const api = createMockSystemAPI({ userRead: vi.fn().mockRejectedValue(new Error('fail')) })
      const { resolveUser } = mountResolver({ systemAPI: api })
      expect(await resolveUser('u1')).toBeNull()
    })
  })

  describe('resolveUsers (batch)', () => {
    it('calls store.resolveUsers with the given IDs', async () => {
      const store = makeUserStore()
      const { resolveUsers } = mountResolver({ userStore: store })
      await resolveUsers(['a', 'b'])
      expect(store.resolveUsers).toHaveBeenCalledWith(['a', 'b'])
    })

    it('no-ops when userStore is null', async () => {
      const { resolveUsers } = mountResolver({ userStore: null })
      await expect(resolveUsers(['a'])).resolves.toBeUndefined()
    })

    it('no-ops when userIDs is empty', async () => {
      const store = makeUserStore()
      const { resolveUsers } = mountResolver({ userStore: store })
      await resolveUsers([])
      expect(store.resolveUsers).not.toHaveBeenCalled()
    })
  })
})

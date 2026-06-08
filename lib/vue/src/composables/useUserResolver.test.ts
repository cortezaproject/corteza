import { describe, it, expect, vi, beforeEach } from 'vitest'
import { defineComponent } from 'vue'
import { flushPromises } from '@vue/test-utils'
import {
  mountWithContext,
  makeUser,
  createMockSystemAPI,
  createTestPinia,
} from '@planetcrust/human-test-utils'
import { useUserResolver } from './useUserResolver'
import { useUserStore } from '../stores/useUserStore'

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
      createTestPinia()
      const { formatUser } = mountResolver()
      expect(formatUser(null)).toBe('')
    })

    it('prefers name over handle and email', () => {
      createTestPinia()
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: '1', name: 'Alice', handle: 'alice', email: 'a@b.com' })).toBe('Alice')
    })

    it('falls back to handle when no name', () => {
      createTestPinia()
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: '1', handle: 'alice' })).toBe('alice')
    })

    it('falls back to email when no name or handle', () => {
      createTestPinia()
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: '1', email: 'a@b.com' })).toBe('a@b.com')
    })

    it('falls back to userID as last resort', () => {
      createTestPinia()
      const { formatUser } = mountResolver()
      expect(formatUser({ userID: 'uid-42' })).toBe('uid-42')
    })
  })

  describe('findCached', () => {
    it('returns null for empty userID', () => {
      createTestPinia()
      const { findCached } = mountResolver()
      expect(findCached('')).toBeNull()
    })

    it('returns cached user when present', () => {
      const api = createMockSystemAPI()
      createTestPinia({ '$SystemAPI': api })
      const user = makeUser({ userID: '90001', name: 'Alice' })
      useUserStore().storeUsers([user])

      const { findCached } = mountResolver({ systemAPI: api })
      expect(findCached('90001')?.name).toBe('Alice')
    })

    it('returns null when user not in cache', () => {
      const api = createMockSystemAPI()
      createTestPinia({ '$SystemAPI': api })

      const { findCached } = mountResolver({ systemAPI: api })
      expect(findCached('90099')).toBeNull()
    })
  })

  describe('resolveUser', () => {
    it('returns null for empty userID', async () => {
      createTestPinia()
      const { resolveUser } = mountResolver()
      expect(await resolveUser('')).toBeNull()
    })

    it('returns cached user without hitting API', async () => {
      const api = createMockSystemAPI()
      createTestPinia({ '$SystemAPI': api })
      const user = makeUser({ userID: '90001' })
      useUserStore().storeUsers([user])

      const { resolveUser } = mountResolver({ systemAPI: api })
      const result = await resolveUser('90001')
      expect(result?.userID).toBe('90001')
      expect(api.userRead).not.toHaveBeenCalled()
    })

    it('fetches from API when not cached', async () => {
      const user = makeUser({ userID: '90002', name: 'Bob' })
      const api = createMockSystemAPI({ userRead: vi.fn().mockResolvedValue(user) })
      createTestPinia({ '$SystemAPI': api })

      const { resolveUser } = mountResolver({ systemAPI: api })
      const result = await resolveUser('90002')
      await flushPromises()

      expect(api.userRead).toHaveBeenCalledWith({ userID: '90002' })
      expect(result?.userID).toBe('90002')
    })

    it('caches fetched user into store', async () => {
      const user = makeUser({ userID: '90003' })
      const api = createMockSystemAPI({ userRead: vi.fn().mockResolvedValue(user) })
      createTestPinia({ '$SystemAPI': api })

      const { resolveUser, findCached } = mountResolver({ systemAPI: api })
      await resolveUser('90003')
      expect(findCached('90003')?.userID).toBe('90003')
    })

    it('returns null when API not available', async () => {
      createTestPinia()
      const { resolveUser } = mountResolver({ systemAPI: null })
      expect(await resolveUser('u1')).toBeNull()
    })

    it('returns null when API throws', async () => {
      const api = createMockSystemAPI({ userRead: vi.fn().mockRejectedValue(new Error('fail')) })
      createTestPinia({ '$SystemAPI': api })

      const { resolveUser } = mountResolver({ systemAPI: api })
      expect(await resolveUser('u1')).toBeNull()
    })
  })

  describe('resolveUsers (batch)', () => {
    it('calls API userList with the given IDs', async () => {
      const api = createMockSystemAPI()
      createTestPinia({ '$SystemAPI': api })

      const { resolveUsers } = mountResolver({ systemAPI: api })
      await resolveUsers(['a', 'b'])
      await flushPromises()

      expect(api.userList).toHaveBeenCalledWith({ userID: ['a', 'b'] })
    })

    it('no-ops when userIDs is empty', async () => {
      const api = createMockSystemAPI()
      createTestPinia({ '$SystemAPI': api })

      const { resolveUsers } = mountResolver({ systemAPI: api })
      await resolveUsers([])
      expect(api.userList).not.toHaveBeenCalled()
    })

    it('skips IDs already in cache', async () => {
      const api = createMockSystemAPI()
      createTestPinia({ '$SystemAPI': api })
      const user = makeUser({ userID: '90004' })
      useUserStore().storeUsers([user])

      const { resolveUsers } = mountResolver({ systemAPI: api })
      await resolveUsers(['90004'])
      expect(api.userList).not.toHaveBeenCalled()
    })
  })
})

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia } from 'pinia'
import { createTestPinia } from '@planetcrust/human-test-utils'
import { useRBACStore } from './useRBAC'

function makeAPI(rules: Array<{ resource: string; operation: string; allow: boolean }> = []) {
  return { permissionsEffective: vi.fn().mockResolvedValue(rules) }
}

describe('useRBACStore', () => {
  beforeEach(() => {
    createTestPinia()
  })

  describe('initial state', () => {
    it('starts not loaded with empty rules', () => {
      const store = useRBACStore()
      expect(store.loaded).toBe(false)
      expect(store.rules).toHaveLength(0)
    })
  })

  describe('load()', () => {
    it('marks loaded after fetch', async () => {
      const store = useRBACStore()
      await store.load([makeAPI()])
      expect(store.loaded).toBe(true)
    })

    it('stores rules from all provided APIs', async () => {
      const store = useRBACStore()
      const api1 = makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }])
      const api2 = makeAPI([{ resource: 'corteza::compose/namespace', operation: 'create', allow: false }])
      await store.load([api1, api2])
      expect(store.rules).toHaveLength(2)
    })

    it('filters out rules not starting with corteza:: prefix', async () => {
      const store = useRBACStore()
      const api = makeAPI([
        { resource: 'corteza::system/user', operation: 'read', allow: true },
        { resource: 'other::something', operation: 'read', allow: true },
      ])
      await store.load([api])
      expect(store.rules).toHaveLength(1)
      expect(store.rules[0].resource).toBe('corteza::system/user')
    })

    it('ignores failed API calls (returns empty for that API)', async () => {
      const store = useRBACStore()
      const good = makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }])
      const bad = { permissionsEffective: vi.fn().mockRejectedValue(new Error('fail')) }
      await store.load([good, bad])
      expect(store.rules).toHaveLength(1)
    })

    it('resets loaded=false at start of load', async () => {
      const store = useRBACStore()
      await store.load([makeAPI()])
      expect(store.loaded).toBe(true)

      let loadedDuringFetch = true
      const slowAPI = {
        permissionsEffective: vi.fn().mockImplementation(async () => {
          loadedDuringFetch = store.loaded
          return []
        }),
      }
      await store.load([slowAPI])
      expect(loadedDuringFetch).toBe(false)
    })
  })

  describe('can()', () => {
    it('returns true when matching allow rule exists', async () => {
      const store = useRBACStore()
      await store.load([makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }])])
      expect(store.can('system/user', 'read')).toBe(true)
    })

    it('returns false when rule is deny', async () => {
      const store = useRBACStore()
      await store.load([makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: false }])])
      expect(store.can('system/user', 'read')).toBe(false)
    })

    it('returns false when no matching rule', async () => {
      const store = useRBACStore()
      await store.load([makeAPI()])
      expect(store.can('system/user', 'delete')).toBe(false)
    })

    it('is case-sensitive on resource and operation', async () => {
      const store = useRBACStore()
      await store.load([makeAPI([{ resource: 'corteza::system/user', operation: 'READ', allow: true }])])
      expect(store.can('system/user', 'read')).toBe(false)
    })
  })

  describe('clear()', () => {
    it('empties rules and resets loaded', async () => {
      const store = useRBACStore()
      await store.load([makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }])])
      store.clear()
      expect(store.rules).toHaveLength(0)
      expect(store.loaded).toBe(false)
    })
  })
})

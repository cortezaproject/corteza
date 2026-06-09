import { describe, it, expect, vi } from 'vitest'
import { createTestPinia } from '@planetcrust/human-test-utils'
import { useRBACStore } from './useRBAC'

function makeAPI(rules: Array<{ resource: string; operation: string; allow: boolean }> = []) {
  return { permissionsEffective: vi.fn().mockResolvedValue(rules) }
}

// The store injects $SystemAPI / $AutomationAPI / $ComposeAPI at setup time, so
// tests register them through the throwaway app hosted by createTestPinia.
function setup(provides: Record<string, unknown> = {}) {
  createTestPinia(provides)
  return useRBACStore()
}

describe('useRBACStore', () => {
  describe('initial state', () => {
    it('starts not loaded with empty rules', () => {
      const store = setup()
      expect(store.loaded).toBe(false)
      expect(store.rules).toHaveLength(0)
    })
  })

  describe('load()', () => {
    it('marks loaded after fetch', async () => {
      const store = setup({ $SystemAPI: makeAPI() })
      await store.load()
      expect(store.loaded).toBe(true)
    })

    it('stores rules from all provided APIs', async () => {
      const store = setup({
        $SystemAPI: makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }]),
        $AutomationAPI: makeAPI([{ resource: 'corteza::compose/namespace', operation: 'create', allow: false }]),
      })
      await store.load()
      expect(store.rules).toHaveLength(2)
    })

    it('filters out rules not starting with corteza:: prefix', async () => {
      const store = setup({
        $SystemAPI: makeAPI([
          { resource: 'corteza::system/user', operation: 'read', allow: true },
          { resource: 'other::something', operation: 'read', allow: true },
        ]),
      })
      await store.load()
      expect(store.rules).toHaveLength(1)
      expect(store.rules[0].resource).toBe('corteza::system/user')
    })

    it('ignores failed API calls (returns empty for that API)', async () => {
      const store = setup({
        $SystemAPI: makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }]),
        $AutomationAPI: { permissionsEffective: vi.fn().mockRejectedValue(new Error('fail')) },
      })
      await store.load()
      expect(store.rules).toHaveLength(1)
    })

    it('resets loaded=false at start of load', async () => {
      let loadedDuringFetch = true
      const slowAPI = {
        permissionsEffective: vi.fn().mockImplementation(async () => {
          loadedDuringFetch = store.loaded
          return []
        }),
      }
      const store = setup({ $SystemAPI: slowAPI })
      await store.load()
      expect(store.loaded).toBe(true)

      await store.load()
      expect(loadedDuringFetch).toBe(false)
    })
  })

  describe('can()', () => {
    it('returns true when matching allow rule exists', async () => {
      const store = setup({
        $SystemAPI: makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }]),
      })
      await store.load()
      expect(store.can('system/user', 'read')).toBe(true)
    })

    it('returns false when rule is deny', async () => {
      const store = setup({
        $SystemAPI: makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: false }]),
      })
      await store.load()
      expect(store.can('system/user', 'read')).toBe(false)
    })

    it('returns false when no matching rule', async () => {
      const store = setup({ $SystemAPI: makeAPI() })
      await store.load()
      expect(store.can('system/user', 'delete')).toBe(false)
    })

    it('is case-sensitive on resource and operation', async () => {
      const store = setup({
        $SystemAPI: makeAPI([{ resource: 'corteza::system/user', operation: 'READ', allow: true }]),
      })
      await store.load()
      expect(store.can('system/user', 'read')).toBe(false)
    })
  })

  describe('clear()', () => {
    it('empties rules and resets loaded', async () => {
      const store = setup({
        $SystemAPI: makeAPI([{ resource: 'corteza::system/user', operation: 'read', allow: true }]),
      })
      await store.load()
      store.clear()
      expect(store.rules).toHaveLength(0)
      expect(store.loaded).toBe(false)
    })
  })
})

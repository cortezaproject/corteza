import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createTestPinia, createMockComposeAPI, makeNamespace } from '@planetcrust/human-test-utils'
import { useNamespaceStore } from '@planetcrust/human-vue'

describe('useNamespaceStore', () => {
  let api: ReturnType<typeof createMockComposeAPI>

  beforeEach(() => {
    api = createMockComposeAPI()
    createTestPinia({ '$ComposeAPI': api })
  })

  describe('load()', () => {
    it('calls namespaceList and populates set', async () => {
      const ns = makeNamespace({ namespaceID: '100', name: 'Main' })
      api.namespaceList.mockResolvedValue({ set: [ns], filter: { total: 1 } })

      const store = useNamespaceStore()
      await store.load()

      expect(api.namespaceList).toHaveBeenCalledWith({})
      expect(store.set).toHaveLength(1)
      expect(store.set[0].namespaceID).toBe('100')
    })

    it('always refetches on every call', async () => {
      const ns1 = makeNamespace({ namespaceID: '200' })
      const ns2 = makeNamespace({ namespaceID: '201' })
      api.namespaceList.mockResolvedValue({ set: [ns1, ns2], filter: {} })

      const store = useNamespaceStore()
      await store.load()
      await store.load()

      expect(api.namespaceList).toHaveBeenCalledTimes(2)
    })
  })

  describe('getByID', () => {
    it('finds namespace by ID after load', async () => {
      const ns = makeNamespace({ namespaceID: '400' })
      api.namespaceList.mockResolvedValue({ set: [ns], filter: {} })

      const store = useNamespaceStore()
      await store.load()

      expect(store.getByID('400')).toBeDefined()
      expect(store.getByID('400')?.namespaceID).toBe('400')
    })

    it('returns undefined for unknown ID', () => {
      const store = useNamespaceStore()
      expect(store.getByID('9999')).toBeUndefined()
    })
  })

  describe('findByID()', () => {
    it('returns cached namespace without API call', async () => {
      const ns = makeNamespace({ namespaceID: '500' })
      api.namespaceList.mockResolvedValue({ set: [ns], filter: {} })

      const store = useNamespaceStore()
      await store.load()

      const result = await store.findByID({ namespaceID: '500' })
      expect(result?.namespaceID).toBe('500')
      expect(api.namespaceRead).not.toHaveBeenCalled()
    })

    it('fetches from API when not cached', async () => {
      const ns = makeNamespace({ namespaceID: '600' })
      api.namespaceRead.mockResolvedValue(ns)

      const store = useNamespaceStore()
      const result = await store.findByID({ namespaceID: '600' })

      expect(api.namespaceRead).toHaveBeenCalledWith({ namespaceID: '600' })
      expect(result?.namespaceID).toBe('600')
    })

    it('returns an independent copy (not the frozen store item)', async () => {
      const ns = makeNamespace({ namespaceID: '550', name: 'Original' })
      api.namespaceList.mockResolvedValue({ set: [ns], filter: {} })

      const store = useNamespaceStore()
      await store.load()

      const result = await store.findByID({ namespaceID: '550' })
      result.name = 'Modified'

      expect(store.getByID('550')?.name).toBe('Original')
    })
  })

  describe('create()', () => {
    it('calls namespaceCreate and adds result to set', async () => {
      const ns = makeNamespace({ namespaceID: '700', name: 'New' })
      api.namespaceCreate.mockResolvedValue(ns)

      const store = useNamespaceStore()
      const result = await store.create({ name: 'New' })

      expect(api.namespaceCreate).toHaveBeenCalledWith({ name: 'New' })
      expect(result.namespaceID).toBe('700')
      expect(store.getByID('700')).toBeDefined()
    })
  })

  describe('update()', () => {
    it('calls namespaceUpdate and replaces item in set', async () => {
      const ns = makeNamespace({ namespaceID: '800', name: 'Original' })
      api.namespaceList.mockResolvedValue({ set: [ns], filter: {} })
      const updated = makeNamespace({ namespaceID: '800', name: 'Updated' })
      api.namespaceUpdate.mockResolvedValue(updated)

      const store = useNamespaceStore()
      await store.load()
      await store.update({ namespaceID: '800', name: 'Updated' })

      expect(store.getByID('800')?.name).toBe('Updated')
    })
  })

  describe('delete()', () => {
    it('removes namespace from set after delete', async () => {
      const ns = makeNamespace({ namespaceID: '900' })
      api.namespaceList.mockResolvedValue({ set: [ns], filter: {} })
      api.namespaceDelete.mockResolvedValue({})

      const store = useNamespaceStore()
      await store.load()
      expect(store.set).toHaveLength(1)

      await store.delete({ namespaceID: '900' })
      expect(store.set).toHaveLength(0)
    })
  })

  describe('clearSet()', () => {
    it('empties the set', async () => {
      const ns = makeNamespace({ namespaceID: '1000' })
      api.namespaceList.mockResolvedValue({ set: [ns], filter: {} })

      const store = useNamespaceStore()
      await store.load()
      store.clearSet()
      expect(store.set).toHaveLength(0)
    })
  })
})

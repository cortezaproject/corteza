import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createTestPinia, createMockComposeAPI, makeModule } from '@planetcrust/human-test-utils'
import { useModuleStore } from '@planetcrust/human-vue'

describe('useModuleStore', () => {
  const NS = '10001'
  let api: ReturnType<typeof createMockComposeAPI>

  beforeEach(() => {
    api = createMockComposeAPI()
    createTestPinia({ '$ComposeAPI': api })
  })

  describe('load()', () => {
    it('calls moduleList and populates set', async () => {
      const mod = makeModule({ moduleID: '1', namespaceID: NS, name: 'Contacts' })
      api.moduleList.mockResolvedValue({ set: [mod], filter: { total: 1 } })

      const store = useModuleStore()
      await store.load({ namespaceID: NS })

      expect(api.moduleList).toHaveBeenCalledWith({ namespaceID: NS, sort: 'name ASC' })
      expect(store.set).toHaveLength(1)
      expect(store.set[0].moduleID).toBe('1')
    })

    it('always refetches on every call', async () => {
      api.moduleList.mockResolvedValue({ set: [makeModule({ namespaceID: NS })], filter: {} })

      const store = useModuleStore()
      await store.load({ namespaceID: NS })
      await store.load({ namespaceID: NS })

      expect(api.moduleList).toHaveBeenCalledTimes(2)
    })
  })

  describe('loadFor()', () => {
    it('fetches on first call', async () => {
      api.moduleList.mockResolvedValue({ set: [makeModule({ namespaceID: NS })], filter: {} })

      const store = useModuleStore()
      await store.loadFor(NS)

      expect(api.moduleList).toHaveBeenCalledTimes(1)
    })

    it('cache-first: skips API on second call', async () => {
      api.moduleList.mockResolvedValue({ set: [makeModule({ namespaceID: NS })], filter: {} })

      const store = useModuleStore()
      await store.loadFor(NS)
      await store.loadFor(NS)

      expect(api.moduleList).toHaveBeenCalledTimes(1)
    })
  })

  describe('findByID()', () => {
    it('returns cached module without API call', async () => {
      const mod = makeModule({ moduleID: '50', namespaceID: NS })
      api.moduleList.mockResolvedValue({ set: [mod], filter: {} })

      const store = useModuleStore()
      await store.load({ namespaceID: NS })

      const result = await store.findByID({ namespaceID: NS, moduleID: '50' })
      expect(result?.moduleID).toBe('50')
      expect(api.moduleRead).not.toHaveBeenCalled()
    })

    it('fetches from API when not cached', async () => {
      const mod = makeModule({ moduleID: '60', namespaceID: NS })
      api.moduleRead.mockResolvedValue(mod)

      const store = useModuleStore()
      const result = await store.findByID({ namespaceID: NS, moduleID: '60' })

      expect(api.moduleRead).toHaveBeenCalledWith({ namespaceID: NS, moduleID: '60' })
      expect(result?.moduleID).toBe('60')
    })

    it('returns an independent copy (not the frozen store item)', async () => {
      const mod = makeModule({ moduleID: '70', namespaceID: NS, name: 'Original' })
      api.moduleList.mockResolvedValue({ set: [mod], filter: {} })

      const store = useModuleStore()
      await store.load({ namespaceID: NS })

      const result = await store.findByID({ namespaceID: NS, moduleID: '70' })
      result.name = 'Modified'

      expect(store.getByID('70')?.name).toBe('Original')
    })
  })

  describe('create()', () => {
    it('calls moduleCreate and adds result to cache', async () => {
      const mod = makeModule({ moduleID: '100', namespaceID: NS, name: 'New' })
      api.moduleCreate.mockResolvedValue(mod)

      const store = useModuleStore()
      store['state'] // ensure store active
      const result = await store.create({ namespaceID: NS, name: 'New', fields: [] })

      expect(api.moduleCreate).toHaveBeenCalled()
      expect(result.moduleID).toBe('100')
    })
  })

  describe('update()', () => {
    it('calls moduleUpdate and replaces item in cache', async () => {
      const mod = makeModule({ moduleID: '200', namespaceID: NS, name: 'Original' })
      api.moduleList.mockResolvedValue({ set: [mod], filter: {} })
      const updated = makeModule({ moduleID: '200', namespaceID: NS, name: 'Updated' })
      api.moduleUpdate.mockResolvedValue(updated)

      const store = useModuleStore()
      await store.load({ namespaceID: NS })
      await store.update({ moduleID: '200', namespaceID: NS, name: 'Updated' })

      expect(store.getByID('200')?.name).toBe('Updated')
    })
  })

  describe('delete()', () => {
    it('removes module from cache after delete', async () => {
      const mod = makeModule({ moduleID: '300', namespaceID: NS })
      api.moduleList.mockResolvedValue({ set: [mod], filter: {} })
      api.moduleDelete.mockResolvedValue({})

      const store = useModuleStore()
      await store.load({ namespaceID: NS })
      expect(store.set).toHaveLength(1)

      await store.delete({ moduleID: '300', namespaceID: NS })
      expect(store.set).toHaveLength(0)
    })
  })

  describe('clearSet()', () => {
    it('empties the active namespace cache', async () => {
      const mod = makeModule({ moduleID: '400', namespaceID: NS })
      api.moduleList.mockResolvedValue({ set: [mod], filter: {} })

      const store = useModuleStore()
      await store.load({ namespaceID: NS })
      store.clearSet()
      expect(store.set).toHaveLength(0)
    })
  })
})

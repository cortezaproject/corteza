import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createTestPinia } from '@planetcrust/human-test-utils'
import { useRecordStore, useModuleStore } from '@planetcrust/human-vue'

vi.mock('@planetcrust/human-js', () => ({
  automation: {
    Function: class MockFunction {
      constructor(data: any) {
        Object.assign(this, data)
      }
    },
  },
  compose: {
    Record: class MockRecord {
      recordID: string
      namespaceID: string
      moduleID: string
      values: any
      ownedBy: string
      meta: any
      constructor(_mod: any, data: any) {
        this.recordID = data.recordID || ''
        this.namespaceID = data.namespaceID || ''
        this.moduleID = data.moduleID || ''
        this.values = data.values || {}
        this.ownedBy = data.ownedBy || ''
        this.meta = data.meta || {}
      }
      serializeValues() {
        return []
      }
    },
    Module: class MockModule {
      moduleID: string
      constructor(data: any) {
        Object.assign(this, data)
      }
    },
  },
}))

const NS_ID = '10001'
const MOD_ID = '20001'
const MOD = { moduleID: MOD_ID, namespaceID: NS_ID, fields: [] }

function makeAPI(overrides: Record<string, any> = {}) {
  return {
    recordList: vi.fn().mockResolvedValue({ set: [], filter: { total: 0 } }),
    recordRead: vi
      .fn()
      .mockResolvedValue({ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }),
    recordCreate: vi
      .fn()
      .mockResolvedValue({ recordID: '30002', namespaceID: NS_ID, moduleID: MOD_ID }),
    recordUpdate: vi
      .fn()
      .mockResolvedValue({ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }),
    recordDelete: vi.fn().mockResolvedValue({}),
    ...overrides,
  }
}

function setup(api = makeAPI()) {
  createTestPinia({ $ComposeAPI: api })
  const moduleStore = useModuleStore()
  moduleStore.updateSet([MOD])
  return { store: useRecordStore(), moduleStore, api }
}

function setupNoModule(api = makeAPI()) {
  createTestPinia({ $ComposeAPI: api })
  return { store: useRecordStore(), api }
}

describe('useRecordStore', () => {
  beforeEach(() => {
    createTestPinia({})
  })

  describe('list()', () => {
    it('throws when moduleID missing', async () => {
      const { store } = setup()
      await expect(store.list({ namespaceID: NS_ID } as any)).rejects.toThrow()
    })

    it('throws when module not in store', async () => {
      const { store } = setupNoModule()
      await expect(store.list({ namespaceID: NS_ID, moduleID: '99999' })).rejects.toThrow(
        'Module 99999 not found',
      )
    })

    it('populates records cache and returns set', async () => {
      const { store } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [
              { recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID },
              { recordID: '30002', namespaceID: NS_ID, moduleID: MOD_ID },
            ],
            filter: { total: 2 },
          }),
        }),
      )
      const { set } = await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })

      expect(set).toHaveLength(2)
      expect(store.getByID('30001')).toBeTruthy()
      expect(store.getByID('30002')).toBeTruthy()
    })

    it('resets loading to false after fetch', async () => {
      const { store } = setup()
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      expect(store.loading).toBe(false)
    })
  })

  describe('findByID()', () => {
    it('returns cached record without API call', async () => {
      const { store, api } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
        }),
      )
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      api.recordRead.mockClear()

      const rec = await store.findByID({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001' })
      expect(rec.recordID).toBe('30001')
      expect(api.recordRead).not.toHaveBeenCalled()
    })

    it('fetches from API when not cached', async () => {
      const { store, api } = setup(
        makeAPI({
          recordRead: vi
            .fn()
            .mockResolvedValue({ recordID: '99001', namespaceID: NS_ID, moduleID: MOD_ID }),
        }),
      )

      const rec = await store.findByID({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '99001' })
      expect(api.recordRead).toHaveBeenCalledWith(
        { namespaceID: NS_ID, moduleID: MOD_ID, recordID: '99001' },
        { signal: undefined },
      )
      expect(rec.recordID).toBe('99001')
    })

    it('force bypasses the cache and refetches', async () => {
      const { store, api } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
          recordRead: vi
            .fn()
            .mockResolvedValue({ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }),
        }),
      )
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      api.recordRead.mockClear()

      await store.findByID({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001', force: true })
      expect(api.recordRead).toHaveBeenCalledWith(
        { namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001' },
        { signal: undefined },
      )
    })

    it('forwards an AbortSignal to recordRead for cancellation', async () => {
      const { store, api } = setup(
        makeAPI({
          recordRead: vi
            .fn()
            .mockResolvedValue({ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }),
        }),
      )
      const ac = new AbortController()

      await store.findByID({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordID: '30001',
        signal: ac.signal,
      })
      expect(api.recordRead).toHaveBeenCalledWith(
        { namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001' },
        { signal: ac.signal },
      )
    })

    it('throws when module not in store', async () => {
      const { store } = setupNoModule()
      await expect(
        store.findByID({ namespaceID: NS_ID, moduleID: '99999', recordID: '30001' }),
      ).rejects.toThrow('Module 99999 not found')
    })
  })

  describe('getByID()', () => {
    it('returns null when not in any cache', () => {
      const { store } = setupNoModule()
      expect(store.getByID('unknown')).toBeNull()
    })

    it('returns record from records map', async () => {
      const { store } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
        }),
      )
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      expect(store.getByID('30001')).toBeTruthy()
    })

    it('returns record from labelCache when not in records', async () => {
      const { store } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '40001', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
        }),
      )

      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40001'],
      })
      expect(store.getByID('40001')).toBeTruthy()
    })
  })

  describe('resolveRecordLabels()', () => {
    it('no-ops when empty recordIDs', async () => {
      const { store, api } = setup()
      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: [] })
      expect(api.recordList).not.toHaveBeenCalled()
      expect(api.recordRead).not.toHaveBeenCalled()
    })

    it('skips IDs already in records cache', async () => {
      const { store, api } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
        }),
      )
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      api.recordList.mockClear()
      api.recordRead.mockClear()

      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['30001'],
      })
      expect(api.recordList).not.toHaveBeenCalled()
      expect(api.recordRead).not.toHaveBeenCalled()
    })

    it('batches missing IDs into a single recordList call', async () => {
      const { store, api } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [
              { recordID: '40042', namespaceID: NS_ID, moduleID: MOD_ID },
              { recordID: '40043', namespaceID: NS_ID, moduleID: MOD_ID },
            ],
            filter: {},
          }),
        }),
      )

      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40042', '40043'],
      })

      expect(api.recordList).toHaveBeenCalledTimes(1)
      expect(api.recordList).toHaveBeenCalledWith(
        expect.objectContaining({
          namespaceID: NS_ID,
          moduleID: MOD_ID,
          query: 'recordID = 40042 OR recordID = 40043',
          incTotal: false,
        }),
      )
      expect(api.recordRead).not.toHaveBeenCalled()
      expect(store.getByID('40042')).toBeTruthy()
      expect(store.getByID('40043')).toBeTruthy()
    })

    it('coalesces same-tick requests for the same ID into one call', async () => {
      const { store, api } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '40044', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
        }),
      )

      const p1 = store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40044'],
      })
      const p2 = store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40044'],
      })
      await Promise.all([p1, p2])

      expect(api.recordList).toHaveBeenCalledTimes(1)
      expect(store.getByID('40044')).toBeTruthy()
    })

    it('skips IDs already in labelCache on second call', async () => {
      const { store, api } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '40045', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
        }),
      )

      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40045'],
      })
      api.recordList.mockClear()
      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40045'],
      })
      expect(api.recordList).not.toHaveBeenCalled()
    })

    it('falls back to per-ID read for IDs the batch omits', async () => {
      const { store, api } = setup(
        makeAPI({
          // Batch returns nothing for the requested ID...
          recordList: vi.fn().mockResolvedValue({ set: [], filter: {} }),
          // ...so it should be read individually.
          recordRead: vi
            .fn()
            .mockResolvedValue({ recordID: '40046', namespaceID: NS_ID, moduleID: MOD_ID }),
        }),
      )

      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40046'],
      })

      expect(api.recordList).toHaveBeenCalled()
      expect(api.recordRead).toHaveBeenCalledWith({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordID: '40046',
      })
      expect(store.getByID('40046')).toBeTruthy()
    })

    it('handles batch + fallback failure gracefully (no throw)', async () => {
      const { store } = setup(
        makeAPI({
          recordList: vi.fn().mockRejectedValue(new Error('boom')),
          recordRead: vi.fn().mockRejectedValue(new Error('not found')),
        }),
      )

      await expect(
        store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['40047'] }),
      ).resolves.toBeUndefined()
      expect(store.getByID('40047')).toBeNull()
    })

    it('ignores non-numeric IDs (cannot be batched safely)', async () => {
      const { store, api } = setup()
      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['bad-id'],
      })
      expect(api.recordList).not.toHaveBeenCalled()
      expect(api.recordRead).not.toHaveBeenCalled()
      expect(store.getByID('bad-id')).toBeNull()
    })
  })

  describe('setNavigationIDs() / getNextAndPrev()', () => {
    it('returns prev and next for middle item', () => {
      const { store } = setupNoModule()
      store.setNavigationIDs(['a', 'b', 'c'])
      expect(store.getNextAndPrev('b')).toEqual({ prev: 'a', next: 'c' })
    })

    it('returns undefined prev for first item', () => {
      const { store } = setupNoModule()
      store.setNavigationIDs(['a', 'b'])
      expect(store.getNextAndPrev('a').prev).toBeUndefined()
    })

    it('returns undefined next for last item', () => {
      const { store } = setupNoModule()
      store.setNavigationIDs(['a', 'b'])
      expect(store.getNextAndPrev('b').next).toBeUndefined()
    })

    it('returns both undefined when ID not in list', () => {
      const { store } = setupNoModule()
      store.setNavigationIDs(['a', 'b'])
      expect(store.getNextAndPrev('z')).toEqual({ prev: undefined, next: undefined })
    })
  })

  describe('clearAll()', () => {
    it('clears both caches and resets state', async () => {
      const { store } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '40050', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
        }),
      )

      await store.resolveRecordLabels({
        namespaceID: NS_ID,
        moduleID: MOD_ID,
        recordIDs: ['40050'],
      })
      store.setNavigationIDs(['x', 'y'])
      expect(store.getByID('40050')).toBeTruthy()

      store.clearAll()
      expect(store.getByID('40050')).toBeNull()
      expect(store.loading).toBe(false)
      expect(store.pending).toBe(false)
      expect(store.paginationRecordIDs).toHaveLength(0)
    })
  })

  describe('delete()', () => {
    it('removes record from cache after successful delete', async () => {
      const { store } = setup(
        makeAPI({
          recordList: vi.fn().mockResolvedValue({
            set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
            filter: {},
          }),
          recordDelete: vi.fn().mockResolvedValue({}),
        }),
      )
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      expect(store.getByID('30001')).toBeTruthy()

      await store.delete({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001' })
      expect(store.getByID('30001')).toBeNull()
    })
  })
})

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createTestPinia } from '@planetcrust/human-test-utils'
import { useRecordStore } from './record'
import { useModuleStore } from './module'

vi.mock('@planetcrust/human-js', () => ({
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
      serializeValues() { return [] }
    },
    Module: class MockModule {
      moduleID: string
      constructor(data: any) { Object.assign(this, data) }
    },
  },
}))

const NS_ID = '10001'
const MOD_ID = '20001'
const MOD = { moduleID: MOD_ID, namespaceID: NS_ID, fields: [] }

function makeAPI(overrides: Record<string, any> = {}) {
  return {
    recordList: vi.fn().mockResolvedValue({ set: [], filter: { total: 0 } }),
    recordRead: vi.fn().mockResolvedValue({ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }),
    recordCreate: vi.fn().mockResolvedValue({ recordID: '30002', namespaceID: NS_ID, moduleID: MOD_ID }),
    recordUpdate: vi.fn().mockResolvedValue({ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }),
    recordDelete: vi.fn().mockResolvedValue({}),
    ...overrides,
  }
}

function setup(api?: any) {
  const provides: Record<string, any> = api ? { '$ComposeAPI': api } : {}
  createTestPinia(provides)
  const moduleStore = useModuleStore()
  moduleStore.updateSet([MOD])
  return { store: useRecordStore(), moduleStore }
}

describe('useRecordStore', () => {
  describe('list()', () => {
    it('throws without $ComposeAPI', async () => {
      createTestPinia({})
      const store = useRecordStore()
      await expect(store.list({ namespaceID: NS_ID, moduleID: MOD_ID })).rejects.toThrow('ComposeAPI not available')
    })

    it('throws when moduleID missing', async () => {
      const { store } = setup(makeAPI())
      await expect(store.list({ namespaceID: NS_ID } as any)).rejects.toThrow()
    })

    it('throws when module not in store', async () => {
      createTestPinia({ '$ComposeAPI': makeAPI() })
      const store = useRecordStore()
      await expect(store.list({ namespaceID: NS_ID, moduleID: '99999' })).rejects.toThrow('Module 99999 not found')
    })

    it('populates records cache and returns set', async () => {
      const api = makeAPI({
        recordList: vi.fn().mockResolvedValue({
          set: [
            { recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID },
            { recordID: '30002', namespaceID: NS_ID, moduleID: MOD_ID },
          ],
          filter: { total: 2 },
        }),
      })
      const { store } = setup(api)
      const { set } = await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })

      expect(set).toHaveLength(2)
      expect(store.getByID('30001')).toBeTruthy()
      expect(store.getByID('30002')).toBeTruthy()
    })

    it('resets loading to false after fetch', async () => {
      const { store } = setup(makeAPI())
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      expect(store.loading).toBe(false)
    })
  })

  describe('findByID()', () => {
    it('returns cached record without API call', async () => {
      const api = makeAPI({
        recordList: vi.fn().mockResolvedValue({
          set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
          filter: {},
        }),
      })
      const { store } = setup(api)
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      api.recordRead.mockClear()

      const rec = await store.findByID({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001' })
      expect(rec.recordID).toBe('30001')
      expect(api.recordRead).not.toHaveBeenCalled()
    })

    it('fetches from API when not cached', async () => {
      const api = makeAPI({
        recordRead: vi.fn().mockResolvedValue({ recordID: '99001', namespaceID: NS_ID, moduleID: MOD_ID }),
      })
      const { store } = setup(api)

      const rec = await store.findByID({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '99001' })
      expect(api.recordRead).toHaveBeenCalledWith({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '99001' })
      expect(rec.recordID).toBe('99001')
    })

    it('force=true bypasses cache', async () => {
      const api = makeAPI({
        recordList: vi.fn().mockResolvedValue({
          set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
          filter: {},
        }),
        recordRead: vi.fn().mockResolvedValue({ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }),
      })
      const { store } = setup(api)
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })

      await store.findByID({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001', force: true })
      expect(api.recordRead).toHaveBeenCalled()
    })

    it('throws when module not in store', async () => {
      createTestPinia({ '$ComposeAPI': makeAPI() })
      const store = useRecordStore()
      await expect(
        store.findByID({ namespaceID: NS_ID, moduleID: '99999', recordID: '30001' }),
      ).rejects.toThrow('Module 99999 not found')
    })
  })

  describe('getByID()', () => {
    it('returns null when not in any cache', () => {
      createTestPinia({})
      const store = useRecordStore()
      expect(store.getByID('unknown')).toBeNull()
    })

    it('returns record from records map', async () => {
      const api = makeAPI({
        recordList: vi.fn().mockResolvedValue({
          set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
          filter: {},
        }),
      })
      const { store } = setup(api)
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      expect(store.getByID('30001')).toBeTruthy()
    })

    it('returns record from labelCache when not in records', async () => {
      const api = makeAPI({
        recordRead: vi.fn().mockResolvedValue({ recordID: 'lbl-1', namespaceID: NS_ID, moduleID: MOD_ID }),
      })
      createTestPinia({ '$ComposeAPI': api })
      const store = useRecordStore()

      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['lbl-1'] })
      expect(store.getByID('lbl-1')).toBeTruthy()
    })
  })

  describe('resolveRecordLabels()', () => {
    it('no-ops when empty recordIDs', async () => {
      const api = makeAPI()
      createTestPinia({ '$ComposeAPI': api })
      const store = useRecordStore()
      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: [] })
      expect(api.recordRead).not.toHaveBeenCalled()
    })

    it('skips IDs already in records cache', async () => {
      const api = makeAPI({
        recordList: vi.fn().mockResolvedValue({
          set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
          filter: {},
        }),
      })
      const { store } = setup(api)
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      api.recordRead.mockClear()

      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['30001'] })
      expect(api.recordRead).not.toHaveBeenCalled()
    })

    it('fetches missing IDs into labelCache', async () => {
      const api = makeAPI({
        recordRead: vi.fn().mockResolvedValue({ recordID: 'lbl-42', namespaceID: NS_ID, moduleID: MOD_ID }),
      })
      createTestPinia({ '$ComposeAPI': api })
      const store = useRecordStore()

      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['lbl-42'] })
      expect(api.recordRead).toHaveBeenCalledWith({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: 'lbl-42' })
      expect(store.getByID('lbl-42')).toBeTruthy()
    })

    it('skips IDs already in labelCache on second call', async () => {
      const api = makeAPI({
        recordRead: vi.fn().mockResolvedValue({ recordID: 'lc-1', namespaceID: NS_ID, moduleID: MOD_ID }),
      })
      createTestPinia({ '$ComposeAPI': api })
      const store = useRecordStore()

      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['lc-1'] })
      api.recordRead.mockClear()
      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['lc-1'] })
      expect(api.recordRead).not.toHaveBeenCalled()
    })

    it('handles per-record API errors gracefully (no throw)', async () => {
      const api = makeAPI({
        recordRead: vi.fn().mockRejectedValue(new Error('not found')),
      })
      createTestPinia({ '$ComposeAPI': api })
      const store = useRecordStore()

      await expect(
        store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['bad-id'] }),
      ).resolves.toBeUndefined()
      expect(store.getByID('bad-id')).toBeNull()
    })
  })

  describe('setNavigationIDs() / getNextAndPrev()', () => {
    it('returns prev and next for middle item', () => {
      createTestPinia({})
      const store = useRecordStore()
      store.setNavigationIDs(['a', 'b', 'c'])
      expect(store.getNextAndPrev('b')).toEqual({ prev: 'a', next: 'c' })
    })

    it('returns undefined prev for first item', () => {
      createTestPinia({})
      const store = useRecordStore()
      store.setNavigationIDs(['a', 'b'])
      expect(store.getNextAndPrev('a').prev).toBeUndefined()
    })

    it('returns undefined next for last item', () => {
      createTestPinia({})
      const store = useRecordStore()
      store.setNavigationIDs(['a', 'b'])
      expect(store.getNextAndPrev('b').next).toBeUndefined()
    })

    it('returns both undefined when ID not in list', () => {
      createTestPinia({})
      const store = useRecordStore()
      store.setNavigationIDs(['a', 'b'])
      expect(store.getNextAndPrev('z')).toEqual({ prev: undefined, next: undefined })
    })
  })

  describe('clearAll()', () => {
    it('clears both caches and resets state', async () => {
      const api = makeAPI({
        recordRead: vi.fn().mockResolvedValue({ recordID: 'clr-1', namespaceID: NS_ID, moduleID: MOD_ID }),
      })
      createTestPinia({ '$ComposeAPI': api })
      const store = useRecordStore()

      await store.resolveRecordLabels({ namespaceID: NS_ID, moduleID: MOD_ID, recordIDs: ['clr-1'] })
      store.setNavigationIDs(['x', 'y'])
      expect(store.getByID('clr-1')).toBeTruthy()

      store.clearAll()
      expect(store.getByID('clr-1')).toBeNull()
      expect(store.loading).toBe(false)
      expect(store.pending).toBe(false)
      expect(store.paginationRecordIDs).toHaveLength(0)
    })
  })

  describe('delete()', () => {
    it('removes record from cache after successful delete', async () => {
      const api = makeAPI({
        recordList: vi.fn().mockResolvedValue({
          set: [{ recordID: '30001', namespaceID: NS_ID, moduleID: MOD_ID }],
          filter: {},
        }),
        recordDelete: vi.fn().mockResolvedValue({}),
      })
      const { store } = setup(api)
      await store.list({ namespaceID: NS_ID, moduleID: MOD_ID })
      expect(store.getByID('30001')).toBeTruthy()

      await store.delete({ namespaceID: NS_ID, moduleID: MOD_ID, recordID: '30001' })
      expect(store.getByID('30001')).toBeNull()
    })
  })
})

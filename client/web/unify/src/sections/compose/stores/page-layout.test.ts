import { describe, it, expect, beforeEach } from 'vitest'
import { createTestPinia, createMockComposeAPI } from '@planetcrust/human-test-utils'
import { usePageLayoutStore } from '@planetcrust/human-vue'

const NS = '10001'
const PAGE = '20001'

const layout = (pageLayoutID: string, pageID = PAGE) => ({
  pageLayoutID,
  pageID,
  namespaceID: NS,
  handle: `l${pageLayoutID}`,
  blocks: [],
})

describe('usePageLayoutStore', () => {
  let api: ReturnType<typeof createMockComposeAPI>

  beforeEach(() => {
    api = createMockComposeAPI()
    createTestPinia({ $ComposeAPI: api })
  })

  describe('findByPageID()', () => {
    it('serves the cache without an API call', async () => {
      api.pageLayoutListNamespace.mockResolvedValue({ set: [layout('1')], filter: {} })
      const store = usePageLayoutStore()
      await store.load({ namespaceID: NS })

      const found = await store.findByPageID({ namespaceID: NS, pageID: PAGE })

      expect(found).toHaveLength(1)
      expect(api.pageLayoutList).not.toHaveBeenCalled()
    })

    it('force refetches even with a populated cache', async () => {
      api.pageLayoutListNamespace.mockResolvedValue({ set: [layout('1')], filter: {} })
      api.pageLayoutList.mockResolvedValue({ set: [layout('1'), layout('2')], filter: {} })
      const store = usePageLayoutStore()
      await store.load({ namespaceID: NS })

      const found = await store.findByPageID({ namespaceID: NS, pageID: PAGE, force: true })

      expect(api.pageLayoutList).toHaveBeenCalledWith({
        namespaceID: NS,
        pageID: PAGE,
        sort: 'weight ASC',
      })
      expect(found.map(l => l.pageLayoutID)).toEqual(['1', '2'])
      expect(store.getByPageID(PAGE)).toHaveLength(2)
    })

    it('force drops the page layouts the server no longer lists', async () => {
      api.pageLayoutListNamespace.mockResolvedValue({
        set: [layout('1'), layout('2'), layout('9', '20009')],
        filter: {},
      })
      api.pageLayoutList.mockResolvedValue({ set: [layout('2')], filter: {} })
      const store = usePageLayoutStore()
      await store.load({ namespaceID: NS })

      await store.findByPageID({ namespaceID: NS, pageID: PAGE, force: true })

      expect(store.getByPageID(PAGE).map(l => l.pageLayoutID)).toEqual(['2'])
      // another page's layouts are none of this fetch's business
      expect(store.getByPageID('20009')).toHaveLength(1)
    })
  })

  describe('findByID()', () => {
    it('serves the cache without an API call', async () => {
      api.pageLayoutListNamespace.mockResolvedValue({ set: [layout('1')], filter: {} })
      const store = usePageLayoutStore()
      await store.load({ namespaceID: NS })

      await store.findByID({ namespaceID: NS, pageID: PAGE, pageLayoutID: '1' })

      expect(api.pageLayoutRead).not.toHaveBeenCalled()
    })

    it('force refetches even when cached', async () => {
      api.pageLayoutListNamespace.mockResolvedValue({ set: [layout('1')], filter: {} })
      api.pageLayoutRead.mockResolvedValue({ ...layout('1'), handle: 'renamed' })
      const store = usePageLayoutStore()
      await store.load({ namespaceID: NS })

      const found = await store.findByID({
        namespaceID: NS,
        pageID: PAGE,
        pageLayoutID: '1',
        force: true,
      })

      expect(api.pageLayoutRead).toHaveBeenCalledTimes(1)
      expect(found.handle).toBe('renamed')
    })
  })
})

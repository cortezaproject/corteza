import { describe, it, expect, beforeEach } from 'vitest'
import { createTestPinia, createMockComposeAPI } from '@planetcrust/human-test-utils'
import { usePageStore, usePageLayoutStore } from '@planetcrust/human-vue'

const NS = '10001'
const PAGE = '20001'

const page = (pageID = PAGE, overrides = {}) => ({
  pageID,
  namespaceID: NS,
  title: `Page ${pageID}`,
  blocks: [],
  ...overrides,
})

describe('usePageStore', () => {
  let api: ReturnType<typeof createMockComposeAPI>

  beforeEach(() => {
    api = createMockComposeAPI()
    createTestPinia({ $ComposeAPI: api })
  })

  describe('findByID()', () => {
    it('serves the cache without an API call', async () => {
      api.pageList.mockResolvedValue({ set: [page()], filter: {} })
      const store = usePageStore()
      await store.load({ namespaceID: NS })

      await store.findByID({ namespaceID: NS, pageID: PAGE })

      expect(api.pageRead).not.toHaveBeenCalled()
    })

    it('force refetches even when cached', async () => {
      api.pageList.mockResolvedValue({ set: [page()], filter: {} })
      api.pageRead.mockResolvedValue(page(PAGE, { title: 'Renamed elsewhere' }))
      const store = usePageStore()
      await store.load({ namespaceID: NS })

      const found = await store.findByID({ namespaceID: NS, pageID: PAGE, force: true })

      expect(api.pageRead).toHaveBeenCalledWith({ namespaceID: NS, pageID: PAGE })
      expect(found.title).toBe('Renamed elsewhere')
      expect(store.getByID(PAGE).title).toBe('Renamed elsewhere')
    })
  })

  describe('create()', () => {
    it('pulls the layout the server made with the page', async () => {
      api.pageCreate.mockResolvedValue(page())
      api.pageLayoutList.mockResolvedValue({
        set: [{ pageLayoutID: '1', pageID: PAGE, namespaceID: NS, handle: 'primary', blocks: [] }],
        filter: {},
      })

      const created = await usePageStore().create(page())

      expect(created.pageID).toBe(PAGE)
      expect(api.pageLayoutList).toHaveBeenCalledWith({
        namespaceID: NS,
        pageID: PAGE,
        sort: 'weight ASC',
      })
      expect(usePageLayoutStore().getByPageID(PAGE)).toHaveLength(1)
    })

    it('still returns the page when its layouts cannot be fetched', async () => {
      api.pageCreate.mockResolvedValue(page())
      api.pageLayoutList.mockRejectedValue(new Error('nope'))

      const created = await usePageStore().create(page())

      expect(created.pageID).toBe(PAGE)
    })
  })
})

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { defineComponent, nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { useResourceList } from './useResourceList'

function makeAPI(items: unknown[] = [], extra: Record<string, unknown> = {}) {
  return vi.fn().mockReturnValue({
    response: vi.fn().mockResolvedValue({
      set: items,
      filter: { total: items.length, nextPage: '', prevPage: '', incTotal: true, ...extra },
    }),
    cancel: vi.fn(),
  })
}

function makeRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<router-view/>' } },
      { path: '/list', component: { template: '<div/>' }, name: 'list' },
    ],
  })
}

let wrappers: ReturnType<typeof mount>[] = []

function mountList(apiFn: ReturnType<typeof makeAPI>, options = {}) {
  const router = makeRouter()
  let listRef: ReturnType<typeof useResourceList> | undefined

  const Comp = defineComponent({
    setup() {
      listRef = useResourceList(apiFn, options)
      return listRef
    },
    template: '<div/>',
  })

  const wrapper = mount(Comp, {
    global: { plugins: [router] },
    attachTo: document.body,
  })
  wrappers.push(wrapper)
  return { wrapper, router, get list() { return listRef! } }
}

describe('useResourceList', () => {
  afterEach(() => {
    wrappers.forEach(w => w.unmount())
    wrappers = []
    vi.restoreAllMocks()
    document.body.innerHTML = ''
  })

  describe('initial fetch', () => {
    it('calls apiFn on mount and populates items', async () => {
      const api = makeAPI([{ id: 1 }, { id: 2 }])
      const { list } = mountList(api)
      await flushPromises()
      expect(api).toHaveBeenCalled()
      expect(list.items.value).toHaveLength(2)
    })

    it('starts with loading true, ends with loading false', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      expect(list.loading.value).toBe(true)
      await flushPromises()
      expect(list.loading.value).toBe(false)
    })

    it('does not fetch when immediate=false', async () => {
      const api = makeAPI([])
      const { list } = mountList(api, { immediate: false })
      await flushPromises()
      expect(api).not.toHaveBeenCalled()
      expect(list.loading.value).toBe(false)
    })
  })

  describe('fetchItems()', () => {
    it('can be called manually to refresh', async () => {
      const api = makeAPI([{ id: 1 }])
      const { list } = mountList(api, { immediate: false })
      await flushPromises()
      expect(list.items.value).toHaveLength(0)

      await list.fetchItems(false)
      await flushPromises()
      expect(list.items.value).toHaveLength(1)
      expect(api).toHaveBeenCalledTimes(1)
    })
  })

  describe('pagination state', () => {
    it('applies default limit of 100', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      await flushPromises()
      const [params] = api.mock.calls[0]
      expect(params.limit).toBe(100)
    })

    it('respects custom initial limit', async () => {
      const api = makeAPI([])
      const { list } = mountList(api, { pagination: { limit: 25 } })
      await flushPromises()
      const [params] = api.mock.calls[0]
      expect(params.limit).toBe(25)
    })

    it('updates total from filter response', async () => {
      const api = vi.fn().mockReturnValue({
        response: vi.fn().mockResolvedValue({
          set: [1, 2, 3],
          filter: { total: 99, incTotal: true, nextPage: '', prevPage: '' },
        }),
        cancel: vi.fn(),
      })
      const { list } = mountList(api)
      await flushPromises()
      expect(list.pagination.total).toBe(99)
    })
  })

  describe('handleSort()', () => {
    it('sets sortBy and fetches', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      await flushPromises()
      api.mockClear()

      list.handleSort({ sortField: 'name', sortOrder: 1 })
      await flushPromises()

      expect(list.sorting.sortBy).toBe('name')
      expect(api).toHaveBeenCalled()
    })

    it('toggles sortDesc when sorting same field again', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      await flushPromises()

      list.handleSort({ sortField: 'name', sortOrder: 1 })
      await flushPromises()
      const first = list.sorting.sortDesc

      list.handleSort({ sortField: 'name', sortOrder: 1 })
      await flushPromises()
      expect(list.sorting.sortDesc).toBe(!first)
    })

    it('resets to ascending when sorting a new field', async () => {
      const api = makeAPI([])
      const { list } = mountList(api, { sorting: { sortBy: 'name', sortDesc: true } })
      await flushPromises()

      list.handleSort({ sortField: 'email', sortOrder: 1 })
      await flushPromises()

      expect(list.sorting.sortBy).toBe('email')
      expect(list.sorting.sortDesc).toBe(false)
    })

    it('ignores event with no sortField', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      await flushPromises()
      api.mockClear()

      list.handleSort({})
      await flushPromises()
      expect(api).not.toHaveBeenCalled()
    })
  })

  describe('handlePageChange()', () => {
    it('fetches with the given cursor and updates page', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      await flushPromises()
      api.mockClear()

      list.handlePageChange({ pageCursor: 'cursor-xyz', page: 2 })
      await flushPromises()

      expect(list.pagination.page).toBe(2)
      // Verify the API was called with the cursor
      expect(api).toHaveBeenCalled()
      const [params] = api.mock.calls[api.mock.calls.length - 1]
      expect(params.pageCursor).toBe('cursor-xyz')
    })

    it('updates limit if provided', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      await flushPromises()

      list.handlePageChange({ pageCursor: '', page: 1, limit: 50 })
      await flushPromises()

      expect(list.pagination.limit).toBe(50)
    })
  })

  describe('filterList()', () => {
    it('resets page to 1 and fetches', async () => {
      const api = makeAPI([])
      const { list } = mountList(api)
      await flushPromises()

      list.pagination.page = 5
      api.mockClear()

      list.filterList()
      await flushPromises()

      expect(list.pagination.page).toBe(1)
      // fetchItems clears pageCursor to undefined after a successful fetch
      expect(list.pagination.pageCursor).toBeUndefined()
      expect(api).toHaveBeenCalled()
    })
  })

  describe('abortRequests()', () => {
    it('cancels in-flight requests', async () => {
      const cancel = vi.fn()
      const api = vi.fn().mockReturnValue({
        response: vi.fn().mockImplementation(() => new Promise(() => {})), // never resolves
        cancel,
      })
      const { list } = mountList(api)
      await nextTick()

      list.abortRequests()
      expect(cancel).toHaveBeenCalled()
    })
  })

  describe('error state', () => {
    it('sets error when apiFn rejects', async () => {
      const api = vi.fn().mockReturnValue({
        response: vi.fn().mockRejectedValue(new Error('fetch failed')),
        cancel: vi.fn(),
      })
      const { list } = mountList(api)
      await flushPromises()
      expect(list.error.value).toBeInstanceOf(Error)
      expect(list.loading.value).toBe(false)
    })
  })
})

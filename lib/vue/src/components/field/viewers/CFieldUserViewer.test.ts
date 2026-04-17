import { describe, it, expect, vi } from 'vitest'
import { flushPromises } from '@vue/test-utils'
import { mountWithContext } from '@planetcrust/human-test-utils'
import CFieldUserViewer from './CFieldUserViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'assignee', isSystem: false, isMulti: false, options: {}, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

function makeUserStore(users: Array<{ userID: string; name?: string; handle?: string; email?: string }> = []) {
  const map = new Map(users.map(u => [u.userID, u]))
  return {
    findByID: (id: string) => map.get(id) || null,
    resolveUsers: vi.fn().mockResolvedValue(undefined),
  }
}

describe('CFieldUserViewer', () => {
  describe('no value', () => {
    it('renders empty when no value in record', async () => {
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: makeUserStore() }, {
        props: { field: field(), record: record() },
      })
      await flushPromises()
      expect(wrapper.findAll('span')).toHaveLength(0)
    })
  })

  describe('single value — user not cached', () => {
    it('shows userID as fallback when user not in store', async () => {
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: makeUserStore() }, {
        props: { field: field(), record: record({ assignee: '9001' }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain('9001')
    })
  })

  describe('single value — user cached', () => {
    it('shows formatted name from store', async () => {
      const store = makeUserStore([{ userID: '9001', name: 'Alice Smith' }])
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: store }, {
        props: { field: field(), record: record({ assignee: '9001' }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain('Alice Smith')
    })

    it('shows handle when name is absent', async () => {
      const store = makeUserStore([{ userID: '9002', handle: 'alice_h' }])
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: store }, {
        props: { field: field(), record: record({ assignee: '9002' }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain('alice_h')
    })
  })

  describe('multi value', () => {
    it('renders one span per user ID', async () => {
      const store = makeUserStore([
        { userID: '9001', name: 'Alice' },
        { userID: '9002', name: 'Bob' },
      ])
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: store }, {
        props: { field: field({ isMulti: true }), record: record({ assignee: ['9001', '9002'] }) },
      })
      await flushPromises()
      const spans = wrapper.findAll('span')
      expect(spans).toHaveLength(2)
    })

    it('uses default delimiter between users', async () => {
      const store = makeUserStore([
        { userID: '9001', name: 'Alice' },
        { userID: '9002', name: 'Bob' },
      ])
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: store }, {
        props: { field: field({ isMulti: true }), record: record({ assignee: ['9001', '9002'] }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain(', ')
    })

    it('uses custom multiDelimiter', async () => {
      const store = makeUserStore([
        { userID: '9001', name: 'Alice' },
        { userID: '9002', name: 'Bob' },
      ])
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: store }, {
        props: {
          field: field({ isMulti: true, options: { multiDelimiter: ' / ' } }),
          record: record({ assignee: ['9001', '9002'] }),
        },
      })
      await flushPromises()
      expect(wrapper.text()).toContain(' / ')
    })
  })

  describe('resolveUsers', () => {
    it('calls userStore.resolveUsers on mount with userIDs', async () => {
      const store = makeUserStore()
      mountWithContext(CFieldUserViewer, { userStore: store }, {
        props: { field: field(), record: record({ assignee: '9001' }) },
      })
      await flushPromises()
      expect(store.resolveUsers).toHaveBeenCalledWith(['9001'])
    })

    it('does not call resolveUsers when no value', async () => {
      const store = makeUserStore()
      mountWithContext(CFieldUserViewer, { userStore: store }, {
        props: { field: field(), record: record() },
      })
      await flushPromises()
      expect(store.resolveUsers).not.toHaveBeenCalled()
    })
  })

  describe('no userStore', () => {
    it('shows userID as fallback with null store', async () => {
      const wrapper = mountWithContext(CFieldUserViewer, { userStore: null }, {
        props: { field: field(), record: record({ assignee: '9001' }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain('9001')
    })
  })
})

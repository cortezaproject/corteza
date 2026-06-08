import { describe, it, expect, vi, beforeEach } from 'vitest'
import { flushPromises } from '@vue/test-utils'
import { mountWithContext, makeUser, createMockSystemAPI, createTestPinia } from '@planetcrust/human-test-utils'
import { useUserStore } from '../../../stores/useUserStore'
import CFieldUserViewer from './CFieldUserViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'assignee', isSystem: false, isMulti: false, options: {}, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

describe('CFieldUserViewer', () => {
  let api: ReturnType<typeof createMockSystemAPI>

  beforeEach(() => {
    api = createMockSystemAPI()
    createTestPinia({ '$SystemAPI': api })
  })

  describe('no value', () => {
    it('renders empty when no value in record', async () => {
      const wrapper = mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field(), record: record() },
      })
      await flushPromises()
      expect(wrapper.findAll('span')).toHaveLength(0)
    })
  })

  describe('single value — user not cached', () => {
    it('shows userID as fallback when user not in store', async () => {
      const wrapper = mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field(), record: record({ assignee: '9001' }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain('9001')
    })
  })

  describe('single value — user cached', () => {
    it('shows formatted name from store', async () => {
      useUserStore().storeUsers([makeUser({ userID: '9001', name: 'Alice Smith' })])
      const wrapper = mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field(), record: record({ assignee: '9001' }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain('Alice Smith')
    })

    it('shows handle when name is absent', async () => {
      useUserStore().storeUsers([makeUser({ userID: '9002', name: '', handle: 'alice_h' })])
      const wrapper = mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field(), record: record({ assignee: '9002' }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain('alice_h')
    })
  })

  describe('multi value', () => {
    it('renders one span per user ID', async () => {
      useUserStore().storeUsers([
        makeUser({ userID: '9001', name: 'Alice' }),
        makeUser({ userID: '9002', name: 'Bob' }),
      ])
      const wrapper = mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field({ isMulti: true }), record: record({ assignee: ['9001', '9002'] }) },
      })
      await flushPromises()
      const spans = wrapper.findAll('span')
      expect(spans).toHaveLength(2)
    })

    it('uses default delimiter between users', async () => {
      useUserStore().storeUsers([
        makeUser({ userID: '9001', name: 'Alice' }),
        makeUser({ userID: '9002', name: 'Bob' }),
      ])
      const wrapper = mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field({ isMulti: true }), record: record({ assignee: ['9001', '9002'] }) },
      })
      await flushPromises()
      expect(wrapper.text()).toContain(', ')
    })

    it('uses custom multiDelimiter', async () => {
      useUserStore().storeUsers([
        makeUser({ userID: '9001', name: 'Alice' }),
        makeUser({ userID: '9002', name: 'Bob' }),
      ])
      const wrapper = mountWithContext(CFieldUserViewer, { systemAPI: api }, {
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
    it('calls userList on mount when user not cached', async () => {
      mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field(), record: record({ assignee: '9001' }) },
      })
      await flushPromises()
      expect(api.userList).toHaveBeenCalledWith({ userID: ['9001'] })
    })

    it('does not call userList when no value', async () => {
      mountWithContext(CFieldUserViewer, { systemAPI: api }, {
        props: { field: field(), record: record() },
      })
      await flushPromises()
      expect(api.userList).not.toHaveBeenCalled()
    })
  })
})

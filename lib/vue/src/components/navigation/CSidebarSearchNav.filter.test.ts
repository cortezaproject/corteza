import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import CSidebarSearchNav from './CSidebarSearchNav.vue'

// Searching a nav has to leave a tree behind: a match keeps the groups it lives
// under, and a group that matches keeps what is in it. Anything else renders
// orphans, or a group whose name the user just typed with nothing under it.

const router = createRouter({
  history: createMemoryHistory(),
  routes: [{ path: '/', component: { template: '<div />' } }],
})

type Item = { id: string; parentId: string; label: string; search?: string }

const ITEMS: Item[] = [
  { id: 'agents', parentId: '0', label: 'Agents' },
  { id: 'a1', parentId: 'agents', label: 'Invoice reader', search: 'invoice_reader' },
  { id: 'a2', parentId: 'agents', label: 'Support triage', search: 'support_triage' },
  { id: 'sessions', parentId: '0', label: 'Sessions' },
]

function mountNav(items: Item[] = ITEMS) {
  return mount(CSidebarSearchNav, {
    props: {
      items,
      idKey: 'id',
      parentKey: 'parentId',
      labelKey: 'label',
      searchKey: 'search',
      expandAll: true,
      noResultsLabel: 'No results',
    },
    global: { plugins: [router] },
  })
}

async function search(wrapper: ReturnType<typeof mountNav>, query: string) {
  await wrapper.find('input').setValue(query)
  await flushPromises()
}

describe('CSidebarSearchNav', () => {
  it('shows everything until something is typed', async () => {
    const wrapper = mountNav()
    await flushPromises()

    expect(wrapper.text()).toContain('Invoice reader')
    expect(wrapper.text()).toContain('Support triage')
    expect(wrapper.text()).toContain('Sessions')
  })

  it('keeps a match under the group it lives in', async () => {
    const wrapper = mountNav()
    await search(wrapper, 'triage')

    expect(wrapper.text()).toContain('Agents')
    expect(wrapper.text()).toContain('Support triage')
    expect(wrapper.text()).not.toContain('Invoice reader')
    expect(wrapper.text()).not.toContain('Sessions')
  })

  it('keeps the whole group when the group name matches', async () => {
    const wrapper = mountNav()
    await search(wrapper, 'agents')

    expect(wrapper.text()).toContain('Invoice reader')
    expect(wrapper.text()).toContain('Support triage')
    expect(wrapper.text()).not.toContain('Sessions')
  })

  it('matches the extra search text, not only the label', async () => {
    const wrapper = mountNav()
    await search(wrapper, 'invoice_reader')

    expect(wrapper.text()).toContain('Invoice reader')
    expect(wrapper.text()).not.toContain('Support triage')
  })

  it('says so when nothing matches', async () => {
    const wrapper = mountNav()
    await search(wrapper, 'nothing here')

    expect(wrapper.find('[data-testid="sidebar-nav-no-results"]').text()).toBe('No results')
    expect(wrapper.text()).not.toContain('Agents')
  })
})

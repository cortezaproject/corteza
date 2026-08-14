import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import CSidebarNav from './CSidebarNav.vue'

// `expandAll` decides which groups start open. Navs are fed by stores that are
// still loading at setup, and by callers that turn the prop on and off (search),
// so it has to be answered whenever either changes — reading it once leaves the
// groups shut and the nav looks empty.

const router = createRouter({
  history: createMemoryHistory(),
  routes: [{ path: '/', component: { template: '<div />' } }],
})

type Item = { id: string; parentId: string; label: string }

const GROUP: Item = { id: 'g', parentId: '0', label: 'Group' }
const CHILD: Item = { id: 'c', parentId: 'g', label: 'Child' }

function mountNav(items: Item[], expandAll = true) {
  return mount(CSidebarNav, {
    props: { items, idKey: 'id', parentKey: 'parentId', labelKey: 'label', expandAll },
    global: { plugins: [router] },
  })
}

describe('CSidebarNav expand-all', () => {
  it('expands a group whose items arrive after setup', async () => {
    const wrapper = mountNav([])
    await flushPromises()
    expect(wrapper.text()).not.toContain('Child')

    await wrapper.setProps({ items: [GROUP, CHILD] })
    await flushPromises()

    expect(wrapper.text()).toContain('Child')
  })

  it('leaves a group the user collapsed alone when the items change', async () => {
    const wrapper = mountNav([GROUP, CHILD])
    await flushPromises()
    expect(wrapper.text()).toContain('Child')

    await wrapper.findAll('button')[0].trigger('click')
    await flushPromises()
    expect(wrapper.text()).not.toContain('Child')

    await wrapper.setProps({ items: [GROUP, CHILD, { ...CHILD, id: 'c2', label: 'Second' }] })
    await flushPromises()

    expect(wrapper.text()).not.toContain('Child')
  })

  it('expands again when the caller turns expand-all back on', async () => {
    const wrapper = mountNav([GROUP, CHILD], false)
    await flushPromises()
    expect(wrapper.text()).not.toContain('Child')

    await wrapper.setProps({ expandAll: true })
    await flushPromises()

    expect(wrapper.text()).toContain('Child')
  })
})

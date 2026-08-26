import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive } from 'vue'

// The admin panel is the namespace's administration, and the server's `manage`
// op is what opens it. Offering it to a user who does not hold that op costs a
// tree of links that all refuse on arrival, so the whole section is absent
// rather than disabled — and when it is present it says what it is.

const route = reactive({ name: 'page', params: { slug: 'ns' } })

vi.mock('vue-router', () => ({
  useRoute: () => route,
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

let namespace

const namespaceStore = { getByUrlPart: urlPart => (urlPart === 'ns' ? namespace : null) }
const moduleStore = { set: [{ moduleID: 'M1', name: 'Leads' }] }
const pageStore = {
  set: [{ pageID: 'P1', selfID: '0', title: 'Home', visible: true, weight: 1 }],
}

vi.mock('@planetcrust/human-vue', () => ({
  useModuleStore: () => moduleStore,
  useNamespaceStore: () => namespaceStore,
  usePageStore: () => pageStore,
  components: {
    // Renders its items' labels so a test can see which nav got which items.
    CSidebarNav: {
      name: 'CSidebarNav',
      props: {
        items: { type: Array, default: () => [] },
        labelKey: { type: String, default: '' },
      },
      template: '<ul><li v-for="i in items" :key="i[labelKey]">{{ i[labelKey] }}</li></ul>',
    },
    CInputSearch: { name: 'CInputSearch', template: '<div />' },
  },
}))

import CSidebarNavigation from './CSidebarNavigation.vue'

let wrapper

function mountNav() {
  wrapper = mount(CSidebarNavigation, {
    global: {
      mocks: { $t: k => k },
      provide: { $ComposeAPI: { baseURL: '' } },
    },
  })
  return wrapper
}

beforeEach(() => {
  route.name = 'page'
  route.params = { slug: 'ns' }
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('Compose sidebar admin panel', () => {
  it('withholds it from a user without manage on the namespace', () => {
    namespace = { namespaceID: 'N1', slug: 'ns', canManageNamespace: false }

    mountNav()

    expect(wrapper.find('[data-testid="sidebar-admin-panel"]').exists()).toBe(false)
    expect(wrapper.text()).not.toContain('sidebar.modules')
    expect(wrapper.text()).not.toContain('sidebar.charts')
    // The module tree goes with it — its names are the admin panel's content.
    expect(wrapper.text()).not.toContain('Leads')
  })

  it('withholds it when the namespace is not in the store yet', () => {
    namespace = null

    mountNav()

    expect(wrapper.find('[data-testid="sidebar-admin-panel"]').exists()).toBe(false)
  })

  it('shows it, named, to a user with manage', () => {
    namespace = { namespaceID: 'N1', slug: 'ns', canManageNamespace: true }

    mountNav()

    const panel = wrapper.find('[data-testid="sidebar-admin-panel"]')
    expect(panel.exists()).toBe(true)
    expect(panel.find('h2').text()).toBe('sidebar.adminPanel')
    expect(panel.text()).toContain('sidebar.modules')
    expect(panel.text()).toContain('sidebar.pages')
    expect(panel.text()).toContain('sidebar.charts')
    expect(panel.text()).toContain('Leads')
  })

  it('leaves the public page tree alone either way', () => {
    namespace = { namespaceID: 'N1', slug: 'ns', canManageNamespace: false }
    mountNav()
    expect(wrapper.text()).toContain('Home')
  })

  it('scrolls the pages and the admin panel as one region', () => {
    namespace = { namespaceID: 'N1', slug: 'ns', canManageNamespace: true }

    mountNav()

    const scrollers = wrapper.findAll('.overflow-auto')
    expect(scrollers).toHaveLength(1)
    // The panel travels with the pages rather than sitting in a band of its own
    // below a second scrollbar...
    expect(scrollers[0].find('[data-testid="sidebar-admin-panel"]').exists()).toBe(true)
    // ...and the search box stays outside it, so it does not scroll away.
    expect(scrollers[0].findComponent({ name: 'CInputSearch' }).exists()).toBe(false)
  })

  it('drops the panel to the foot when there is nothing to scroll', () => {
    namespace = { namespaceID: 'N1', slug: 'ns', canManageNamespace: true }

    mountNav()

    // Flex does the work, not `position: sticky` — the column fills the
    // scroller and the pages take the slack, which pushes the panel down. Once
    // there is more than fits, the column grows and the panel simply scrolls.
    const panel = wrapper.find('[data-testid="sidebar-admin-panel"]')
    expect(panel.classes()).not.toContain('sticky')

    const column = wrapper.find('.overflow-auto > div')
    expect(column.classes()).toContain('min-h-full')
    const pages = column.find('.flex-1')
    expect(pages.exists()).toBe(true)
    expect(pages.find('[data-testid="sidebar-admin-panel"]').exists()).toBe(false)
  })

  it('names the panel above the rule, not below it', () => {
    namespace = { namespaceID: 'N1', slug: 'ns', canManageNamespace: true }

    mountNav()

    const panel = wrapper.find('[data-testid="sidebar-admin-panel"]')
    const children = [...panel.element.children]
    expect(children[0].tagName).toBe('H2')
    expect(children[0].textContent.trim()).toBe('sidebar.adminPanel')
    // the rule sits on the block holding the nav, under the heading
    expect(children[1].className).toContain('border-t')
    expect(children[1].querySelector('ul')).toBeTruthy()
  })
})

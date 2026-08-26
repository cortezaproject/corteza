import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive } from 'vue'

// The shell wraps this body in a scroller of its own. This column has to fit it
// exactly, so that one never engages and the namespace switcher stays put — a
// child that claims a full column on top of the switcher is taller than the
// scroller, and the sidebar ends up with two scrollbars. The nav takes what the
// switcher leaves and scrolls what is inside it.

const route = reactive({ meta: {} })

vi.mock('vue-router', () => ({ useRoute: () => route }))

vi.mock('@/sections/compose/components/CSidebarNamespaceNav.vue', () => ({
  default: { name: 'CSidebarNamespaceNav', template: '<div />' },
}))
vi.mock('@/sections/compose/components/CSidebarNamespaceSwitcher.vue', () => ({
  default: { name: 'CSidebarNamespaceSwitcher', template: '<div />' },
}))
vi.mock('@/sections/compose/components/CSidebarNavigation.vue', () => ({
  default: { name: 'CSidebarNavigation', template: '<div />' },
}))

import ComposeSidebar from './ComposeSidebar.vue'

let wrapper

function mountSidebar(sidebar) {
  route.meta = sidebar ? { sidebar } : {}
  wrapper = mount(ComposeSidebar, { global: { mocks: { $route: route } } })
  return wrapper
}

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

describe('Compose sidebar column', () => {
  it('fits the shell scroller exactly', () => {
    mountSidebar()

    expect(wrapper.classes()).toContain('h-full')
    expect(wrapper.classes()).not.toContain('min-h-full')
  })

  it('lets the nav take the height the switcher leaves', () => {
    mountSidebar()

    const nav = wrapper.findComponent({ name: 'CSidebarNavigation' })
    expect(nav.exists()).toBe(true)
    expect(nav.classes()).toContain('flex-1')
    expect(nav.classes()).toContain('min-h-0')
    expect(nav.classes()).not.toContain('h-full')
  })

  it('does the same for the namespace nav', () => {
    mountSidebar('namespaces')

    const nav = wrapper.findComponent({ name: 'CSidebarNamespaceNav' })
    expect(nav.exists()).toBe(true)
    expect(nav.classes()).toContain('flex-1')
    expect(nav.classes()).toContain('min-h-0')
    expect(nav.classes()).not.toContain('h-full')
  })
})

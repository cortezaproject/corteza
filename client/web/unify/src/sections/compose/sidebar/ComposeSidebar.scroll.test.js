import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive } from 'vue'

// The shell wraps this body in one scroller of its own. A child that asks for a
// full column on top of the namespace switcher is taller than that scroller, so
// the sidebar ends up with two scrollbars — the shell's and the nav's. The nav
// takes what the switcher leaves instead.

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

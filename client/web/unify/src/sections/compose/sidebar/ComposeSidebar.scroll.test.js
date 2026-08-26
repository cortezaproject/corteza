import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { reactive } from 'vue'

// The shell wraps this body in the one scroller the sidebar has, so everything
// here scrolls together — switcher, search, pages. The column must be free to
// grow past that scroller (`min-h-full`, never `h-full`, which would cap it and
// leave the nav to scroll on its own) while still filling it when there is
// little to show, since that is what pushes the admin panel to the foot.

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
  it('fills the shell scroller without being capped by it', () => {
    mountSidebar()

    expect(wrapper.classes()).toContain('min-h-full')
    expect(wrapper.classes()).not.toContain('h-full')
  })

  it('lets the nav take the height the switcher leaves', () => {
    mountSidebar()

    const nav = wrapper.findComponent({ name: 'CSidebarNavigation' })
    expect(nav.exists()).toBe(true)
    expect(nav.classes()).toContain('flex-1')
    expect(nav.classes()).not.toContain('h-full')
  })

  it('does the same for the namespace nav', () => {
    mountSidebar('namespaces')

    const nav = wrapper.findComponent({ name: 'CSidebarNamespaceNav' })
    expect(nav.exists()).toBe(true)
    expect(nav.classes()).toContain('flex-1')
    expect(nav.classes()).not.toContain('h-full')
  })
})

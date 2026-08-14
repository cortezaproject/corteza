import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'

vi.mock('./PageBlock.vue', () => ({
  default: { template: '<div><slot /></div>' },
}))

// The pages the namespace holds: one plain page with two children, and a
// record page. Both the sub-pages menu and the record-page route read these.
const PAGES = [
  { pageID: 'P1', title: 'Target', selfID: '0', moduleID: '0' },
  { pageID: 'P2', title: 'Sub One', selfID: 'P1', moduleID: '0' },
  { pageID: 'P3', title: 'Sub Two', selfID: 'P1', moduleID: '0' },
  { pageID: 'P9', title: 'Thing Record', selfID: '0', moduleID: 'M1' },
]

vi.mock('@planetcrust/human-vue', () => ({
  usePageStore: () => ({
    set: PAGES,
    getByID: id => PAGES.find(p => p.pageID === id),
  }),
}))

const pushed = []
let currentPageID = 'P0'

vi.mock('vue-router', () => ({
  useRoute: () => ({ params: { slug: 'ns', pageID: currentPageID } }),
  useRouter: () => ({ push: to => pushed.push(to) }),
}))

import NavigationBlock from './NavigationBlock.vue'

const composeItem = (label, pageID, extra = {}) => ({
  type: 'compose',
  options: { enabled: true, item: { label, pageID, ...extra } },
})

const ButtonStub = {
  name: 'ButtonStub',
  props: ['label'],
  template: '<button>{{ label }}</button>',
}
const MenuStub = { name: 'MenuStub', props: ['model'], template: '<div class="menu" />' }
const LinkStub = { name: 'LinkStub', props: ['to'], template: '<a><slot /></a>' }

function mountBlock(display, navigationItems) {
  return mount(NavigationBlock, {
    props: {
      block: { options: { display, navigationItems } },
      namespace: { namespaceID: 'N1' },
    },
    global: {
      stubs: { Button: ButtonStub, Menu: MenuStub, 'router-link': LinkStub },
      mocks: { $t: k => k },
    },
  })
}

beforeEach(() => {
  pushed.length = 0
  currentPageID = 'P0'
})

const rowClasses = w => w.find('.flex.h-full').classes()
const linkClasses = w => w.find('a').classes()

describe('NavigationBlock appearance', () => {
  it('draws tabs as a full-height row on a rule', () => {
    const w = mountBlock({ appearance: 'tabs' }, [composeItem('Target', 'P1')])

    expect(rowClasses(w)).toEqual(expect.arrayContaining(['items-stretch', 'border-b']))
    expect(linkClasses(w)).toEqual(expect.arrayContaining(['border-b-2', 'px-3', 'py-2']))
    expect(linkClasses(w)).not.toContain('rounded-full')
  })

  it('draws pills as rounded items, centred in the block', () => {
    const w = mountBlock({ appearance: 'pills' }, [composeItem('Target', 'P1')])

    expect(rowClasses(w)).toEqual(expect.arrayContaining(['items-center']))
    expect(rowClasses(w)).not.toContain('border-b')
    expect(linkClasses(w)).toContain('rounded-full')
  })

  it('draws small as the same row, tightened', () => {
    const w = mountBlock({ appearance: 'small' }, [composeItem('Target', 'P1')])

    expect(linkClasses(w)).toEqual(expect.arrayContaining(['px-2', 'py-1', 'text-xs']))
  })

  it('gives tabs and pills different classes', () => {
    const tabs = linkClasses(mountBlock({ appearance: 'tabs' }, [composeItem('T', 'P1')]))
    const pills = linkClasses(mountBlock({ appearance: 'pills' }, [composeItem('T', 'P1')]))

    expect(tabs).not.toEqual(pills)
  })
})

describe('NavigationBlock active item', () => {
  it('marks the item for the page being shown', () => {
    currentPageID = 'P1'
    const w = mountBlock({ appearance: 'pills' }, [composeItem('Target', 'P1')])

    expect(linkClasses(w)).toEqual(expect.arrayContaining(['bg-primary', 'text-primary-contrast']))
  })

  it('leaves the other items alone', () => {
    currentPageID = 'P9'
    const w = mountBlock({ appearance: 'pills' }, [composeItem('Target', 'P1')])

    expect(linkClasses(w)).not.toContain('bg-primary')
  })

  it('never marks a URL item, whatever page is shown', () => {
    currentPageID = 'P1'
    const w = mountBlock({ appearance: 'tabs' }, [
      {
        type: 'url',
        options: { enabled: true, item: { label: 'Out', url: 'https://example.com' } },
      },
    ])

    expect(linkClasses(w)).not.toContain('border-primary')
  })
})

describe('NavigationBlock compose page routing', () => {
  it('links a plain page by page route', () => {
    const w = mountBlock({}, [composeItem('Target', 'P1')])

    expect(w.findComponent(LinkStub).props('to')).toEqual({
      name: 'page',
      params: { pageID: 'P1' },
      query: {},
    })
  })

  it('carries a chosen layout as a query', () => {
    const w = mountBlock({}, [composeItem('Target', 'P1', { pageLayoutID: 'L7' })])

    expect(w.findComponent(LinkStub).props('to').query).toEqual({ layoutID: 'L7' })
  })

  it('opens a record page on a blank record rather than on the page itself', () => {
    const w = mountBlock({}, [composeItem('Thing', 'P9')])

    expect(w.findComponent(LinkStub).props('to')).toEqual({
      name: 'page.record',
      params: { pageID: 'P9', recordID: '0' },
      query: {},
    })
  })
})

describe('NavigationBlock sub-pages dropdown', () => {
  const menuModel = w => w.findComponent(MenuStub).props('model')

  it('lists the page itself and then its children', () => {
    const w = mountBlock({}, [composeItem('Target', 'P1', { displaySubPages: true })])

    expect(menuModel(w).map(i => i.label ?? '---')).toEqual(['Target', '---', 'Sub One', 'Sub Two'])
  })

  it('navigates to the child a menu entry names', () => {
    const w = mountBlock({}, [composeItem('Target', 'P1', { displaySubPages: true })])
    menuModel(w).at(-1).command()

    expect(pushed).toEqual([{ name: 'page', params: { pageID: 'P3' }, query: {} }])
  })

  it('stays a plain link when the page has no children', () => {
    const w = mountBlock({}, [composeItem('Thing', 'P9', { displaySubPages: true })])

    expect(w.findComponent(MenuStub).exists()).toBe(false)
    expect(w.find('a').exists()).toBe(true)
  })

  it('marks the item while one of its children is being shown', () => {
    currentPageID = 'P2'
    const w = mountBlock({ appearance: 'pills' }, [
      composeItem('Target', 'P1', { displaySubPages: true }),
    ])

    expect(w.find('.relative').classes()).toContain('bg-primary')
  })
})

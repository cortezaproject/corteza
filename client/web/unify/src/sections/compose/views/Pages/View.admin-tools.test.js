import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive } from 'vue'

// The page builder and the page editor are screens of the namespace's admin
// panel, so the topbar offers them only to a user holding the namespace's
// `manage` op as well as the page's own update right. The page's right alone
// used to be enough, which let someone with no claim on the admin panel walk
// straight into it.

const route = reactive({ name: 'page', params: { slug: 'ns', pageID: 'P1' }, query: {} })

const router = { push: vi.fn(), replace: vi.fn() }

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
}))

let page
let layouts

const pageStore = { getByID: () => page }
const pageLayoutStore = { getByPageID: () => layouts }

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@planetcrust/human-vue', () => ({
  usePageStore: () => pageStore,
  usePageLayoutStore: () => pageLayoutStore,
  useHistoryBack: () => vi.fn(),
}))

vi.mock('@planetcrust/human-js', () => ({
  NoID: '0',
  compose: { PageBlockMaker: b => b, interpolateTemplate: t => t },
}))

vi.mock('@/sections/compose/composables/usePageVisibility', () => ({
  fetchBlockID: b => b.blockID,
  refuseOnce: () => true,
  clearRefusal: () => {},
  usePageVisibility: () => ({
    buildExpressionVariables: () => ({}),
    determineLayout: () => Promise.resolve(layouts[0]),
    evaluateBlocks: () => Promise.resolve(new Set()),
  }),
}))

vi.mock('@/sections/compose/composables/useResourceTranslations', () => ({
  useResourceTranslations: () => ({ showTranslatorButton: false }),
}))

vi.mock('@/sections/compose/components/PageBlocks/Grid.vue', () => ({
  default: { template: '<div />' },
}))

vi.mock('@/sections/compose/components/Admin/Page/PageTranslator.vue', () => ({
  default: { template: '<div />' },
}))

import View from './View.vue'

let wrapper

const PAGE_BUILDER = 'page.block.general.label.pageBuilder'
const ADD_LAYOUT = 'page.page-layout.add'

async function mountView(namespace) {
  wrapper = mount(View, {
    props: { namespace },
    global: {
      stubs: { teleport: true, Teleport: true, Message: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $SystemAPI: {},
        $Auth: { user: { userID: 'U1', name: 'Ada' } },
      },
      renderStubDefaultSlot: false,
    },
    shallow: true,
  })

  await flushPromises()
  return wrapper
}

beforeEach(() => {
  page = {
    pageID: 'P1',
    namespaceID: 'N1',
    title: 'Plain page title',
    canUpdatePage: true,
    blocks: [],
  }
  layouts = [{ pageLayoutID: 'L1', blocks: [], config: {}, meta: {} }]
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
  vi.clearAllMocks()
})

describe('Page View admin tools', () => {
  it('offers the builder to a user with manage and the page right', async () => {
    await mountView({ namespaceID: 'N1', canManageNamespace: true })

    expect(wrapper.html()).toContain(PAGE_BUILDER)
  })

  it('withholds it from a user without manage on the namespace', async () => {
    await mountView({ namespaceID: 'N1', canManageNamespace: false })

    expect(wrapper.html()).not.toContain(PAGE_BUILDER)
  })

  it('withholds it from a user who cannot update the page', async () => {
    page = { ...page, canUpdatePage: false }

    await mountView({ namespaceID: 'N1', canManageNamespace: true })

    expect(wrapper.html()).not.toContain(PAGE_BUILDER)
  })

  it('withholds the empty-state "add a layout" button on the same terms', async () => {
    layouts = []

    await mountView({ namespaceID: 'N1', canManageNamespace: false })

    expect(wrapper.html()).not.toContain(ADD_LAYOUT)

    wrapper.unmount()
    await mountView({ namespaceID: 'N1', canManageNamespace: true })

    expect(wrapper.html()).toContain(ADD_LAYOUT)
  })
})

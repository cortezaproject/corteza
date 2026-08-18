import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createTestPinia } from '@planetcrust/human-test-utils'
import { reactive } from 'vue'

// The builder renders a layout, so it has to fetch one. Reading the namespace's
// cached set instead leaves it on its "no layout yet" empty state for every page
// created since that set was loaded — including the record pages the module
// editor's own buttons make, whose primary layout the server creates server-side.

const route = reactive({
  name: 'admin.pages.builder',
  params: { slug: 'ns', pageID: 'P1' },
  query: {},
})

const router = { push: vi.fn(), replace: vi.fn() }

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
  onBeforeRouteLeave: () => {},
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

vi.mock('primevue/useconfirm', () => ({
  useConfirm: () => ({ require: vi.fn() }),
}))

let pages
let layoutsByPage

const pageStore = {
  getByID: vi.fn(id => pages[id]),
  findByID: vi.fn(({ pageID }) => Promise.resolve(pages[pageID])),
  update: vi.fn(p => Promise.resolve(p)),
}
const pageLayoutStore = {
  getByPageID: vi.fn(() => []),
  findByPageID: vi.fn(({ pageID }) => Promise.resolve(layoutsByPage[pageID] || [])),
  create: vi.fn(),
  update: vi.fn(),
  delete: vi.fn(),
}
const moduleStore = { getByID: () => null }

vi.mock('@planetcrust/human-vue', () => ({
  usePageStore: () => pageStore,
  usePageLayoutStore: () => pageLayoutStore,
  useModuleStore: () => moduleStore,
  useHistoryBack: () => vi.fn(),
  useDraftGuard: () => ({
    isDirty: { value: false },
    capture: vi.fn(),
    reset: vi.fn(),
    markSaved: vi.fn(),
  }),
  useRBACStore: () => ({ can: () => true, canGlobal: () => true }),
}))

vi.mock('@planetcrust/human-js', () => ({
  compose: {
    Page: class {
      constructor(p = {}) {
        Object.assign(this, p)
      }
    },
    PageBlockMaker: i => ({ ...i }),
  },
  NoID: '0',
}))

vi.mock('@/sections/compose/components/PageBlocks/Grid.vue', () => ({
  default: { props: ['blocks'], template: '<div />' },
}))

import Builder from './Builder.vue'

// Globally registered in the app, so `shallow` has nothing to stub them from
// and each warns once per mount — noise that buries the test's own output.
const GLOBAL_COMPONENTS = Object.fromEntries(
  [
    'Button',
    'ButtonGroup',
    'CEditorActions',
    'CExpressionHint',
    'CFormGroup',
    'CInputExpression',
    'CInputRole',
    'Checkbox',
    'Dialog',
    'Divider',
    'Fieldset',
    'InputNumber',
    'InputText',
    'ProgressSpinner',
    'Select',
    'Tab',
    'TabList',
    'TabPanel',
    'TabPanels',
    'Tabs',
  ].map(name => [name, true]),
)

let wrapper

async function mountBuilder() {
  wrapper = mount(Builder, {
    props: { namespace: { namespaceID: 'N1' } },
    global: {
      stubs: { teleport: true, Teleport: true, ...GLOBAL_COMPONENTS },
      directives: { tooltip: {}, focus: {} },
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastDanger: vi.fn(), toastErrorHandler: () => vi.fn() },
        $ComposeAPI: {},
        $SystemAPI: {},
        $Settings: { get: () => undefined },
        $Auth: { user: { userID: 'U1', roles: [] } },
        $eventBus: null,
      },
      renderStubDefaultSlot: false,
    },
    shallow: true,
  })
  await flushPromises()
  return wrapper
}

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

beforeEach(() => {
  createTestPinia({ $ComposeAPI: {}, $SystemAPI: {} })
  route.params = { slug: 'ns', pageID: 'P1' }
  route.query = {}
  pageStore.findByID.mockClear()
  pageLayoutStore.findByPageID.mockClear()
  pageLayoutStore.getByPageID.mockClear()
  pages = { P1: { pageID: 'P1', title: 'A page', blocks: [], moduleID: '0' } }
  layoutsByPage = { P1: [{ pageLayoutID: 'L1', pageID: 'P1', blocks: [] }] }
})

describe('page builder load', () => {
  it('fetches the page rather than reading the namespace cache', async () => {
    await mountBuilder()

    expect(pageStore.findByID).toHaveBeenCalledWith({
      namespaceID: 'N1',
      pageID: 'P1',
      force: true,
    })
  })

  it('fetches the layouts rather than reading the namespace cache', async () => {
    await mountBuilder()

    expect(pageLayoutStore.findByPageID).toHaveBeenCalledWith({
      namespaceID: 'N1',
      pageID: 'P1',
      force: true,
    })
    expect(pageLayoutStore.getByPageID).not.toHaveBeenCalled()
  })

  it('opens a layout the namespace cache never held', async () => {
    // What the reported defect looks like: the layout exists server-side, made
    // alongside the page, and the cached set predates it.
    pageLayoutStore.getByPageID.mockReturnValue([])

    await mountBuilder()

    expect(wrapper.text()).not.toContain('page.build.noLayout')
  })

  it('shows the empty state only when the page genuinely has no layout', async () => {
    layoutsByPage = { P1: [] }

    await mountBuilder()

    expect(wrapper.text()).toContain('page.build.noLayout')
  })
})

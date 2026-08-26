import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { computed, reactive, toValue } from 'vue'

// The record page's topbar offers two ways into the namespace's admin panel:
// the module editor and the page builder. Each takes the namespace's `manage`
// op plus its own resource's right — the module's for the module editor, the
// page's for the builder — so neither is ever offered to someone the server
// will turn away, and neither stands in for the other.

const route = reactive({
  name: 'page.record',
  params: { slug: 'ns', pageID: 'P1', recordID: 'R1' },
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

function makeRecord(recordID, values) {
  const rec = { recordID, values, ownedBy: 'U9', serialize: () => ({ recordID, values }) }
  rec.clone = () => makeRecord(recordID, { ...values })
  return rec
}

const records = { R1: makeRecord('R1', { name: 'Acme Corp' }) }

const pageStore = { getByID: () => page }
const pageLayoutStore = { getByPageID: () => layouts }
const moduleStore = { getByID: () => module_ }
const recordStore = {
  paginationRecordIDs: [],
  findByID: vi.fn(({ recordID }) => Promise.resolve(records[recordID])),
  update: vi.fn(rec => Promise.resolve(rec)),
  create: vi.fn(rec => Promise.resolve(rec)),
}

vi.mock('@planetcrust/human-vue', () => ({
  useDeferredBusy: src => computed(() => !!toValue(src)),
  useModuleStore: () => moduleStore,
  usePageStore: () => pageStore,
  usePageLayoutStore: () => pageLayoutStore,
  useRecordStore: () => recordStore,
  useHistoryBack: () => vi.fn(),
  components: { CInputDelete: { template: '<div />' } },
}))

vi.mock('@planetcrust/human-js', () => ({
  NoID: '0',
  compose: {
    Record: class {
      constructor(mod, opts = {}) {
        this.values = {}
        Object.assign(this, opts)
      }
      setValue(name, value) {
        this.values[name] = value
      }
      serialize() {
        return { values: this.values }
      }
    },
    interpolateTemplate: t => t,
  },
  validator: { IsEmpty: () => false },
}))

vi.mock('@/sections/compose/components/PageBlocks/Grid.vue', () => ({
  default: { template: '<div />' },
}))

import RecordView from './RecordView.vue'

let page
let layouts
let module_
let wrapper

const MODULE_EDIT = 'page.moduleEdit'
const PAGE_BUILDER = 'page.block.general.label.pageBuilder'

async function mountView(namespace) {
  wrapper = mount(RecordView, {
    props: { namespace },
    global: {
      stubs: { teleport: true, Teleport: true },
      directives: { tooltip: {}, focus: {} },
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastErrorHandler: () => vi.fn() },
        $ComposeAPI: {},
        $SystemAPI: { expressionEvaluate: vi.fn(() => Promise.resolve({ results: {} })) },
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

beforeEach(() => {
  page = {
    pageID: 'P1',
    namespaceID: 'N1',
    moduleID: 'M1',
    title: 'Plain page title',
    isRecordPage: true,
    canUpdatePage: true,
    blocks: [],
  }
  layouts = [{ pageLayoutID: 'L1', blocks: [], config: {}, meta: {} }]
  module_ = { moduleID: 'M1', namespaceID: 'N1', fields: [{ name: 'name' }], canUpdateModule: true }
  route.params = { slug: 'ns', pageID: 'P1', recordID: 'R1' }
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
  vi.clearAllMocks()
})

describe('RecordView admin tools', () => {
  it('offers both to a user with manage and both resource rights', async () => {
    await mountView({ namespaceID: 'N1', canManageNamespace: true })

    expect(wrapper.html()).toContain(MODULE_EDIT)
    expect(wrapper.html()).toContain(PAGE_BUILDER)
  })

  it('withholds both from a user without manage, whatever the resources say', async () => {
    await mountView({ namespaceID: 'N1', canManageNamespace: false })

    expect(wrapper.html()).not.toContain(MODULE_EDIT)
    expect(wrapper.html()).not.toContain(PAGE_BUILDER)
  })

  it('withholds the module editor from a user who cannot update the module', async () => {
    module_ = { ...module_, canUpdateModule: false }

    await mountView({ namespaceID: 'N1', canManageNamespace: true })

    expect(wrapper.html()).not.toContain(MODULE_EDIT)
    // The page's own right is untouched by the module's.
    expect(wrapper.html()).toContain(PAGE_BUILDER)
  })

  it('withholds the builder from a user who cannot update the page', async () => {
    page = { ...page, canUpdatePage: false }

    await mountView({ namespaceID: 'N1', canManageNamespace: true })

    expect(wrapper.html()).not.toContain(PAGE_BUILDER)
    expect(wrapper.html()).toContain(MODULE_EDIT)
  })

  it('withholds the module editor on a page that is not a record page', async () => {
    page = { ...page, isRecordPage: false }

    await mountView({ namespaceID: 'N1', canManageNamespace: true })

    expect(wrapper.html()).not.toContain(MODULE_EDIT)
  })
})

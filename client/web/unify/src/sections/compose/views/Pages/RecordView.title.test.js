import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive } from 'vue'

// A page layout may override the record page's title with its own, interpolated
// against the open record (`config.useTitle`). The feature existed in the Vue 2
// compose app and was lost in the Vue 3 port — these pin the restored behaviour,
// including every fallback to the plain page title.

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
  const rec = {
    recordID,
    values,
    ownedBy: 'U9',
    serialize: () => ({ recordID, values }),
  }
  rec.clone = () => makeRecord(recordID, { ...values })
  return rec
}

const records = { R1: makeRecord('R1', { name: 'Acme Corp' }) }

const pageStore = { getByID: () => page }
const pageLayoutStore = { getByPageID: () => layouts }
const moduleStore = {
  getByID: () => ({ moduleID: 'M1', namespaceID: 'N1', fields: [{ name: 'name' }] }),
}
const recordStore = {
  paginationRecordIDs: [],
  findByID: vi.fn(({ recordID }) => Promise.resolve(records[recordID])),
  update: vi.fn(rec => Promise.resolve(rec)),
  create: vi.fn(rec => Promise.resolve(rec)),
}

vi.mock('@planetcrust/human-vue', () => ({
  useModuleStore: () => moduleStore,
  usePageStore: () => pageStore,
  usePageLayoutStore: () => pageLayoutStore,
  useRecordStore: () => recordStore,
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
    // The real implementation, so the test exercises real template semantics
    // (including throwing on a missing record) rather than a stand-in.
    interpolateTemplate(template, { record, user, recordID, ownerID, userID }) {
      const evaluate = new Function(
        'record',
        'user',
        'recordID',
        'ownerID',
        'userID',
        'return `' + template + '`',
      )
      return evaluate(record, user, recordID, ownerID, userID)
    },
  },
  validator: { IsEmpty: () => false },
}))

vi.mock('@/sections/compose/components/PageBlocks/Grid.vue', () => ({
  default: { template: '<div />' },
}))

import RecordView from './RecordView.vue'

let page
let layouts
let wrapper

function layoutWith(config, meta = {}) {
  return [{ pageLayoutID: 'L1', blocks: [], config, meta }]
}

async function mountView() {
  wrapper = mount(RecordView, {
    props: { namespace: { namespaceID: 'N1' } },
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
    blocks: [],
  }
  layouts = layoutWith({})
  route.params = { slug: 'ns', pageID: 'P1', recordID: 'R1' }
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
  vi.clearAllMocks()
})

describe('RecordView title', () => {
  it('uses the page title when the layout does not override it', async () => {
    layouts = layoutWith({ useTitle: false }, { title: 'Layout name' })
    await mountView()

    expect(wrapper.text()).toContain('Plain page title')
    expect(wrapper.text()).not.toContain('Layout name')
  })

  it('uses the layout title when useTitle is on', async () => {
    layouts = layoutWith({ useTitle: true }, { title: 'Custom title' })
    await mountView()

    expect(wrapper.text()).toContain('Custom title')
    expect(wrapper.text()).not.toContain('Plain page title')
  })

  it('interpolates the layout title against the open record', async () => {
    layouts = layoutWith({ useTitle: true }, { title: '${record.values.name} — invoice' })
    await mountView()

    expect(wrapper.text()).toContain('Acme Corp — invoice')
  })

  it('resolves recordID, ownerID and userID in the layout title', async () => {
    layouts = layoutWith({ useTitle: true }, { title: '${recordID}/${ownerID}/${userID}' })
    await mountView()

    expect(wrapper.text()).toContain('R1/U9/U1')
  })

  it('falls back to the page title when useTitle is on but no title is set', async () => {
    layouts = layoutWith({ useTitle: true }, { title: '' })
    await mountView()

    expect(wrapper.text()).toContain('Plain page title')
  })

  it('falls back to the page title when the template cannot be evaluated', async () => {
    layouts = layoutWith({ useTitle: true }, { title: '${record.values.name' })
    await mountView()

    expect(wrapper.text()).toContain('Plain page title')
  })
})

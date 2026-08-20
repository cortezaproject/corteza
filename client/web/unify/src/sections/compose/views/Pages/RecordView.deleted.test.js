import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { computed, reactive, toValue } from 'vue'

// A soft-deleted record still reads back, so the record page opens it and used
// to render it as a live one: no statement that it was deleted, a delete button
// on something already deleted, and no way to put it back.

const route = reactive({
  name: 'page.record',
  params: { slug: 'ns', pageID: 'P1', recordID: 'R1' },
  query: {},
})

const router = { push: vi.fn(), replace: vi.fn() }
const goBack = vi.fn()

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
  onBeforeRouteLeave: () => {},
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

let record

const pageStore = { getByID: () => page }
const pageLayoutStore = { getByPageID: () => layouts }
let recordModule
const moduleStore = { getByID: () => recordModule }
const recordStore = {
  paginationRecordIDs: [],
  findByID: vi.fn(() => Promise.resolve(record)),
  update: vi.fn(rec => Promise.resolve(rec)),
  create: vi.fn(rec => Promise.resolve(rec)),
  delete: vi.fn(() => Promise.resolve(true)),
  undelete: vi.fn(() => Promise.resolve(true)),
}

vi.mock('@planetcrust/human-vue', () => ({
  useDeferredBusy: src => computed(() => !!toValue(src)),
  useModuleStore: () => moduleStore,
  usePageStore: () => pageStore,
  usePageLayoutStore: () => pageLayoutStore,
  useRecordStore: () => recordStore,
  useHistoryBack: () => goBack,
  components: { CInputDelete: { props: ['label'], template: '<div />' } },
}))

vi.mock('@planetcrust/human-js', () => ({
  compose: {
    Record: class {
      constructor(mod, opts = {}) {
        this.values = {}
        Object.assign(this, opts)
      }
      serialize() {
        return { values: this.values }
      }
    },
    PageBlockMaker: i => ({ ...i }),
  },
  validator: { IsEmpty: v => v === undefined || v === null || v === '' },
  NoID: '0',
}))

vi.mock('@/sections/compose/components/PageBlocks/Grid.vue', () => ({
  default: { props: ['blocks', 'namespace', 'page', 'record', 'loading'], template: '<div />' },
}))

import RecordView from './RecordView.vue'

let page
let layouts
let wrapper

function makeRecord(overrides = {}) {
  const rec = {
    recordID: 'R1',
    moduleID: 'M1',
    namespaceID: 'N1',
    values: { status: 'open' },
    ownedBy: 'U1',
    canUpdateRecord: true,
    canDeleteRecord: true,
    canUndeleteRecord: true,
    serialize: () => ({ recordID: 'R1' }),
    ...overrides,
  }
  rec.clone = () => makeRecord(overrides)
  return rec
}

// Toolbar buttons, by the label they carry
function buttonLabels() {
  return wrapper
    .findAll('button')
    .map(b => b.attributes('label'))
    .filter(Boolean)
}

async function mountView() {
  wrapper = mount(RecordView, {
    props: { namespace: { namespaceID: 'N1' } },
    global: {
      stubs: { teleport: true, Teleport: true },
      directives: { tooltip: {}, focus: {} },
      mocks: { $t: k => k },
      provide: {
        $toast: { toastSuccess: vi.fn(), toastWarning: vi.fn(), toastErrorHandler: () => vi.fn() },
        $ComposeAPI: {},
        $SystemAPI: {
          expressionEvaluate: vi.fn(({ expressions }) =>
            Promise.resolve(Object.fromEntries(Object.keys(expressions).map(k => [k, true]))),
          ),
        },
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
  route.params = { slug: 'ns', pageID: 'P1', recordID: 'R1' }
  route.query = {}
  recordModule = {
    moduleID: 'M1',
    namespaceID: 'N1',
    fields: [{ fieldID: 'F-status', name: 'status' }],
  }
  record = makeRecord()
  recordStore.findByID.mockClear()
  recordStore.undelete.mockClear()
  page = {
    pageID: 'P1',
    title: 'Record page',
    isRecordPage: true,
    moduleID: 'M1',
    blocks: [{ blockID: 'B1', kind: 'Record', xywh: [0, 0, 12, 12] }],
    canUpdatePage: false,
  }
  layouts = [{ pageLayoutID: 'L1', blocks: [{ blockID: 'B1', xywh: [0, 0, 12, 12] }], config: {} }]
})

describe('RecordView on a deleted record', () => {
  it('says the record was deleted', async () => {
    record = makeRecord({ deletedAt: '2026-08-20T11:24:49Z' })
    await mountView()

    expect(wrapper.find('message').exists()).toBe(true)
  })

  it('says nothing of the kind on a live record', async () => {
    await mountView()

    expect(wrapper.find('message').exists()).toBe(false)
  })

  it('offers restore in place of delete, and nothing that writes to the record', async () => {
    record = makeRecord({ deletedAt: '2026-08-20T11:24:49Z' })
    await mountView()

    expect(buttonLabels()).toContain('general.label.restore')
    expect(buttonLabels()).not.toContain('general.label.edit')
    expect(buttonLabels()).not.toContain('general.label.saveAsCopy')
    expect(buttonLabels()).not.toContain('general.label.add')
    expect(wrapper.find('c-input-delete-stub').exists()).toBe(false)
  })

  it('leaves a live record its full toolbar', async () => {
    await mountView()

    expect(buttonLabels()).toContain('general.label.edit')
    expect(buttonLabels()).toContain('general.label.saveAsCopy')
    expect(buttonLabels()).toContain('general.label.add')
    expect(buttonLabels()).not.toContain('general.label.restore')
    expect(wrapper.find('c-input-delete-stub').exists()).toBe(true)
  })

  it('withholds restore from someone who may not undelete', async () => {
    record = makeRecord({ deletedAt: '2026-08-20T11:24:49Z', canUndeleteRecord: false })
    await mountView()

    expect(buttonLabels()).not.toContain('general.label.restore')
  })

  it('does not open in edit mode for an ?edit=1 that outlived the delete', async () => {
    record = makeRecord({ deletedAt: '2026-08-20T11:24:49Z' })
    route.query = { edit: '1' }
    await mountView()

    expect(buttonLabels()).not.toContain('general.label.save')
    expect(buttonLabels()).toContain('general.label.restore')
  })

  it('restores through the store and reads the record back', async () => {
    record = makeRecord({ deletedAt: '2026-08-20T11:24:49Z' })
    await mountView()
    const readsBefore = recordStore.findByID.mock.calls.length

    const restore = wrapper
      .findAll('button')
      .find(b => b.attributes('label') === 'general.label.restore')
    await restore.trigger('click')
    await flushPromises()

    expect(recordStore.undelete).toHaveBeenCalledWith({
      namespaceID: 'N1',
      moduleID: 'M1',
      recordID: 'R1',
    })
    expect(recordStore.findByID.mock.calls.length).toBeGreaterThan(readsBefore)
  })
})

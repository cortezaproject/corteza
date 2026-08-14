import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive } from 'vue'

// Layout conditions are evaluated against the record, so the layout must be
// picked after the record is resolved — and re-picked whenever the record or
// the mode changes. These are the paths that silently skip loadPage().

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
    ownedBy: 'U1',
    serialize: () => ({ recordID, values }),
  }
  rec.clone = () => makeRecord(recordID, { ...values })
  return rec
}

const records = {
  R1: makeRecord('R1', { status: 'open' }),
  R2: makeRecord('R2', { status: 'closed' }),
}

const pageStore = { getByID: () => page }
const pageLayoutStore = { getByPageID: () => layouts }
const moduleStore = {
  getByID: () => ({ moduleID: 'M1', namespaceID: 'N1', fields: [{ name: 'status' }] }),
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
  useHistoryBack: () => vi.fn(),
  components: { CInputDelete: { template: '<div />' } },
}))

vi.mock('@planetcrust/human-js', () => {
  // Stands in for the page-block registry: a kind whose class carries a method
  // the renderer calls, so a block that lost its prototype is visible in a test.
  class PageBlockRecordRevisions {
    fetch() {
      return Promise.resolve([])
    }
  }

  return {
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
      PageBlockMaker: i =>
        i.kind === 'RecordRevisions' ? Object.assign(new PageBlockRecordRevisions(), i) : { ...i },
    },
    validator: { IsEmpty: () => false },
  }
})

vi.mock('@/sections/compose/components/PageBlocks/Grid.vue', () => ({
  default: { props: ['blocks', 'namespace', 'page', 'record'], template: '<div />' },
}))

import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'
import RecordView from './RecordView.vue'

let page
let layouts
let expressionEvaluate
let wrapper

// Variables of every layout-condition evaluation, in order
function layoutEvaluations() {
  return expressionEvaluate.mock.calls
    .filter(([{ expressions }]) => Object.keys(expressions).some(id => id.startsWith('L')))
    .map(([{ variables }]) => variables)
}

// Variables of every block-condition evaluation, in order
function blockEvaluations() {
  return expressionEvaluate.mock.calls
    .filter(([{ expressions }]) => Object.keys(expressions).some(id => id.startsWith('B')))
    .map(([{ variables }]) => variables)
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
        $SystemAPI: { expressionEvaluate },
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

// Mounted views keep reacting to the shared route object; leaving them alive
// would let earlier tests answer later route changes.
afterEach(() => {
  wrapper?.unmount()
  wrapper = null
})

beforeEach(() => {
  route.params = { slug: 'ns', pageID: 'P1', recordID: 'R1' }
  route.query = {}
  expressionEvaluate = vi.fn(({ expressions }) =>
    Promise.resolve(Object.fromEntries(Object.keys(expressions).map(k => [k, true]))),
  )
  recordStore.findByID.mockClear()
  page = {
    pageID: 'P1',
    title: 'Record page',
    isRecordPage: true,
    moduleID: 'M1',
    blocks: [],
    canUpdatePage: false,
  }
  layouts = [
    {
      pageLayoutID: 'L1',
      blocks: [],
      config: { visibility: { expression: 'record.values.status == "open"', roles: [] } },
    },
  ]
})

describe('RecordView layout resolution', () => {
  it('evaluates layout conditions against the loaded record, not an empty one', async () => {
    await mountView()

    const evaluations = layoutEvaluations()
    expect(evaluations).toHaveLength(1)
    expect(evaluations[0].record).toEqual({ recordID: 'R1', values: { status: 'open' } })
  })

  it('re-evaluates when fast-swapping to another record on the same page', async () => {
    await mountView()
    expect(layoutEvaluations()).toHaveLength(1)

    // Same page, another record: this path calls loadRecord() and skips loadPage()
    route.params = { ...route.params, recordID: 'R2' }
    await flushPromises()

    const evaluations = layoutEvaluations()
    expect(evaluations).toHaveLength(2)
    expect(evaluations[1].record).toEqual({ recordID: 'R2', values: { status: 'closed' } })
  })

  it('re-evaluates when switching from view to edit', async () => {
    await mountView()
    expect(layoutEvaluations()[0].isView).toBe(true)

    route.query = { edit: '1' }
    await flushPromises()

    const evaluations = layoutEvaluations()
    expect(recordStore.findByID).toHaveBeenCalledTimes(1)
    expect(evaluations).toHaveLength(2)
    expect(evaluations[1].isEdit).toBe(true)
    expect(evaluations[1].record).toEqual({ recordID: 'R1', values: { status: 'open' } })
  })

  it('re-evaluates after the record is saved', async () => {
    await mountView()
    route.query = { edit: '1' }
    await flushPromises()

    const before = layoutEvaluations().length
    wrapper.vm.record.values.status = 'closed'
    await wrapper.vm.handleSave({ valid: true })
    await flushPromises()

    const evaluations = layoutEvaluations()
    expect(evaluations.length).toBeGreaterThan(before)
    expect(evaluations.at(-1).record.values.status).toBe('closed')
  })

  describe('with a layout requested via ?layoutID', () => {
    beforeEach(() => {
      layouts = [
        {
          pageLayoutID: 'L1',
          blocks: [],
          config: { visibility: { expression: 'record.values.status == "open"', roles: [] } },
        },
        {
          pageLayoutID: 'L2',
          blocks: [],
          config: { visibility: { expression: 'record.values.status == "closed"', roles: [] } },
        },
      ]
      route.query = { layoutID: 'L2' }
    })

    it('opens it when its own condition passes', async () => {
      expressionEvaluate.mockImplementation(() => Promise.resolve({ L1: true, L2: true }))
      await mountView()

      expect(wrapper.vm.layout.pageLayoutID).toBe('L2')
    })

    it('falls back to the normal layout order when its condition fails', async () => {
      expressionEvaluate.mockImplementation(() => Promise.resolve({ L1: true, L2: false }))
      await mountView()

      expect(wrapper.vm.layout.pageLayoutID).toBe('L1')
    })
  })

  it('leaves the layout alone but re-evaluates blocks as values are edited', async () => {
    page.blocks = [
      {
        blockID: 'B1',
        kind: 'Content',
        meta: { visibility: { expression: 'record.values.status == "open"', roles: [] } },
      },
    ]
    await mountView()
    expect(layoutEvaluations()).toHaveLength(1)
    expect(blockEvaluations()).toHaveLength(1)

    wrapper.vm.record.values.status = 'closed'
    // Block visibility is debounced (300ms) so typing doesn't fire per keystroke
    await new Promise(resolve => setTimeout(resolve, 350))
    await flushPromises()

    const blocks = blockEvaluations()
    expect(blocks).toHaveLength(2)
    expect(blocks[1].record.values.status).toBe('closed')
    expect(layoutEvaluations()).toHaveLength(1)
  })

  it('builds layout blocks through PageBlockMaker so they keep their class methods', async () => {
    // A spread carries options but drops the prototype, and the block renderers
    // call methods on it (block.fetch, block.reorderViews).
    page.blocks = [{ blockID: 'B1', kind: 'RecordRevisions', options: { preload: true } }]
    layouts[0].blocks = [{ blockID: 'B1', xywh: [0, 0, 24, 20] }]

    await mountView()

    const [block] = wrapper.findComponent(Grid).props('blocks')
    expect(typeof block.fetch).toBe('function')
    expect(block.options).toEqual({ preload: true })
    expect(block.xywh).toEqual([0, 0, 24, 20])
  })

  it('does not re-evaluate while field values are edited', async () => {
    await mountView()
    expect(layoutEvaluations()).toHaveLength(1)

    wrapper.vm.record.values.status = 'closed'
    await flushPromises()

    expect(layoutEvaluations()).toHaveLength(1)
  })
})

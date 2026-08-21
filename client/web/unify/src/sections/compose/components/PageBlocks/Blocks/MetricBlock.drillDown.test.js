import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// A metric's filter is authored as a template — `ownedBy = ${userID}`. The
// filter it drills into a record list with has to be the evaluated one: the
// list interpolates the prefilter it was configured with, never the one it is
// handed, and an unevaluated ${...} reaches the server as an illegal token.

vi.mock('vue-router', () => ({ useRoute: () => ({ name: 'page', params: {} }) }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('./PageBlock.vue', () => ({
  default: { props: ['block', 'record'], template: '<div><slot /></div>' },
}))
// The dialog's list, captured so its options can be asserted on. It is loaded
// lazily by MetricBlock, so the module is mocked whole — a name stub never
// reaches an async component, and the real one wants a page around it.
const drilled = vi.hoisted(() => ({ block: null }))
vi.mock('./RecordListBlock.vue', () => ({
  __esModule: true,
  default: {
    name: 'RecordListBlock',
    props: ['block'],
    setup(props) {
      drilled.block = props.block
      return () => null
    },
  },
}))

vi.mock('@planetcrust/human-vue', () => ({
  components: {},
}))

import MetricBlock from './MetricBlock.vue'

const USER_ID = '506815840423837697'
const LIST_BLOCK = '_list'

let emitted
let eventBus
let recordReport

function metricBlock(filter, extra = {}) {
  return {
    blockID: 'B1',
    kind: 'Metric',
    title: '',
    options: {
      metrics: [
        {
          label: 'Owned by me',
          moduleID: '510151835201437697',
          metricField: 'count',
          operation: 'count',
          filter,
          ...extra,
          drillDown: { enabled: true, blockID: LIST_BLOCK, recordListOptions: { fields: [] } },
        },
      ],
    },
    // The component reads its value through the block class
    fetch: vi.fn(({ m }) => {
      seenFilters.push(m.filter)
      return Promise.resolve([{ value: 4 }])
    }),
  }
}

let seenFilters

beforeEach(() => {
  emitted = []
  seenFilters = []
  eventBus = { on: () => () => {}, emit: (name, payload) => emitted.push({ name, payload }) }
  recordReport = vi.fn(() => Promise.resolve([{ count: 4 }]))
})

function mountBlock(filter, extra) {
  return mount(MetricBlock, {
    props: {
      block: metricBlock(filter, extra),
      namespace: { namespaceID: 'N1' },
      page: { pageID: '0' },
    },
    global: {
      stubs: {
        ProgressSpinner: true,
        Button: true,
        // Renders its slot, so the drilled-into list actually mounts
        Dialog: { template: '<div class="dialog"><slot /></div>' },
      },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $ComposeAPI: { recordReport },
        $Auth: { user: { userID: USER_ID } },
        $eventBus: eventBus,
      },
    },
  })
}

const drill = () => emitted.find(e => e.name === `drill-down-recordList:${LIST_BLOCK}`)

// The block holds its spinner for 300ms after the values land
async function settle(w) {
  await flushPromises()
  await new Promise(resolve => setTimeout(resolve, 350))
  await flushPromises()
  return w
}

describe('MetricBlock drill-down', () => {
  it('drills with the evaluated filter, not the template', async () => {
    const w = await settle(mountBlock('ownedBy = ${userID}'))

    await w.find('.metric-item').trigger('click')

    expect(drill()).toBeTruthy()
    expect(drill().payload.prefilter).toBe(`(ownedBy = ${USER_ID})`)
    expect(drill().payload.prefilter).not.toContain('${')
  })

  it('reads its own value through the same evaluated filter', async () => {
    mountBlock('ownedBy = ${userID}')
    await flushPromises()

    expect(seenFilters).toEqual([`ownedBy = ${USER_ID}`])
  })

  it('passes a plain filter through untouched', async () => {
    const w = await settle(mountBlock('createdAt >= DATE_SUB(NOW(), INTERVAL 7 DAY)'))

    await w.find('.metric-item').trigger('click')

    expect(drill().payload.prefilter).toBe('(createdAt >= DATE_SUB(NOW(), INTERVAL 7 DAY))')
  })

  it('sends no filter at all for the metric that counts everything', async () => {
    const w = await settle(mountBlock(''))

    await w.find('.metric-item').trigger('click')

    // An empty prefilter is what clears the table's drill-down filter
    expect(drill().payload.prefilter).toBe('')
  })
})

// Without a drillDown.blockID there is no list to filter, so the records open
// in the block's own dialog instead.
function modalBlock(extra = {}) {
  const b = metricBlock('', extra)
  b.options.metrics[0].drillDown = { enabled: true, recordListOptions: { fields: [] } }
  return b
}

function mountModal(extra) {
  drilled.block = null
  return mount(MetricBlock, {
    props: { block: modalBlock(extra), namespace: { namespaceID: 'N1' }, page: { pageID: '0' } },
    global: {
      stubs: {
        ProgressSpinner: true,
        Button: true,
        Dialog: {
          props: ['header'],
          template:
            '<div class="dialog"><span class="dialog-header">{{ header }}</span><slot /></div>',
        },
      },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $ComposeAPI: { recordReport },
        $Auth: { user: { userID: USER_ID } },
        $eventBus: eventBus,
      },
    },
  })
}

// The dialog's list is loaded lazily, so it mounts a tick after the click
async function openDialog(w) {
  await w.find('.metric-item').trigger('click')
  await flushPromises()
  return w
}

describe('MetricBlock drill-down dialog', () => {
  it('names the dialog once when the block has no title of its own', async () => {
    const w = await openDialog(await settle(mountModal()))

    // The fallback title is the metric label, which is also the drilled value
    expect(w.find('.dialog-header').text()).toBe('Owned by me')
  })

  it('opens the records in a dialog when no list is named', async () => {
    await openDialog(await settle(mountModal()))

    expect(drill()).toBeUndefined()
    expect(drilled.block).toBeTruthy()
    expect(drilled.block.options.moduleID).toBe('510151835201437697')
  })

  it('opens a deleted-records metric on the deleted records', async () => {
    await openDialog(await settle(mountModal({ deleted: 2 })))

    const { options } = drilled.block
    expect(options.showDeletedRecordsInitially).toBe(true)
    expect(options.showDeletedRecordsOption).toBe(true)
    // A record made here would land among the existing ones
    expect(options.hideAddButton).toBe(true)
  })

  it('leaves an ordinary metric drilling into existing records', async () => {
    await openDialog(await settle(mountModal()))

    const { options } = drilled.block
    expect(options.showDeletedRecordsInitially).toBe(false)
    expect(options.showDeletedRecordsOption).toBe(false)
    expect(options.hideAddButton).toBe(false)
  })
})

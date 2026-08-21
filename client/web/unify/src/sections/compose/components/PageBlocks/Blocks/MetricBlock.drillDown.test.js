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
vi.mock('./RecordListBlock.vue', () => ({ default: { template: '<div />' } }))

vi.mock('@planetcrust/human-vue', () => ({
  components: {},
}))

import MetricBlock from './MetricBlock.vue'

const USER_ID = '506815840423837697'
const LIST_BLOCK = '_list'

let emitted
let eventBus
let recordReport

function metricBlock(filter) {
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

function mountBlock(filter) {
  return mount(MetricBlock, {
    props: {
      block: metricBlock(filter),
      namespace: { namespaceID: 'N1' },
      page: { pageID: '0' },
    },
    global: {
      stubs: { ProgressSpinner: true, Dialog: true, Button: true },
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

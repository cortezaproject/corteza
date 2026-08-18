import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

vi.mock('./PageBlock.vue', () => ({
  default: { template: '<div><slot /></div>' },
}))

vi.mock('./Metric/MetricItem.vue', () => ({
  default: { template: '<div />' },
}))

import MetricBlock from './MetricBlock.vue'

// A metric's transform formula is author-typed and reaches `new Function` — so
// like its filter, it is interpolated before it gets there. `v * ${...}` left
// uninterpolated is a syntax error, not a number.

const USER = { userID: '42', name: 'Ada' }
const RECORD = { recordID: '7', ownedBy: '3', values: { rate: 1.5 } }

let fetched

function mountBlock(metrics, record) {
  fetched = []

  return mount(MetricBlock, {
    props: {
      block: {
        blockID: 'B1',
        title: '',
        options: { metrics },
        fetch: ({ m }) => {
          fetched.push(m)
          return Promise.resolve([{ label: '', value: 1 }])
        },
      },
      namespace: { namespaceID: 'N1' },
      record,
    },
    global: {
      stubs: { ProgressSpinner: true, Dialog: true, Button: true },
      directives: { tooltip: {} },
      provide: { $ComposeAPI: {}, $Auth: { user: USER }, $eventBus: null },
      mocks: { $t: k => k },
    },
  })
}

const metric = extra => ({ moduleID: 'M1', field: 'v', operation: 'sum', ...extra })

beforeEach(() => {
  vi.spyOn(console, 'warn').mockImplementation(() => {})
})

describe('MetricBlock interpolation', () => {
  it('interpolates the transform formula against the record', async () => {
    mountBlock([metric({ transformFx: 'v * ${record.values.rate}' })], RECORD)
    await flushPromises()

    expect(fetched[0].transformFx).toBe('v * 1.5')
  })

  it('interpolates the filter and the formula together', async () => {
    mountBlock(
      [metric({ filter: 'recordID = ${recordID}', transformFx: 'v + ${ownerID}' })],
      RECORD,
    )
    await flushPromises()

    expect(fetched[0].filter).toBe('recordID = 7')
    expect(fetched[0].transformFx).toBe('v + 3')
  })

  it('skips a metric whose formula needs a record there is none of', async () => {
    mountBlock([metric({ transformFx: 'v * ${record.values.rate}' })])
    await flushPromises()

    expect(fetched).toHaveLength(0)
  })

  it('leaves a plain formula alone', async () => {
    mountBlock([metric({ transformFx: 'v / 100' })], RECORD)
    await flushPromises()

    expect(fetched[0].transformFx).toBe('v / 100')
  })
})

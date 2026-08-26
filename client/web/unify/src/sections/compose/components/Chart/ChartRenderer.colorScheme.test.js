import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// The colour tables lib/js compiles in are only half a chart's palette; the
// other half is whatever an admin defined for this instance. Without them a
// chart naming a custom scheme resolves to no palette, and echarts draws its
// legend and none of its series — a blank panel that reads as "no data".

const madeWith = []

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

vi.mock('@planetcrust/human-vue', () => ({
  useModuleStore: () => ({ getByID: () => undefined }),
  useRecordStore: () => ({ getByID: () => undefined, resolveRecordLabels: vi.fn() }),
  components: { CChart: { name: 'CChart', props: ['chart'], template: '<div />' } },
}))

vi.mock('../../lib/charts', () => ({
  chartConstructor: () => ({
    isValid: () => true,
    fetchReports: () => Promise.resolve({ labels: [], datasets: [] }),
    makeOptions: data => {
      madeWith.push(data)
      return {}
    },
  }),
}))

import ChartRenderer from './ChartRenderer.vue'

const schemes = [{ id: 'custom-1', name: 'Brand', colors: ['#FF00AA'] }]

function render($Settings) {
  return mount(ChartRenderer, {
    props: {
      chart: { chartID: 'C1', config: { colorScheme: 'custom-1', reports: [{ moduleID: 'M1' }] } },
      reporter: () => Promise.resolve([]),
    },
    global: {
      stubs: { ProgressSpinner: true },
      provide: { $Settings },
    },
  })
}

beforeEach(() => {
  madeWith.length = 0
})

describe('ChartRenderer', () => {
  it("hands the chart the instance's own colour schemes", async () => {
    render({ get: (k, d) => (k === 'ui.charts.colorSchemes' ? schemes : d) })
    await flushPromises()

    expect(madeWith).toHaveLength(1)
    expect(madeWith[0].customColorSchemes).toEqual(schemes)
  })

  // A block can render before settings have answered, and getColorschemeColors
  // reaches straight for .find() on whatever it is given.
  it('hands it an empty list where settings never loaded', async () => {
    render(undefined)
    await flushPromises()

    expect(madeWith[0].customColorSchemes).toEqual([])
  })
})

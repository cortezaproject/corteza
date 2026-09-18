import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// echarts draws on a canvas, and a canvas takes a concrete family or nothing:
// a font naming 'inherit' is rejected whole, size included, and the text falls
// back to the canvas default of 10px sans-serif.

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

function render() {
  return mount(ChartRenderer, {
    props: {
      chart: { chartID: 'C1', config: { reports: [{ moduleID: 'M1' }] } },
      reporter: () => Promise.resolve([]),
    },
    global: { stubs: { ProgressSpinner: true } },
  })
}

beforeEach(() => {
  madeWith.length = 0
})

afterEach(() => {
  document.body.style.fontFamily = ''
})

describe('ChartRenderer font', () => {
  it("hands the chart the page's own font family", async () => {
    document.body.style.fontFamily = 'Poppins, sans-serif'
    render()
    await flushPromises()

    expect(madeWith[0].themeVariables['font-regular']).toBe('Poppins, sans-serif')
  })

  it('falls back to a family a canvas accepts', async () => {
    render()
    await flushPromises()

    expect(madeWith[0].themeVariables['font-regular']).toBe('sans-serif')
  })
})

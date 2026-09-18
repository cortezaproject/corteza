import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'

// How many metrics a report takes is up to the chart type's editor: bar, line,
// pie, doughnut, scatter and radar take any number, funnel and gauge take one.
// Each editor is mounted the way the chart editor mounts it — chart and
// modules, nothing else.

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k }),
}))

vi.mock('@/sections/compose/composables/useExpressionScope', () => ({
  useExpressionScope: () => ({ scope: ref({}) }),
}))

import { compose } from '@planetcrust/human-js'
import { GenericChart, RadarChart, FunnelChart, GaugeChart } from './index.js'

const PanelStub = {
  name: 'Panel',
  template: '<div><slot name="icons" /><slot /></div>',
}

const ButtonStub = {
  name: 'Button',
  props: { label: { type: String, default: '' } },
  emits: ['click'],
  template: '<button @click="$emit(\'click\')">{{ label }}</button>',
}

function mountEditor(editor, chart, metricCount) {
  const report = ref({
    ...chart.defReport(),
    moduleID: 'M1',
    metrics: Array.from({ length: metricCount }, () => chart.defMetric()),
  })
  const wrapper = mount(editor, {
    props: {
      chart,
      modules: [{ moduleID: 'M1', fields: [] }],
    },
    global: {
      provide: { reportDraft: report },
      stubs: { Panel: PanelStub, Button: ButtonStub, NumberFormatting: true },
      mocks: { $t: k => k },
      config: { warnHandler: () => {} },
    },
  })
  return { wrapper, report }
}

const addButton = wrapper =>
  wrapper.findAll('button').find(b => b.text().includes('chart.edit.metric.add'))

describe('metrics per chart type', () => {
  it.each([
    ['GenericChart', GenericChart, compose.Chart],
    ['RadarChart', RadarChart, compose.RadarChart],
  ])('%s offers another metric past the first', async (_, editor, Chart) => {
    const { wrapper, report } = mountEditor(editor, new Chart(), 2)

    const add = addButton(wrapper)
    expect(add).toBeDefined()

    await add.trigger('click')
    expect(report.value.metrics).toHaveLength(3)
  })

  it.each([
    ['FunnelChart', FunnelChart, compose.FunnelChart],
    ['GaugeChart', GaugeChart, compose.GaugeChart],
  ])('%s stops at one metric', (_, editor, Chart) => {
    const { wrapper } = mountEditor(editor, new Chart(), 1)

    expect(addButton(wrapper)).toBeUndefined()
  })
})

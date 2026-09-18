import { expect } from 'chai'

import Chart from './chart'
import FunnelChart from './funnel'
import GaugeChart from './gauge'
import RadarChart from './radar'
import { chartFontSize } from './util'

/**
 * echarts hard-codes 12px on axis labels, legends and gauge text, and none of
 * those follow the option's global textStyle. Every piece of chart text has to
 * name the size itself or it stays at echarts' default.
 */

const rows = [
  { dimension_0: 'Proposal', count: 3, sum_amount: 300 },
  { dimension_0: 'Negotiation', count: 2, sum_amount: 200 },
]

const config = (type: string, extra: Record<string, unknown> = {}) => ({
  config: {
    reports: [
      {
        moduleID: '1',
        dimensions: [{ field: 'stage', modifier: '(no grouping / buckets)', ...extra }],
        metrics: [{ field: 'count', type }],
        yAxis: { label: 'Deals' },
      },
    ],
  },
})

const optionsOf = async (chart: any) => {
  const data = await chart.fetchReports({ reporter: () => Promise.resolve(rows) })
  return chart.makeOptions({ ...data, themeVariables: { 'font-regular': 'Poppins, sans-serif' } })
}

describe('chart text size', () => {
  it('sizes a bar chart’s text, legend and both axes', async () => {
    const o = await optionsOf(new Chart(config('bar') as any))

    expect(o.textStyle).to.include({ fontSize: chartFontSize, fontFamily: 'Poppins, sans-serif' })
    expect(o.legend.textStyle.fontSize).to.equal(chartFontSize)
    expect(o.legend.pageTextStyle.fontSize).to.equal(chartFontSize)
    expect(o.xAxis[0].axisLabel.fontSize).to.equal(chartFontSize)
    expect(o.yAxis[0].axisLabel.fontSize).to.equal(chartFontSize)
  })

  it('keeps a rotated axis rotated at the same size', async () => {
    const o = await optionsOf(new Chart(config('bar', { rotateLabel: '45' }) as any))

    expect(o.xAxis[0].axisLabel).to.include({ rotate: '45', fontSize: chartFontSize })
  })

  it('sizes a funnel’s and a radar’s text and legend', async () => {
    for (const chart of [
      new FunnelChart(config('funnel') as any),
      new RadarChart(config('radar') as any),
    ]) {
      const o = await optionsOf(chart)

      expect(o.textStyle.fontSize).to.equal(chartFontSize)
      expect(o.legend.textStyle.fontSize).to.equal(chartFontSize)
    }
  })

  it('reads a gauge’s value larger than its title', async () => {
    const o = await optionsOf(
      new GaugeChart(config('gauge', { meta: { steps: [{ value: 5 }, { value: 10 }] } }) as any),
    )
    const [gauge] = o.series

    expect(o.textStyle.fontSize).to.equal(chartFontSize)
    expect(gauge.title.fontSize).to.equal(chartFontSize)
    expect(gauge.detail.fontSize).to.be.above(chartFontSize)
  })
})

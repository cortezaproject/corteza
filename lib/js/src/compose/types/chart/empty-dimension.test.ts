import { expect } from 'chai'

import Chart from './chart'
import GaugeChart from './gauge'
import { formatChartValue } from './util'

/**
 * Records that hold no value for the dimension field aggregate into their own
 * group, which the report returns with a null dimension value (and a null
 * aggregate, when the metric is one). Nothing downstream may turn that null into
 * a zero, and nothing may render it as "NaN".
 *
 * Regression cover for crusttech/human#39: a doughnut over an all-empty numeric
 * column drew its slice labelled "NaN", while "Default value" and "Skip missing
 * values" both did nothing because the empty group had already become a real 0.
 */

// Drives BaseChart.processReporterResults through the one public entry point it
// has, standing in for the API call with a fixed rowset.
const reportWith = (chart: any, rows: Array<Record<string, unknown>>) =>
  chart.fetchReports({ reporter: () => Promise.resolve(rows) })

const doughnut = (dimension: Record<string, unknown>) =>
  new Chart({
    config: {
      reports: [
        {
          moduleID: '1',
          dimensions: [{ field: 'amount', modifier: '(no grouping / buckets)', ...dimension }],
          metrics: [{ field: 'amount', aggregate: 'SUM', type: 'doughnut' }],
        },
      ],
    },
  } as any)

describe('empty dimension groups', () => {
  describe('formatChartValue', () => {
    it('should never render the literal string NaN', () => {
      for (const v of ['', '   ', null, undefined]) {
        expect(formatChartValue(v as any)).to.not.contain('NaN')
      }
    })

    it('should not render null as the string "null"', () => {
      expect(formatChartValue(null as any)).to.not.contain('null')
    })

    it('should still format a real zero as zero', () => {
      expect(formatChartValue(0)).to.contain('0')
      expect(formatChartValue(0)).to.not.contain('NaN')
    })

    it('should still format ordinary numbers', () => {
      expect(formatChartValue(1234)).to.contain('1,234')
      expect(formatChartValue('42')).to.contain('42')
    })

    it('should pass non-numeric text through untouched', () => {
      expect(formatChartValue('Alpha')).to.contain('Alpha')
    })

    it('should honour prefix and suffix on a zero', () => {
      const out = formatChartValue(0, { prefix: '$', suffix: 'EUR' })
      expect(out).to.contain('$')
      expect(out).to.contain('EUR')
      expect(out).to.not.contain('NaN')
    })
  })

  describe('label selection', () => {
    it('should label an empty group with the sentinel the renderer translates', async () => {
      const data: any = await reportWith(doughnut({}), [
        { dimension_0: null, sum_amount: null, count: 3 },
      ])

      expect(data.labels).to.deep.equal(['undefined'])
    })

    it('should label an empty group with the configured default value', async () => {
      const data: any = await reportWith(doughnut({ default: 'No amount' }), [
        { dimension_0: null, sum_amount: null, count: 3 },
      ])

      expect(data.labels).to.deep.equal(['No amount'])
    })

    it('should keep a real zero as its own label, not the default', async () => {
      const data: any = await reportWith(doughnut({ default: 'No amount' }), [
        { dimension_0: null, sum_amount: null, count: 3 },
        { dimension_0: 0, sum_amount: 0, count: 1 },
        { dimension_0: 5, sum_amount: 5, count: 1 },
      ])

      expect(data.labels).to.deep.equal(['No amount', 0, 5])
    })

    it('should drop the empty group when skipMissing is set, keeping a real zero', async () => {
      const data: any = await reportWith(doughnut({ skipMissing: true }), [
        { dimension_0: null, sum_amount: null, count: 3 },
        { dimension_0: 0, sum_amount: 0, count: 1 },
        { dimension_0: 5, sum_amount: 5, count: 1 },
      ])

      expect(data.labels).to.deep.equal([0, 5])
    })
  })

  describe('metric values', () => {
    it('should not put the dimension default text into the data array', async () => {
      const data: any = await reportWith(doughnut({ default: 'No amount' }), [
        { dimension_0: null, sum_amount: null, count: 3 },
      ])

      expect(data.datasets[0].data).to.deep.equal([null])
    })

    it('should keep a real zero metric as zero', async () => {
      const data: any = await reportWith(doughnut({}), [
        { dimension_0: 0, sum_amount: 0, count: 1 },
        { dimension_0: 5, sum_amount: 5, count: 1 },
      ])

      expect(data.datasets[0].data).to.deep.equal([0, 5])
    })
  })

  describe('gauge totals', () => {
    // parseFloat(null) is NaN while isNaN(null) is false, so an empty group used
    // to make the whole gauge total NaN.
    it('should ignore an empty group instead of poisoning the total', async () => {
      const chart = new GaugeChart({
        config: {
          reports: [
            {
              moduleID: '1',
              dimensions: [
                {
                  field: 'amount',
                  modifier: '(no grouping / buckets)',
                  meta: { steps: [{ label: 'low', value: '10' }] },
                },
              ],
              metrics: [{ field: 'amount', aggregate: 'SUM', type: 'gauge' }],
            },
          ],
        },
      } as any)

      const data: any = await reportWith(chart, [
        { dimension_0: null, sum_amount: null, count: 3 },
        { dimension_0: 5, sum_amount: 5, count: 1 },
      ])

      expect(data.datasets[0].value).to.equal('5.000')
    })
  })
})

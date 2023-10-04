import { BaseChart, PartialChart } from './base'
import {
  Metric,
  Report,
  Dimension,
  ChartType,
  TemporalDataPoint,
} from './util'
import { getColorschemeColors } from '../../../shared'

export default class GaugeChart extends BaseChart {
  constructor (def: PartialChart = {}) {
    super(def)

    // Assure required fields
    for (const v of (this.config.reports || []) as Array<Report>) {
      for (const d of (v.dimensions || []) as Array<Dimension>) {
        // Since gauge produces one value we want one dataset, deletedAt is the same for all existing records
        d.field = 'deletedAt'

        if (!d.meta) {
          d.meta = {}
        }

        if (!d.meta.steps) {
          d.meta.steps = []
        }
      }
    }
  }

  // Gauge charts (at the moment) support only 1 report per chart
  async fetchReports (a: any) {
    return super.fetchReports(a).then((rr: any) => {
      return rr[0]
    })
  }

  processLabels (ll: Array<string>, d: Dimension) {
    return (d.meta?.steps || []).map(({ label }: any) => label)
  }

  makeDataset (m: Metric, d: Dimension, data: Array<number|TemporalDataPoint>, alias: string) {
    const steps = (d.meta?.steps || [])

    data = this.datasetPostProc(data, m)

    const value = data.reduce((acc: any, cur: any) => {
      return !isNaN(cur) ? acc + parseFloat(cur) : acc
    }, 0).toFixed(3)

    const max = Math.max(...steps.map(({ value }: any) => parseFloat(value)))

    const sortedSteps = [...steps].sort((a: any, b: any) => {
      return parseFloat(b.value) - parseFloat(a.value)
    })

    const { label: name } = sortedSteps.reduce((acc: any, cur: any) => {
      const curValue = parseFloat(cur.value)
      return value < curValue ? cur : acc
    }, sortedSteps[0] || {})

    return {
      steps,
      name,
      max,
      value,
      startAngle: m.startAngle,
      endAngle: m.endAngle,
      tooltip: {
        fixed: m.fixTooltips,
      },
    }
  }

  makeOptions (data: any) {
    const { colorScheme, noAnimation = false, toolbox } = this.config
    const { saveAsImage } = toolbox || {}

    const { datasets = [] } = data
    const { steps = [], name, value, max, tooltip, startAngle, endAngle } = datasets.find(({ value }: any) => value) || datasets[0]
    const colors = getColorschemeColors(colorScheme, data.customColorSchemes)

    const color = steps.map((s: any, i: number) => {
      return [s.value / max, colors[i]]
    })

    return {
      animation: !noAnimation,
      textStyle: {
        fontFamily: 'Poppins-Regular',
      },
      toolbox: {
        feature: {
          saveAsImage: saveAsImage ? {
            name: this.name
          } : undefined,
        },
        top: 15,
        right: 5,
      },
      grid: {
        bottom: 0,
      },
      series: [
        {
          type: 'gauge',
          startAngle,
          endAngle,
          min: 0,
          max,
          splitNumber: 5,
          radius: '100%',
          center: ['50%', '50%'],
          pointer: {
            width: 5,
            length: '75%',
            itemStyle: {
              color: '#464646',
            },
          },
          splitLine: {
            distance: 0,
            length: 0,
            lineStyle: {
              color: '#FFFFFF',
            },
          },
          axisLine: {
            lineStyle: {
              width: 30,
              color,
            },
          },
          axisTick: {
            show: false,
            distance: -30,
          },
          axisLabel: {
            show: false,
            distance: 60,
          },
          title: {
            fontSize: 14,
            show: tooltip.fixed,
            offsetCenter: [0, '30%'],
          },
          detail: {
            fontSize: 13,
            offsetCenter: [0, '55%'],
            valueAnimation: true,
          },
          data: [
            {
              name,
              value,
            },
          ],
        },
      ],
    }
  }

  baseChartType (): string {
    return 'gauge'
  }

  defMetrics (): Metric {
    return Object.assign({}, {
      type: ChartType.gauge,
      startAngle: 200,
      endAngle: -20
    })
  }

  /**
   * Checks validity of dimensions.
   * If invalid it throws an error
   */
  dimCheck ({ meta }: Dimension): void | Error {
    if ((meta?.steps || []).length === 0) {
      throw new Error('notification.chart.invalidConfig.missingDimensionsSteps')
    }
  }

  /**
   * Since gauge charts always define one type, this check can be simplified
   */
  mtrCheck ({ field, aggregate }: Metric): void | Error {
    if (!field) {
      throw new Error('notification.chart.invalidConfig.missingMetricsField')
    }
    if (field !== 'count' && !aggregate) {
      throw new Error('notification.chart.invalidConfig.missingMetricsAggregate')
    }
  }
}

import { BaseChart, PartialChart } from './base'
import {
  Metric,
  Report,
  Dimension,
  ChartType,
  makeDataLabel,
  calculatePercentages,
} from './util'
import { makeTipper } from './chartjs/plugins'
import { defaultBGColor } from './common'

const ChartJS = require('chart.js')

/**
 * Gauge chart provides the definitions for the chartjs-plugin-funnel plugin.
 */
export default class GaugeChart extends BaseChart {
  constructor (def: PartialChart = {}) {
    super(def)

    // Assure required fields
    for (const v of (this.config.reports || []) as Array<Report>) {
      for (const d of (v.dimensions || []) as Array<Dimension>) {
        if (!d.meta) {
          d.meta = {}
        }

        if (!d.meta.steps) {
          d.meta.steps = []
        }
      }
    }
  }

  /**
   * Since gauge charts always define one type, this check can be simplified
   */
  mtrCheck ({ field, aggregate }: Metric) {
    if (!field) {
      throw new Error('notification.chart.invalidConfig.missingMetricsField')
    }
    if (field !== 'count' && !aggregate) {
      throw new Error('notification.chart.invalidConfig.missingMetricsAggregate')
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

  makeDataset (m: Metric, d: Dimension, data: Array<number|any>, alias: string) {
    const steps = (d.meta?.steps || [])

    return {
      value: data.reduce((acc, cur) => {
        return !isNaN(cur) ? acc + parseFloat(cur) : acc
      }, 0),
      data: steps.map(({ value }: any) => parseFloat(value)),
      backgroundColor: steps.map(({ color }: any) => color),
      tooltips: {
        enabled: true,
        labelCallback: m.fixTooltips ? this.makeLabel : this.makeTooltip,
      },
    }
  }

  makeOptions (data: any) {
    const rep = this.config.reports?.[0]
    const { metrics: [metric] = [] } = rep || {}

    const options: any = {
      needle: {
        radiusPercentage: 2,
        widthPercentage: 3.5,
        lengthPercentage: 70,
        color: metric.backgroundColor || defaultBGColor,
      },
      tooltips: {
        enabled: true,
        displayColors: false,
        callbacks: {
          title: () => '',
          label: ({ datasetIndex, index }: any, { datasets, labels }: any) => {
            const dataset = datasets[datasetIndex]
            return dataset.tooltips.labelCallback({ datasetIndex, index }, { datasets, labels })
          },
        },
        titleFontFamily: "'Poppins-Regular'",
        bodyFontFamily: "'Poppins-Regular'",
        footerFontFamily: "'Poppins-Regular'",
      },
      valueLabel: {
        display: true,
        color: 'rgba(0, 0, 0, 1)',
        backgroundColor: 'rgba(255,255,255,0.5)',
        borderRadius: 5,
        bottomMarginPercentage: 5,
        padding: {
          top: 10,
          bottom: 10,
        },
      },
      legend: {
        labels: {
          // This more specific font property overrides the global property
          fontFamily: "'Poppins-Regular'",
        }
      },
    }

    if (this.config.colorScheme) {
      options.plugins = {
        colorschemes: {
          scheme: this.config.colorScheme,
          // this is a bit of a hack to make the plugin work on each dataset value
          // we should improve this at a later point in time, but is ok for now.
          custom: (e: Array<string>) => {
            const cls = [...e]
            while (cls.length < data.datasets[0].data.length) {
              cls.push(...e)
            }
            data.datasets[0].backgroundColor = cls.slice(0, data.datasets[0].data.length)
            return e
          },
        },
      }
    }
    return options
  }

  private makeTooltip ({ datasetIndex, index }: any, { datasets, labels }: any): any {
    const dataset = datasets[datasetIndex]

    return makeDataLabel({
      prefix: labels[index],
      value: dataset.data[index],
    })
  }

  private makeLabel ({ datasetIndex, index }: any, { datasets, labels }: any): any {
    const dataset = datasets[datasetIndex]

    return makeDataLabel({
      prefix: labels[index],
      value: dataset.data[index],
    })
  }

  /**
   * @note Gauge chart requires the use of chartjs-gauge.
   * I was unable to make this work if the plugin was provided from this object,
   * so the plugin is registered on the webapp.
   * We should fix this at a later point in time...
   */
  plugins () {
    const mm: Array<Metric> = []

    for (const r of (this.config.reports || []) as Array<Report>) {
      mm.push(...(r.metrics || []) as Array<Metric>)
    }

    const rr: Array<any> = []
    if (mm.find(({ fixTooltips }) => fixTooltips)) {
      rr.push(makeTipper(ChartJS.Tooltip, {}))
    }
    return rr
  }

  baseChartType (datasets: Array<any>) {
    return 'gauge'
  }

  defMetrics (): Metric {
    return Object.assign({}, { type: ChartType.gauge })
  }
}

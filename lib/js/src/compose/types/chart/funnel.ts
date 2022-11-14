import { BaseChart, PartialChart } from './base'
import {
  Dimension,
  Metric,
  Report,
  ChartType,
  makeDataLabel,
  calculatePercentages,
} from './util'

import { defaultBGColor } from './common'
import { makeTipper } from './chartjs/plugins'
const ChartJS = require('chart.js')

/**
 * Funnel chart provides the definitions for the chartjs-plugin-funnel plugin.
 */
export default class FunnelChart extends BaseChart {
  constructor (def: PartialChart = {}) {
    super(def)

    // Assure required fields; this helps with backwards compatibility
    for (const v of (this.config.reports || []) as Array<Report>) {
      for (const d of (v.dimensions || []) as Array<Dimension>) {
        if (!d.meta) {
          d.meta = {}
        }

        if (!d.meta.fields) {
          d.meta.fields = []
        }
      }

      for (const m of (v.metrics || []) as Array<Metric>) {
        if (m.cumulative === undefined) {
          m.cumulative = true
        }
      }
    }
  }

  /**
   * Since funnel charts always define one type, this check can be simplified
   */
  mtrCheck ({ field, aggregate }: Metric) {
    if (!field) {
      throw new Error('notification.chart.invalidConfig.missingMetricsField')
    }
    if (field !== 'count' && !aggregate) {
      throw new Error('notification.chart.invalidConfig.missingMetricsAggregate')
    }
  }

  /**
   * Extend this method to include filtering for just specific values.
   * For example:
   * We wish to show only new and converted leads.
   */
  formatReporterParams (r: Report) {
    const base = super.formatReporterParams(r)
    const ff = base.filter

    let df = ''
    if (r.dimensions && r.dimensions[0]) {
      const rd = r.dimensions[0]
      if (r.dimensions[0].meta) {
        const fields = r.dimensions[0].meta.fields || []
        df = fields.map(({ value }: any) => `${rd.field || ''}='${value}'`)
          .join(' OR ')
      }
    }

    if (ff && df) {
      base.filter = `(${base.filter}) AND (${df})`
    } else if (!ff && df) {
      base.filter = df
    }

    return base
  }

  // Funnel chart creates a metric including all reports, so this step is deferred to there
  makeDataset (m: Metric, d: Dimension, data: Array<number|any>, alias: string) {
    const ds: any = { data }
    return ds
  }

  // No much configurations available for funnel charts
  makeOptions (data: any) {
    const options: any = {
      sort: 'desc',
      maintainAspectRatio: false,
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
        datalabels: {
          display: false,
        },
      }
    }

    options.tooltips = {
      enabled: true,
      displayColors: false,
      callbacks: {
        label: this.makeLabel,
      },
      titleFontFamily: "'Poppins-Regular'",
      bodyFontFamily: "'Poppins-Regular'",
      footerFontFamily: "'Poppins-Regular'",
    }
    return options
  }

  private makeLabel ({ datasetIndex, index }: any, { datasets }: any): any {
    const dataset = datasets[datasetIndex]

    // We use org data here to get actual percentages and not cumulative percentages
    const percentages = calculatePercentages(
      [...dataset.orgData],
      2,
      true,
      dataset.cumulative,
    )

    return makeDataLabel({
      value: dataset.data[index],
      relativeValue: percentages[index],
    })
  }

  /**
   * @note Funel chart requires the use of chartjs-plugin-funnel.
   * I was unable to make this work if the plugin was provided from this object,
   * so the plugin is registered on the webapp.
   * We should fix this at a later point in time...
   */
  plugins (mm: Array<Metric>) {
    return [makeTipper(ChartJS.Tooltip, {})]
  }

  baseChartType (datasets: Array<any>) {
    return 'funnel'
  }

  /**
   * Includes a few additional post processing steps:
   * * generate a set of labels based on all reports, all data sets,
   * * generates a set of data based on all reports, all data sets,
   */
  async fetchReports (a: any) {
    const rr = await super.fetchReports(a) as any
    const values = []

    // Above provided data sets might not have their labels/values ordered
    // correctly
    const valMap: any = {}
    // Map values to their labels
    for (let ri = 0; ri < rr.length; ri++) {
      const r = rr[ri]
      r.labels.forEach((l: string, i: number) => {
        valMap[l] = r.datasets[0].data[i]
      })

      // Construct labels & data based on provided reports
      const report = this.config.reports?.[ri]
      const d = report?.dimensions?.[0] as Dimension

      for (const { value } of d.meta?.fields || []) {
        values.push({
          // Use value for label and resolve it on FE (i18n)
          label: value,
          data: valMap[value] || 0,
        })
      }
    }

    // We are rendering the chart upside down
    // (by default it renders in ASC, but we want DESC)
    const labels: any[] = []
    const data: any[] = []

    values.sort((a, b) => a.data - b.data).forEach(v => {
      labels.push(v.label)
      data.push(v.data)
    })

    // Determine color to render for specific value
    const colorMap: { [_: string]: string } = {}
    this.config.reports?.forEach(r => {
      for (const { value, color } of r.dimensions?.[0].meta?.fields) {
        colorMap[value] = color
      }
    })

    // Get cumulative data but also keep original for tooltips
    const orgData = [...data]
    if (this.isCumulative()) {
      for (let i = 1; i < data.length; i++) {
        data[i] += data[i - 1]
      }
    }

    return {
      labels,
      datasets: [{
        data,
        orgData,
        backgroundColor: labels.map(l => colorMap[l] || defaultBGColor),
        cumulative: this.isCumulative(),
      }],
    }
  }

  isCumulative (): boolean {
    // Cumulative true by default
    // Find false value
    let cumulative = true
    const { reports = [] } = this.config

    reports.forEach(({ metrics = [] }) => {
      if (cumulative && !metrics[0].cumulative) {
        cumulative = false
      }
    })

    return cumulative
  }

  defMetrics (): Metric {
    return Object.assign({}, { type: ChartType.funnel })
  }

  defDimension (): Dimension {
    return Object.assign({}, { conditions: {}, meta: { fields: [] } })
  }
}

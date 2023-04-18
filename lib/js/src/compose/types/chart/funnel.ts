import { BaseChart, PartialChart } from './base'
import {
  Dimension,
  Metric,
  Report,
  ChartType,
} from './util'
import { getColorschemeColors } from '../../../shared'

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
    return {
      type: m.type,
      label: m.label || m.field,
      data,
      tooltip: {
        fixed: !!m.fixTooltips,
        relative: !!m.relativeValue,
      },
    }
  }

  makeOptions (data: any) {
    const { reports = [], colorScheme, noAnimation = false, toolbox } = this.config
    const { saveAsImage } = toolbox || {}

    const { labels, datasets = [], tooltip } = data
    const { legend: l } = reports[0] || {}
    const colors = getColorschemeColors(colorScheme)

    const tooltipFormatter = `{b}<br />{c} ${tooltip.relative ? ' ({d}%)' : ''}`
    const labelFormatter = `{c}${tooltip.relative ? ' ({d}%)' : ''}`

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
      tooltip: {
        show: true,
        trigger: 'item',
        formatter: tooltipFormatter,
        appendToBody: true,
      },
      legend: {
        show: !l?.isHidden,
        type: l?.isScrollable ? 'scroll' : 'plain',
        top: (l?.position?.isDefault ? undefined : l?.position?.top) || undefined,
        right: (l?.position?.isDefault ? undefined : l?.position?.right) || undefined,
        bottom: (l?.position?.isDefault ? undefined : l?.position?.bottom) || undefined,
        left: (l?.position?.isDefault ? l?.align || 'center' : l?.position?.left) || 'auto',
        orient: l?.orientation || 'horizontal'
      },
      series: datasets.map(({ data }: any) => {
        return {
          type: 'funnel',
          sort: 'descending',
          top: 45,
          bottom: 10,
          left: '5%',
          width: '90%',
          label: {
            show: tooltip.fixed,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
            formatter: labelFormatter,
          },
          emphasis: {
            label: {
              show: false,
              fontSize: 14,
            },
          },
          data: labels.map((name: string, i: number) => {
            return { name, value: data[i], itemStyle: { color: colors[i] } }
          }),
        }
      }),
    }
  }

  baseChartType (): string {
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

    let tooltip = {}

    // Above provided data sets might not have their labels/values ordered
    // correctly
    const valMap: any = {}
    // Map values to their labels
    for (let ri = 0; ri < rr.length; ri++) {
      const r = rr[ri]

      r.labels.forEach((l: string, i: number) => {
        valMap[l] = r.datasets[0].data[i]
      })

      tooltip = { ...tooltip, ...r.datasets[0].tooltip }

      // Construct labels & data based on provided reports
      const report = this.config.reports?.[ri]
      const d = report?.dimensions?.[0] as Dimension

      let { fields = [] } = d.meta || {}
      fields = fields.length ? fields : r.labels

      for (let label of fields) {
        const value = typeof label === 'object' ? label.value : label
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
    if (this.isCumulative()) {
      for (let i = 1; i < data.length; i++) {
        data[i] += data[i - 1]
      }
    }

    return {
      labels,
      datasets: [{
        data,
      }],
      tooltip,
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
    return Object.assign({}, {
      type: ChartType.funnel,
      fixTooltips: false,
      relativeValue: true,
    })
  }

  defDimension (): Dimension {
    return Object.assign({}, {
      conditions: {},
      meta: { fields: [] }
    })
  }
}

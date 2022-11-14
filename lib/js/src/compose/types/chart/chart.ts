import { BaseChart } from './base'
import {
  Dimension,
  Metric,
  dimensionFunctions,
  isRadialChart,
  makeDataLabel,
  KV,
  Report,
  TemporalDataPoint,
  calculatePercentages,
} from './util'

import { makeTipper } from './chartjs/plugins'
const ChartJS = require('chart.js')

// The default dataset post processing function to use.
// This one simply returns the current value.
const defaultFx = 'n'

/**
 * Chart represents a generic chart, such as a bar chart, line chart, ...
 */
export default class Chart extends BaseChart {
  // Generic charts (at the moment) support only 1 report per chart
  async fetchReports (a: any) {
    return super.fetchReports(a).then((rr: any) => {
      return rr[0]
    })
  }

  /**
   * The method performs post processing for each value in the given dataset.
   * It works with a simple equation written in javascript (example: n + m).
   * Available variables to use:
   * * n - current value
   * * m - previous value (undefined in case of the first element)
   * * r - entire data array.
   *
   * @param data Array of values in the given data set
   * @param m Metric for the given dataset
   */
  private datasetPostProc (data: Array<number|TemporalDataPoint>, m: Metric): Array<number|TemporalDataPoint> {
    // Define a valid function to evaluate
    let fxRaw = (m.fx || defaultFx).trim()
    if (!fxRaw.startsWith('return')) {
      fxRaw = 'return ' + fxRaw
    }
    const fx = new Function('n', 'm', 'r', fxRaw)

    // Define a new array, so we don't alter the original one.
    const r = [...data]

    // Run postprocessing for all data in the given data set
    // There is a slight difference between temporal data points and categorical data points.
    if (data[0] instanceof Object) {
      // Temporal
      for (let i = 0; i < data.length; i++) {
        const a = data[i] as TemporalDataPoint
        const b = data[i - 1] as TemporalDataPoint|undefined

        const n = a.y
        let m: number|undefined
        if (i > 0) {
          m = b?.y
        }

        a.y = fx(n, m, r)
      }
    } else {
      // Categorical
      for (let i = 0; i < data.length; i++) {
        const n = data[i] as number
        let m: number|undefined
        if (i > 0) {
          m = data[i - 1] as number
        }
        data[i] = fx(n, m, r)
      }
    }

    return data
  }

  makeDataset (m: Metric, d: Dimension, data: Array<number|any>, alias: string) {
    data = this.datasetPostProc(data, m)
    const ds: any = { data }

    // colors
    if (typeof m.backgroundColor === 'string' && !isRadialChart(m)) {
      ds.backgroundColor = 'rgba(' + parseInt(m.backgroundColor.slice(-6, -4), 16) + ',' + parseInt(m.backgroundColor.slice(-4, -2), 16) + ',' + parseInt(m.backgroundColor.slice(-2), 16) + ',0.7)'
      ds.hoverBackgroundColor = m.backgroundColor
    }

    return Object.assign(ds, {
      label: m.label || m.field,
      type: m.type,
      fill: !!m.fill,
      lineTension: m.lineTension || 0,
      tooltips: {
        enabled: true,
        relativeValue: !!m.relativeValue,
        relativePrecision: m.relativePrecision,
        labelCallback: m.fixTooltips ? this.makeLabel : this.makeTooltip,
      },
    })
  }

  makeOptions () {
    const options: any = {
      // Allow chart to consume entire container
      responsive: true,
      maintainAspectRatio: false,
      animation: {
        duration: 500,
      },
      legend: {
        position: 'top',
        labels: {
          // This more specific font property overrides the global property
          fontFamily: "'Poppins-Regular'",
        },
      },
    }

    if (this.config.colorScheme) {
      options.plugins = {
        colorschemes: {
          scheme: this.config.colorScheme,
        },
      }
    }

    (this.config.reports || []).forEach(r => {
      if (!options.scales) {
        options.scales = { xAxes: [], yAxes: [] }
      }

      // can't disable tooltips on dataset level, so this is required
      options.tooltips = {
        filter: ({ datasetIndex }: any, { datasets }: any) => {
          // enabled can be undefined, so it must be checked against false
          return ((datasets[datasetIndex] || {}).tooltips || {}).enabled !== false
        },

        callbacks: {
          label: ({ datasetIndex, index }: any, { datasets, labels }: any) => {
            const dataset = datasets[datasetIndex]
            return dataset.tooltips.labelCallback({ datasetIndex, index }, { datasets, labels })
          },
        },
        titleFontFamily: "'Poppins-Regular'",
        bodyFontFamily: "'Poppins-Regular'",
        footerFontFamily: "'Poppins-Regular'",
        displayColors: false,
      }

      if (r.metrics?.find((m: Metric) => !isRadialChart(m as KV))) {
        options.scales.xAxes = r.dimensions?.map((d: Dimension, i: number) => {
          const ticks = {
            autoSkip: !!d.autoSkip,
          }
          const timeDimensionUnit = (dimensionFunctions.lookup(d) || {}).time

          if (timeDimensionUnit) {
            return {
              type: 'time',
              time: timeDimensionUnit,
              ticks,
            }
          } else {
            return {
              ticks,
            }
          }
        })
      }

      for (const m of r.metrics || []) {
        if (m.legendPosition) {
          options.legend.position = m.legendPosition
          break
        }
      }

      options.scales.yAxes = this.makeYAxis(r)
    })
    return options
  }

  private makeTooltip ({ datasetIndex, index }: any, { datasets, labels }: any): any {
    const dataset = datasets[datasetIndex]

    const percentages = calculatePercentages(
      [...dataset.data],
      dataset.tooltips.relativePrecision,
      dataset.tooltips.relativeValue,
    )

    return makeDataLabel({
      prefix: labels[index],
      value: dataset.data[index],
      relativeValue: dataset.tooltips.relativeValue ? percentages[index] : undefined,
    })
  }

  private makeLabel ({ datasetIndex, index }: any, { datasets }: any): any {
    const dataset = datasets[datasetIndex]

    const percentages = calculatePercentages(
      [...dataset.data],
      dataset.tooltips.relativePrecision,
      dataset.tooltips.relativeValue,
    )

    return makeDataLabel({
      value: dataset.data[index],
      relativeValue: dataset.tooltips.relativeValue ? percentages[index] : undefined,
    })
  }

  private makeYAxis (r: Report) {
    if (r.yAxis) {
      return [{
        display: !r.metrics?.find((m: KV) => isRadialChart(m)),
        type: r.yAxis.axisType || 'linear',
        position: r.yAxis.axisPosition || 'left',
        scaleLabel: {
          display: true,
          labelString: r.yAxis.label || undefined,
        },
        ticks: {
          beginAtZero: !!r.yAxis.beginAtZero,
          min: r.yAxis.min ? parseFloat(r.yAxis.min) : undefined,
          max: r.yAxis.max ? parseFloat(r.yAxis.max) : undefined,
        },
      }]
    } else {
      const m: Metric = r.metrics?.[0] || {}
      return [{
        display: !isRadialChart(m as KV),
        type: m.axisType || 'linear',
        position: m.axisPosition || 'left',
        scaleLabel: {
          display: true,
          labelString: m.label || m.field,
        },
        ticks: {
          beginAtZero: !!m.beginAtZero,
        },
      }]
    }
  }

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
    return datasets[0].type
  }
}

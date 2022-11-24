import { BaseChart } from './base'
import {
  Dimension,
  Metric,
  dimensionFunctions,
  TemporalDataPoint,
} from './util'
import { getColorschemeColors } from '../../../shared'

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

    return {
      type: m.type,
      label: m.label || m.field,
      data,
      fill: !!m.fill,
      tooltip: {
        fixed: m.fixTooltips,
        relative: !!m.relativeValue,
      },
    }
  }

  makeOptions (data: any): any {
    const { reports = [], colorScheme } = this.config

    const options: any = {
      series: [],
      xAxis: [],
      yAxis: [],
      tooltip: {
        show: true,
        appendToBody: true,
        position: 'inside',
      },
    }

    const { labels, datasets = [] } = data
    const {
      dimensions: [dimension] = [],
      yAxis, metrics: [metric] = [],
      offset,
      tooltip: t,
      legend: l,
    } = reports[0] || {}

    const hasAxis = datasets.some(({ type }: any) => ['bar', 'line'].includes(type))
    const timeDimension = (dimensionFunctions.lookup(dimension) || {}).time

    if (hasAxis) {
      if (yAxis) {
        const {
          label: yLabel,
          axisType: yType = 'linear',
          axisPosition: position = 'left',
          labelPosition = 'end',
          beginAtZero,
          min,
          max,
        } = yAxis


        const tempYAxis = {
          name: yLabel,
          type: yType === 'linear' ? 'value' : 'log',
          position,
          nameGap: labelPosition === 'center' ? 25 : 7,
          nameLocation: labelPosition,
          min: beginAtZero ? 0 : min || undefined,
          max: max || undefined,
          axisLabel: {
            interval: 0,
            overflow: 'truncate',
            hideOverlap: true,
            rotate: yAxis.rotateLabel,
          },
          axisLine: {
            show: true,
            onZero: false,
          },
          nameTextStyle: {
            align: labelPosition === 'center' ? 'center' : position,
            padding: labelPosition !== 'center' ? (position === 'left' ? [0, 0, 2, -20] : [0, -20, 2, 0]) : undefined,
          },
        }

        // If we provide undefined, log scale breaks
        if (tempYAxis.type === 'log') {
          delete tempYAxis.min
          delete tempYAxis.max
        }

        options.yAxis = [tempYAxis]
      }
    }

    options.series = datasets.map(({ type, label, data, fill, tooltip }: any, index: number) => {
      const { fixed, relative } = tooltip

      const tooltipFormatter = t?.formatting ? t.formatting : `{a}<br />{b} : {c}${relative ? ' ({d}%)' : ''}`
      const labelFormatter = `{c}${relative ? ' ({d}%)' : ''}`

      if (['pie', 'doughnut'].includes(type)) {
        const startRadius = type === 'doughnut' ? 40 : 0

        options.tooltip.trigger = 'item'

        let lbl =  {}

        if (t?.labelsNextToPartition) {
          lbl = {
            show: true,
            overflow: 'truncate',
          }
        } else {
          lbl = {
            show: fixed,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
          }
        }


        return {
          z: index,
          name: label,
          type: 'pie',
          radius: [`${startRadius}%`, '80%'],
          center: ['50%', '55%'],
          tooltip: {
            trigger: 'item',
            formatter: tooltipFormatter,
          },
          label: {
            ...lbl,
            formatter: labelFormatter,
          },
          itemStyle: {
            borderRadius: 5,
            borderColor: '#fff',
            borderWidth: 1,
          },
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)',
            },
          },
          data: labels.map((name: string, i: number) => {
            return { name, value: data[i] }
          }),
          top: offset?.isDefault ? undefined : offset?.top,
          right: offset?.isDefault ? undefined : offset?.right,
          bottom: offset?.isDefault ? undefined : offset?.bottom,
          left: offset?.isDefault ? undefined : offset?.left,
        }
      } else if (['bar', 'line'].includes(type)) {
        options.tooltip.trigger = 'axis'

        if (!options.xAxis.length) {
          options.xAxis.push({
            nameLocation: 'center',
            type: 'category',
            data: labels,
            axisLabel: {
              interval: 0,
              overflow: 'truncate',
              hideOverlap: true,
              rotate: dimension.rotateLabel,
            },
          })
        }

        options.grid = {
          top: offset?.isDefault ? 45 : offset?.top,
          right: offset?.isDefault ? 25 : offset?.right,
          bottom: offset?.isDefault ? 25 : offset?.bottom,
          left: offset?.isDefault ? 25 : offset?.left,
          containLabel: true,
        }

        return {
          z: index,
          name: label,
          type: type,
          smooth: true,
          areaStyle: {
            opacity: fill ? 0.7 : 0,
          },
          label: {
            show: fixed,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
            formatter: labelFormatter,
          },
          data,
        }
      }
    })

    return {
      color: getColorschemeColors(colorScheme),
      textStyle: {
        fontFamily: 'Poppins-Regular',
      },
      legend: {
        show: !l?.isHidden,
        type: l?.isScrollable ? 'scroll' : 'plain',
        top: (l?.position?.isDefault ? l?.position?.top : undefined) || undefined,
        right: (l?.position?.isDefault ? l?.position?.right : undefined) || undefined,
        bottom: (l?.position?.isDefault ? l?.position?.bottom : undefined) || undefined,
        left: (l?.position?.isDefault ? l?.position?.left : l?.align || 'center') || 'auto',
        orient: l?.orientation || 'horizontal'
      },
      ...options,
    }
  }

  baseChartType (datasets: Array<any>): string {
    return datasets[0].type
  }
}

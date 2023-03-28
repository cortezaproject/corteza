import { BaseChart } from './base'
import {
  Dimension,
  Metric,
  dimensionFunctions,
  TemporalDataPoint,
} from './util'
import { getColorschemeColors } from '../../../shared'

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

  makeDataset (m: Metric, d: Dimension, data: Array<number|TemporalDataPoint>, alias: string) {
    data = this.datasetPostProc(data, m)

    return {
      type: m.type,
      label: m.label || m.field,
      data,
      fill: !!m.fill,
      tooltip: {
        fixed: m.fixTooltips,
        relative: m.relativeValue && !['bar', 'line'].includes(m.type as string),
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
          nameGap: labelPosition === 'center' ? 30 : 7,
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
            padding: labelPosition !== 'center' ? (position === 'left' ? [0, 0, 2, -3] : [0, -3, 2, 0]) : undefined,
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

      // We should render the first metric in the dataset as the last
      const z = (datasets.length - 1) - index

      if (['pie', 'doughnut'].includes(type)) {
        const startRadius = type === 'doughnut' ? 40 : 0

        options.tooltip.trigger = 'item'

        let lbl:any =  {
          rotate: dimension.rotateLabel ? +dimension.rotateLabel: 0
        }

        if (t?.labelsNextToPartition) {
          lbl = {
            ...lbl,
            show: true,
            overflow: 'truncate',
          }
        } else {
          lbl = {
            ...lbl,
            show: fixed,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
          }
        }

        return {
          z,
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
          top: offset?.isDefault ? 50 : offset?.top,
          right: offset?.isDefault ? 30 : offset?.right,
          bottom: offset?.isDefault ? 20 : offset?.bottom,
          left: offset?.isDefault ? 30 : offset?.left,
          containLabel: true,
        }

        return {
          z,
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
        top: (l?.position?.isDefault ? undefined : l?.position?.top) || undefined,
        right: (l?.position?.isDefault ? undefined : l?.position?.right) || undefined,
        bottom: (l?.position?.isDefault ? undefined : l?.position?.bottom) || undefined,
        left: (l?.position?.isDefault ? l?.align || 'center' : l?.position?.left) || 'auto',
        orient: l?.orientation || 'horizontal'
      },
      ...options,
    }
  }

  baseChartType (datasets: Array<any>): string {
    return datasets[0].type
  }
}

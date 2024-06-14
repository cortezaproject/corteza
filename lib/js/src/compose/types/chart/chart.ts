import { BaseChart } from './base'
import {
  Dimension,
  Metric,
  TemporalDataPoint,
  formatChartValue,
  formatChartTooltip,
  TooltipParams,
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
      fill: m.fill,
      smooth: m.smooth,
      step: m.step ? 'middle' : undefined,
      roseType: m.rose ? 'radius' : undefined,
      symbol: m.symbol,
      stack: m.stack,
      tooltip: {
        fixed: m.fixTooltips,
        relative: m.relativeValue && !['bar', 'line'].includes(m.type as string),
      },
      formatting: m.formatting,
    }
  }

  makeOptions (data: any): any {
    const { reports = [], colorScheme, noAnimation = false, toolbox } = this.config
    const { saveAsImage, timeline = '' } = toolbox || {}

    const options: any = {
      animation: !noAnimation,
      series: [],
      xAxis: [],
      yAxis: [],
      tooltip: {
        show: true,
        appendToBody: true,
        position: 'inside',
      },
    }

    const { labels, datasets = [], themeVariables = {} } = data
    const {
      dimensions: [dimension] = [],
      yAxis,
      offset,
      tooltip: t,
      legend: l,
    } = reports[0] || {}

    const hasAxis = datasets.some(({ type }: any) => ['bar', 'line', 'scatter'].includes(type))
    let horizontal = false

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

        horizontal = !!yAxis.horizontal

        const xAxis = {
          nameLocation: 'center',
          type: dimension.timeLabels ? 'time' : 'category',
          axisLabel: {
            interval: 0,
            overflow: 'break',
            hideOverlap: true,
            rotate: dimension.rotateLabel,
          },
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
          },
        }

        const tempYAxis = {
          name: yLabel,
          type: yType === 'linear' ? 'value' : 'log',
          position,
          nameLocation: labelPosition,
          min: beginAtZero ? 0 : Number(min) || undefined,
          max: Number(max) || undefined,
          axisLabel: {
            interval: 0,
            overflow: 'break',
            hideOverlap: true,
            rotate: yAxis.rotateLabel,
            formatter: (value: string | number): string => formatChartValue(value, yAxis.formatting),
          },
          axisLine: {
            show: false,
            onZero: false,
          },
          splitLine: {
            lineStyle: {
              color: [themeVariables['extra-light']],
            },
          },
          nameTextStyle: {
            align: labelPosition === 'center' ? 'center' : position,
          },
        }

        // If we provide undefined, log scale breaks
        if (tempYAxis.type === 'log') {
          delete tempYAxis.min
          delete tempYAxis.max
        }

        if (horizontal) {
          options.xAxis = [tempYAxis]
          options.yAxis = [xAxis]
        } else {
          options.xAxis = [xAxis]
          options.yAxis = [tempYAxis]
        }
      }
    }

    options.series = datasets.map(({ formatting, type, label, data, stack, tooltip, fill, smooth, step, roseType, symbol }: any, index: number) => {
      const { fixed, relative } = tooltip

      // We should render the first metric in the dataset as the last
      const z = (datasets.length - 1) - index

      if (['pie', 'doughnut'].includes(type)) {
        const startRadius = type === 'doughnut' ? 40 : 0
        const endRadius = 80
        const radiusLength = (endRadius - startRadius) / (datasets.length || 1)

        const sr = startRadius + (index * radiusLength)
        const er = startRadius + ((index + 1) * radiusLength)

        options.tooltip.trigger = 'item'

        let lbl :any = {
          rotate: dimension.rotateLabel ? +dimension.rotateLabel : 0,
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
          stack,
          name: label,
          type: 'pie',
          roseType,
          radius: [`${sr}%`, `${er}%`],
          center: ['50%', '55%'],
          tooltip: {
            trigger: 'item',
            appendToBody: true,
            formatter: (params: TooltipParams): string => {
              const v = formatChartValue(params.value || '', formatting)

              if (t?.formatting) {
                return formatChartTooltip(t?.formatting, params)
              }

              return `${params.seriesName}<br>${params.marker}${params.name}<span style="float: right; margin-left: 20px">${v}${relative ? ' (' + params.percent + '%)' : ''}</span>`
            },
          },
          label: {
            ...lbl,
            formatter: (params: TooltipParams): string => formatChartValue(params.value || '', formatting),
          },
          itemStyle: {
            borderRadius: 5,
            borderColor: themeVariables.white,
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
      } else if (['bar', 'line', 'scatter'].includes(type)) {
        options.tooltip.trigger = 'axis'

        const defaultOffset = {
          top: 65,
          right: timeline.includes('x') ? 40 : 30,
          bottom: timeline.includes('x') ? 60 : 20,
          left: 30,
        }

        options.grid = {
          top: offset?.isDefault ? defaultOffset.top : offset?.top,
          right: offset?.isDefault ? defaultOffset.right : offset?.right,
          bottom: offset?.isDefault ? defaultOffset.bottom : offset?.bottom,
          left: offset?.isDefault ? defaultOffset.left : offset?.left,
          containLabel: true,
        }

        if (horizontal) {
          data = labels.map((name: string, i: number) => {
            return [data[i], name]
          })
        } else {
          data = labels.map((name: string, i: number) => {
            return [name, data[i]]
          })
        }

        return {
          z,
          stack,
          name: label,
          type: type,
          smooth,
          step,
          areaStyle: {
            opacity: fill ? 0.7 : 0,
          },
          symbol,
          symbolSize: type === 'scatter' ? 16 : 10,
          tooltip: {
            appendToBody: true,
            // pass trigger type to determine if valueFormatter or formatter will be used
            trigger: t?.formatting ? 'item' : 'axis',
            // we can either
            // add formatting to the value and apply tooltip if trigger: 'item'
            // display the same tooltip format name <br/> seriesName value if trigger: 'axis'

            // works when trigger is set to axis
            valueFormatter: (value: string | number): string => formatChartValue(value, formatting),
            // works when trigger is set to item
            formatter: (params: { seriesName?: string; name?: string;value: Array<any>, percent: string | number, marker: string;}): string => {
              const { value = [], percent = '' } = params

              const formattedValue = formatChartValue(value[1], formatting)

              if (t?.formatting) {
                return formatChartTooltip(t?.formatting, { ...params, value: value[1], percent })
              }

              return `${params.seriesName}<br>${params.marker}${params.name}<span style="float: right; margin-left: 20px">${formattedValue}${relative ? ' (' + params.percent + '%)' : ''}</span>`
            },
          },
          label: {
            show: fixed,
            position: 'inside',
            align: 'center',
            verticalAlign: 'middle',
            tooltip: {
              trigger: 'axis',
            },
            formatter: (params: { seriesName: string, name: string, value: Array<any>, percent: string | number }): string => {
              const { value = [], percent = '' } = params

              return `${formatChartValue(value[1], formatting)}${relative ? ` (${percent}%)` : ''}`
            },
          },
          data,
        }
      }
    })

    const dataZoom = timeline ? [
      {
        show: timeline.includes('x'),
        type: 'slider',
        height: 30,
      },
      {
        show: timeline.includes('y'),
        type: 'slider',
        width: 15,
        yAxisIndex: 0,
      },
    ] : undefined

    return {
      color: getColorschemeColors(colorScheme, data.customColorSchemes),
      textStyle: {
        fontFamily: themeVariables['font-regular'],
        overflow: 'break',
        color: themeVariables.black,
      },
      toolbox: {
        feature: {
          saveAsImage: saveAsImage ? {
            name: this.name,
          } : undefined,
        },
        top: 15,
        right: 5,
      },
      dataZoom,
      legend: {
        show: !l?.isHidden,
        type: l?.isScrollable ? 'scroll' : 'plain',
        top: (l?.position?.isDefault ? undefined : l?.position?.top) || undefined,
        right: (l?.position?.isDefault ? undefined : l?.position?.right) || undefined,
        bottom: (l?.position?.isDefault ? undefined : l?.position?.bottom) || undefined,
        left: (l?.position?.isDefault ? l?.align || 'center' : l?.position?.left) || 'auto',
        orient: l?.orientation || 'horizontal',
        textStyle: {
          color: themeVariables.black,
        },
        pageTextStyle: {
          color: themeVariables.black,
        },
        pageIconColor: themeVariables.black,
        pageIconInactiveColor: themeVariables.light,
      },
      ...options,
    }
  }

  defMetric (): Metric {
    return Object.assign(super.defMetric(), {
      smooth: true,
      fill: false,
      rose: false,
      symbol: 'circle',
    })
  }

  baseChartType (datasets: Array<any>): string {
    return datasets[0].type
  }
}

import { BaseChart } from './base'
import {
  Dimension,
  Metric,
  ChartType,
  formatChartValue,
  TooltipParams,
} from './util'
import { getColorschemeColors } from '../../../shared'

export default class RadarChart extends BaseChart {
  mtrCheck ({ field, aggregate }: Metric) {
    if (!field) {
      throw new Error('notification.chart.invalidConfig.missingMetricsField')
    }
    if (field !== 'count' && !aggregate) {
      throw new Error('notification.chart.invalidConfig.missingMetricsAggregate')
    }
  }

  makeDataset (m: Metric, d: Dimension, data: Array<number|any>) {
    return {
      type: m.type,
      label: m.label || m.field,
      data,
      formatting: m.formatting,
    }
  }

  makeOptions (data: any) {
    const { reports = [], colorScheme, noAnimation = false, toolbox } = this.config
    const { saveAsImage } = toolbox || {}
    const { labels, datasets = [], dimension = {}, themeVariables = {} } = data

    const {
      legend: l,
    } = reports[0] || {}

    const { formatting } = datasets[0] || {}

    let min: number = 0
    let max: number = Math.max()
    const seriesData: any[] = []

    datasets.forEach(({ data: value, label: name }: any) => {
      value.forEach((v: number) => {
        if (v < min) min = v
        if (v > max) max = v
      })

      seriesData.push({ value, name })
    })

    return {
      color: getColorschemeColors(colorScheme, data.customColorSchemes),
      animation: !noAnimation,
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
      },
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
      tooltip: {
        show: true,
        position: 'top',
        appendToBody: true,
        valueFormatter: (value: string | number): string => formatChartValue(value, formatting),
      },
      radar: {
        shape: dimension.shape,
        indicator: labels.map((name: string) => {
          return { name, min, max }
        }),
        center: ['50%', '55%'],
      },
      series: {
        type: 'radar',
        label: {
          show: dimension.fixTooltips,
          formatter: (params: { value: string | number }): string => {
            const { value = '' } = params
            return formatChartValue(value, formatting)
          },
        },
        data: seriesData,
      },
    }
  }

  baseChartType (): string {
    return 'radar'
  }

  async fetchReports (a: any) {
    return super.fetchReports(a).then((rr: any) => {
      return rr[0]
    })
  }

  defMetric (): Metric {
    return Object.assign(super.defMetric(), {
      type: ChartType.radar,
    })
  }

  defDimension (): Dimension {
    return Object.assign({}, {
      shape: 'polygon',
      fixTooltips: false,
      conditions: {},
      meta: {},
    })
  }
}

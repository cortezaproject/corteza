import { BaseChart } from './base'
import {
  Dimension,
  Metric,
  ChartType,
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
    }
  }

  makeOptions (data: any) {
    const { reports = [], colorScheme, noAnimation = false, toolbox } = this.config
    const { saveAsImage } = toolbox || {}
    const { labels, datasets = [], dimension = {} } = data
    const { legend: l } = reports[0] || {}

    const labelFormatter = '{c}'

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
        fontFamily: 'Poppins-Regular',
      },
      toolbox: {
        feature: {
          saveAsImage: saveAsImage ? {
            name: this.name
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
        orient: l?.orientation || 'horizontal'
      },
      tooltip: {
        show: true,
        position: 'top',
        appendToBody: true,
      },
      radar: {
        shape: dimension.shape,
        indicator: labels.map((name: string) => {
          return { name, min, max }
        }),
        center: ['50%', '55%']
      },
      series: {
        type: 'radar',
        label: {
          show: dimension.fixTooltips,
          formatter: labelFormatter,
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

  defMetrics (): Metric {
    return Object.assign({}, {
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

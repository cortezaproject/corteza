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

    return {
      color: getColorschemeColors(colorScheme),
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
          return { name }
        }),
        center: ['50%', '55%']
      },
      series: {
        type: 'radar',
        label: {
          show: dimension.fixTooltips,
          formatter: labelFormatter,
        },
        data: datasets.map(({ data: value, label: name }: any, index: number) => {
          return { value, name }
        }),
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

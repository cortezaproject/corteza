import { BaseChart } from './base'
import { Dimension, Metric, ChartType } from './util'
export default class RadarChart extends BaseChart {
  mtrCheck({ field, aggregate }: Metric): void
  makeDataset(
    m: Metric,
    d: Dimension,
    data: Array<number | any>,
  ): {
    type: ChartType | undefined
    label: any
    data: any[]
    formatting: import('./util').FormatData
  }
  makeOptions(data: any): {
    color: string[]
    animation: boolean
    textStyle: {
      fontFamily: any
      overflow: string
      color: any
    }
    toolbox: {
      feature: {
        saveAsImage:
          | {
              name: string
            }
          | undefined
      }
    }
    legend: {
      show: boolean
      type: string
      top: string | undefined
      right: string | undefined
      bottom: string | undefined
      left: string
      orient: string
      textStyle: {
        color: any
      }
      pageTextStyle: {
        color: any
      }
      pageIconColor: any
      pageIconInactiveColor: any
    }
    tooltip: {
      show: boolean
      position: string
      appendToBody: boolean
      valueFormatter: (value: string | number) => string
    }
    radar: {
      shape: any
      indicator: any
      center: string[]
    }
    series: {
      type: string
      label: {
        show: any
        formatter: (params: { value: string | number }) => string
      }
      data: any[]
    }
  }
  baseChartType(): string
  fetchReports(a: any): Promise<any>
  defMetric(): Metric
  defDimension(): Dimension
}

import { BaseChart, PartialChart } from './base'
import { Metric, Dimension, TemporalDataPoint } from './util'
export default class GaugeChart extends BaseChart {
  constructor(def?: PartialChart)
  fetchReports(a: any): Promise<any>
  processLabels(ll: Array<string>, d: Dimension): any
  makeDataset(
    m: Metric,
    d: Dimension,
    data: Array<number | TemporalDataPoint>,
    alias: string,
  ): {
    steps: any
    name: any
    max: number
    value: any
    startAngle: any
    endAngle: any
    tooltip: {
      fixed: boolean | undefined
    }
    formatting: import('./util').FormatData
  }
  makeOptions(data: any): {
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
      top: number
      right: number
    }
    grid: {
      bottom: number
    }
    series: {
      type: string
      startAngle: any
      endAngle: any
      min: number
      max: any
      splitNumber: number
      radius: string
      center: string[]
      pointer: {
        width: number
        length: string
        itemStyle: {
          color: any
        }
      }
      splitLine: {
        distance: number
        length: number
        lineStyle: {
          color: any
        }
      }
      axisLine: {
        lineStyle: {
          width: number
          color: any
        }
      }
      axisTick: {
        show: boolean
        distance: number
      }
      axisLabel: {
        show: boolean
        distance: number
      }
      title: {
        fontSize: number
        show: any
        offsetCenter: (string | number)[]
        color: any
      }
      detail: {
        fontSize: number
        offsetCenter: (string | number)[]
        valueAnimation: boolean
        color: any
        formatter: (value: string | number) => string
      }
      data: {
        name: any
        value: any
      }[]
    }[]
  }
  baseChartType(): string
  defMetric(): Metric
  /**
   * Checks validity of dimensions.
   * If invalid it throws an error
   */
  dimCheck({ meta }: Dimension): void | Error
  /**
   * Since gauge charts always define one type, this check can be simplified
   */
  mtrCheck({ field, aggregate }: Metric): void | Error
}

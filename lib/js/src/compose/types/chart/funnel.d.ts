import { BaseChart, PartialChart } from './base'
import { Dimension, Metric, Report, ChartType, TooltipParams } from './util'
export default class FunnelChart extends BaseChart {
  constructor(def?: PartialChart)
  /**
   * Since funnel charts always define one type, this check can be simplified
   */
  mtrCheck({ field, aggregate }: Metric): void
  /**
   * Extend this method to include filtering for just specific values.
   * For example:
   * We wish to show only new and converted leads.
   */
  formatReporterParams(r: Report): {
    moduleID: string | null | undefined
    filter: string | null | undefined
    metrics: string | undefined
    dimensions: any
  }
  makeDataset(
    m: Metric,
    d: Dimension,
    data: Array<number | any>,
    alias: string,
  ): {
    type: ChartType | undefined
    label: any
    data: any[]
    tooltip: {
      fixed: boolean
      relative: boolean
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
    tooltip: {
      trigger: string
      formatter: (params: TooltipParams) => string
      appendToBody: boolean
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
    series: any
  }
  baseChartType(): string
  /**
   * Includes a few additional post processing steps:
   * * generate a set of labels based on all reports, all data sets,
   * * generates a set of data based on all reports, all data sets,
   */
  fetchReports(a: any): Promise<{
    labels: any[]
    datasets: {
      label: string
      data: any[]
      formatting: Record<string, unknown>
    }[]
    tooltip: Record<string, unknown>
  }>
  isCumulative(): boolean
  defMetric(): Metric
  defDimension(): Dimension
}

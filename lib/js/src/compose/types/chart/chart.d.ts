import { BaseChart } from './base'
import { Dimension, Metric, TemporalDataPoint } from './util'
/**
 * Chart represents a generic chart, such as a bar chart, line chart, ...
 */
export default class Chart extends BaseChart {
  fetchReports(a: any): Promise<any>
  makeDataset(
    m: Metric,
    d: Dimension,
    data: Array<number | TemporalDataPoint>,
    alias: string,
  ): {
    type: import('./util').ChartType | undefined
    label: any
    data: (number | TemporalDataPoint)[]
    fill: any
    smooth: any
    step: string | undefined
    roseType: string | undefined
    symbol: string | undefined
    stack: any
    tooltip: {
      fixed: boolean | undefined
      relative: boolean | undefined
    }
    formatting: import('./util').FormatData
  }
  makeOptions(data: any): any
  defMetric(): Metric
  baseChartType(datasets: Array<any>): string
}

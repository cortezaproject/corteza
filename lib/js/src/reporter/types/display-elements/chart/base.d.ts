import { DisplayElement, DisplayElementInput } from '../base'
import { DefinitionOptions, FrameDefinition } from '../../frame'
export type PartialChartOptions = Partial<ChartOptions>
interface XAxisOptions {
  type: string
  label?: string
  unit?: string
  skipMissing: boolean
  defaultValue?: any
  labelRotation: number
}
interface YAxisOptions {
  label?: string
  labelPosition?: string
  labelRotation: number
  type?: string
  position?: string
  beginAtZero?: boolean
  stepSize?: string
  min?: string
  max?: string
}
interface Legend {
  hide: boolean
  orientation: string
  align: string
  scrollable: boolean
  position: Offset
}
interface Tooltips {
  showAlways: boolean
}
interface Offset {
  default: boolean
  top?: string
  right?: string
  bottom?: string
  left?: string
}
export declare class ChartOptions {
  title: string
  type: string
  colorScheme: string
  noAnimation: boolean
  source: string
  datasources: Array<FrameDefinition>
  xAxis: XAxisOptions
  yAxis: YAxisOptions
  legend: Legend
  tooltips: Tooltips
  offset: Offset
  constructor(o?: PartialChartOptions)
}
export declare const ChartOptionsRegistry: Map<string, typeof ChartOptions>
export declare function ChartOptionsMaker<T extends ChartOptions>(options: Partial<ChartOptions>): T
export declare class DisplayElementChart extends DisplayElement {
  readonly kind = 'Chart'
  options: ChartOptions
  constructor(i?: DisplayElementInput)
  applyOptions(o?: PartialChartOptions): void
  reportDefinitions(definition?: DefinitionOptions): {
    dataframes: Array<FrameDefinition>
  }
}
export {}

export declare const rgbaRegex: RegExp
export declare const toRGBA: ([r, g, b, a]: Array<number>) => string
export declare enum ChartType {
  pie = 'pie',
  bar = 'bar',
  line = 'line',
  doughnut = 'doughnut',
  funnel = 'funnel',
  gauge = 'gauge',
  radar = 'radar',
  scatter = 'scatter',
}
export interface TemporalDataPoint {
  t: Date
  y: number
}
export interface KV {
  [_: string]: any
}
export interface FormatData {
  format?: string
  prefix?: string
  suffix?: string
  presetFormat?: string
}
export interface Tooltip {
  formatting?: string
  labelsNextToPartition?: boolean
}
export interface TooltipParams {
  seriesName?: string
  name?: string
  value?: string | number
  percent?: string | number
  marker?: string
}
export interface Dimension {
  meta?: KV
  conditions: object
  field?: string
  modifier?: string
  default?: string
  skipMissing?: boolean
  timeLabels?: boolean
  autoSkip?: boolean
  rotateLabel?: number
}
export interface Metric {
  axisType?: string
  field?: string
  fixTooltips?: boolean
  relativeValue?: boolean
  cumulative?: boolean
  type?: ChartType
  alias?: string
  aggregate?: string
  modifier?: string
  fx?: string
  backgroundColor?: string
  symbol?: string
  formatting: FormatData
  [_: string]: any
}
export interface YAxis {
  axisPosition?: string
  axisType?: string
  beginAtZero?: boolean
  label?: string
  labelPosition?: string
  min?: string
  max?: string
  rotateLabel?: number
  horizontal?: boolean
  formatting: FormatData
}
export interface ChartOffset {
  top?: string
  right?: string
  bottom?: string
  left?: string
  isDefault?: boolean
}
export interface Position {
  isDefault?: boolean
  top?: string
  right?: string
  bottom?: string
  left?: string
}
export interface Legend {
  isHidden?: boolean
  orientation?: string
  align?: string
  isScrollable?: boolean
  isDefault?: boolean
  position?: Position
}
export interface Report {
  moduleID?: string | null
  filter?: string | null
  dimensions?: Array<Dimension>
  metrics?: Array<Metric>
  yAxis?: YAxis
  tooltip?: Tooltip
  legend?: Legend
  offset?: ChartOffset
}
export interface ChartToolbox {
  saveAsImage: boolean
  timeline: string
}
export interface ChartConfig {
  reports?: Array<Report>
  colorScheme?: string
  noAnimation?: boolean
  toolbox?: ChartToolbox
}
export declare const aggregateFunctions: {
  value: string
  text: string
}[]
interface DimensionFunction {
  text: string
  value: string
  convert: (f: string) => string
}
export declare class DimensionFunctions<T> extends Array<T> {
  private constructor()
  static create<T>(): DimensionFunctions<T>
  lookup(d: any): any
  convert(d: any): any
}
export declare const dimensionFunctions: DimensionFunctions<DimensionFunction>
export declare const predefinedFilters: {
  value: string
  text: string
}[]
export declare const isRadialChart: ({ type }: KV) => boolean
export declare const hasRelativeDisplay: ({ type }: KV) => boolean
export declare const makeAlias: ({ alias, aggregate, modifier, field }: Partial<Metric>) => string
export declare function formatChartValue(value: string | number, formatting?: FormatData): string
export declare function formatChartTooltip(tooltip: string, params: TooltipParams): string
export declare function defFormatData(): FormatData
declare const chartUtil: {
  dimensionFunctions: DimensionFunctions<DimensionFunction>
  hasRelativeDisplay: ({ type }: KV) => boolean
  aggregateFunctions: {
    value: string
    text: string
  }[]
  predefinedFilters: {
    value: string
    text: string
  }[]
  ChartType: typeof ChartType
}
export { chartUtil }

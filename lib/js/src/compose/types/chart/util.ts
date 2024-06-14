import numeral from 'numeral'
import * as fmt from '../../../formatting'

export const rgbaRegex = /^rgba\((\d+),.*?(\d+),.*?(\d+),.*?(\d*\.?\d*)\)$/

const ln = (n: number) => Math.round(n < 0 ? 255 + n : (n > 255) ? n - 255 : n)
export const toRGBA = ([r, g, b, a]: Array<number>) =>
  `rgba(${ln(r)}, ${ln(g)}, ${ln(b)}, ${a})`

export enum ChartType {
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
  t: Date;
  y: number;
}

export interface KV {
  [_: string]: any;
}

export interface FormatData {
  format?: string,
  prefix?: string,
  suffix?: string,
  presetFormat?: string,
}

export interface Tooltip {
  formatting?: string;
  labelsNextToPartition?: boolean;
}

export interface TooltipParams {
  seriesName?: string;
  name?: string;
  value?: string | number;
  percent?: string | number;
  marker?: string;
}

export interface Dimension {
  meta?: KV;
  conditions: object;
  field?: string;
  modifier?: string;
  default?: string;
  skipMissing?: boolean;
  timeLabels?: boolean;
  autoSkip?: boolean;
  rotateLabel?: number;
}

export interface Metric {
  axisType?: string;
  field?: string;
  fixTooltips?: boolean;
  relativeValue?: boolean;
  cumulative?: boolean;
  type?: ChartType;
  alias?: string;
  aggregate?: string;
  modifier?: string;
  fx?: string;
  backgroundColor?: string;
  symbol?: string;
  formatting: FormatData;
  [_: string]: any;
}

export interface YAxis {
  axisPosition?: string;
  axisType?: string;
  beginAtZero?: boolean;
  label?: string;
  labelPosition?: string;
  min?: string;
  max?: string;
  rotateLabel?: number;
  horizontal?: boolean;
  formatting: FormatData;
}

export interface ChartOffset {
  top?: string;
  right?: string;
  bottom?: string;
  left?: string;
  isDefault?: boolean;
}

export interface Position {
  isDefault?: boolean;
  top?: string;
  right?: string;
  bottom?: string;
  left?: string;
}

export interface Legend {
  isHidden?: boolean;
  orientation?: string;
  align?: string;
  isScrollable?: boolean;
  isDefault?: boolean;
  position?: Position;
}

export interface Report {
  moduleID?: string|null;
  filter?: string|null;
  dimensions?: Array<Dimension>;
  metrics?: Array<Metric>;
  yAxis?: YAxis;
  tooltip?: Tooltip;
  legend?: Legend;
  offset?: ChartOffset;
}

export interface ChartToolbox {
  saveAsImage: boolean;
  timeline: string;
}

export interface ChartConfig {
  reports?: Array<Report>;
  colorScheme?: string;
  noAnimation?: boolean;
  toolbox?: ChartToolbox;
}

export const aggregateFunctions = [
  {
    value: 'SUM',
    text: 'sum',
  },
  {
    value: 'MAX',
    text: 'max',
  },
  {
    value: 'MIN',
    text: 'min',
  },
  {
    value: 'AVG',
    text: 'avg',
  },
  {
    value: 'STD',
    text: 'std',
  },
]

interface DimensionFunction {
  text: string;
  value: string;
  convert: (f: string) => string;
}

export class DimensionFunctions<T> extends Array<T> {
  private constructor (items?: Array<T>) {
    super(...(items || []))
  }

  static create<T> (): DimensionFunctions<T> {
    return Object.create(DimensionFunctions.prototype)
  }

  public lookup (d: any): any {
    return this.find((f: any) => d.modifier === f.value)
  }

  public convert (d: any): any {
    return (this.lookup(d) || {}).convert(d.field)
  }
}

export const dimensionFunctions: DimensionFunctions<DimensionFunction> = DimensionFunctions.create<DimensionFunction>()
dimensionFunctions.push(...[
  {
    text: 'none',
    value: '(no grouping / buckets)',
    convert: (f: string) => f,
  },

  {
    text: 'date',
    value: 'DATE',
    convert: (f: string) => `DATE(${f})`,
  },

  {
    text: 'week',
    value: 'WEEK',
    convert: (f: string) => `WEEK(${f})`,
  },

  {
    text: 'month',
    value: 'MONTH',
    convert: (f: string) => `DATE_FORMAT(${f}, '%Y-%m-01')`,
  },

  {
    text: 'quarter',
    value: 'QUARTER',
    convert: (f: string) => `QUARTER(${f})`,
  },

  {
    text: 'year',
    value: 'YEAR',
    convert: (f: string) => `DATE_FORMAT(${f}, '%Y-01-01')`,
  },
])

export const predefinedFilters = [
  {
    value: 'YEAR(createdAt) = YEAR(NOW())',
    text: 'recordsCreatedThisYear',
  },
  {
    value: 'YEAR(createdAt) = YEAR(NOW()) - 1',
    text: 'recordsCreatedLastYear',
  },

  {
    value: 'YEAR(createdAt) = YEAR(NOW()) AND QUARTER(createdAt) = QUARTER(NOW())',
    text: 'recordsCreatedThisQuarter',
  },
  {
    value: 'YEAR(createdAt) = YEAR(NOW()) AND QUARTER(createdAt) = QUARTER(DATE_SUB(NOW(), INTERVAL 3 MONTH)',
    text: 'recordsCreatedLastQuarter',
  },

  {
    value: 'DATE_FORMAT(createdAt, \'%Y-%m\') = DATE_FORMAT(NOW(), \'%Y-%m\')',
    text: 'recordsCreatedThisMonth',
  },
  {
    value: 'DATE_FORMAT(createdAt, \'%Y-%m\') = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), \'%Y-%m\')',
    text: 'recordsCreatedLastMonth',
  },
]

dimensionFunctions.lookup = d => dimensionFunctions.find(f => d.modifier === f.value) || dimensionFunctions[0]
dimensionFunctions.convert = d => dimensionFunctions.lookup(d).convert(d.field)

export const isRadialChart = ({ type }: KV) => type === 'doughnut' || type === 'pie'
export const hasRelativeDisplay = ({ type }: KV) => isRadialChart({ type })

// Makes a standardized alias from modifier or dimension report option
export const makeAlias = ({ alias, aggregate, modifier, field }: Partial<Metric>) => alias || `${aggregate || modifier || 'none'}_${field}`.toLocaleLowerCase()

export function formatChartValue (value: string | number, formatting?: FormatData): string {
  let n: number | string = 0 || ''
  // if value contains alphabetic chars parseFloat() will return NaN
  // and n will equal 0
  const containsAlphabeticChars = isNaN(Number(value))
  let result = ''

  if (!containsAlphabeticChars) {
    switch (typeof value) {
      case 'string':
        n = parseFloat(value)
        break
      case 'number':
        n = value
        break
      default:
        n = 0
    }

    if (formatting?.format) {
      result = numeral(n).format(formatting.format)
    } else {
      result = fmt.number(n)
    }
  }

  if (formatting?.presetFormat === 'accounting') {
    result = fmt.accountingNumber(Number(n))
  }

  return ` ${formatting?.prefix ?? ''} ${result || value} ${formatting?.suffix ?? ''}`
}

export function formatChartTooltip (tooltip: string, params: TooltipParams): string {
  const { seriesName = '', name = '', value = '', percent = '' } = params

  return tooltip
    .replace('{a}', seriesName)
    .replace('{b}', name)
    .replace('{c}', value.toString())
    .replace('{d}', percent.toString())
}

export function defFormatData (): FormatData {
  return Object.assign({}, {
    presetFormat: 'custom',
    prefix: '',
    suffix: '',
    format: '',
  })
}

const chartUtil = {
  dimensionFunctions,
  hasRelativeDisplay,
  aggregateFunctions,
  predefinedFilters,
  ChartType,
}

export {
  chartUtil,
}

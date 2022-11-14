import { DisplayElement, DisplayElementInput, Registry } from '../base'
import { DefinitionOptions, FrameDefinition } from '../../frame'
import { Apply } from '../../../../cast'

const kind = 'Chart'

export type PartialChartOptions = Partial<ChartOptions>

interface XAxisOptions {
  type: string;
  label?: string;
  unit?: string;
  skipMissing: boolean;
  defaultValue?: any;
}

interface YAxisOptions {
  label?: string;
  type?: string;
  position?: string;
  beginAtZero?: boolean;
  stepSize?: string;
  min?: string;
  max?: string;
}

export class ChartOptions {
  public title = ''
  public type = 'bar'
  public colorScheme = ''
  public source = ''
  public datasources: Array<FrameDefinition> = []
  public showTooltips = true
  public showLegend = true

  public xAxis: XAxisOptions = {
    type: '',
    skipMissing: true,
  }

  public yAxis: YAxisOptions = {
    type: 'linear',
    position: 'left',
    beginAtZero: true,
  }

  constructor (o: PartialChartOptions = {}) {
    if (!o) return

    Apply(this, o, String, 'title', 'type', 'colorScheme', 'source')
    Apply(this, o, Boolean, 'showTooltips', 'showLegend')

    if (o.datasources) {
      this.datasources = o.datasources
    }

    if (o.xAxis) {
      this.xAxis = { ...this.xAxis, ...o.xAxis }
    }

    if (o.yAxis) {
      this.yAxis = { ...this.yAxis, ...o.yAxis }
    }
  }
}

export const ChartOptionsRegistry = new Map<string, typeof ChartOptions>()

export function ChartOptionsMaker<T extends ChartOptions> (options: Partial<ChartOptions>): T {
  const { type } = options

  if (type) {
    const ChartOptionsTemp = ChartOptionsRegistry.get(type)
    if (ChartOptionsTemp === undefined) {
      throw new Error(`unknown chart type '${type}'`)
    }

    if (options instanceof ChartOptions) {
      // Get rid of the references
      options = JSON.parse(JSON.stringify(options))
    }

    return new ChartOptionsTemp(options) as T
  } else {
    throw new Error('no chart type')
  }
}

export class DisplayElementChart extends DisplayElement {
  readonly kind = kind

  options: ChartOptions = ChartOptionsMaker({ type: 'bar' })

  constructor (i?: DisplayElementInput) {
    super(i)
    this.applyOptions(i?.options as Partial<ChartOptions>)
  }

  applyOptions (o?: PartialChartOptions): void {
    if (!o) return

    this.options = ChartOptionsMaker(o)
  }

  reportDefinitions (definition: DefinitionOptions = {}): { dataframes: Array<FrameDefinition> } {
    if (typeof this.options.source === 'object') {
      // @todo allow implicit sources
      throw new Error('chart source must be provided as a reference')
    }

    const dataframes: Array<FrameDefinition> = []

    this.options.datasources.forEach(({ name = '', filter, sort }) => {
      const df: FrameDefinition = {
        name: this.elementID,
        source: this.options.source,
        ref: name,
        filter,
        sort,
      }

      const relatedDefinition = definition[name]

      if (relatedDefinition) {
        df.sort = (relatedDefinition.sort ? relatedDefinition.sort : sort) || undefined

        if (relatedDefinition.filter && relatedDefinition.filter?.ref) {
          // If element and scenario have filter AND them together
          if (filter && filter.ref) {
            df.filter = {
              ref: 'and',
              args: [
                filter,
                relatedDefinition.filter,
              ],
            }
          } else {
            df.filter = relatedDefinition.filter
          }
        }
      }

      dataframes.push(df)
    })

    return { dataframes }
  }
}

Registry.set(kind, DisplayElementChart)

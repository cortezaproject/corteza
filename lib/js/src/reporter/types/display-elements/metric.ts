import { DisplayElement, DisplayElementInput, Registry } from './base'
import { FrameDefinition, DefinitionOptions } from '../frame'
import { Apply } from '../../../cast'

const kind = 'Metric'

interface Options {
  source?: string;
  datasources: Array<FrameDefinition>;

  valueColumn: string;

  format: string;
  prefix: string;
  suffix: string;

  color: string;
  backgroundColor: string;
}

const defaults: Readonly<Options> = Object.freeze({
  source: '',
  datasources: [],

  valueColumn: '',

  format: '',
  prefix: '',
  suffix: '',

  color: '#000000',
  backgroundColor: '#ffffff',
})

export class DisplayElementMetric extends DisplayElement {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: DisplayElementInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'source', 'valueColumn', 'format', 'prefix', 'suffix', 'color', 'backgroundColor')

    if (o.datasources) {
      this.options.datasources = o.datasources || []
    }
  }

  reportDefinitions (definition: DefinitionOptions = {}): { dataframes: Array<FrameDefinition> } {
    if (typeof this.options.source === 'object') {
      // @todo allow implicit sources
      throw new Error('metric source must be provided as a reference')
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

Registry.set(kind, DisplayElementMetric)

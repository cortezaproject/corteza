import { DisplayElement, DisplayElementInput } from './base'
import { FrameDefinition, FrameColumn, DefinitionOptions } from '../frame'
interface TableColumns {
  [key: string]: Array<FrameColumn>
}
interface Options {
  source?: string
  datasources: Array<FrameDefinition>
  columns?: TableColumns
  striped: boolean
  bordered: boolean
  borderless: boolean
  small: boolean
  hover: boolean
  dark: boolean
  fixed: boolean
  responsive: boolean
  noCollapse: boolean
  headVariant: string | null
  tableVariant: string
}
export declare class DisplayElementTable extends DisplayElement {
  readonly kind = 'Table'
  options: Options
  constructor(i?: DisplayElementInput)
  applyOptions(o?: Partial<Options>): void
  reportDefinitions(definition?: DefinitionOptions): {
    dataframes: Array<FrameDefinition>
  }
}
export {}

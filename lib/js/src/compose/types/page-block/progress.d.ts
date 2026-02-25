import { PageBlock, PageBlockInput } from './base'
import { Compose as ComposeAPI } from '../../../api-clients'
interface ValueOptions {
  default: number
  moduleID: string
  filter: string
  field: string
  operation: string
}
interface Threshold {
  value: number
  variant: string
}
interface DisplayOptions {
  showValue: boolean
  showRelative: boolean
  showProgress: boolean
  animated: boolean
  variant: string
  thresholds: Threshold[]
}
interface Options {
  value: ValueOptions
  minValue: ValueOptions
  maxValue: ValueOptions
  display: DisplayOptions
  refreshRate: number
  showRefresh: boolean
  magnifyOption: string
}
export declare class PageBlockProgress extends PageBlock {
  readonly kind = 'Progress'
  options: Options
  constructor(i?: PageBlockInput)
  applyOptions(o?: Partial<Options>): void
  /**
   * Helper function to fetch and parse reporter's reports.
   */
  fetch(additionalOptions: Options, api: ComposeAPI, namespaceID: string): Promise<object>
}
export {}

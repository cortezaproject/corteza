import { PageBlock, PageBlockInput } from './base'
interface Options {
  reportID: string
  scenarioID: string
  elementID: string
  refreshRate: number
  showRefresh: boolean
  magnifyOption: string
}
export declare class PageBlockReport extends PageBlock {
  readonly kind = 'Report'
  options: Options
  constructor(i?: PageBlockInput)
  applyOptions(o?: Partial<Options>): void
}
export {}

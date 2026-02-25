import { PageBlock, PageBlockInput } from './base'
import { PageBlockWrap } from './types'
interface Options {
  srcField: string
  src: string
  wrap: PageBlockWrap
  refreshRate: number
  showRefresh: boolean
  magnifyOption: string
}
export declare class PageBlockIFrame extends PageBlock {
  readonly kind = 'IFrame'
  options: Options
  constructor(i?: PageBlockInput)
  applyOptions(o?: Partial<Options>): void
}
export {}

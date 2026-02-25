import { PageBlock, PageBlockInput } from './base'
interface Options {
  body: string
  magnifyOption: string
}
export declare class PageBlockContent extends PageBlock {
  readonly kind = 'Content'
  options: Options
  constructor(i?: PageBlockInput)
  applyOptions(o?: Partial<Options>): void
}
export {}

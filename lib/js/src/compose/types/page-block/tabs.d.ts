import { PageBlock, PageBlockInput } from './base'
interface Style {
  appearance: string
  alignment: string
  justify: string
  orientation: string
  position: string
}
interface Tab {
  blockID: string
  title: string
}
interface Options {
  style: Style
  tabs: Tab[]
  magnifyOption: string
}
export declare class PageBlockTab extends PageBlock {
  readonly kind = 'Tabs'
  options: Options
  constructor(i?: PageBlockInput)
  applyOptions(o?: Partial<Options>): void
}
export {}

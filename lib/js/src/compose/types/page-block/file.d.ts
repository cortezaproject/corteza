import { PageBlock, PageBlockInput } from './base'
interface Options {
  mode: string
  attachments: string[]
  hideFileName: boolean
  height: string
  width: string
  maxHeight: string
  maxWidth: string
  borderRadius: string
  margin: string
  backgroundColor: string
  magnifyOption: string
  clickToView?: boolean
  enableDownload?: boolean
}
export declare class PageBlockFile extends PageBlock {
  readonly kind = 'File'
  options: Options
  constructor(i?: PageBlockInput)
  applyOptions(o?: Partial<Options>): void
}
export {}

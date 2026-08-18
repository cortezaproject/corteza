import { PageBlock, PageBlockInput, Registry, BlockIssue } from './base'
import { Apply } from '../../../cast'

const kind = 'IFrame'
interface Options {
  srcField: string
  src: string
  refreshRate: number
  showRefresh: boolean
  magnifyOption: string
}

const defaults: Readonly<Options> = Object.freeze({
  srcField: '',
  src: '',
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
})

export class PageBlockIFrame extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options & { url?: string }>): void {
    if (!o) return

    // Handle legacy 'url' key: map to 'src' if 'src' is not provided
    if (!o.src && (o as Record<string, unknown>).url) {
      o = { ...o, src: (o as Record<string, unknown>).url as string }
    }

    Apply(this.options, o, String, 'srcField', 'src', 'magnifyOption')
    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')
  }

  validate(): Array<BlockIssue> {
    const ee = super.validate()
    if (!this.options.src && !this.options.srcField) {
      ee.push({ option: 'src', labelKey: 'block.issue.iframeSource' })
    }
    return ee
  }
}

Registry.set(kind, PageBlockIFrame)

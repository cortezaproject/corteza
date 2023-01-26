import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'
import { PageBlockWrap } from './types'

const kind = 'IFrame'
interface Options {
  srcField: string;
  src: string;
  wrap: PageBlockWrap;
  refreshRate: number;
  showRefresh: boolean;
  magnifyOption: string;
}

const defaults: Readonly<Options> = Object.freeze({
  srcField: '',
  src: '',
  wrap: 'Plain',
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
})

export class PageBlockIFrame extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'srcField', 'src', 'wrap', 'magnifyOption')
    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')
  }
}

Registry.set(kind, PageBlockIFrame)

import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'Content'

interface Options {
  body: string;
  magnifyOption: string
}

const defaults: Readonly<Options> = Object.freeze({
  body: '',
  magnifyOption: ''
})

export class PageBlockContent extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'body', 'magnifyOption')
  }
}

Registry.set(kind, PageBlockContent)

import { PageBlock, PageBlockInput, Registry, BlockIssue } from './base'
import { Apply } from '../../../cast'

const kind = 'Content'

interface Options {
  body: string
  magnifyOption: string
}

const defaults: Readonly<Options> = Object.freeze({
  body: '',
  magnifyOption: '',
})

export class PageBlockContent extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'body', 'magnifyOption')
  }

  validate(): Array<BlockIssue> {
    const ee = super.validate()
    if (!this.options.body) {
      ee.push({ option: 'body', labelKey: 'block.issue.contentBody' })
    }
    return ee
  }
}

Registry.set(kind, PageBlockContent)

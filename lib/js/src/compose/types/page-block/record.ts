import { PageBlock, PageBlockInput, Registry } from './base'

const kind = 'Record'

interface Options {
  fields: unknown[];
}

const defaults: Readonly<Options> = Object.freeze({
  fields: [],
})

export class PageBlockRecord extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    if (o.fields) {
      this.options.fields = o.fields
    }
  }
}

Registry.set(kind, PageBlockRecord)

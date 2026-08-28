import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'Group'

interface GroupBlock {
  blockID: string;
  xywh: number[];
}

interface Options {
  blocks: GroupBlock[];
  padding: number;
}

const defaults: Readonly<Options> = Object.freeze({
  blocks: [],
  padding: 0,
})

export class PageBlockGroup extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    if (o.blocks) {
      this.options.blocks = o.blocks
    }

    if (o.padding !== undefined) {
      this.options.padding = o.padding
    }
  }
}

Registry.set(kind, PageBlockGroup)

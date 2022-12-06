import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'Nylas'

interface Options {
  kind: string,
  componentID: string;
}

const defaults: Readonly<Options> = Object.freeze({
  kind: 'Composer', // Default kind of nylas component
  componentID: '',
})

export class PageBlockNylas extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, String, 'kind', 'componentID')
  }
}

Registry.set(kind, PageBlockNylas)

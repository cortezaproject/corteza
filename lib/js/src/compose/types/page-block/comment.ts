import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'

const kind = 'Comment'
interface Options {
  moduleID: string;
  filter: string;
  titleField: string;
  contentField: string;
  referenceField: string;
  sortDirection: string;
  refreshRate: number;
}

const defaults: Readonly<Options> = Object.freeze({
  moduleID: NoID,
  filter: '',
  titleField: '',
  contentField: '',
  sortDirection: '',
  referenceField: '',
  refreshRate: 0,
})

export class PageBlockComment extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return
    Apply(this.options, o, CortezaID, 'moduleID')
    Apply(this.options, o, String, 'titleField', 'contentField', 'referenceField', 'filter', 'sortDirection')
    Apply(this.options, o, Number, 'refreshRate')
  }
}

Registry.set(kind, PageBlockComment)

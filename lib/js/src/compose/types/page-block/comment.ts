import { Apply, CortezaID, NoID } from '../../../cast';
import { PageBlock, PageBlockInput, Registry } from './base';

const kind = 'Comment'
interface Options {
  moduleID: string;
  filter: string;
  titleField: string;
  contentField: string;
  replyField: string;
  referenceField: string;
  sortDirection: string;
  refreshRate: number;
  showRefresh: boolean;
  magnifyOption: string;
  attachmentField: string;
}

const defaults: Readonly<Options> = Object.freeze({
  moduleID: NoID,
  filter: '',
  titleField: '',
  contentField: '',
  replyField: '',
  sortDirection: 'asc',
  referenceField: '',
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
  attachmentField: '',
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
    Apply(this.options, o, String, 'titleField', 'contentField', 'replyField', 'referenceField', 'attachmentField', 'filter', 'sortDirection', 'magnifyOption')
    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')
  }
}

Registry.set(kind, PageBlockComment)

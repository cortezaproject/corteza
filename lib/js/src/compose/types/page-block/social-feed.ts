import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'

const kind = 'SocialFeed'
interface Options {
  moduleID: string;
  fields: unknown[];
  profileSourceField: string;
  profileUrl: string;
  refreshRate: number;
}

const defaults: Readonly<Options> = Object.freeze({
  moduleID: NoID,
  fields: [],
  profileSourceField: '',
  profileUrl: '',
  refreshRate: 0,
})

export class PageBlockSocialFeed extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, CortezaID, 'moduleID')
    Apply(this.options, o, String, 'profileSourceField', 'profileUrl')
    Apply(this.options, o, Number, 'refreshRate')

    if (o.fields) {
      this.options.fields = o.fields
    }
  }
}

Registry.set(kind, PageBlockSocialFeed)

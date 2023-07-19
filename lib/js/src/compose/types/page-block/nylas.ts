import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'Nylas'

interface Options {
  kind: string,
  componentID: string;
  accessTokenRequired: boolean;
  magnifyOption: string;
  prefill: Prefill;
}

interface Prefill {
  to: string;
  subject: string;
  body: string;
  queryString: string;
}

const defaults: Readonly<Options> = Object.freeze({
  kind: 'Composer', // Default kind of nylas component
  componentID: '',
  accessTokenRequired: true,
  magnifyOption: '',
  prefill: {
    to: '',
    subject: '',
    body: '',
    queryString: ''
  }
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

    Apply(this.options, o, String, 'kind', 'componentID', 'magnifyOption')
    Apply(this.options, o, Boolean, 'accessTokenRequired')

    this.options.prefill = { ...defaults.prefill, ...o.prefill }
  }
}

Registry.set(kind, PageBlockNylas)

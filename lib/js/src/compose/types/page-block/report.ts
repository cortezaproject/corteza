import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'

const kind = 'Report'

interface Options {
  reportID: string;
  scenarioID: string;
  elementID: string;
  refreshRate: number;
  showRefresh: boolean;
  magnifyOption: string;
}

const defaults: Readonly<Options> = Object.freeze({
  reportID: NoID,
  scenarioID: NoID,
  elementID: NoID,
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
})

export class PageBlockReport extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, CortezaID, 'reportID', 'scenarioID', 'elementID')
    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')
    Apply(this.options, o, String, 'magnifyOption')
  }
}

Registry.set(kind, PageBlockReport)

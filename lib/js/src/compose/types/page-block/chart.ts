import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'

const kind = 'Chart'

interface Options {
  chartID: string;
  refreshRate: number;
}

const defaults: Readonly<Options> = Object.freeze({
  chartID: NoID,
  refreshRate: 0,
})

export class PageBlockChart extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, CortezaID, 'chartID')
    Apply(this.options, o, Number, 'refreshRate')
  }
}

Registry.set(kind, PageBlockChart)

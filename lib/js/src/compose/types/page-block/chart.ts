import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, CortezaID, NoID } from '../../../cast'

const kind = 'Chart'

interface DrillDown {
  enabled: boolean;
  blockID?: string;
}

interface Options {
  chartID: string;
  refreshRate: number;
  showRefresh: boolean;
  magnifyOption: string;
  drillDown: DrillDown;
}

const defaults: Readonly<Options> = Object.freeze({
  chartID: '',
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
  drillDown: {
    enabled: false,
    blockID: ''
  }
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

    o.chartID = o.chartID === NoID ? '' : o.chartID

    Apply(this.options, o, String, 'chartID', 'magnifyOption')
    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')

    if (o.drillDown) {
      this.options.drillDown = o.drillDown
    }
  }
}

Registry.set(kind, PageBlockChart)

import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply, NoID } from '../../../cast'
import { Options as PageBlockRecordListOptions } from './record-list'
import { cloneDeep, merge } from 'lodash'

const kind = 'Chart'

interface DrillDown {
  enabled: boolean;
  blockID: string;
  recordListOptions: Partial<PageBlockRecordListOptions>;
}

interface Options {
  chartID: string;
  refreshRate: number;
  showRefresh: boolean;
  magnifyOption: string;
  drillDown: DrillDown;
  liveFilterEnabled: boolean;
}

const defaults: Readonly<Options> = Object.freeze({
  chartID: '',
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
  liveFilterEnabled: false,
  drillDown: {
    enabled: false,
    blockID: '',
    recordListOptions: {
      fields: [],
    },
  },
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
    Apply(this.options, o, Boolean, 'showRefresh', 'liveFilterEnabled')

    if (o.drillDown) {
      this.options.drillDown = merge({}, defaults.drillDown, o.drillDown)
    }
  }

  resetDrillDown (): void {
    this.options.drillDown = cloneDeep(defaults.drillDown)
  }
}

Registry.set(kind, PageBlockChart)

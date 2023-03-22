import { PageBlock, PageBlockInput, Registry } from './base'
import { dimensionFunctions } from '../chart/util'
import { CortezaID, Apply } from '../../../cast'

const kind = 'Metric'

type Reporter = (p: ReporterParams) => Promise<any>

interface ReporterParams {
  moduleID: string;
  filter?: string;
  metrics?: string;
  dimensions: string;
}

interface Style {
  color: string;
  backgroundColor: string;
  fontSize: string;
}

interface Metric {
  label: string;
  moduleID: string;
  dimensionField: string;
  dateFormat?: string;
  filter?: string;
  bucketSize?: string;
  metricField: string;
  operation: string;
  numberFormat?: string;
  prefix?: string;
  suffix?: string;
  transformFx?: string;

  // @todo allow conditional styles; eg. if value is < 10 render with bold red text
  labelStyle?: Style;
  valueStyle?: Style;
}

interface Options {
  metrics: Array<Metric>;
  refreshRate: number;
  showRefresh: boolean;
  magnifyOption: string;
}

const defaults: Readonly<Options> = Object.freeze({
  metrics: [],
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
})

export class PageBlockMetric extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return
    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')
    Apply(this.options, o, String, 'magnifyOption')
    if (o.metrics) {
      this.options.metrics = o.metrics
    }
  }

  /**
   * Helper function to fetch and parse reporter's reports.
   */
  async fetch ({ m }: { m: Metric }, reporter: Reporter): Promise<object> {
    const w = await reporter(this.formatParams(m))
    const datasets = w.map((r: any) => r.rp !== undefined ? r.rp : r.count)

    let rtr: number
    if (m.operation === 'max') {
      rtr = datasets.sort((a: number, b: number) => b - a)[0]
    } else if (m.operation === 'min') {
      rtr = datasets.sort((a: number, b: number) => a - b)[0]
    } else if (m.operation === 'avg') {
      rtr = datasets.reduce((acc: number, cur: number) => acc + cur, 0) / datasets.length
    } else {
      rtr = datasets.reduce((acc: number, cur: number) => acc + cur, 0)
    }

    if (m.transformFx) {
      rtr = (new Function('v', `return ${m.transformFx}`))(rtr)
    }

    return [{ value: rtr }]
  }

  /**
   * Helper to construct reporter's params
   */
  private formatParams ({ moduleID, filter, metricField, operation = '' }: Metric): ReporterParams {
    let metrics = ''

    if (operation && metricField && metricField !== 'count') {
      metrics = `${operation}(${metricField}) AS rp`
    }

    return {
      moduleID,
      filter,
      metrics,
      dimensions: dimensionFunctions.convert({ modifier: 'YEAR', field: 'createdAt' }),
    }
  }
}

Registry.set(kind, PageBlockMetric)

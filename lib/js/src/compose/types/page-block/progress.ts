import { PageBlock, PageBlockInput, Registry } from './base'
import { dimensionFunctions } from '../chart/util'
import { Compose as ComposeAPI } from '../../../api-clients'

const kind = 'Progress'

interface ValueOptions {
  moduleID: string;
  filter: string;
  field: string;
  operation: string;
}

interface Threshold {
  value: number;
  variant: string;
}

interface DisplayOptions {
  showValue: boolean;
  showRelative: boolean;
  showProgress: boolean;
  animated: boolean;
  variant: string;
  thresholds: Threshold[];
}

interface Options {
  value: ValueOptions;
  maxValue: ValueOptions;
  display: DisplayOptions;
}

const defaults: Readonly<Options> = Object.freeze({
  value: {
    moduleID: '',
    filter: '',
    field: '',
    operation: '',
  },

  maxValue: {
    moduleID: '',
    filter: '',
    field: '',
    operation: '',
  },

  display: {
    showValue: true,
    showRelative: true,
    showProgress: false,
    animated: false,
    variant: 'success',
    thresholds: [],
  },
})

export class PageBlockProgress extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    if (o.value) {
      this.options.value = o.value
    }

    if (o.maxValue) {
      this.options.maxValue = o.maxValue
    }

    if (o.display) {
      this.options.display = o.display
    }
  }

  /**
   * Helper function to fetch and parse reporter's reports.
   */
  fetch (options: Options, api: ComposeAPI, namespaceID: string): Promise<object> {
    const dimensions = dimensionFunctions.convert({ modifier: 'YEAR', field: 'createdAt' })

    let metrics = ''

    // Construct value report
    const { field: valueField, operation: valueOperation = '' } = this.options.value
    if (valueOperation && valueField && valueField !== 'count') {
      metrics = `${valueOperation}(${valueField}) AS rp`
    }
    const valueReport = api.recordReport({ namespaceID, metrics, dimensions, ...this.options.value, ...options.value })

    // Construct max value report
    metrics = ''
    const { field: maxValueField, operation: maxValueOperation = '' } = this.options.maxValue
    if (maxValueOperation && maxValueField && maxValueField !== 'count') {
      metrics = `${maxValueOperation}(${maxValueField}) AS rp`
    }
    const maxValueReport = api.recordReport({ namespaceID, metrics, dimensions, ...this.options.maxValue, ...options.maxValue })

    return Promise.all([valueReport, maxValueReport]).then(([v, m]: Array<any>) => {
      let value: number
      let max: number

      let datasets = v.map((r: any) => r.rp !== undefined ? r.rp : r.count)
      if (valueOperation === 'max') {
        value = datasets.sort((a: number, b: number) => b - a)[0]
      } else if (valueOperation === 'min') {
        value = datasets.sort((a: number, b: number) => a - b)[0]
      } else if (valueOperation === 'avg') {
        value = datasets.reduce((acc: number, cur: number) => acc + cur, 0) / datasets.length
      } else {
        value = datasets.reduce((acc: number, cur: number) => acc + cur, 0)
      }

      datasets = m.map((r: any) => r.rp !== undefined ? r.rp : r.count)
      if (maxValueOperation === 'max') {
        max = datasets.sort((a: number, b: number) => b - a)[0]
      } else if (maxValueOperation === 'min') {
        max = datasets.sort((a: number, b: number) => a - b)[0]
      } else if (maxValueOperation === 'avg') {
        max = datasets.reduce((acc: number, cur: number) => acc + cur, 0) / datasets.length
      } else {
        max = datasets.reduce((acc: number, cur: number) => acc + cur, 0)
      }

      return { value, max }
    })
  }
}

Registry.set(kind, PageBlockProgress)

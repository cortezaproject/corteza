import _ from 'lodash'
import { PageBlock, PageBlockInput, Registry } from './base'
import { dimensionFunctions } from '../chart/util'
import { Compose as ComposeAPI } from '../../../api-clients'
import { Apply } from '../../../cast'

const kind = 'Progress'

interface ValueOptions {
  default: number;
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
  minValue: ValueOptions;
  maxValue: ValueOptions;
  display: DisplayOptions;
  refreshRate: number;
  showRefresh: boolean;
  magnifyOption: string;
}

const defaults: Readonly<Options> = Object.freeze({
  value: {
    default: 0,
    moduleID: '',
    filter: '',
    field: '',
    operation: '',
  },

  minValue: {
    default: 0,
    moduleID: '',
    filter: '',
    field: '',
    operation: '',
  },

  maxValue: {
    default: 100,
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
  refreshRate: 0,
  showRefresh: false,
  magnifyOption: '',
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

    Apply(this.options, o, Number, 'refreshRate')
    Apply(this.options, o, Boolean, 'showRefresh')
    Apply(this.options, o, String, 'magnifyOption')

    if (o.value) {
      this.options.value = { ...this.options.value, ...o.value }
    }

    if (o.minValue) {
      this.options.minValue = { ...this.options.minValue, ...o.minValue }
    }

    if (o.maxValue) {
      this.options.maxValue = { ...this.options.maxValue, ...o.maxValue }
    }

    if (o.display) {
      this.options.display = { ...this.options.display, ...o.display }
    }
  }

  /**
   * Helper function to fetch and parse reporter's reports.
   */
  fetch (options: Options, api: ComposeAPI, namespaceID: string): Promise<object> {
    options = _.merge(this.options, options)
    const reports = []
    const dimensions = dimensionFunctions.convert({ modifier: 'YEAR', field: 'createdAt' })

    let metrics = ''

    // Construct value report
    const { field: valueField, operation: valueOperation = '' } = this.options.value

    if (options.value.moduleID && valueField) {
      if (valueOperation && valueField !== 'count') {
        metrics = `${valueOperation}(${valueField}) AS rp`
      }

      reports.push(api.recordReport({ namespaceID, metrics, dimensions, ...options.value }))
    } else {
      reports.push(new Promise(resolve => resolve(options.value.default)))
    }

    // Construct minValue report
    const { field: minValueField, operation: minValueOperation = '' } = this.options.minValue

    if (options.minValue.moduleID && minValueField) {
      metrics = ''
      if (minValueOperation && minValueField !== 'count') {
        metrics = `${minValueOperation}(${minValueField}) AS rp`
      }

      reports.push(api.recordReport({ namespaceID, metrics, dimensions, ...options.minValue }))
    } else {
      reports.push(new Promise(resolve => resolve(options.minValue.default)))
    }

    // Construct minValue report
    const { field: maxValueField, operation: maxValueOperation = '' } = this.options.maxValue

    if (options.maxValue.moduleID && maxValueField) {
      metrics = ''
      if (maxValueOperation && maxValueField !== 'count') {
        metrics = `${maxValueOperation}(${maxValueField}) AS rp`
      }

      reports.push(api.recordReport({ namespaceID, metrics, dimensions, ...options.maxValue }))
    } else {
      reports.push(new Promise(resolve => resolve(options.maxValue.default)))
    }

    return Promise.all(reports).then(([value, min, max]: Array<any>) => {
      if (Array.isArray(value)) {
        const datasets = value.map((r: any) => r.rp !== undefined ? r.rp : r.count)
        if (valueOperation === 'max') {
          value = datasets.sort((a: number, b: number) => b - a)[0]
        } else if (valueOperation === 'min') {
          value = datasets.sort((a: number, b: number) => a - b)[0]
        } else if (valueOperation === 'avg') {
          value = datasets.reduce((acc: number, cur: number) => acc + cur, 0) / datasets.length
        } else {
          value = datasets.reduce((acc: number, cur: number) => acc + cur, 0)
        }
      }

      if (Array.isArray(min)) {
        const datasets = min.map((r: any) => r.rp !== undefined ? r.rp : r.count)
        if (minValueOperation === 'max') {
          min = datasets.sort((a: number, b: number) => b - a)[0]
        } else if (minValueOperation === 'min') {
          min = datasets.sort((a: number, b: number) => a - b)[0]
        } else if (minValueOperation === 'avg') {
          min = datasets.reduce((acc: number, cur: number) => acc + cur, 0) / datasets.length
        } else {
          min = datasets.reduce((acc: number, cur: number) => acc + cur, 0)
        }
      }

      if (Array.isArray(max)) {
        const datasets = max.map((r: any) => r.rp !== undefined ? r.rp : r.count)
        if (maxValueOperation === 'max') {
          max = datasets.sort((a: number, b: number) => b - a)[0]
        } else if (maxValueOperation === 'min') {
          max = datasets.sort((a: number, b: number) => a - b)[0]
        } else if (maxValueOperation === 'avg') {
          max = datasets.reduce((acc: number, cur: number) => acc + cur, 0) / datasets.length
        } else {
          max = datasets.reduce((acc: number, cur: number) => acc + cur, 0)
        }
      }

      return { value, min, max }
    })
  }
}

Registry.set(kind, PageBlockProgress)

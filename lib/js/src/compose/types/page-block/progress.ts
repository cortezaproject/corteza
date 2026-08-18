import { Compose as ComposeAPI } from '../../../api-clients'
import { Apply } from '../../../cast'
import { dimensionFunctions } from '../chart/util'
import { PageBlock, PageBlockInput, Registry, BlockIssue, isUnsetID } from './base'

const kind = 'Progress'

interface ValueOptions {
  default: number
  moduleID: string
  filter: string
  field: string
  operation: string
}

interface Threshold {
  value: number
  variant: string
}

interface DisplayOptions {
  showValue: boolean
  showRelative: boolean
  showProgress: boolean
  variant: string
  thresholds: Threshold[]
}

interface Options {
  value: ValueOptions
  minValue: ValueOptions
  maxValue: ValueOptions
  display: DisplayOptions
  refreshRate: number
  showRefresh: boolean
  magnifyOption: string
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

  constructor(i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions(o?: Partial<Options>): void {
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

    // Default field to 'count' when a module is selected — keeps fetch consistent with UI
    for (const section of ['value', 'minValue', 'maxValue'] as const) {
      const s = this.options[section]
      if (s.moduleID && !s.field) s.field = 'count'
    }

    // Migrate legacy top-level `thresholds` (older configurator wrote here) → display.thresholds
    const legacyThresholds = (o as any)?.thresholds
    if (Array.isArray(legacyThresholds) && !this.options.display.thresholds?.length) {
      this.options.display.thresholds = legacyThresholds
    }

    // Migrate legacy 'light' variant (no PrimeVue Button equivalent) → 'secondary'
    if (this.options.display.variant === 'light') {
      this.options.display.variant = 'secondary'
    }
    if (Array.isArray(this.options.display.thresholds)) {
      for (const th of this.options.display.thresholds) {
        if (th.variant === 'light') th.variant = 'secondary'
      }
    }
  }

  /**
   * Helper function to fetch and parse reporter's reports.
   */
  fetch(additionalOptions: Options, api: ComposeAPI, namespaceID: string): Promise<object> {
    const reports = []
    const dimensions = dimensionFunctions.convert({ modifier: 'YEAR', field: 'createdAt' })

    let metrics = ''

    // Construct value report
    const { field: valueField, operation: valueOperation = '' } = this.options.value

    if (this.options.value.moduleID && valueField) {
      if (valueOperation && valueField !== 'count') {
        metrics = `${valueOperation}(${valueField}) AS rp`
      }

      reports.push(
        api.recordReport({
          namespaceID,
          metrics,
          dimensions,
          ...this.options.value,
          ...additionalOptions.value,
        }),
      )
    } else {
      reports.push(new Promise(resolve => resolve(this.options.value.default)))
    }

    // Construct minValue report
    const { field: minValueField, operation: minValueOperation = '' } = this.options.minValue

    if (this.options.minValue.moduleID && minValueField) {
      metrics = ''
      if (minValueOperation && minValueField !== 'count') {
        metrics = `${minValueOperation}(${minValueField}) AS rp`
      }

      reports.push(
        api.recordReport({
          namespaceID,
          metrics,
          dimensions,
          ...this.options.minValue,
          ...additionalOptions.minValue,
        }),
      )
    } else {
      reports.push(new Promise(resolve => resolve(this.options.minValue.default)))
    }

    // Construct minValue report
    const { field: maxValueField, operation: maxValueOperation = '' } = this.options.maxValue

    if (this.options.maxValue.moduleID && maxValueField) {
      metrics = ''
      if (maxValueOperation && maxValueField !== 'count') {
        metrics = `${maxValueOperation}(${maxValueField}) AS rp`
      }

      reports.push(
        api.recordReport({
          namespaceID,
          metrics,
          dimensions,
          ...this.options.maxValue,
          ...additionalOptions.maxValue,
        }),
      )
    } else {
      reports.push(new Promise(resolve => resolve(this.options.maxValue.default)))
    }

    return Promise.all(reports).then(([value, min, max]: Array<any>) => {
      if (Array.isArray(value)) {
        const datasets = value.map((r: any) => (r.rp !== undefined ? r.rp : r.count))
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
        const datasets = min.map((r: any) => (r.rp !== undefined ? r.rp : r.count))
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
        const datasets = max.map((r: any) => (r.rp !== undefined ? r.rp : r.count))
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

  validate(): Array<BlockIssue> {
    const ee = super.validate()

    // A bar with no source reads as a fixed default, which is a value nobody set.
    const { value } = this.options
    if (isUnsetID(value.moduleID) && value.default === undefined) {
      ee.push({ option: 'value', labelKey: 'block.issue.progressValue' })
    }

    return ee
  }
}

Registry.set(kind, PageBlockProgress)

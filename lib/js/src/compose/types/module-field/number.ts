import numeral from 'numeral'
import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'
import * as fmt from '../../../formatting'

const kind = 'Number'

interface Threshold {
  value: number;
  variant: string;
}

interface NumberOptions extends Options {
  presetFormat: string;
  format: string;
  prefix: string;
  suffix: string;
  precision: number;
  multiDelimiter: string;
  display: string;
  min: number;
  max: number;
  showValue: boolean;
  showRelative: boolean;
  showProgress: boolean;
  animated: boolean;
  variant: string;
  thresholds: Threshold[];
}

const defaults = (): Readonly<NumberOptions> => Object.freeze({
  ...defaultOptions(),
  presetFormat: 'currencyFormat',
  precision: 3,
  multiDelimiter: '\n',
  display: 'number', // Either number or progress (progress bar)
  // Number display options
  format: '',
  prefix: '',
  suffix: '',
  // Progress bar display options
  min: 0,
  max: 100,
  showValue: true,
  showRelative: true,
  showProgress: false,
  animated: false,
  variant: 'success',
  thresholds: [],
})

export class ModuleFieldNumber extends ModuleField {
  readonly kind = kind

  options: NumberOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldNumber>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<NumberOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'format', 'prefix', 'suffix', 'multiDelimiter', 'display', 'variant', 'presetFormat')
    Apply(this.options, o, Number, 'precision', 'min', 'max')
    Apply(this.options, o, Boolean, 'showValue', 'showRelative', 'showProgress', 'animated')

    if (o.thresholds) {
      this.options.thresholds = o.thresholds
    }
  }

  formatValue (value: string): string {
    const o = this.options
    let n: number

    switch (typeof value) {
      case 'string':
        n = parseFloat(value)
        break
      case 'number':
        n = value
        break
      default:
        n = 0
    }

    let out = `${n}`

    if (o.format && o.format.length > 0) {
      out = numeral(n).format(o.format)
    } else {
      out = fmt.number(n)
    }

    if (o.presetFormat === 'accountingNumber') {
      out = formatChartValueAsAccountingNumber(Number(n))
    }

    return '' + o.prefix + (out || n) + o.suffix
  }
}

export function formatChartValueAsAccountingNumber (value: number): string {
  let result = ''

  if (value < 0) {
    result = `(${Math.abs(value)})`
  } else if (value === 0) {
    result = '-'
  }

  return result
}

Registry.set(kind, ModuleFieldNumber)

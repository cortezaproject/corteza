import numeral from 'numeral'
import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'
import * as fmt from '../../../formatting'

const kind = 'Number'

interface NumberOptions extends Options {
  format: string;
  prefix: string;
  suffix: string;
  precision: number;
  multiDelimiter: string;
}

const defaults = (): Readonly<NumberOptions> => Object.freeze({
  ...defaultOptions(),
  format: '',
  prefix: '',
  suffix: '',
  precision: 3,
  multiDelimiter: '\n',
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

    Apply(this.options, o, String, 'format', 'prefix', 'suffix', 'multiDelimiter')
    Apply(this.options, o, Number, 'precision')
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

    return '' + o.prefix + out + o.suffix
  }
}

Registry.set(kind, ModuleFieldNumber)

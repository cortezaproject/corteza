// @todo option to allow multiple entries
// @todo option to allow duplicates
import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'
import { AreStrings } from '../../../guards'

const kind = 'Select'

interface SelectOptions extends Options {
  options: Array<string | { value: string; text?: string }>;
  selectType: string;
  multiDelimiter: string;
}

const defaults = (): Readonly<SelectOptions> => Object.freeze({
  ...defaultOptions(),
  options: [],
  selectType: 'default',
  multiDelimiter: '\n',
})

export class ModuleFieldSelect extends ModuleField {
  readonly kind = kind

  options: SelectOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldSelect>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<SelectOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'selectType', 'multiDelimiter')

    if (o.options) {
      let opt: Array<{ value: string; text: string }> = []

      if (AreStrings(o.options)) {
        opt = o.options
          .map(value => ({ value, text: value }))
      } else {
        opt = (o.options as Array<{ value: string; text?: string }>)
          .map(({ value, text }) => ({ value, text: text || value }))
      }

      this.options.options = opt
    }
  }
}

Registry.set(kind, ModuleFieldSelect)

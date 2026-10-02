import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'
import { AreStrings } from '../../../guards'

const kind = 'Select'

interface SelectOptionStyle {
  textColor?: string
  backgroundColor?: string
}

interface SelectOption {
  value: string;
  text: string;
  style: SelectOptionStyle;
}

interface SelectOptions extends Options {
  options: Array<SelectOption>;
  selectType: string;
  displayType: 'text' | 'badge';
  multiDelimiter: string;
  isUniqueMultiValue: boolean;
}

const defaults = (): Readonly<SelectOptions> => Object.freeze({
  ...defaultOptions(),
  options: [],
  selectType: 'default',
  multiDelimiter: '\n',
  isUniqueMultiValue: false,
  displayType: 'text',
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

    Apply(this.options, o, String, 'selectType', 'multiDelimiter', 'displayType')
    Apply(this.options, o, Boolean, 'isUniqueMultiValue')

    if (o.options) {
      let opt: Array<SelectOption> = []

      if (AreStrings(o.options)) {
        opt = o.options.map((value: string) => this.createSelectOption({ value, text: value }))
      } else {
        opt = o.options.map(o => this.createSelectOption(o))
      }

      this.options.options = opt
    }
  }

  createSelectOption ({ value = '', text = '', style = {} }: Partial<SelectOption> = {}): SelectOption {
    const { textColor = '', backgroundColor = '' } = style || {}
    return {
      value,
      text: text || value,
      style: {
        textColor,
        backgroundColor,
      },
    }
  }
}

Registry.set(kind, ModuleFieldSelect)

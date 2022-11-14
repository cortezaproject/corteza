import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'

const kind = 'String'

interface StringOptions extends Options {
  multiLine: boolean;
  useRichTextEditor: boolean;
  multiDelimiter: string;
}

const defaults = (): Readonly<StringOptions> => Object.freeze({
  ...defaultOptions(),
  multiLine: false,
  useRichTextEditor: false,
  multiDelimiter: '\n',
})

export class ModuleFieldString extends ModuleField {
  readonly kind = kind

  options: StringOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldString>) {
    super(i)

    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<StringOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'multiDelimiter')
    Apply(this.options, o, Boolean, 'multiLine', 'useRichTextEditor')
  }
}

Registry.set(kind, ModuleFieldString)

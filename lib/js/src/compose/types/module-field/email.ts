// @todo option to allow multiple entries
// @todo option to allow duplicates
// @todo option to allow only whitelisted domains
import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'

const kind = 'Email'

interface EmailOptions extends Options {
  outputPlain: boolean;
  multiDelimiter: string;
}

const defaults = (): Readonly<EmailOptions> => Object.freeze({
  ...defaultOptions(),
  outputPlain: true,
  multiDelimiter: '\n',
})

export class ModuleFieldEmail extends ModuleField {
  readonly kind = kind

  options: EmailOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldEmail>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<EmailOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'multiDelimiter')
    Apply(this.options, o, Boolean, 'outputPlain')
  }
}

Registry.set(kind, ModuleFieldEmail)

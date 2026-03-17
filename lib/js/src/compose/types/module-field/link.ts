import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'

const kind = 'Link'

interface LinkOptions extends Options {
  dynamicMode: boolean;
  linkType: string;
  customProtocol: string;
  tempProtocol: string;
  trimFragment: boolean;
  trimQuery: boolean;
  trimPath: boolean;
  onlySecure: boolean;
  outputPlain: boolean;
  multiDelimiter: string;
}

const defaults = (): Readonly<LinkOptions> => Object.freeze({
  ...defaultOptions(),
  dynamicMode: true,
  linkType: '',
  customProtocol: '',
  tempProtocol: '',
  trimFragment: false,
  trimQuery: false,
  trimPath: false,
  onlySecure: false,
  outputPlain: false,
  multiDelimiter: '\n',
})

export class ModuleFieldLink extends ModuleField {
  readonly kind = kind

  options: LinkOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldLink>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<LinkOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'multiDelimiter', 'linkType', 'customProtocol', 'tempProtocol')
    Apply(this.options, o, Boolean, 'dynamicMode', 'trimFragment', 'trimQuery', 'trimPath', 'onlySecure', 'outputPlain')
  }
}

Registry.set(kind, ModuleFieldLink)
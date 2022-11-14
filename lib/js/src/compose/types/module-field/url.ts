// @todo option to allow multiple entries
// @todo option to allow duplicates

import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'

const kind = 'Url'

interface UrlOptions extends Options {
  trimFragment: boolean;
  trimQuery: boolean;
  trimPath: boolean;
  onlySecure: boolean;
  outputPlain: boolean;
  multiDelimiter: string;
}

const defaults = (): Readonly<UrlOptions> => Object.freeze({
  ...defaultOptions(),
  trimFragment: false,
  trimQuery: false,
  trimPath: false,
  onlySecure: false,
  outputPlain: false,
  multiDelimiter: '\n',
})

export class ModuleFieldUrl extends ModuleField {
  readonly kind = kind

  options: UrlOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldUrl>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<UrlOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'multiDelimiter')
    Apply(this.options, o, Boolean, 'trimFragment', 'trimQuery', 'trimPath', 'onlySecure', 'outputPlain')
  }
}

Registry.set(kind, ModuleFieldUrl)

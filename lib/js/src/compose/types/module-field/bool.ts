import { Capabilities, ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'

const kind = 'Bool'

interface BoolOptions extends Options {
  trueLabel: string;
  falseLabel: string;
}

const defaults = (): Readonly<BoolOptions> => Object.freeze({
  ...defaultOptions(),
  trueLabel: '',
  falseLabel: '',
})

export class ModuleFieldBool extends ModuleField {
  readonly kind = kind

  options: BoolOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldBool>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<BoolOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'trueLabel', 'falseLabel')
  }

  /**
   * Per module field type capabilities
   */
  public get cap (): Readonly<Capabilities> {
    return {
      ...super.cap,
      multi: false,
    }
  }
}

Registry.set(kind, ModuleFieldBool)

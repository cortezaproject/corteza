import { Capabilities, ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'

const kind = 'Bool'

interface BoolOptions extends Options {
  trueLabel: string
  falseLabel: string
  switch: boolean
}

const defaults = (): Readonly<BoolOptions> =>
  Object.freeze({
    ...defaultOptions(),
    trueLabel: '',
    falseLabel: '',
    switch: false,
  })

export class ModuleFieldBool extends ModuleField {
  readonly kind = kind

  options: BoolOptions = { ...defaults() }

  constructor(i?: Partial<ModuleFieldBool>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions(o?: Partial<BoolOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'trueLabel', 'falseLabel')
    Apply(this.options, o, Boolean, 'switch')
  }

  /**
   * Per module field type capabilities
   */
  public get cap(): Readonly<Capabilities> {
    return {
      ...super.cap,
      multi: false,

      // A checkbox is ticked or it is not; it has no unanswered state for
      // required to demand an answer to. Marking one required only ever meant
      // "must be ticked", which a validator expression says out loud.
      required: false,
    }
  }
}

Registry.set(kind, ModuleFieldBool)

import { CortezaTypes } from './args-corteza'
import { Caster } from './shared'

/**
 * Handles arguments, passed to the script
 *
 * By convention variables holding "current" resources are prefixed with dollar ($) sign.
 * For example, before/after triggers for record will call registered scripts with $record, $module
 * and $namespace, holding current record, it's module and namespace.
 *
 * All these variables are casted (if passed as an argument) to proper types ($record => Record, $module => Module, ...)
 */
export class Args {
  constructor (args: {[_: string]: unknown}, caster: Caster = CortezaTypes) {
    const cachedArgs: { [_: string]: any } = {}

    for (const arg in args) {
      if (caster && caster.has(arg)) {
        const cast = caster.get(arg)

        if (cast) {
          Object.defineProperty(this, `$${arg}`, {
            get: () => {
              if (!cachedArgs[arg]) {
                cachedArgs[arg] = cast.call(this, args[arg])
              }

              return cachedArgs[arg]
            },
            configurable: false,
            enumerable: true,
          })
        }

        Object.defineProperty(this, `raw${arg.substring(0, 1).toUpperCase()}${arg.substring(1)}`, {
          value: args[arg],
          writable: false,
          configurable: false,
          enumerable: true,
        })
      } else {
        Object.defineProperty(this, arg, {
          value: args[arg],
          writable: false,
          configurable: false,
          enumerable: true,
        })
      }
    }
  }
}

/**
 * Handles arguments, passed to the script but preserves references to the original objects
 *
 * By convention variables holding "current" resources are prefixed with dollar ($) sign.
 * For example, before/after triggers for record will call registered scripts with $record, $module
 * and $namespace, holding current record, it's module and namespace.
 *
 * These variables are not additionally casted, since in order to preserve references they should
 * already be in the correct type.
 */
export class ArgsProxy {
  constructor (args: {[_: string]: unknown}, caster: Caster = CortezaTypes) {
    for (const arg in args) {
      // For consistency only prefix args with & and raw that have a defined caster
      if (caster && caster.has(arg)) {
        Object.defineProperty(this, `$${arg}`, {
          get: () => args[arg],
          configurable: false,
          enumerable: true,
        })

        Object.defineProperty(this, `raw${arg.substring(0, 1).toUpperCase()}${arg.substring(1)}`, {
          value: args[arg],
          writable: false,
          configurable: false,
          enumerable: true,
        })
      } else {
        Object.defineProperty(this, arg, {
          value: args[arg],
          writable: false,
          configurable: false,
          enumerable: true,
        })
      }
    }
  }
}

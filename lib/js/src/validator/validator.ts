import { IsOf } from '../guards'
import { merge } from 'lodash'

interface Meta { [key: string]: unknown }

export class ValidatorError {
  /**
   * kind - key used for translation
   */
  readonly kind: string

  /**
   * Plain error message
   */
  readonly message: string

  /**
   * Any additional meta data that can be used to expand (translated) message,
   * or to group, categorize validator results
   */
  readonly meta: Meta = {}

  constructor (i: string | { kind: string; message?: string; meta?: Meta }) {
    if (typeof i === 'string') {
      this.kind = i
      this.message = i
    } else {
      this.kind = i.kind
      this.message = i.message || i.kind
      if (i.meta) {
        this.meta = merge({}, this.meta, i.meta)
      }
    }
  }
}

const ValidatorFalseDefaultError = Object.freeze(new ValidatorError({
  message: 'Internal error',
  kind: 'internal',
}))

interface ValidatorRawResult {
  kind: string;
  message?: string;
  meta?: Meta;
}

interface ValidatorResultGetter {
  get (): ValidatorError[];
}

/**
 * Supporting as much as we can so that we can make script-developer's life as easy as possible
 */
export type ValidatorResult =
  ValidatorResultGetter |

  // raw pojo results { err: 'error' }
  ValidatorRawResult |

  // error(s)
  ValidatorError[] | ValidatorError |

  // when true its converted to ValidatorFalseDefaultError by Validated class
  boolean |

  // boolean false or any of the rest will result in no error
  null | undefined | void

export function NormalizeValidatorResults (...r: ValidatorResult[]): ValidatorError[] {
  const out: ValidatorError[] = []

  r.forEach(r => {
    if (r === undefined || r === null || r === true) {
      // considering these validation results as valid
      return
    }

    if (Array.isArray(r)) {
      // Expand & normalize each item...
      out.push(...(NormalizeValidatorResults(...(r as ValidatorError[]))))
      return
    }

    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    if (r instanceof Validated) {
      out.push(...r.get())
      return
    }

    if (r instanceof ValidatorError) {
      out.push(r)
      return
    }

    if (IsOf<ValidatorRawResult>(r, 'kind') && r.kind) {
      out.push(new ValidatorError(r))
      return
    }

    // Catch-all for non object errors
    out.push(ValidatorFalseDefaultError)
  })

  return out
}

interface FilterValidatedFn {
  (w: ValidatorError): boolean;
}

/**
 * Holds an manipulates set of errors
 */
export class Validated {
  protected set: ValidatorError[] = []

  constructor (...r: ValidatorResult[]) {
    this.push(...r)
  }

  public get (): ValidatorError[] {
    return this.set
  }

  public get length (): number {
    return this.set.length
  }

  public valid (): boolean {
    return this.length === 0
  }

  public push (...r: ValidatorResult[]): void {
    this.set.push(...NormalizeValidatorResults(...r))
  }

  public applyMeta (meta: Meta): void {
    this.set = this.set.map(r => {
      const appliedMeta = { ...r, meta: { ...r.meta, ...meta } }

      if (r instanceof ValidatorError) {
        return new ValidatorError(appliedMeta)
      }

      return new ValidatorError(appliedMeta)
    })
  }

  public filter (fn: FilterValidatedFn): Validated {
    return new Validated(this.set.filter(fn))
  }

  /**
   * Filters by meta keys
   *
   * If only key is given it returns entries that have meta with that key
   *
   * @param {string} key
   * @param {unknown} value
   */
  public filterByMeta (key: string, value?: unknown): Validated {
    return this.filter(
      (err) => (value === undefined ? err.meta[key] !== undefined : err.meta[key] === value),
    )
  }
}

export interface ValidatorFn<T> {
  (this: T, ...args: unknown[]): ValidatorResult;
}

export function IsEmpty (v: unknown): boolean {
  if (!v || (IsOf(v, 'length') && v.length && v.length === 0)) {
    return true
  }

  if (Array.isArray(v)) {
    return v.every(i => !i)
  }

  return !v
}

/**
 * Checks if values are equal
 * @param {string|string[]} v1 Value in question
 * @param {string|string[]} v2 Value to compare to
 * @returns {boolean}
 */
export function AreEqual (v1: string|string[], v2: string|string[]): boolean {
  if (Array.isArray(v1)) {
    if (!Array.isArray(v2) || v1.length !== v2.length) {
      return false
    }

    return !!v1.find((v, i) => v !== v2[i])
  } else {
    return v1 === v2
  }
}

/**
 * Validator is record validation tool that registers and runs record & field validators
 *
 * Record and field validators are functions that
 */
export class Validator<T> {
  /**
   * Validators
   */
  protected registered: ValidatorFn<T>[] = []

  constructor (...vfn: ValidatorFn<T>[]) {
    if (vfn) {
      this.registered.push(...vfn)
    }
  }

  public push (...vfn: ValidatorFn<T>[]): void {
    this.registered.push(...vfn)
  }

  public run (target: T, ...args: unknown[]): Validated {
    return new Validated(...this.registered.map(vfn => vfn.call(target, ...args)))
  }
}

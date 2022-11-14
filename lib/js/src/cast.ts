// const iso8601check = /^([\\+-]?\d{4}(?!\d{2}\b))((-?)((0[1-9]|1[0-2])(\3([12]\d|0[1-9]|3[01]))?|W([0-4]\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24:?00)([.,]\d+(?!:))?)?(\17[0-5]\d([.,]\d+)?)?([zZ]|([\\+-])([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$/

const uint64zeropad = '00000000000000000000'

/**
 * Reasons behind this small snippet:
 *  - backend is using uint64 as prefered type for handling CortezaID (of all things)
 *  - JavaScript can not (without external help) deal with uint64
 *  - Backend's JSON marshaller converts uint64 to string
 *  - Backend's JSON unmarshaller raises error when given anything but string with a number inside
 *
 *  Until this last thing is fixed or has a proper workaround, we're stuck with this.
 */
export const NoID = '0'

export function ISO8601Date (ts?: unknown): Date|undefined {
  if (ts instanceof Date) {
    return ts
  }

  if (typeof ts === 'string' || typeof ts === 'number') {
    return new Date(ts)
  }

  return undefined
}

interface Caster<T> {
  (input: unknown): T;
}

/**
 * Is native class?
 *
 * @param thing
 * @returns {boolean}
 */
// function isNativeClass (o: unknown): boolean {
//   return typeof o === 'function' &&
//     Object.prototype.hasOwnProperty.call(o, 'prototype') &&
//     !Object.prototype.hasOwnProperty.call(o, 'arguments')
// }

/**
 * Casts value to <type> or returns default
 */
export function PropCast<T> (type: Caster<T>, o: {[_: string]: unknown}|undefined, prop: string): T|undefined {
  if (o === undefined || Object.prototype.hasOwnProperty.call(o, prop)) {
    return undefined
  }

  return type(o.prop)
}

/**
 * Tests if a given value looks like corteza ID
 * @param ID
 * @constructor
 */
export function IsCortezaID (ID: unknown): boolean {
  if (typeof ID !== 'string') {
    return false
  }

  if (!/^\d+$/.test(ID)) {
    return false
  }

  return true
}

/**
 * @return {string}
 */
export function CortezaID (value: unknown): string {
  if (!value) {
    return NoID
  }

  if (typeof value === 'number') {
    return value.toString()
  }

  if (IsCortezaID(value)) {
    return value as string
  }

  throw new Error('Invalid CortezaID value')
}

/**
 * Apply caster interface that satisfies basic casting functions + String, Number etc...
 */
interface ApplyCaster {
  (val: unknown): unknown;
}

/**
 * Apply takes all given props, their values (from src) and assignes them to props (on dst)
 *
 * A casting function can be used (see ApplyCaster) to modify the values before assigning them
 */
export function Apply<DST, SRC, T extends keyof DST> (dst: DST, src: SRC, cast: ApplyCaster|keyof DST, ...props: (keyof DST)[]): void {
  if (typeof cast !== 'function') {
    // Handle case where we do not use caster
    props.unshift(cast)

    // and use String as a caster, effectively forcing
    // cast-to-string on all applied values to
    cast = String
  }

  if (typeof src !== 'object') {
    return
  }

  props.forEach(prop => {
    // prop must exist on dst
    if (!Object.prototype.hasOwnProperty.call(dst, prop)) {
      return
    }

    // prop must exist on src
    if (!Object.prototype.hasOwnProperty.call(src, prop)) {
      return
    }

    // sProp is prop from source
    const sProp = (prop as unknown) as keyof SRC

    // value on src should be defined
    if (src[sProp] === undefined || src[sProp] === null) {
      return
    }

    // Cast value from src to type of value from prop on dst
    let val = (src[sProp] as unknown) as DST[T]

    // Run value through cast fn
    val = ((cast as ApplyCaster)(val) as unknown) as DST[T]
    if (val === undefined) {
      return
    }

    // Assign (valid ony) value to dst
    dst[prop] = val
  })
}

export function ApplyWhitelisted<DST, SRC, WL, T extends keyof DST> (dst: DST, src: SRC, whitelist: (DST[T])[], ...props: (keyof DST)[]): void {
  if (typeof src !== 'object') {
    return
  }

  props.forEach(prop => {
    // prop must exist on dst
    if (!Object.prototype.hasOwnProperty.call(dst, prop)) {
      return
    }

    // prop must exist on src
    if (!Object.prototype.hasOwnProperty.call(src, prop)) {
      return
    }

    // sProp is prop from source
    const sProp = (prop as unknown) as keyof SRC

    // value on src should be defined
    if (src[sProp] === undefined) {
      return
    }

    // Cast value from src to type of value from prop on dst
    const val = (src[sProp] as unknown) as DST[T]

    if (whitelist.includes(val)) {
      dst[prop] = val
    }
  })
}

export function makeIDSortable (ID?: string): string {
  // We're using uint64 for CortezaID and JavaScript does not know how to handle this type
  // natively. We get the value from backend as string anyway and we need to prefix
  // it with '0' to ensure string sorting does what we need it to.
  return uint64zeropad.substr((ID || '').length) + (ID || '')
}

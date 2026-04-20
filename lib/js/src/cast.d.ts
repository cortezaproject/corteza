/**
 * Reasons behind this small snippet:
 *  - backend is using uint64 as prefered type for handling HumanID (of all things)
 *  - JavaScript can not (without external help) deal with uint64
 *  - Backend's JSON marshaller converts uint64 to string
 *  - Backend's JSON unmarshaller raises error when given anything but string with a number inside
 *
 *  Until this last thing is fixed or has a proper workaround, we're stuck with this.
 */
export declare const NoID = '0'
export declare function ISO8601Date(ts?: unknown): Date | undefined
interface Caster<T> {
  (input: unknown): T
}
/**
 * Is native class?
 *
 * @param thing
 * @returns {boolean}
 */
/**
 * Casts value to <type> or returns default
 */
export declare function PropCast<T>(
  type: Caster<T>,
  o:
    | {
        [_: string]: unknown
      }
    | undefined,
  prop: string,
): T | undefined
/**
 * Tests if a given value looks like a human ID
 * @param ID
 * @constructor
 */
export declare function IsHumanID(ID: unknown): boolean
/**
 * @return {string}
 */
export declare function HumanID(value: unknown): string
/**
 * Apply caster interface that satisfies basic casting functions + String, Number etc...
 */
interface ApplyCaster {
  (val: unknown): unknown
}
/**
 * Apply takes all given props, their values (from src) and assignes them to props (on dst)
 *
 * A casting function can be used (see ApplyCaster) to modify the values before assigning them
 */
export declare function Apply<DST, SRC, _T extends keyof DST>(
  dst: DST,
  src: SRC,
  cast: ApplyCaster | keyof DST,
  ...props: (keyof DST)[]
): void
export declare function ApplyWhitelisted<DST, SRC, _WL, T extends keyof DST>(
  dst: DST,
  src: SRC,
  whitelist: DST[T][],
  ...props: (keyof DST)[]
): void
export declare function makeIDSortable(ID?: string): string
export {}

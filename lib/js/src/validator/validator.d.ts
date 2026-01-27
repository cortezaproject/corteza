interface Meta {
    [key: string]: unknown;
}
export declare class ValidatorError {
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
  readonly meta: Meta
  constructor(i: string | {
        kind: string;
        message?: string;
        meta?: Meta;
    });
}
interface ValidatorRawResult {
    kind: string;
    message?: string;
    meta?: Meta;
}
interface ValidatorResultGetter {
    get(): ValidatorError[];
}
/**
 * Supporting as much as we can so that we can make script-developer's life as easy as possible
 */
export type ValidatorResult = ValidatorResultGetter | ValidatorRawResult | ValidatorError[] | ValidatorError | boolean | null | undefined | void;
export declare function NormalizeValidatorResults(...r: ValidatorResult[]): ValidatorError[];
interface FilterValidatedFn {
    (w: ValidatorError): boolean;
}
/**
 * Holds an manipulates set of errors
 */
export declare class Validated {
  protected set: ValidatorError[];
  constructor(...r: ValidatorResult[]);
  get(): ValidatorError[];
  get length(): number;
  valid(): boolean;
  push(...r: ValidatorResult[]): void;
  applyMeta(meta: Meta): void;
  filter(fn: FilterValidatedFn): Validated;
  /**
     * Filters by meta keys
     *
     * If only key is given it returns entries that have meta with that key
     *
     * @param {string} key
     * @param {unknown} value
     */
  filterByMeta(key: string, value?: unknown): Validated;
}
export interface ValidatorFn<T> {
    (this: T, ...args: unknown[]): ValidatorResult;
}
export declare function IsEmpty(v: unknown): boolean;
/**
 * Checks if values are equal
 * @param {string|string[]} v1 Value in question
 * @param {string|string[]} v2 Value to compare to
 * @returns {boolean}
 */
export declare function AreEqual(v1: string | string[], v2: string | string[]): boolean;
/**
 * Validator is record validation tool that registers and runs record & field validators
 *
 * Record and field validators are functions that
 */
export declare class Validator<T> {
  /**
     * Validators
     */
  protected registered: ValidatorFn<T>[]
  constructor(...vfn: ValidatorFn<T>[]);
  push(...vfn: ValidatorFn<T>[]): void;
  run(target: T, ...args: unknown[]): Validated;
}
export {}

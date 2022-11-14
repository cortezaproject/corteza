import { User } from '../system'

interface GenericCtor<T> {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  new (...args: any[]): T;
}

/**
 * Generic type caster
 *
 * Takes argument (ref to class) and returns a function that will initialize class of that type
 */
export function GenericCaster<T> (C: GenericCtor<T>): GenericGetterFn<T|undefined> {
  return function (val: unknown): T|undefined {
    if (!val || typeof val !== 'object') {
      return undefined
    }

    return new C(val)
  }
}

/**
 * Generic type caster with Object.freeze
 *
 * Takes argument (ref to class) and returns a function that will initialize class of that type
 */
export function GenericCasterFreezer<T> (C: GenericCtor<T>): GenericGetterFn<Readonly<T>|undefined> {
  return function (val: unknown): Readonly<T>|undefined {
    if (!val || typeof val !== 'object') {
      return undefined
    }

    return Object.freeze(new C(val))
  }
}

export interface BaseArgs {
  $invoker: User;
  authToken: string;
}

export interface GenericGetterFn<T> {
  (val: unknown): T;
}

export interface GetterFn {
  (key: unknown): unknown;
}

export type Caster = Map<string, GetterFn>

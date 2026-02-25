import { User } from '../system'
interface GenericCtor<T> {
  new (...args: any[]): T
}
/**
 * Generic type caster
 *
 * Takes argument (ref to class) and returns a function that will initialize class of that type
 */
export declare function GenericCaster<T>(C: GenericCtor<T>): GenericGetterFn<T | undefined>
/**
 * Generic type caster with Object.freeze
 *
 * Takes argument (ref to class) and returns a function that will initialize class of that type
 */
export declare function GenericCasterFreezer<T>(
  C: GenericCtor<T>,
): GenericGetterFn<Readonly<T> | undefined>
export interface BaseArgs {
  $invoker: User
  authToken: string
}
export interface GenericGetterFn<T> {
  (val: unknown): T
}
export interface GetterFn {
  (key: unknown): unknown
}
export type Caster = Map<string, GetterFn>
export {}

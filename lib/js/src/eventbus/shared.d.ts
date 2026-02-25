import { ConstraintMatcher } from './constraints'
export interface Constraint {
  name?: string
  op?: string
  value: string[]
}
export declare const onManual = 'onManual'
export interface HandlerFn {
  (ev: Event): Promise<unknown>
}
export interface Trigger {
  eventTypes: string[]
  resourceTypes: string[]
  weight?: number
  constraints?: Constraint[]
  scriptName?: string
}
interface SortableScript {
  weight: number
}
export declare function scriptSorter(a: SortableScript, b: SortableScript): number
interface EventMatcher {
  (c: ConstraintMatcher): boolean
}
interface EventArgs {
  [_: string]: unknown
}
export interface Event {
  resourceType: string
  eventType: string
  match?: EventMatcher
  args?: EventArgs
}
interface ResourceTypeGetter {
  resourceType: string
}
export declare function GenericEventMaker<T extends ResourceTypeGetter>(
  t: T,
  eventType: string,
  match: EventMatcher,
  args: EventArgs,
): Event
export {}

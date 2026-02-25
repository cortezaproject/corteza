import { ConstraintMatcher } from './constraints'
import { Event, HandlerFn, Trigger } from './shared'
export declare function DummyHandler(): Promise<undefined>
export declare class Handler {
  readonly resourceTypes: string[]
  readonly eventTypes: string[]
  readonly constraints: ConstraintMatcher[]
  readonly weight: number
  readonly handle: HandlerFn
  readonly scriptName?: string
  constructor(h: HandlerFn, t: Trigger)
  /**
   * Match this handler with a given event - type, resource, constraints + scriptName when ManualEvent
   *
   * @param {Event} ev
   * @return bool
   */
  Match(ev: Event, script?: string): boolean
  Handle(ev: Event): Promise<unknown>
}

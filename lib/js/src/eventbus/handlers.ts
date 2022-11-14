import { ConstraintMaker, ConstraintMatcher } from './constraints'
import { Event, HandlerFn, onManual, Trigger } from './shared'

// Dummy handler, can be used for tests
export async function DummyHandler (): Promise<undefined> { return undefined }

export class Handler {
  readonly resourceTypes: string[];
  readonly eventTypes: string[];
  readonly constraints: ConstraintMatcher[];
  readonly weight: number;
  readonly handle: HandlerFn;
  readonly scriptName?: string

  constructor (h: HandlerFn, t: Trigger) {
    this.handle = h
    this.eventTypes = t.eventTypes
    this.resourceTypes = t.resourceTypes
    this.weight = t.weight || 0
    // @todo parse constraints to constraint matchers
    this.constraints = t.constraints ? t.constraints.map(ConstraintMaker) : []
    this.scriptName = t.scriptName
  }

  /**
   * Match this handler with a given event - type, resource, constraints + scriptName when ManualEvent
   *
   * @param {Event} ev
   * @return bool
   */
  Match (ev: Event, script?: string): boolean {
    if (!this.eventTypes.includes(ev.eventType)) {
      return false
    }

    if (!this.resourceTypes.includes(ev.resourceType)) {
      return false
    }

    if (script && this.scriptName !== script) {
      return false
    }

    if (ev.match) {
      // Event should match all trigger's constraints
      for (const c of this.constraints) {
        if (!ev.match(c)) {
          return false
        }
      }
    }

    return true
  }

  Handle (ev: Event): Promise<unknown> {
    return this.handle(ev)
  }
}

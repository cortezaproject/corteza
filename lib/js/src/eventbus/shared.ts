import { ConstraintMatcher } from './constraints'

export interface Constraint {
  name?:
      string;
  op?:
      string;
  value:
      string[];
}

export const onManual = 'onManual'

export interface HandlerFn {
  (ev: Event): Promise<unknown>;
}

export interface Trigger {
  eventTypes: string[];
  resourceTypes: string[];
  weight?: number;
  constraints?: Constraint[];
  scriptName?: string;
}

interface SortableScript {
  weight: number;
}

export function scriptSorter (a: SortableScript, b: SortableScript): number {
  return a.weight - b.weight
}

interface EventMatcher {
  (c: ConstraintMatcher): boolean;
}

interface EventArgs { [_: string]: unknown }

export interface Event {
  resourceType: string;
  eventType: string;
  match?: EventMatcher;
  args?: EventArgs;
}

interface ResourceTypeGetter { resourceType: string }

export function GenericEventMaker<T extends ResourceTypeGetter> (t: T, eventType: string, match: EventMatcher, args: EventArgs): Event {
  return {
    resourceType: t.resourceType,
    eventType,
    match,
    args,
  }
}

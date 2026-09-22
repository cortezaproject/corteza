import { Constraint } from '@planetcrust/human-js/src/eventbus/shared'
import { UIProp } from './types'

const defaultResource = 'system'

function isIterable<T>(o: unknown): o is Iterable<T> {
  return (
    typeof o === 'object' && typeof (o as Record<symbol, unknown>)[Symbol.iterator] === 'function'
  )
}

function distinct(arr: string[]): string[] {
  return arr.filter((n, i) => arr.indexOf(n) === i)
}

// prefixes and captializes list of given event names
function eventize(prefix: string, ee: string[]): string[] {
  return ee.map(e => prefix + e.substring(0, 1).toUpperCase() + e.substring(1))
}

export class Trigger {
  readonly eventTypes: string[] = []

  readonly resourceTypes: string[] = []

  readonly constraints: Constraint[] = []

  readonly uiProps: UIProp[] = []

  constructor(t?: Partial<Trigger>) {
    if (!t) {
      return
    }

    this.eventTypes = t.eventTypes || []
    this.resourceTypes = t.resourceTypes || []
    this.constraints = t.constraints || []
    this.uiProps = t.uiProps || []
  }

  on(...events: string[]): Trigger {
    const t = this ?? new Trigger()

    return new Trigger({
      ...t,
      eventTypes: distinct([...t.eventTypes, ...eventize('on', events)]),
    })
  }

  before(...events: string[]): Trigger {
    const t = this ?? new Trigger()

    return new Trigger({
      ...t,
      eventTypes: distinct([...t.eventTypes, ...eventize('before', events)]),
    })
  }

  after(...events: string[]): Trigger {
    const t = this ?? new Trigger()

    return new Trigger({
      ...t,
      eventTypes: distinct([...t.eventTypes, ...eventize('after', events)]),
    })
  }

  for(...resources: string[]): Trigger {
    const t = this ?? new Trigger()

    return new Trigger({
      ...t,
      resourceTypes: distinct([...t.resourceTypes, ...resources]),
    })
  }

  at(...timestamps: string[]): Trigger {
    return deferred('timestamp', ...timestamps)
  }

  every(...intervals: string[]): Trigger {
    return deferred('interval', ...intervals)
  }

  where(...aa: unknown[]): Trigger {
    let name: string | undefined
    let op: string | undefined
    let value: string[]

    if (aa.length > 1) {
      // More than 1 argument, assuming first is constraint's name
      name = aa.shift() as string
    }

    if (aa.length > 1) {
      // More than 2 arguments, assuming 2nd is operator
      op = aa.shift() as string
    }

    // There rest of arguments are values
    if (aa.length === 1 && Array.isArray(aa[0])) {
      value = aa[0]
    } else if (aa.length === 1) {
      value = [aa[0] as string]
    } else {
      value = aa as string[]
    }

    const { constraints } = this
    constraints.push({ name, op, value })

    return new Trigger({ ...this, constraints })
  }

  uiProp(name: string, value: string): Trigger {
    const t = this ?? new Trigger()
    return new Trigger({ ...t, uiProps: [...t.uiProps, { name, value }] })
  }
}

/**
 * Makes triggers out of an input params
 */
export default function Make(t: unknown): Trigger[] {
  let tt = []

  if (typeof t === 'function') {
    // Execute trigger callback to convert to array of triggers
    // and overwrite the function with definition
    t = t(new Trigger())

    if (isIterable(t)) {
      t = [...t]
    }
  }

  if (Array.isArray(t)) {
    tt = [...t]
  } else {
    tt = [t]
  }

  return tt.filter(t => t instanceof Trigger)
}

function deferred(name: string, ...value: string[]): Trigger {
  return new Trigger({
    resourceTypes: [defaultResource],
    eventTypes: eventize('on', [name]),
    constraints: [{ name, value }],
  })
}

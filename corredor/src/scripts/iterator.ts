export interface IteratorFilter {
  query: string
  sort: string
  limit: number | string
  offset: number | string

  // Any additional, non-standard filter params
  // like ns and module for record
  [_: string]: number | string
}

export class Iterator {
  public resourceType!: string
  public eventType: 'onManual' | 'onInterval' | 'onTimestamp' = 'onManual'
  public action: 'update' | 'delete' | 'clone' | '' = ''
  public filter!: IteratorFilter
  public deferred: Array<string> = []

  constructor(o: Partial<Iterator>) {
    Object.assign(this, o)
  }

  every(...value: string[]): Iterator {
    this.eventType = 'onInterval'
    this.deferred = [...value]
    return this
  }

  at(...value: string[]): Iterator {
    this.eventType = 'onTimestamp'
    this.deferred = [...value]
    return this
  }
}

/**
 * Iterator construction helper
 *
 * Allows user to construct resource iterator from resource type, action and filter
 */
export function each(o: Partial<Iterator>): Iterator {
  return new Iterator(o)
}

export default function Make(i: unknown): Iterator | undefined {
  if (typeof i === 'function') {
    return i(each)
  } else if (i instanceof Iterator) {
    return i
  }

  return undefined
}

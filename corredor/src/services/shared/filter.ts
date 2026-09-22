import { Script } from '../../types'

interface ListFilter {
  query?: string
  resourceType?: string
  eventTypes?: string[]
}

interface TriggerFilterArgs {
  resourceTypes: string[]
  eventTypes: string[]
}

/**
 * Creates Array.filter() compatible function for filtering script array
 */
export default function (f: ListFilter): { (item: Script): boolean } {
  return (item: Script): boolean => {
    if (f === undefined) {
      // Match all when no filter
      return true
    }

    if (f.resourceType || f.eventTypes) {
      const tt = (item.triggers || []).filter(
        ({ resourceTypes, eventTypes }: TriggerFilterArgs) => {
          if (f.resourceType && f.resourceType.length > 0) {
            // Filter by resource
            if (!resourceTypes || !resourceTypes.includes(f.resourceType)) {
              // No resources found on trigger
              return false
            }
          }

          if (f.eventTypes && f.eventTypes.length > 0) {
            // Filter by events
            if (!eventTypes || f.eventTypes.find(fe => eventTypes.includes(fe)) === undefined) {
              return false
            }
          }

          return true
        },
      )

      if (tt.length === 0) {
        // Filtering by resource and/or events but
        // none of the triggers matched
        return false
      }
    }

    if (f.query) {
      // Strings to search through
      const str = `${item.name} ${item.label} ${item.description}`

      // search query terms
      for (const t of f.query.split(' ')) {
        if (str.includes(t)) {
          return true
        }
      }

      // none matched, fail
      return false
    }

    // No match
    return true
  }
}

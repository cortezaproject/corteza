// The one-at-a-time status filter of a resource list.
//
// A list shows Active rows or the rows in exactly one state. The API speaks in
// per-state filters ('0' without / '1' including / '2' only), so a status is a
// set of those: Active excludes every state, a state is "only" that one and
// excludes the rest — except Deleted, which keeps its rows whatever else they
// carry, so a deleted user still shows when also suspended.

export type ResourceStatus = 'active' | 'deleted' | 'suspended' | 'archived' | 'disabled'

type StateFilter = Record<string, unknown>

/** The status a filter currently shows; anything not "only" one state is Active. */
export function statusOf(filter: StateFilter, states: string[]): ResourceStatus {
  const only = states.find(s => String(filter[s]) === '2')
  return (only as ResourceStatus) ?? 'active'
}

/** The per-state filter values that show one status. */
export function statusFilter(status: string, states: string[]): Record<string, string> {
  return Object.fromEntries(
    states.map(s => {
      if (s === status) return [s, '2']
      return [s, status === 'deleted' ? '1' : '0']
    }),
  )
}

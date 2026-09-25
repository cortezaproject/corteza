// The "last change" column every resource list ends with.
//
// One key, one value, one sort expression: the cell shows the most recent of
// deletedAt/updatedAt/createdAt and `useResourceList` turns the key into the
// matching COALESCE sort, so the column orders by what it displays.

import { locFullDateTime } from '../filters/date'

/** Column key; `useResourceList` maps it to the COALESCE sort expression. */
export const CHANGED_AT_KEY = 'changedAt'

/** Sort expression the server is given for CHANGED_AT_KEY. */
export const CHANGED_AT_SORT = 'coalesce(deletedAt, updatedAt, createdAt)'

// Go's zero time on the wire. A resource the server assembles rather than
// stores — a catalog connection, a federation node — carries it in place of a
// real createdAt, and it is a value, not a date. The clock and offset vary with
// the zone the server rendered it in ('…T00:00:00Z', '…T00:58:04+00:58'), so
// only the date part identifies it.
const ZERO_DATE = '0001-01-01'

// A row is either the wire object, whose timestamps are strings, or a lib/js
// model, whose constructor casts the same fields to Date (ISO8601Date).
function stamped(value: unknown): value is string | Date {
  if (value instanceof Date) {
    return !Number.isNaN(value.getTime()) && value.getUTCFullYear() > 1
  }

  return typeof value === 'string' && value !== '' && !value.startsWith(ZERO_DATE)
}

/** The most recent of a resource's three timestamps. */
export function changedAt(resource: any): string | Date | undefined {
  for (const value of [resource?.deletedAt, resource?.updatedAt, resource?.createdAt]) {
    if (stamped(value)) {
      return value
    }
  }

  return undefined
}

/**
 * The cell's text, empty when the resource carries no timestamp at all — a
 * catalog entry that is not a stored record yet. The date filters go through
 * moment, which reads a missing value as *now*, so an unguarded call renders
 * today's date for a resource that was never written.
 */
export function changedAtText(resource: any): string {
  const value = changedAt(resource)
  return value ? locFullDateTime(value) : ''
}

/** A row's lifecycle state, shown as a tag in the column; the most final wins. */
export type ResourceState = 'deleted' | 'suspended' | 'archived'

export function resourceState(resource: any): ResourceState | undefined {
  if (stamped(resource?.deletedAt)) return 'deleted'
  if (stamped(resource?.suspendedAt)) return 'suspended'
  if (stamped(resource?.archivedAt)) return 'archived'
  return undefined
}

interface ChangedAtField {
  key: string
  header: string
  sortable: boolean
  class: string
  pt: Record<string, unknown>
  [key: string]: unknown
}

/**
 * Field definition for the column. `header` is passed in so the key stays with
 * the app's locale bundle; `overrides` covers the lists that cannot sort it —
 * `changedAtField(t('general.columns.changedAt'), { sortable: false })`.
 */
export function changedAtField(
  header: string,
  overrides: Record<string, unknown> = {},
): ChangedAtField {
  return {
    key: CHANGED_AT_KEY,
    header,
    sortable: true,
    class: 'text-right',
    pt: {
      columnHeaderContent: 'justify-end',
    },
    ...overrides,
  }
}

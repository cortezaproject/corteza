// Shared date-normalization helper for the dashboard stores (events,
// backlogItems). GovernanceForm's date fields bind PrimeVue date pickers,
// which surface date-typed values as JS Date objects. The backend stores
// dates as plain YYYY-MM-DD strings (the report endpoint's overdue metric
// parses them as ISO), so every create/update call normalizes here rather
// than leaving each store to serialize its own payload.

export function toISODate(v) {
  if (!(v instanceof Date)) return v
  if (Number.isNaN(v.getTime())) return ''
  const y = v.getFullYear()
  const m = String(v.getMonth() + 1).padStart(2, '0')
  const d = String(v.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

// Mutates and returns `body`, converting each listed key from a Date (or
// passing through anything else, e.g. an already-ISO string) to YYYY-MM-DD.
// Keys not present on `body` are left untouched.
export function normalizeDates(body, keys) {
  for (const k of keys) {
    if (k in body) body[k] = toISODate(body[k])
  }
  return body
}

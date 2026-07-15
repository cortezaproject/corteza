// Shared chart palette. Hexes mirror the EventBadge tint families (severity,
// status) so charts and the table badges read as one colour system. Grouped by
// the dimension they colour; use colorFor(variant, label) to resolve.

// Lifecycle ramp: colour tracks progress toward done (urgency lives in the
// separate "overdue" metric, not here). Distinct hues, green = done, and Open
// is off-grey so it never reads as the "unknown/blank" fallback.
export const STATUS_COLORS = {
  Open: '#3b82f6', // blue — new, logged
  'In Progress': '#f59e0b', // amber — active, in-flight
  'Ready to Test': '#8b5cf6', // violet — in review
  Completed: '#10b981', // emerald — done
}

export const SEVERITY_COLORS = {
  Critical: '#ef4444',
  Serious: '#d946ef',
  Major: '#f97316',
  Minor: '#f59e0b',
  Informational: '#0ea5e9',
}

export const RISK_COLORS = {
  Critical: '#ef4444',
  High: '#f97316',
  Medium: '#f59e0b',
  Low: '#10b981',
}

// Per-category accent — matches the badge icon colours in config/categories.
export const CATEGORY_COLORS = {
  incident: '#ef4444',
  feature: '#3b82f6',
  privacy: '#a855f7',
  task: '#10b981',
  review: '#f59e0b',
}

// Canonical display order per dimension, so stacked/legend segments read
// worst→best (or lifecycle order) rather than first-seen.
export const STATUS_ORDER = ['Open', 'In Progress', 'Ready to Test', 'Completed']
export const SEVERITY_ORDER = ['Critical', 'Serious', 'Major', 'Minor', 'Informational']
export const RISK_ORDER = ['Critical', 'High', 'Medium', 'Low']

const ORDERS = { status: STATUS_ORDER, severity: SEVERITY_ORDER, risk: RISK_ORDER }

// Sort index for a label within a variant's canonical order; unknown labels
// sort last (stable, alphabetical among themselves).
export function orderIndex(variant, label) {
  const i = (ORDERS[variant] || []).indexOf(String(label ?? ''))
  return i === -1 ? Number.MAX_SAFE_INTEGER : i
}

// Neutral for unknown/blank labels, and the axis/legend text tone (reads on
// both light and dark surfaces).
export const MUTED = '#94a3b8'

// Categorical fallback for open-ended dimensions (type) and unknown labels.
const FALLBACK = ['#6366f1', '#0ea5e9', '#14b8a6', '#8b5cf6', '#f43f5e', '#f59e0b', '#10b981', '#64748b']

const MAPS = { status: STATUS_COLORS, severity: SEVERITY_COLORS, risk: RISK_COLORS }

// Deterministic fallback colour from a string, stable across renders.
function hashColor(label) {
  let h = 0
  for (let i = 0; i < label.length; i++) h = (h * 31 + label.charCodeAt(i)) >>> 0
  return FALLBACK[h % FALLBACK.length]
}

// Resolve a colour for a dimension value. `variant`: status | severity | risk |
// category | type (or anything else → hashed fallback).
export function colorFor(variant, label) {
  const map = MAPS[variant]
  const key = String(label ?? '')
  if (map) return map[key] || MUTED
  if (variant === 'category') return CATEGORY_COLORS[key] || hashColor(key)
  return hashColor(key)
}

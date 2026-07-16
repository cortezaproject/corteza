// Shared chart palette. Hexes mirror the EventBadge tint families (severity,
// status) so charts and the table badges read as one colour system. Grouped by
// the dimension they colour; use colorFor(variant, label, category) to resolve
// (category is only needed for variant 'type' — see below).
//
// Every palette here has been run through the dataviz skill's
// validate_palette.js (categorical checks: lightness band, chroma floor, CVD
// separation, normal-vision floor, contrast) for BOTH `--mode light` and
// `--mode dark` — a single hex per label is used in both themes, so it must
// clear the *intersection* of the two mode bands. Re-run the validator with
// the same command whenever a hex in this file changes:
//   node <dataviz-skill>/scripts/validate_palette.js "<hex,hex,...>" --mode light
//   node <dataviz-skill>/scripts/validate_palette.js "<hex,hex,...>" --mode dark
import { EVENT_FORMS } from './eventForm'

// Lifecycle ramp: colour tracks progress toward done (urgency lives in the
// separate "overdue" metric, not here). Distinct hues, green = done, and Open
// is off-grey so it never reads as the "unknown/blank" fallback.
// Validated both modes (ALL PASS, no contrast relief needed): the older
// amber/emerald steps (#f59e0b, #10b981) overflowed the dark-mode lightness
// cap, so both sit one family step darker.
export const STATUS_COLORS = {
  Open: '#3b82f6', // blue — new, logged
  'In Progress': '#b45309', // amber — active, in-flight
  'Ready to Test': '#8b5cf6', // violet — in review
  Completed: '#059669', // emerald — done
}

// Severity is ordinal (Critical > Serious > Major > Minor > Informational)
// and is treated as a SEQUENTIAL ramp, not a categorical palette: one red
// family, darker = worse, so the worst→least order is readable at a glance
// and charts stay in the same warm family as the EventBadge pills.
// Informational sits outside the scale (it isn't "a bit less bad than
// Minor") and gets a cool blue instead.
// The trade (chosen deliberately over a validator-passing categorical
// zig-zag): adjacent red steps are told apart by darkness, position and
// labels — severity always renders with a direct count label (rank bars) or
// a legend + tooltip (donuts, stacked trend) — not by hue distance. For
// sequential ramps the validation rule is lightness monotonicity (see the
// dataviz skill's color-formula.md), which holds: L 0.45 < 0.53 < 0.62 <
// 0.72, Informational off-scale.
export const SEVERITY_COLORS = {
  Critical: '#991b1b', // darkest — worst
  Serious: '#c32222',
  Major: '#e04343',
  Minor: '#f07f74', // lightest warm step
  Informational: '#449ecd', // cool — outside the badness scale
}

// Reuses the severity ramp's four warm steps 1:1 (Critical/Serious/Major/
// Minor → Critical/High/Medium/Low) so the two ordinal scales read as one
// system.
export const RISK_COLORS = {
  Critical: '#991b1b',
  High: '#c32222',
  Medium: '#e04343',
  Low: '#f07f74',
}

// Audit-event activity accent (pulse charts in Overview + All Events). Kept
// distinct from the category/status palettes above.
export const EVENTS_COLOR = '#6366f1'

// Per-category accent (identity/nominal — order doesn't carry meaning, only
// distinctness does). incident=red, feature=blue and review=amber keep their
// original anchors; task moved to a darker green and review to a darker amber
// so both clear the dark-mode lightness band (were 0.696/0.769, band caps at
// 0.67); privacy's purple (ΔE 0.9 vs feature's blue under deutan — collapses)
// is replaced with a teal/cyan, re-stepped dark enough to clear both bands
// and stay separated from its green/blue neighbours.
// Validated (display order incident,feature,privacy,task,review): all-PASS
// adjacent, light+dark; worst adjacent CVD ΔE 8.7 (protan), worst normal ΔE
// 15.9 (privacy↔feature). Privacy sits at 2.84:1 contrast on light (WARN/
// relief — mitigated by the category's own title/legend, never colour-only).
export const CATEGORY_COLORS = {
  incident: '#ef4444',
  feature: '#3b82f6',
  privacy: '#00a7b5',
  task: '#007835',
  review: '#d97706',
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

// Categorical fallback for the open-ended "type" dimension and any label that
// isn't in a known list. Fixed hue order (blue, green, magenta, olive/yellow,
// teal, orange, violet, red), each re-stepped into the light/dark
// intersection band (OKLCH L 0.48–0.67) with lightness *alternated* between
// neighbours (on top of the hue jump) for margin — the same technique as
// SEVERITY_COLORS. Assigned by list index (see colorFor), so only ADJACENT
// pairs need to clear the floor (bars/donuts, index-stable order) — not
// all-pairs (that's the scatter/bubble/small-multiples case, and per
// color-formula.md an 8-hue set can only clear all-pairs at 4 slots; not
// needed here since no chart scatters type values against each other).
// Validated: worst adjacent CVD ΔE 13.7, worst adjacent normal ΔE 21.3 (both
// light+dark, well clear of the 8/15 targets).
const FALLBACK = ['#007df4', '#007900', '#f5008a', '#865900', '#00a26f', '#a83a00', '#7e6fff', '#bc001d']

const MAPS = { status: STATUS_COLORS, severity: SEVERITY_COLORS, risk: RISK_COLORS }

// The type-select field key per category (mirrors FIELDS in config/
// categories.js) — the single source for "what are this category's type
// options, in order" is the New Event form schema itself (config/eventForm),
// so the type colour assignment can never drift from what the dropdown shows.
const TYPE_FIELD_KEY = {
  incident: 'incidentType',
  feature: 'featureType',
  privacy: 'requestType',
  task: 'taskType',
  review: 'reviewType',
}

// category -> ordered option labels for its type field. Built once at import
// time from EVENT_FORMS; a label's position here is its stable FALLBACK index.
const TYPE_OPTIONS = Object.fromEntries(
  Object.entries(TYPE_FIELD_KEY).map(([category, key]) => {
    const fields = EVENT_FORMS[category]?.[0]?.fields || []
    return [category, fields.find(f => f.key === key)?.options || []]
  }),
)

// Deterministic fallback colour from a string, stable across renders. Last
// resort only — for a type label that isn't in its category's option list
// (e.g. legacy/free-text data), or any other unlisted dimension value.
function hashColor(label) {
  let h = 0
  for (let i = 0; i < label.length; i++) h = (h * 31 + label.charCodeAt(i)) >>> 0
  return FALLBACK[h % FALLBACK.length]
}

// Resolve a colour for a dimension value. `variant`: status | severity | risk |
// category | type (or anything else → hashed fallback). `category` is only
// consulted for variant 'type' — it picks which category's option list to
// index into (see TYPE_OPTIONS above); omit it for every other variant.
export function colorFor(variant, label, category) {
  const key = String(label ?? '')
  if (!key) return MUTED // unknown/blank always stays muted, never hashed

  const map = MAPS[variant]
  if (map) return map[key] || MUTED
  if (variant === 'category') return CATEGORY_COLORS[key] || hashColor(key)
  if (variant === 'type') {
    const options = TYPE_OPTIONS[category] || []
    const idx = options.indexOf(key)
    return idx === -1 ? hashColor(key) : FALLBACK[idx % FALLBACK.length]
  }
  return hashColor(key)
}

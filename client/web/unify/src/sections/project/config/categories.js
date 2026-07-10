// Per-category dashboard config. Drives CategoryView.vue: which columns the
// CResourceList shows, which bar charts render, which KPIs appear, and the form
// schema reused for the "+ New {type}" dialog. Icons/colors are copied from the
// EVENT_CATEGORIES entries (single source of truth for the visual identity);
// the form schemas are reused as-is from EVENT_FORMS.
import { EVENT_CATEGORIES, EVENT_FORMS } from './eventForm'

// i18n prefix for reused per-field labels (shared with the New Event form).
const f = 'project.dashboard.event.f.'

// Fixed display order of the five categories (nav + any iteration).
export const CATEGORY_ORDER = ['incident', 'feature', 'privacy', 'task', 'review']

// Quick lookup of the EVENT_CATEGORIES entry (icon/color/key) by key.
const cat = key => EVENT_CATEGORIES.find(c => c.key === key)

// Per-category column sets, adapted from the reference demo's events table
// (its Category column is dropped — it's implied by the page). `kind` tells
// CategoryView how to render the cell:
//   id -> muted mono · title -> title + description · type -> neutral tag ·
//   severity/risk -> level badge · status -> status tag · backlog -> BL pills ·
//   date -> localized date · (undefined) -> plain text (owner names).
// Header labels reuse the shared event field labels (event.f.*) where they
// exist, plus a few new dashboard.columns.* keys.
const col = (key, headerKey, kind) => ({ key, headerKey, kind })

// The "rich" event categories (incident/feature/privacy) share the demo's full
// column set. typeKey/ownerKey are the category's own field names.
const richColumns = (typeKey, ownerKey) => [
  col('id', 'project.dashboard.columns.id', 'id'),
  col('title', 'project.dashboard.columns.title', 'title'),
  col(typeKey, `${f}${typeKey}`, 'type'),
  col('severity', 'project.dashboard.columns.severity', 'severity'),
  col('status', `${f}status`, 'status'),
  col('risk', 'project.dashboard.columns.risk', 'risk'),
  col(ownerKey, `${f}${ownerKey}`, 'user'),
  col('changeOwner', `${f}changeOwner`, 'user'),
  col('backlog', 'project.dashboard.columns.backlog', 'backlog'),
  col('changeApprovedBy', `${f}changeApprovedBy`, 'user'),
  col('dateDue', `${f}dateDue`, 'date'),
]

// Explicit column set per category (tasks/reviews carry fewer demo fields).
const COLUMNS = {
  incident: richColumns('incidentType', 'issueOwner'),
  feature: richColumns('featureType', 'featureOwner'),
  privacy: richColumns('requestType', 'requestOwner'),
  task: [
    col('id', 'project.dashboard.columns.id', 'id'),
    col('title', 'project.dashboard.columns.title', 'title'),
    col('taskType', `${f}taskType`, 'type'),
    col('severity', 'project.dashboard.columns.severity', 'severity'),
    col('status', `${f}status`, 'status'),
    col('risk', 'project.dashboard.columns.risk', 'risk'),
    col('owner', `${f}owner`, 'user'),
    col('changeOwner', `${f}changeOwner`, 'user'),
    col('backlog', 'project.dashboard.columns.backlog', 'backlog'),
    col('dateDue', `${f}dateDue`, 'date'),
  ],
  review: [
    col('id', 'project.dashboard.columns.id', 'id'),
    col('title', 'project.dashboard.columns.title', 'title'),
    col('reviewType', `${f}reviewType`, 'type'),
    col('reviewFrequency', `${f}reviewFrequency`),
    col('status', `${f}status`, 'status'),
    col('reviewer', `${f}reviewer`, 'user'),
    col('approvedBy', `${f}approvedBy`, 'user'),
    col('dateDue', `${f}dateDue`, 'date'),
  ],
}

// Charts: always exactly two — a status breakdown plus a breakdown on the
// category's type field (with an apt chart.* titleKey).
const buildCharts = ({ typeKey, typeChartKey }) => [
  { field: 'status', titleKey: 'project.dashboard.chart.byStatus' },
  { field: typeKey, titleKey: typeChartKey },
]

// Leading title badge per category — same shape/idiom as config/kinds.js
// KIND_CONFIG (icon + text/bg/ring, light+dark), so the CategoryView header
// reads like the wizard's step header and the section sidebar.
const BADGES = {
  incident: { icon: 'pi pi-exclamation-triangle', text: 'text-red-600 dark:text-red-400', bg: 'bg-red-50 dark:bg-red-950/40', ring: 'ring-red-200 dark:ring-red-800/60' },
  feature: { icon: 'pi pi-sparkles', text: 'text-blue-600 dark:text-blue-400', bg: 'bg-blue-50 dark:bg-blue-950/40', ring: 'ring-blue-200 dark:ring-blue-800/60' },
  privacy: { icon: 'pi pi-shield', text: 'text-purple-600 dark:text-purple-400', bg: 'bg-purple-50 dark:bg-purple-950/40', ring: 'ring-purple-200 dark:ring-purple-800/60' },
  task: { icon: 'pi pi-check-square', text: 'text-emerald-600 dark:text-emerald-400', bg: 'bg-emerald-50 dark:bg-emerald-950/40', ring: 'ring-emerald-200 dark:ring-emerald-800/60' },
  review: { icon: 'pi pi-sync', text: 'text-amber-600 dark:text-amber-400', bg: 'bg-amber-50 dark:bg-amber-950/40', ring: 'ring-amber-200 dark:ring-amber-800/60' },
}

// Shared KPI definitions (values are computed live from the store per category).
const KPIS = [
  { key: 'total', labelKey: 'project.dashboard.kpi.total' },
  { key: 'open', labelKey: 'project.dashboard.kpi.open' },
  { key: 'overdue', labelKey: 'project.dashboard.kpi.overdue' },
]

// Field-key map per category, feeding the column/chart builders above.
const FIELDS = {
  incident: { typeKey: 'incidentType', ownerKey: 'issueOwner', typeChartKey: 'project.dashboard.chart.byType' },
  feature: { typeKey: 'featureType', ownerKey: 'featureOwner', typeChartKey: 'project.dashboard.chart.byFeatureType' },
  privacy: { typeKey: 'requestType', ownerKey: 'requestOwner', typeChartKey: 'project.dashboard.chart.byRequestType' },
  task: { typeKey: 'taskType', ownerKey: 'owner', typeChartKey: 'project.dashboard.chart.byTaskType' },
  review: { typeKey: 'reviewType', ownerKey: 'reviewer', typeChartKey: 'project.dashboard.chart.byReviewType' },
}

// Assemble the full config for one category from its FIELDS entry.
const build = key => {
  const c = cat(key)
  const fields = FIELDS[key]
  return {
    key,
    icon: c.icon,
    color: c.color,
    badge: BADGES[key],
    titleKey: `project.dashboard.categories.${key}`,
    // Screen description — reuses the New Event category blurbs.
    descKey: `project.dashboard.event.categories.${key}.desc`,
    singularKey: `project.dashboard.categorySingular.${key}`,
    formSchema: EVENT_FORMS[key],
    columns: COLUMNS[key],
    charts: buildCharts(fields),
    kpis: KPIS,
  }
}

// key -> full category dashboard config.
export const CATEGORY_CONFIG = CATEGORY_ORDER.reduce((acc, key) => {
  acc[key] = build(key)
  return acc
}, {})

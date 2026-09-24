// What the record report view reads out of a compose_record_report result, and
// the chart form it picks for it.

import type { ViewField } from './recordLookup'
import { viewDataKey } from './recordLookup'

export type Metric = { key: string; field?: string }

export type ReportView = {
  module: { name: string; handle?: string }
  dimension?: ViewField
  metrics: Metric[]
}

export type Unit = { prefix?: string; suffix?: string; label?: string }

export type Report = {
  rows: Record<string, unknown>[]
  refs: Record<string, string>
  units: Record<string, Unit>
  view?: ReportView
}

export type Form = 'total' | 'line' | 'bar' | 'barHorizontal'

// The key report rows carry each group under, and its value when there is no
// grouping.
export const dimensionKey = 'dimension_0'
const wholeSet = '*'

// Categorical hues in fixed order, validated for light and dark surfaces with
// the dataviz validate_palette.js (all checks pass; three hues sit under 3:1
// on dark, which the table view and direct labels relieve).
export const palette = [
  '#007df4',
  '#007900',
  '#f5008a',
  '#865900',
  '#00a26f',
  '#a83a00',
  '#7e6fff',
  '#bc001d',
]

type ToolResult = {
  structuredContent?: Record<string, unknown>
  _meta?: Record<string, unknown>
}

export function parseReport(result: ToolResult): Report {
  const data = (result.structuredContent ?? {}) as Record<string, unknown>
  return {
    rows: Array.isArray(data.rows) ? (data.rows as Record<string, unknown>[]) : [],
    refs: (data.refs as Record<string, string>) ?? {},
    units: (data.units as Record<string, Unit>) ?? {},
    view: result._meta?.[viewDataKey] as ReportView | undefined,
  }
}

// The metrics to chart: those asked for, or the record count when none were.
export function metricsOf(report: Report): Metric[] {
  return report.view?.metrics?.length ? report.view.metrics : [{ key: 'count' }]
}

export function isTotal(report: Report) {
  return report.rows.length <= 1 && (report.rows[0]?.[dimensionKey] ?? wholeSet) === wholeSet
}

export function formOf(report: Report, labels: string[]): Form {
  if (isTotal(report)) return 'total'
  if (report.view?.dimension?.kind === 'DateTime') return 'line'
  const long = labels.some(l => l.length > 12)
  return long || labels.length > 8 ? 'barHorizontal' : 'bar'
}

type SelectOption = { value: string; text?: string } | string

// The name a person reads for a group: a Select option's text, a referenced
// record's or user's name, or the value as stored. Empty groups get `empty`.
export function groupLabel(value: unknown, report: Report, empty: string): string {
  if (value === null || value === undefined || value === '') return empty

  const v = String(value)
  const dim = report.view?.dimension

  if (dim?.kind === 'Select' && Array.isArray(dim.options)) {
    const opt = (dim.options as SelectOption[]).find(
      o => (typeof o === 'string' ? o : o.value) === v,
    )
    if (opt && typeof opt !== 'string' && opt.text) return opt.text
  }

  return report.refs[v] ?? v
}

export function unitOf(metric: Metric, report: Report): Unit {
  return (metric.field && report.units[metric.field]) || {}
}

export function formatValue(value: unknown, unit: Unit, locale?: string) {
  if (value === null || value === undefined) return ''
  const n = typeof value === 'number' ? value : Number(value)
  const text = Number.isFinite(n)
    ? n.toLocaleString(locale, { maximumFractionDigits: 2 })
    : String(value)
  return `${unit.prefix ?? ''}${text}${unit.suffix ?? ''}`
}

import { describe, expect, it } from 'vitest'
import { viewDataKey } from './recordLookup'
import {
  formOf,
  formatValue,
  groupLabel,
  isTotal,
  metricsOf,
  parseReport,
  unitOf,
} from './recordReport'

const report = (rows: Record<string, unknown>[], view?: object, extra: object = {}) =>
  parseReport({ structuredContent: { rows, ...extra }, _meta: view ? { [viewDataKey]: view } : {} })

const selectView = {
  module: { name: 'Deals' },
  dimension: { name: 'stage', kind: 'Select', options: [{ value: 'won', text: 'Won' }, 'lead'] },
  metrics: [{ key: 'total', field: 'amount' }],
}

describe('metricsOf', () => {
  it('charts the metrics asked for', () => {
    expect(metricsOf(report([], selectView)).map(m => m.key)).toEqual(['total'])
  })

  it('charts the record count when no metric was asked for', () => {
    expect(metricsOf(report([], { module: { name: 'x' }, metrics: [] }))).toEqual([
      { key: 'count' },
    ])
  })
})

describe('formOf', () => {
  it('shows a report over the whole set as a total', () => {
    const r = report([{ dimension_0: '*', count: 9 }])
    expect(isTotal(r)).toBe(true)
    expect(formOf(r, ['*'])).toBe('total')
  })

  it('draws a date dimension as a line', () => {
    const r = report(
      [
        { dimension_0: '2026-08-01', count: 1 },
        { dimension_0: '2026-08-02', count: 2 },
      ],
      {
        ...selectView,
        dimension: { name: 'due', kind: 'DateTime' },
      },
    )
    expect(formOf(r, ['a', 'b'])).toBe('line')
  })

  it('lays bars on their side when labels are long or many', () => {
    const r = report([{ dimension_0: 'a' }, { dimension_0: 'b' }], selectView)
    expect(formOf(r, ['short', 'labels'])).toBe('bar')
    expect(formOf(r, ['a label longer than twelve', 'b'])).toBe('barHorizontal')
    expect(
      formOf(
        r,
        Array.from({ length: 9 }, (_, i) => String(i)),
      ),
    ).toBe('barHorizontal')
  })
})

describe('groupLabel', () => {
  it('reads Select option text, then refs, then the stored value', () => {
    const r = report([], selectView, { refs: { '42': 'Alice Novak' } })
    expect(groupLabel('won', r, '(empty)')).toBe('Won')
    expect(groupLabel('lead', r, '(empty)')).toBe('lead')
    expect(groupLabel('42', r, '(empty)')).toBe('Alice Novak')
    expect(groupLabel(null, r, '(empty)')).toBe('(empty)')
  })
})

describe('units', () => {
  it('finds a metric unit by the field it aggregates', () => {
    const r = report([], selectView, { units: { amount: { prefix: '€ ' } } })
    expect(formatValue(1234.5, unitOf(r.view!.metrics[0], r), 'en')).toBe('€ 1,234.5')
    expect(formatValue(3, unitOf({ key: 'count' }, r), 'en')).toBe('3')
  })
})

import { describe, it, expect, beforeAll } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import CInputModuleField from './CInputModuleField.vue'

beforeAll(() => {
  // PrimeVue overlays bind a matchMedia listener on mount; jsdom has none.
  // @ts-expect-error assigning a stub over a property jsdom does not define
  window.matchMedia = (q: string) => ({
    matches: false,
    media: q,
    onchange: null,
    addEventListener() {},
    removeEventListener() {},
    addListener() {},
    removeListener() {},
    dispatchEvent: () => false,
  })
})

const field = (over: Record<string, unknown>) => ({
  name: 'f',
  label: 'F',
  kind: 'String',
  isMulti: false,
  isQueryable: true,
  ...over,
})

// Stands in for compose.Module: the component only reads `fields` and calls
// `systemFields()`, which is the one list of them.
const moduleWith = (fields: Array<Record<string, unknown>>) => ({
  fields,
  systemFields: () => [
    { name: 'createdAt', label: 'Created at', kind: 'DateTime', isQueryable: true },
    { name: 'revision', label: 'Revision', kind: 'Number', isQueryable: true },
  ],
})

const optionsOf = (props: Record<string, unknown>) => {
  const w = mount(CInputModuleField, {
    props,
    global: { plugins: [createPinia()] },
  })
  return (w.vm as unknown as { options: Array<{ name: string; label: string }> }).options
}

describe('CInputModuleField options', () => {
  it('falls back to the field name when the label is blank', () => {
    const opts = optionsOf({ module: moduleWith([field({ name: 'no_label', label: '' })]) })

    expect(opts.map(o => o.label)).toEqual(['no_label'])
  })

  it('keeps only the requested kinds', () => {
    const opts = optionsOf({
      module: moduleWith([
        field({ name: 'when', kind: 'DateTime' }),
        field({ name: 'who', kind: 'User' }),
        field({ name: 'what', kind: 'String' }),
      ]),
      kinds: ['DateTime', 'User'],
    })

    expect(opts.map(o => o.name).sort()).toEqual(['when', 'who'])
  })

  it('offers system fields only when asked, and after the module’s own', () => {
    const mod = moduleWith([field({ name: 'zzz', label: 'Zzz' })])

    expect(optionsOf({ module: mod }).map(o => o.name)).toEqual(['zzz'])
    expect(optionsOf({ module: mod, includeSystem: true }).map(o => o.name)).toEqual([
      'zzz',
      'createdAt',
      'revision',
    ])
  })

  it('applies a kind filter to system fields too', () => {
    const opts = optionsOf({
      module: moduleWith([field({ name: 'when', kind: 'DateTime' })]),
      kinds: ['DateTime'],
      includeSystem: true,
    })

    expect(opts.map(o => o.name)).toEqual(['when', 'createdAt'])
  })

  it('drops multi-value and non-queryable fields when asked', () => {
    const fields = [
      field({ name: 'single', label: 'Single' }),
      field({ name: 'many', label: 'Many', isMulti: true }),
      field({ name: 'opaque', label: 'Opaque', isQueryable: false }),
    ]

    expect(optionsOf({ module: moduleWith(fields), excludeMulti: true }).map(o => o.name)).toEqual([
      'opaque',
      'single',
    ])
    expect(optionsOf({ module: moduleWith(fields), queryableOnly: true }).map(o => o.name)).toEqual(
      ['many', 'single'],
    )
  })

  it('offers nothing when no module is resolved', () => {
    expect(optionsOf({})).toEqual([])
  })

  it('appends the technical name under showName, and never doubles it up', () => {
    const opts = optionsOf({
      module: moduleWith([
        field({ name: 'total', label: 'Total' }),
        field({ name: 'bare', label: '' }),
      ]),
      showName: true,
    })

    expect(opts.map(o => o.label).sort()).toEqual(['Total (total)', 'bare'])
  })

  it('applies a caller predicate on top of the built-in narrowing', () => {
    const opts = optionsOf({
      module: moduleWith([
        field({ name: 'parent', label: 'Parent', kind: 'Record', options: { moduleID: 'M1' } }),
        field({ name: 'other', label: 'Other', kind: 'Record', options: { moduleID: 'M2' } }),
      ]),
      kinds: ['Record'],
      filter: (f: { options?: { moduleID?: string } }) => f.options?.moduleID === 'M1',
    })

    expect(opts.map(o => o.name)).toEqual(['parent'])
  })

  it('puts extra options above the module’s fields, unfiltered', () => {
    const opts = optionsOf({
      module: moduleWith([field({ name: 'amount', label: 'Amount', kind: 'Number' })]),
      kinds: ['Number'],
      extraOptions: [{ name: 'count', label: 'Count' }],
    })

    expect(opts.map(o => o.name)).toEqual(['count', 'amount'])
  })

  it('still offers extra options with no module', () => {
    expect(
      optionsOf({ extraOptions: [{ name: 'count', label: 'Count' }] }).map(o => o.name),
    ).toEqual(['count'])
  })
})

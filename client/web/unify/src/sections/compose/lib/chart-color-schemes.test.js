import { describe, it, expect } from 'vitest'

import {
  COLOR_SCHEMES_SETTING,
  builtinColorSchemes,
  colorSchemeOptions,
  isCustomScheme,
  newCustomScheme,
  readColorSchemes,
  removeColorScheme,
  upsertColorScheme,
} from './chart-color-schemes'

const tables = {
  tableau: {
    Tableau10: ['#1', '#2'],
    ClassicOrangeBlue13: ['#3'],
  },
}

const custom = [
  { id: 'custom-1', name: 'Brand', colors: ['#FF00AA', '#00FFD5'] },
  { id: 'custom-2', name: 'Muted', colors: ['#999999'] },
]

describe('readColorSchemes', () => {
  it('reads the one setting the palettes live in', () => {
    const $Settings = { get: (k, d) => (k === COLOR_SCHEMES_SETTING ? custom : d) }
    expect(readColorSchemes($Settings)).to.have.lengthOf(2)
  })

  // The webapp renders charts before settings have necessarily answered, and a
  // non-array reaching getColorschemeColors throws inside the render.
  it('is an empty list where settings never loaded', () => {
    expect(readColorSchemes(undefined)).toEqual([])
    expect(readColorSchemes({ get: () => null })).toEqual([])
    expect(readColorSchemes({ get: () => 'nonsense' })).toEqual([])
  })
})

describe('isCustomScheme', () => {
  it('matches the substring getColorschemeColors selects on', () => {
    expect(isCustomScheme('custom-1755000000000')).toBe(true)
    expect(isCustomScheme('tableau.Tableau10')).toBe(false)
    expect(isCustomScheme(undefined)).toBe(false)
  })

  it('accepts every id newCustomScheme can mint', () => {
    expect(isCustomScheme(newCustomScheme(1755000000000).id)).toBe(true)
  })
})

describe('builtinColorSchemes', () => {
  it('names a scheme by family, label and swatch count', () => {
    const [first] = builtinColorSchemes(tables, c => `${c} colors`)
    expect(first.id).toBe('tableau.Tableau10')
    expect(first.name).toBe('Tableau: Tableau (10 colors)')
  })

  it('splits the count off the end of the key, not the middle', () => {
    const found = builtinColorSchemes(tables, c => c).find(
      s => s.id === 'tableau.ClassicOrangeBlue13',
    )
    expect(found.name).toBe('Tableau: ClassicOrangeBlue (13)')
  })

  it('copies the table rather than handing out the live array', () => {
    builtinColorSchemes(tables)[0].colors.push('#nope')
    expect(tables.tableau.Tableau10).toHaveLength(2)
  })
})

describe('colorSchemeOptions', () => {
  it("puts the instance's own schemes ahead of the built-ins", () => {
    const opts = colorSchemeOptions(custom, tables)
    expect(opts.slice(0, 2).map(o => o.id)).toEqual(['custom-1', 'custom-2'])
    expect(opts).toHaveLength(4)
  })

  it('falls back to the id when a stored scheme has no name', () => {
    const [only] = colorSchemeOptions([{ id: 'custom-3', colors: [] }], {})
    expect(only.name).toBe('custom-3')
  })

  it('drops a stored entry with no id, which nothing could ever select', () => {
    expect(colorSchemeOptions([null, { name: 'orphan' }], {})).toEqual([])
  })
})

describe('upsertColorScheme', () => {
  it('appends a scheme the list does not hold', () => {
    const next = upsertColorScheme(custom, { id: 'custom-3', name: 'New', colors: ['#000'] })
    expect(next).toHaveLength(3)
    expect(next[2]).toEqual({ id: 'custom-3', name: 'New', colors: ['#000'] })
  })

  it('replaces in place, keeping the order the picker shows', () => {
    const next = upsertColorScheme(custom, { id: 'custom-1', name: 'Brand v2', colors: ['#111'] })
    expect(next.map(s => s.id)).toEqual(['custom-1', 'custom-2'])
    expect(next[0]).toEqual({ id: 'custom-1', name: 'Brand v2', colors: ['#111'] })
  })

  it('trims the name, so a stray space is not what makes two schemes differ', () => {
    expect(upsertColorScheme([], { id: 'c', name: '  Brand  ', colors: ['#0'] })[0].name).toBe(
      'Brand',
    )
  })

  // A write that fails must leave the picker showing what the server still has.
  it('leaves the list it was given untouched', () => {
    upsertColorScheme(custom, { id: 'custom-1', name: 'Mutated', colors: [] })
    expect(custom[0].name).toBe('Brand')
    expect(custom).toHaveLength(2)
  })
})

describe('removeColorScheme', () => {
  it('drops the named scheme and keeps the rest', () => {
    expect(removeColorScheme(custom, 'custom-1').map(s => s.id)).toEqual(['custom-2'])
  })

  it('is a no-op for an id the list does not hold', () => {
    expect(removeColorScheme(custom, 'custom-9')).toHaveLength(2)
  })
})

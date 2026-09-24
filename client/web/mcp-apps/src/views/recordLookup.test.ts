import { describe, expect, it } from 'vitest'
import { appendPage, columnsOf, parseResult, refLabel, toRow, viewDataKey } from './recordLookup'

const view = {
  module: { name: 'Tasks' },
  fields: [
    { name: 'title', label: 'Title', kind: 'String' },
    { name: 'tags', kind: 'String', multi: true },
    { name: 'contact', kind: 'Record' },
  ],
  links: { '1': 'https://human.example/records/1' },
}

const result = {
  structuredContent: {
    records: [
      {
        recordID: '1',
        values: [
          { name: 'title', value: 'Call Alice' },
          { name: 'tags', value: 'a' },
          { name: 'tags', value: 'b' },
          { name: 'contact', value: '9' },
        ],
      },
    ],
    refs: { '9': 'Alice Novak' },
    nextPageCursor: 'c1',
  },
  _meta: { [viewDataKey]: view },
}

describe('parseResult', () => {
  it('reads records, refs, cursor and the view', () => {
    const page = parseResult(result)

    expect(page.rows).toEqual([
      {
        _dataKey: '1',
        recordID: '1',
        values: { title: 'Call Alice', tags: ['a', 'b'], contact: '9' },
      },
    ])
    expect(page.refs).toEqual({ '9': 'Alice Novak' })
    expect(page.cursor).toBe('c1')
    expect(page.view?.module.name).toBe('Tasks')
  })

  it('reads a lookup by recordID as one row', () => {
    const page = parseResult({
      structuredContent: { recordID: '5', values: [{ name: 'title', value: 'x' }] },
    })
    expect(page.rows.map(r => r.recordID)).toEqual(['5'])
  })

  it('reads the last page as having no cursor', () => {
    expect(
      parseResult({ structuredContent: { records: [], nextPageCursor: null } }).cursor,
    ).toBeUndefined()
  })
})

describe('toRow', () => {
  it('keeps the last value of a single-value field that repeats', () => {
    const row = toRow(
      {
        recordID: '1',
        values: [
          { name: 'title', value: 'a' },
          { name: 'title', value: 'b' },
        ],
      },
      [],
    )
    expect(row.values.title).toBe('b')
  })
})

describe('columnsOf', () => {
  it('uses the module fields when the server described them', () => {
    expect(columnsOf(parseResult(result)).map(f => f.name)).toEqual(['title', 'tags', 'contact'])
  })

  it('falls back to the value names the records carry', () => {
    const page = parseResult({ structuredContent: result.structuredContent })
    expect(columnsOf(page)).toEqual([
      { name: 'title', kind: 'String' },
      { name: 'tags', kind: 'String' },
      { name: 'contact', kind: 'String' },
    ])
  })
})

describe('refLabel', () => {
  it('names each referenced ID and leaves unknown ones as they are', () => {
    expect(refLabel(['9', '8'], { '9': 'Alice Novak' })).toBe('Alice Novak, 8')
    expect(refLabel(undefined, {})).toBe('')
  })
})

describe('appendPage', () => {
  it('adds rows, merges refs and links, and takes the next cursor', () => {
    const first = parseResult(result)
    const next = parseResult({
      structuredContent: {
        records: [{ recordID: '2', values: [] }],
        refs: { '7': 'Bojan' },
        nextPageCursor: null,
      },
      _meta: { [viewDataKey]: { ...view, links: { '2': 'https://human.example/records/2' } } },
    })

    const both = appendPage(first, next)

    expect(both.rows.map(r => r.recordID)).toEqual(['1', '2'])
    expect(both.refs).toEqual({ '9': 'Alice Novak', '7': 'Bojan' })
    expect(both.cursor).toBeUndefined()
    expect(Object.keys(both.view?.links ?? {})).toEqual(['1', '2'])
  })
})

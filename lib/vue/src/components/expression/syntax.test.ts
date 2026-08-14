import { describe, it, expect } from 'vitest'
import { buildScope, membersAt, resolvePath } from './catalog'
import { completionAt, lintExpression, scanHoles, tokenRanges } from './syntax'

const orders = {
  moduleID: '1',
  name: 'Order',
  fields: [
    { name: 'status', label: 'Status', kind: 'Select' },
    { name: 'quantity', label: 'Quantity', kind: 'Number' },
  ],
}

const recordScope = buildScope({ recordModule: orders, hasRecord: true })
const listScope = buildScope({ hasRecord: false })

describe('scanHoles', () => {
  it('finds a hole and its bounds', () => {
    const [hole] = scanHoles('id = ${recordID}')
    expect(hole).toMatchObject({ from: 5, to: 16, inner: 'recordID', terminated: true })
  })

  it('reports an unterminated hole rather than swallowing the rest', () => {
    const [hole] = scanHoles('id = ${recordID')
    expect(hole.terminated).toBe(false)
    expect(hole.to).toBe(15)
  })

  it('does not close on a brace inside a string', () => {
    const [hole] = scanHoles("x = ${a === '}' ? b : c}")
    expect(hole.inner).toBe("a === '}' ? b : c")
    expect(hole.terminated).toBe(true)
  })

  it('closes at the matching brace when the hole nests one', () => {
    const [hole] = scanHoles('x = ${ {a: 1}.a }')
    expect(hole.inner).toBe(' {a: 1}.a ')
  })

  it('finds every hole', () => {
    expect(scanHoles('${a} AND ${b}').map(h => h.inner)).toEqual(['a', 'b'])
  })
})

describe('resolvePath', () => {
  it('walks into module fields under record.values', () => {
    expect(resolvePath(recordScope, ['record', 'values', 'status']).status).toBe('found')
  })

  it('names the unknown segment and its parent', () => {
    const res = resolvePath(recordScope, ['record', 'values', 'nope'])
    expect(res).toMatchObject({ status: 'unknown', segment: 'nope' })
  })

  it('reports unverifiable rather than unknown when the module is not loaded', () => {
    const scope = buildScope({ recordModule: null, hasRecord: true })
    expect(resolvePath(scope, ['record', 'values', 'anything']).status).toBe('unverifiable')
  })
})

describe('lintExpression', () => {
  it('passes a correct prefilter', () => {
    expect(lintExpression("status = '${record.values.status}'", recordScope)).toEqual([])
  })

  it('flags an unterminated hole', () => {
    const [d] = lintExpression('id = ${recordID', recordScope)
    expect(d.message).toMatch(/Unterminated/)
  })

  it('flags an empty hole', () => {
    expect(lintExpression('id = ${}', recordScope)[0].message).toBe('Empty ${}')
  })

  it('flags a mistyped variable and suggests the real ones', () => {
    const [d] = lintExpression('id = ${recrdID}', recordScope)
    expect(d.message).toContain("'recrdID' is not an available variable")
    expect(d.message).toContain('recordID')
  })

  it('flags a mistyped module field', () => {
    const [d] = lintExpression('qty > ${record.values.quantitiy}', recordScope)
    expect(d.message).toContain("'quantitiy' is not a member of values")
    expect(d.message).toContain('quantity')
  })

  it('flags a record variable on a page that has no record', () => {
    const [d] = lintExpression('id = ${recordID}', listScope)
    expect(d.message).toContain("'recordID' is not an available variable")
  })

  it('stays silent on an expression it cannot check', () => {
    expect(lintExpression("x = ${record.values.status || 'none'}", recordScope)).toEqual([])
    expect(lintExpression('x = ${someFn(record)}', recordScope)).toEqual([])
  })

  it('accepts a reachable property that is not suggested', () => {
    expect(lintExpression('r = ${record.revision}', recordScope)).toEqual([])
  })

  it('marks the whole hole, not just the bad word', () => {
    const [d] = lintExpression('id = ${recrdID}', recordScope)
    expect([d.from, d.to]).toEqual([5, 15])
  })
})

describe('membersAt', () => {
  it('offers the roots when nothing is typed', () => {
    expect(membersAt(recordScope, []).map(e => e.name)).toEqual([
      'recordID',
      'ownerID',
      'record',
      'userID',
      'user',
    ])
  })

  it('drops record roots when there is no record', () => {
    expect(membersAt(listScope, []).map(e => e.name)).toEqual(['userID', 'user'])
  })

  it('withholds non-suggested members but keeps them resolvable', () => {
    const names = membersAt(recordScope, ['record']).map(e => e.name)
    expect(names).toContain('recordID')
    expect(names).not.toContain('revision')
    expect(resolvePath(recordScope, ['record', 'revision']).status).toBe('found')
  })

  it('offers module fields under record.values', () => {
    expect(membersAt(recordScope, ['record', 'values']).map(e => e.name)).toEqual([
      'status',
      'quantity',
    ])
  })
})

describe('tokenRanges', () => {
  const kinds = (text: string, dialect: 'ql' | 'interpolation' = 'ql') =>
    tokenRanges(text, dialect).map(r => [r.kind, text.slice(r.from, r.to)])

  it('marks holes in either dialect', () => {
    expect(kinds('/page/${recordID}', 'interpolation')).toEqual([['hole', '${recordID}']])
  })

  it('colours QL keywords, strings and numbers outside holes', () => {
    expect(kinds("status = 'Open' AND qty > 10")).toEqual([
      ['string', "'Open'"],
      ['keyword', 'AND'],
      ['number', '10'],
    ])
  })

  it('leaves a keyword inside a hole to the hole', () => {
    expect(kinds('x = ${a AND b}')).toEqual([['hole', '${a AND b}']])
  })

  it('does not colour a keyword inside a string literal', () => {
    expect(kinds("x = 'AND 10'")).toEqual([['string', "'AND 10'"]])
  })

  it('returns ranges in document order and disjoint', () => {
    const ranges = tokenRanges("a = 1 AND b = ${record.values.x} AND c = 'z'", 'ql')
    const sorted = [...ranges].sort((a, b) => a.from - b.from)
    expect(ranges).toEqual(sorted)
    ranges.forEach((r, i) => i && expect(r.from).toBeGreaterThanOrEqual(ranges[i - 1].to))
  })

  it('colours nothing but holes in a plain template', () => {
    expect(kinds("status = 'Open' AND qty > 10", 'interpolation')).toEqual([])
  })
})

describe('completionAt', () => {
  const complete = (text: string, dialect: 'ql' | 'interpolation' = 'ql', fields = orders.fields) =>
    completionAt(text, text.length, recordScope, dialect, fields)

  it('offers roots right after ${', () => {
    expect(complete('id = ${')?.options.map(o => o.label)).toContain('recordID')
  })

  it('narrows by what is typed and replaces only that word', () => {
    const res = complete('id = ${rec')
    expect(res?.options.map(o => o.label)).toEqual(['recordID', 'record'])
    expect(res?.from).toBe(7)
  })

  it('walks into module fields after a dot', () => {
    expect(complete('x = ${record.values.')?.options.map(o => o.label)).toEqual([
      'status',
      'quantity',
    ])
  })

  it('carries the field label into the completion detail', () => {
    const opt = complete('x = ${record.values.stat')?.options[0]
    expect(opt).toMatchObject({ label: 'status', detail: 'Select · Status' })
  })

  it('offers queried module fields and keywords outside a hole in QL', () => {
    const labels = complete('stat')?.options.map(o => o.label)
    expect(labels).toContain('status')
  })

  it('offers keywords outside a hole in QL', () => {
    expect(complete('a = 1 AN')?.options.map(o => o.label)).toEqual(['AND'])
  })

  it('offers nothing outside a hole in a plain template', () => {
    expect(complete('some text', 'interpolation')).toBeNull()
  })

  it('still offers scope inside a hole in a plain template', () => {
    expect(complete('/page/${record', 'interpolation')?.options.map(o => o.label)).toEqual([
      'recordID',
      'record',
    ])
  })

  it('does not offer scope members once the hole is closed', () => {
    // `rec` here is a bare QL identifier, not a variable reference.
    const labels =
      completionAt('${recordID} AND rec', 19, recordScope, 'ql', [])?.options.map(o => o.label) ||
      []
    expect(labels).not.toContain('recordID')
    expect(labels).not.toContain('record')
  })
})

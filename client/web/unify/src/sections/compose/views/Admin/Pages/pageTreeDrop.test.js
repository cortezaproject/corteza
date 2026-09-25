import { describe, it, expect } from 'vitest'
import { projectDrop, dropPlan } from './pageTreeDrop'

// Rows of a tree drawn 50px tall with a 12px gap, 27px per level:
//   A (0)
//     B (1)
//       C (2)
//     D (1)
//   E (0)
//   F (0)
const geometry = { indent: 27, baseX: 100 }
const row = (key, parentKey, depth, i) => ({
  key,
  parentKey,
  depth,
  top: 10 + i * 62,
  bottom: 60 + i * 62,
})
const rows = [
  row('A', '0', 0, 0),
  row('B', 'A', 1, 1),
  row('C', 'B', 2, 2),
  row('D', 'A', 1, 3),
  row('E', '0', 0, 4),
  row('F', '0', 0, 5),
]
const gapAfter = i => 60 + i * 62 + 6
const xAt = depth => geometry.baseX + depth * geometry.indent + 3

describe('projectDrop', () => {
  it('lands before the first row only at the top level', () => {
    expect(projectDrop(rows, xAt(2), 0, geometry)).toMatchObject({
      parentKey: '0',
      afterKey: null,
      depth: 0,
    })
  })

  it('after a row, the pointer picks the depth between the rows around it', () => {
    // after C, between C (depth 2) and D (depth 1): 1..3 are possible
    expect(projectDrop(rows, xAt(3), gapAfter(2), geometry)).toMatchObject({
      parentKey: 'C',
      afterKey: null,
      depth: 3,
    })
    expect(projectDrop(rows, xAt(2), gapAfter(2), geometry)).toMatchObject({
      parentKey: 'B',
      afterKey: 'C',
      depth: 2,
    })
    expect(projectDrop(rows, xAt(1), gapAfter(2), geometry)).toMatchObject({
      parentKey: 'A',
      afterKey: 'B',
      depth: 1,
    })
    // no shallower than D, which would be cut off from A
    expect(projectDrop(rows, xAt(0), gapAfter(2), geometry)).toMatchObject({
      parentKey: 'A',
      afterKey: 'B',
      depth: 1,
    })
  })

  it('cannot land between a parent and its first sub-page', () => {
    expect(projectDrop(rows, xAt(0), gapAfter(0), geometry)).toMatchObject({
      parentKey: 'A',
      afterKey: null,
      depth: 1,
    })
  })

  it('a leaf takes a sub-page when the pointer pulls right', () => {
    expect(projectDrop(rows, xAt(1), gapAfter(5), geometry)).toMatchObject({
      parentKey: 'F',
      afterKey: null,
      depth: 1,
    })
    expect(projectDrop(rows, xAt(9), gapAfter(5), geometry)).toMatchObject({
      parentKey: 'F',
      afterKey: null,
      depth: 1,
    })
    expect(projectDrop(rows, xAt(0), gapAfter(5), geometry)).toMatchObject({
      parentKey: '0',
      afterKey: 'F',
      depth: 0,
    })
  })

  it("on a row's middle band the page goes into that row, first", () => {
    const rowMid = i => 10 + i * 62 + 25
    expect(projectDrop(rows, xAt(0), rowMid(4), geometry)).toMatchObject({
      parentKey: 'E',
      afterKey: null,
      depth: 1,
      intoKey: 'E',
    })
    expect(projectDrop(rows, xAt(0), rowMid(0), geometry)).toMatchObject({
      parentKey: 'A',
      afterKey: null,
      depth: 1,
      intoKey: 'A',
    })
    // the top and bottom quarters are beside it, at the depth x asks for
    expect(projectDrop(rows, xAt(0), rowMid(4) - 20, geometry)).toMatchObject({
      parentKey: '0',
      afterKey: 'A',
      depth: 0,
      intoKey: null,
    })
    expect(projectDrop(rows, xAt(0), rowMid(4) + 20, geometry)).toMatchObject({
      parentKey: '0',
      afterKey: 'E',
      depth: 0,
      intoKey: null,
    })
  })

  it('the line sits in the gap after the row the pointer passed', () => {
    const spot = projectDrop(rows, xAt(0), gapAfter(0) - 20, geometry)
    expect(spot.afterKey).toBe(null)
    expect(spot.parentKey).toBe('A')
    expect(spot.lineY).toBe((60 + 72) / 2)
    expect(projectDrop(rows, xAt(0), 999, geometry).lineY).toBe(60 + 5 * 62 + 6)
  })

  it('an empty tree lands at the top level', () => {
    expect(projectDrop([], 0, 0, geometry)).toMatchObject({
      parentKey: '0',
      afterKey: null,
      depth: 0,
    })
  })
})

describe('dropPlan', () => {
  const children = key => ({ 0: ['A', 'E', 'F'], A: ['B', 'D'], B: ['C'] })[key] ?? []

  it('places the page after its sibling, ignoring where it was', () => {
    expect(dropPlan({ parentKey: 'A', afterKey: 'B' }, 'D', children)).toEqual({
      parentKey: 'A',
      pageIDs: ['B', 'D'],
      index: 1,
    })
    expect(dropPlan({ parentKey: '0', afterKey: 'F' }, 'A', children)).toEqual({
      parentKey: '0',
      pageIDs: ['E', 'F', 'A'],
      index: 2,
    })
  })

  it('makes the page the first sub-page when it follows nothing', () => {
    expect(dropPlan({ parentKey: 'F', afterKey: null }, 'C', children)).toEqual({
      parentKey: 'F',
      pageIDs: ['C'],
      index: 0,
    })
    expect(dropPlan({ parentKey: 'A', afterKey: null }, 'D', children)).toEqual({
      parentKey: 'A',
      pageIDs: ['D', 'B'],
      index: 0,
    })
  })
})

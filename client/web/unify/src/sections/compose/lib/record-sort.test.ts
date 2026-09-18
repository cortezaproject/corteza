import { describe, it, expect } from 'vitest'
import { parseSortExpression, sortExpression } from './record-sort'

describe('sortExpression', () => {
  it('uses the presort until the user sorts', () => {
    expect(sortExpression(null, 'color_count, color, market_value DESC')).toBe(
      'color_count, color, market_value DESC',
    )
    expect(sortExpression(undefined, 'createdAt DESC')).toBe('createdAt DESC')
  })

  it('is no sort once the user has cleared every column, presort or not', () => {
    expect(sortExpression([], 'createdAt DESC')).toBe('')
    expect(sortExpression([])).toBe('')
  })

  it('is empty with neither columns nor a presort', () => {
    expect(sortExpression(null)).toBe('')
    expect(sortExpression(null, null)).toBe('')
  })

  it('replaces the presort with the sorted columns', () => {
    expect(sortExpression([{ field: 'rarity', order: 1 }], 'createdAt DESC')).toBe('rarity ASC')
  })

  it('keeps the columns in the order they were picked, each with its own direction', () => {
    expect(
      sortExpression([
        { field: 'rarity', order: 1 },
        { field: 'market_value', order: -1 },
      ]),
    ).toBe('rarity ASC, market_value DESC')
  })
})

describe('parseSortExpression', () => {
  it('reads each key with its direction, ascending unless DESC', () => {
    expect(parseSortExpression('color_count, color ASC, market_value DESC')).toEqual([
      { field: 'color_count', order: 1 },
      { field: 'color', order: 1 },
      { field: 'market_value', order: -1 },
    ])
  })

  it('takes the direction in any case and ignores extra spacing', () => {
    expect(parseSortExpression('  createdAt   desc ,name')).toEqual([
      { field: 'createdAt', order: -1 },
      { field: 'name', order: 1 },
    ])
  })

  it('is empty for an empty or missing expression', () => {
    expect(parseSortExpression('')).toEqual([])
    expect(parseSortExpression(null)).toEqual([])
    expect(parseSortExpression(' , ')).toEqual([])
  })

  it('round-trips through sortExpression', () => {
    const meta = parseSortExpression('rarity, market_value DESC')
    expect(sortExpression(meta)).toBe('rarity ASC, market_value DESC')
  })
})

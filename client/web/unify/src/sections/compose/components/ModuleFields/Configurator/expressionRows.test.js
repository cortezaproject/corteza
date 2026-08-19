import { describe, it, expect } from 'vitest'
import { incompleteExpressionRows, isBlankExpressionRow } from './expressionRows'

describe('isBlankExpressionRow', () => {
  it('treats nothing, whitespace and a missing value as blank', () => {
    expect(isBlankExpressionRow('')).toBe(true)
    expect(isBlankExpressionRow('   ')).toBe(true)
    expect(isBlankExpressionRow(undefined)).toBe(true)
    expect(isBlankExpressionRow(null)).toBe(true)
  })

  it('accepts an expression', () => {
    expect(isBlankExpressionRow('trim(value)')).toBe(false)
  })
})

describe('incompleteExpressionRows', () => {
  it('passes a field with no expressions at all', () => {
    expect(incompleteExpressionRows(undefined)).toBe(false)
    expect(incompleteExpressionRows({})).toBe(false)
    expect(incompleteExpressionRows({ sanitizers: [], validators: [] })).toBe(false)
  })

  it('passes a fully filled-in set', () => {
    expect(
      incompleteExpressionRows({
        sanitizers: ['trim(value)'],
        validators: [{ test: 'value == ""', error: 'Enter a value' }],
      }),
    ).toBe(false)
  })

  it('catches a sanitizer with no expression', () => {
    expect(incompleteExpressionRows({ sanitizers: ['trim(value)', '  '] })).toBe(true)
  })

  it('catches a validator with no expression', () => {
    expect(incompleteExpressionRows({ validators: [{ test: '', error: 'Enter a value' }] })).toBe(
      true,
    )
  })

  it('catches a validator with no error message', () => {
    expect(incompleteExpressionRows({ validators: [{ test: 'value == ""', error: '' }] })).toBe(
      true,
    )
  })

  it('catches a row added and left untouched', () => {
    expect(incompleteExpressionRows({ validators: [{}] })).toBe(true)
  })
})

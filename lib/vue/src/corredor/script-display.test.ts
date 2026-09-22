import { describe, expect, it } from 'vitest'
import { constraintChips } from './script-display'

describe('constraintChips', () => {
  it('spells the operator the server assumes', () => {
    expect(constraintChips([{ name: 'module', value: ['agent-contact'] }])).toEqual([
      'module = agent-contact',
    ])
  })

  it('keeps an explicit operator and joins the alternatives', () => {
    expect(constraintChips([{ name: 'module', op: 'like', value: ['agent-*', 'crm-*'] }])).toEqual([
      'module like agent-*, crm-*',
    ])
  })

  it('answers empty for no constraints', () => {
    expect(constraintChips()).toEqual([])
    expect(constraintChips(null)).toEqual([])
  })
})

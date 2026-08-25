import { describe, it, expect } from 'vitest'
import { moduleFieldKeyLabel } from './resource-translations'

// The translator shows one row per translation key. A module field's keys are
// dotted paths the person translating never wrote, so each one is named; a key
// the mapper does not know returns '' and the form falls back to its own
// generic formatting.
const t = (k: string, p?: Record<string, string>) => (p ? `${k}:${p.value}` : k)

describe('moduleFieldKeyLabel', () => {
  it('names a select option after the option value it belongs to', () => {
    expect(moduleFieldKeyLabel('meta.options.new.text', t)).toBe('translator.keys.option:new')
  })

  it('keeps an option value that contains dots whole', () => {
    expect(moduleFieldKeyLabel('meta.options.a.b.text', t)).toBe('translator.keys.option:a.b')
  })

  it('names the label, description, hint and bool keys', () => {
    expect(moduleFieldKeyLabel('label', t)).toBe('translator.keys.label')
    expect(moduleFieldKeyLabel('meta.description.view', t)).toBe('translator.keys.description-view')
    expect(moduleFieldKeyLabel('meta.description.edit', t)).toBe('translator.keys.description-edit')
    expect(moduleFieldKeyLabel('meta.hint.view', t)).toBe('translator.keys.hint-view')
    expect(moduleFieldKeyLabel('meta.hint.edit', t)).toBe('translator.keys.hint-edit')
    expect(moduleFieldKeyLabel('meta.bool.true.label', t)).toBe('translator.keys.bool-true')
    expect(moduleFieldKeyLabel('meta.bool.false.label', t)).toBe('translator.keys.bool-false')
  })

  it('names a validator error whatever the validator id is', () => {
    expect(moduleFieldKeyLabel('expression.validator.12345.error', t)).toBe(
      'translator.keys.validator-error',
    )
  })

  it("returns '' for a key that is not a field's own", () => {
    expect(moduleFieldKeyLabel('name', t)).toBe('')
    expect(moduleFieldKeyLabel('title', t)).toBe('')
    expect(moduleFieldKeyLabel('pageBlock.3.title', t)).toBe('')
  })
})

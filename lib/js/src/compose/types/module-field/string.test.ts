import { expect } from 'chai'
import { ModuleFieldString } from './string'

describe('string field options', () => {
  it('sanitizes XSS by default', () => {
    const f = new ModuleFieldString({ name: 'f' })
    expect(f.options.sanitizeXSS).to.eq(true)
  })

  it('keeps an explicit opt-out', () => {
    const f = new ModuleFieldString({ name: 'f', options: { sanitizeXSS: false } })
    expect(f.options.sanitizeXSS).to.eq(false)
  })
})

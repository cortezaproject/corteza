import { describe, it, expect } from 'vitest'
import { coveredByFamily, defaultModeFor, modeOf, splitGrants, tally } from './toolAccess'

const read = { name: 'compose_record_lookup', groups: ['usage'], risk: 'read' }
const write = { name: 'compose_record_create', groups: ['usage'], risk: 'write' }
const destroy = { name: 'compose_record_delete', groups: ['usage'], risk: 'destructive' }

function grants(...gg: any[]) {
  return splitGrants(gg)
}

describe('toolAccess', () => {
  // The runtime resolves a grant that names no mode when it runs the agent, so
  // the editor has to resolve it the same way to show the agent back.
  it('resolves a grant that names no mode the way the runtime does', () => {
    expect(defaultModeFor('read')).toBe('always')
    expect(defaultModeFor('write')).toBe('ask')
    expect(defaultModeFor('destructive')).toBe('ask')
    expect(defaultModeFor(undefined)).toBe('always')
  })

  it('reads a tool nothing grants as blocked', () => {
    const { named, families } = grants()
    expect(modeOf(read, named, families)).toBe('deny')
  })

  it('reads a named grant at the mode it states, or its risk default', () => {
    const { named, families } = grants(
      { name: read.name },
      { name: write.name, permission: 'always' },
    )
    expect(modeOf(read, named, families)).toBe('always')
    expect(modeOf(write, named, families)).toBe('always')
  })

  // A ceiling admits its own level and everything below it.
  it('covers a tool at or below a family ceiling, and no higher', () => {
    const families = [{ group: 'usage', maxRisk: 'write' }]
    expect(coveredByFamily(read, families)).toBe(true)
    expect(coveredByFamily(write, families)).toBe(true)
    expect(coveredByFamily(destroy, families)).toBe(false)
    expect(coveredByFamily({ groups: ['configuring'], risk: 'read' }, families)).toBe(false)
  })

  it('gives a covered tool its risk default without an entry of its own', () => {
    const { named, families } = grants({ group: 'usage', maxRisk: 'write' })
    expect(modeOf(read, named, families)).toBe('always')
    expect(modeOf(write, named, families)).toBe('ask')
    expect(modeOf(destroy, named, families)).toBe('deny')
  })

  // "All data tools, but never delete" is a family grant plus a named entry.
  it('lets a named block override the family that covers it', () => {
    const { named, families } = grants(
      { group: 'usage', maxRisk: 'destructive' },
      { name: destroy.name, permission: 'deny' },
    )
    expect(modeOf(destroy, named, families)).toBe('deny')
    expect(modeOf(read, named, families)).toBe('always')
  })

  it('counts each mode, leaving out the ones nothing sits in', () => {
    const { named, families } = grants({ name: read.name }, { name: write.name })
    expect(tally([read, write, destroy], named, families)).toEqual([
      { mode: 'always', n: 1 },
      { mode: 'ask', n: 1 },
      { mode: 'deny', n: 1 },
    ])
  })

  it('counts nothing granted as all blocked', () => {
    const { named, families } = grants()
    expect(tally([read, write], named, families)).toEqual([{ mode: 'deny', n: 2 }])
  })
})

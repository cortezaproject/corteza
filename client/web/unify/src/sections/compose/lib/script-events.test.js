import { describe, it, expect } from 'vitest'
import { eventbus } from '@planetcrust/human-js'
import { isScriptAbort, scriptConstraintMatcher } from './script-events'

const namespace = { slug: 'agent-sandbox', name: 'Agent sandbox' }
const module = { handle: 'agent-contact', name: 'Contact' }

const match = scriptConstraintMatcher({ namespace, module })

function constraint(name, value, op) {
  return eventbus.ConstraintMaker({ name, value: [value], op })
}

describe('scriptConstraintMatcher', () => {
  it('matches a namespace by slug and by name', () => {
    expect(match(constraint('namespace', 'agent-sandbox'))).toBe(true)
    expect(match(constraint('namespace.slug', 'agent-sandbox'))).toBe(true)
    expect(match(constraint('namespace.name', 'Agent sandbox'))).toBe(true)
    expect(match(constraint('namespace', 'something-else'))).toBe(false)
  })

  it('matches a module by handle and by name', () => {
    expect(match(constraint('module', 'agent-contact'))).toBe(true)
    expect(match(constraint('module.handle', 'agent-contact'))).toBe(true)
    expect(match(constraint('module.name', 'Contact'))).toBe(true)
    expect(match(constraint('module', 'agent-account'))).toBe(false)
  })

  it('honours the constraint operator', () => {
    expect(match(constraint('module', 'agent-*', 'like'))).toBe(true)
    expect(match(constraint('module', 'agent-contact', '!='))).toBe(false)
  })

  it('refuses a constraint on anything it does not carry', () => {
    expect(match(constraint('record.values.status', 'active'))).toBe(false)
  })

  it('refuses every constraint when neither resource is known', () => {
    const blind = scriptConstraintMatcher()

    expect(blind(constraint('namespace', 'agent-sandbox'))).toBe(false)
    expect(blind(constraint('module', 'agent-contact'))).toBe(false)
  })
})

describe('isScriptAbort', () => {
  it('knows a script refusing the action from a script that broke', () => {
    // corredor.Exec rejects with Aborted when a script returns false; the bus
    // lower-cases it when a handler does
    expect(isScriptAbort(new Error('Aborted'))).toBe(true)
    expect(isScriptAbort(new Error('aborted'))).toBe(true)
    expect(isScriptAbort(new Error('x is not a function'))).toBe(false)
    expect(isScriptAbort(undefined)).toBe(false)
  })
})

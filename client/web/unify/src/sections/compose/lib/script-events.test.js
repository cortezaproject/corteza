import { describe, it, expect } from 'vitest'
import { eventbus } from '@planetcrust/human-js'
import {
  isScriptAbort,
  pageFitsResourceType,
  pageScriptResourceTypes,
  scriptConstraintMatcher,
  triggerCanApply,
} from './script-events'

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

describe('pageScriptResourceTypes', () => {
  it('lets a record page satisfy every compose resource', () => {
    expect(pageScriptResourceTypes({ isRecordPage: true })).toEqual([
      'compose',
      'compose:namespace',
      'compose:page',
      'compose:module',
      'compose:record',
    ])
  })

  it('leaves a page with no record unable to satisfy one', () => {
    expect(pageScriptResourceTypes({ isRecordPage: false })).not.toContain('compose:record')
    expect(pageScriptResourceTypes()).not.toContain('compose:record')
  })
})

describe('pageFitsResourceType', () => {
  const recordPage = { isRecordPage: true }
  const dashboardPage = { isRecordPage: false }

  it('fits a record-bound script to a record page only', () => {
    expect(pageFitsResourceType(recordPage, 'compose:record')).toBe(true)
    expect(pageFitsResourceType(dashboardPage, 'compose:record')).toBe(false)
  })

  it('fits the resources every page carries', () => {
    expect(pageFitsResourceType(dashboardPage, 'compose:module')).toBe(true)
    expect(pageFitsResourceType(dashboardPage, 'compose:namespace')).toBe(true)
    expect(pageFitsResourceType(dashboardPage, 'compose:page')).toBe(true)
    expect(pageFitsResourceType(dashboardPage, 'compose')).toBe(true)
  })

  it('refuses a resource no compose page carries', () => {
    expect(pageFitsResourceType(recordPage, 'system:user')).toBe(false)
  })

  it('fits a trigger that names no resource at all', () => {
    expect(pageFitsResourceType(recordPage, undefined)).toBe(true)
  })
})

describe('triggerCanApply', () => {
  const here = { namespace, module }

  it('applies when every constraint it can settle matches', () => {
    expect(
      triggerCanApply(
        {
          constraints: [
            { name: 'module', value: ['agent-contact'] },
            { name: 'namespace', value: ['agent-sandbox'] },
          ],
        },
        here,
      ),
    ).toBe(true)
  })

  it('never applies when a constraint it can settle does not match', () => {
    expect(
      triggerCanApply({ constraints: [{ name: 'module', value: ['agent-deal'] }] }, here),
    ).toBe(false)
    expect(
      triggerCanApply({ constraints: [{ name: 'namespace', value: ['other-space'] }] }, here),
    ).toBe(false)
  })

  it('honours the constraint operator', () => {
    expect(
      triggerCanApply({ constraints: [{ name: 'module', op: 'like', value: ['agent-*'] }] }, here),
    ).toBe(true)
    expect(
      triggerCanApply(
        { constraints: [{ name: 'module', op: '!=', value: ['agent-contact'] }] },
        here,
      ),
    ).toBe(false)
  })

  it('leaves a constraint on anything the page cannot settle alone', () => {
    expect(
      triggerCanApply({ constraints: [{ name: 'record.values.status', value: ['active'] }] }, here),
    ).toBe(true)
    expect(triggerCanApply({ constraints: [{ name: 'module', value: ['agent-deal'] }] }, {})).toBe(
      true,
    )
    expect(
      triggerCanApply({ constraints: [{ name: 'module', value: ['x'] }] }, { namespace }),
    ).toBe(true)
  })

  it('leaves a constraint it cannot read alone', () => {
    expect(triggerCanApply({ constraints: [{ name: 'module' }] }, here)).toBe(true)
    expect(
      triggerCanApply({ constraints: [{ name: 'module', op: 'nonsense', value: ['x'] }] }, here),
    ).toBe(true)
  })

  it('applies a trigger with no constraints anywhere', () => {
    expect(triggerCanApply({ constraints: null }, here)).toBe(true)
    expect(triggerCanApply(undefined, here)).toBe(true)
  })
})

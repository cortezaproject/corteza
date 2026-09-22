import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'
import { sections } from './index'

// The registry applications the server provisions, by `unify.url`. The section
// gate resolves against these, so a section naming a url that is not here is
// gated on an application that will never exist — permanently unreachable.
const PROVISIONED = readFileSync(
  join(__dirname, '../../../../../server/provision/101_applications/0200_applications.yaml'),
  'utf8',
)
  .split('\n')
  .map(line => /^\s*url:\s*(\S+)\s*$/.exec(line))
  .filter(Boolean)
  .map(m => m[1].replace(/^\/|\/$/g, ''))

describe('section access contract', () => {
  it('every section declares the application that gates it', () => {
    const undeclared = sections.filter(s => s.app === undefined).map(s => s.id)
    expect(undeclared).toEqual([])
  })

  it('names a provisioned application, or null for the one ungated section', () => {
    const unknown = sections
      .filter(s => s.app !== null)
      .filter(s => !PROVISIONED.includes(String(s.app).replace(/^\/|\/$/g, '')))
      .map(s => ({ id: s.id, app: s.app }))
    expect(unknown).toEqual([])
  })

  it('leaves exactly one section reachable without a permission', () => {
    // More than one ungated section is a hole; none is a lock-out, since a user
    // with no applications would have nowhere to be redirected to. A `perApp`
    // section declares `app: null` too, and is gated on the application named
    // in its route rather than on nothing.
    expect(sections.filter(s => s.app === null && !s.perApp).map(s => s.id)).toEqual(['home'])
  })

  it('gates a per-application section on the route, not on an application of its own', () => {
    // `perApp` with a url of its own would be two gates on one entry, and the
    // section's would be the one that answers.
    expect(sections.filter(s => s.perApp && s.app !== null).map(s => s.id)).toEqual([])
  })
})

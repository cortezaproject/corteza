import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import yaml from 'js-yaml'
import { describe, expect, it } from 'vitest'
import one from './index'

// The registry application `one` names, as provisioned.
const app = yaml
  .load(
    readFileSync(
      join(__dirname, '../../../../../../server/provision/101_applications/0200_applications.yaml'),
      'utf8',
    ),
  )
  .applications.find(a => (a.unify?.url || '').replace(/^\/|\/$/g, '') === 'one')

describe('one section', () => {
  it('is gated on the provisioned `one` application', () => {
    expect(one.app).toBe('one/')
    expect(app).toBeDefined()
  })

  // The section exists for an installation carrying the application over from
  // before the webapps merged. A new one gets home as its landing page, so the
  // catalog ships switched off and out of the menu, and turning it on is a
  // deliberate act.
  it('ships disabled and unlisted', () => {
    expect(app.enabled).toBe(false)
    expect(app.unify.listed).toBe(false)
  })
})

import { expect } from 'chai'
import { canAccessC311Route, c311DataState, c311StateForError, hasC311Capability, hasC311Route, hasC311Scope } from './c311'
import { sanitizeC311Draft } from '../mixins/c311-dirty-guard.js'
import { createDefaultFixtureSet } from '../../../../lib/js/src/311/fixtures'
import { APPLICATION_ROLES } from '../../../../lib/js/src/311/enums'

const session = {
  authenticated: true,
  actor: {
    actor_id: 'actor-fixture-001',
    capabilities: ['staff_request_queue'],
    scopes: ['service_requests.write'],
    available_routes: ['staff_request_queue'],
  },
}

describe('C311 shared frontend helpers', () => {
  it('maps the required eight visible data states', () => {
    expect(c311DataState({ loading: true })).to.equal('loading')
    expect(c311DataState({ items: [] })).to.equal('empty')
    expect(c311DataState({ items: [{ id: '1' }] })).to.equal('populated')
    expect(c311StateForError({ status: 403 })).to.equal('forbidden')
    expect(c311StateForError({ status: 404 })).to.equal('not-found')
    expect(c311StateForError({ status: 422 })).to.equal('validation-error')
    expect(c311StateForError({ retryable: true })).to.equal('retryable-error')
    expect(c311StateForError({ status: 500, retryable: true })).to.equal('terminal-error')
    expect(c311StateForError({ status: 503 })).to.equal('retryable-error')
    expect(c311StateForError({ code: 'OPERATION_FAILED' })).to.equal('terminal-error')
    expect(c311StateForError({ status: 401 })).to.equal('forbidden')
  })

  it('checks capabilities, scopes and authentication before entering a route', () => {
    expect(hasC311Capability(session, 'staff_request_queue')).to.equal(true)
    expect(hasC311Scope(session, 'service_requests.write')).to.equal(true)
    expect(hasC311Route(session, 'staff_request_queue')).to.equal(true)
    expect(canAccessC311Route(session, { requiresAuth: true, route: 'staff_request_queue', capabilities: ['staff_request_queue'], scopes: ['service_requests.write'] })).to.equal(true)
    expect(canAccessC311Route({ authenticated: true, actor: { capabilities: ['portal_my_requests'], scopes: ['service_requests.write'], available_routes: [] } }, { public: true, requiresAuth: true, route: 'portal_my_requests', capabilities: ['portal_my_requests'], scopes: ['service_requests.write'] })).to.equal(false)
    expect(canAccessC311Route(session, { requiresAuth: true, capabilities: ['admin_branding_get'] })).to.equal(false)
    expect(canAccessC311Route(session, { requiresAuth: true, route: 'report_catalogue' })).to.equal(false)
    expect(canAccessC311Route(session, { requiresAuth: true, scopes: ['workflow.execute'] })).to.equal(false)
    expect(canAccessC311Route({ authenticated: false }, { requiresAuth: true })).to.equal(false)
  })

  it('removes sensitive values before persisting an unsaved draft', () => {
    expect(sanitizeC311Draft({ summary: 'draft', password: 'secret', nested: { accessToken: 'token', value: 'keep' } })).to.deep.equal({
      summary: 'draft',
      nested: { value: 'keep' },
    })
  })

  it('keeps every FE-00 role from entering its denied route, capability or scope', () => {
    const fixtures = createDefaultFixtureSet()
    for (const role of APPLICATION_ROLES) {
      const fixture = fixtures.role_fixtures[role]
      expect(fixture, role).to.exist
      const session = fixture.session
      expect(hasC311Capability(session, fixture.denied_capability)).to.equal(false, role)
      expect(hasC311Route(session, fixture.denied_route)).to.equal(false, role)
      expect(hasC311Scope(session, fixture.denied_scope)).to.equal(false, role)
      expect(canAccessC311Route(session, {
        requiresAuth: true,
        route: fixture.denied_route,
        capabilities: [fixture.denied_capability],
        scopes: [fixture.denied_scope],
      })).to.equal(false, role)
    }
  })
})

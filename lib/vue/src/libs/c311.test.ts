import { expect } from 'chai'
import { canAccessC311Route, c311DataState, c311StateForError, hasC311Capability } from './c311'
import { sanitizeC311Draft } from '../mixins/c311-dirty-guard.js'

const session = {
  authenticated: true,
  actor: {
    actor_id: 'actor-fixture-001',
    capabilities: ['staff_request_queue'],
    scopes: ['staff_request_queue'],
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
  })

  it('checks capabilities, scopes and authentication before entering a route', () => {
    expect(hasC311Capability(session, 'staff_request_queue')).to.equal(true)
    expect(canAccessC311Route(session, { requiresAuth: true, capabilities: ['staff_request_queue'] })).to.equal(true)
    expect(canAccessC311Route(session, { requiresAuth: true, capabilities: ['admin_branding_get'] })).to.equal(false)
    expect(canAccessC311Route({ authenticated: false }, { requiresAuth: true })).to.equal(false)
  })

  it('removes sensitive values before persisting an unsaved draft', () => {
    expect(sanitizeC311Draft({ summary: 'draft', password: 'secret', nested: { accessToken: 'token', value: 'keep' } })).to.deep.equal({
      summary: 'draft',
      nested: { value: 'keep' },
    })
  })
})

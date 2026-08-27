import { expect } from 'chai'
import { installC311RouterGuards } from './c311-router'

describe('C311 router guards', () => {
  function setup (session: any, expired = false, error: any = null) {
    let guard: any
    const router = {
      beforeEach (handler: any) { guard = handler },
      afterEach () {},
    }
    const runtime = {
      async loadSession () { return session },
      isExpired () { return expired },
      error,
    }
    installC311RouterGuards(router, () => runtime as any)
    return guard
  }

  it('redirects anonymous and expired sessions to the shared 401 route', async () => {
    for (const runtime of [setup({ authenticated: false }), setup({ authenticated: true }, true)]) {
      let redirect: any
      await runtime({
        path: '/c311/staff',
        fullPath: '/c311/staff?view=mine',
        matched: [{ meta: { c311: { requiresAuth: true } } }],
      }, {}, value => { redirect = value })

      expect(redirect.name).to.equal('c311.unauthorized')
      expect(redirect.query.returnTo).to.equal(encodeURIComponent('/c311/staff?view=mine'))
    }
  })

  it('redirects authenticated actors without capability to the shared 403 route', async () => {
    const guard = setup({ authenticated: true, actor: { capabilities: [] } })
    let redirect: any
    await guard({
      path: '/c311/staff/reports',
      fullPath: '/c311/staff/reports',
      matched: [{ meta: { c311: { requiresAuth: true, capabilities: ['report_catalogue'] } } }],
    }, {}, value => { redirect = value })

    expect(redirect.name).to.equal('c311.forbidden')
    expect(redirect.query.returnTo).to.equal(encodeURIComponent('/c311/staff/reports'))
  })

  it('does not turn a session service outage into a false 401 redirect', async () => {
    const guard = setup(null, false, { status: 503, code: 'TEMPORARILY_UNAVAILABLE' })
    let nextCalled = false
    await guard({
      path: '/c311/staff',
      fullPath: '/c311/staff',
      matched: [{ meta: { c311: { requiresAuth: true } } }],
    }, {}, () => { nextCalled = true })

    expect(nextCalled).to.equal(true)
  })
})

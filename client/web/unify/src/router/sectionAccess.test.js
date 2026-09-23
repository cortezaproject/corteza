import { describe, expect, it, vi } from 'vitest'
import { makeSectionAccessGuard } from './sectionAccess'

const SECTIONS = {
  home: { id: 'home', app: null },
  admin: { id: 'admin', app: 'admin/' },
  compose: { id: 'compose', app: 'compose/' },
  // Gated on the application named in the route, not on one of its own.
  app: { id: 'app', app: null, perApp: true },
  // A section that forgot to declare one.
  rogue: { id: 'rogue' },
}

// The applications the route-named lookup can find, by id.
const CUSTOM_APPS = {
  mine: { unify: { kind: 'custom' }, canAccessApplication: true },
  theirs: { unify: { kind: 'custom' }, canAccessApplication: false },
  binned: {
    unify: { kind: 'custom' },
    canAccessApplication: true,
    deletedAt: '2026-09-22T10:00:00Z',
  },
  section: { unify: { kind: 'section' }, canAccessApplication: true },
}

function guardWith(allowed, { ready = vi.fn().mockResolvedValue(undefined) } = {}) {
  const canAccessApp = vi.fn(url => allowed.includes(url))
  const findByID = vi.fn(async id => {
    // The list is read-filtered, so one the user cannot read is not there.
    if (!CUSTOM_APPS[id]) throw new Error('not found')
    return CUSTOM_APPS[id]
  })
  const applications = { ready, canAccessApp, findByID }
  const guard = makeSectionAccessGuard({
    useApplications: () => applications,
    sectionById: id => SECTIONS[id] || null,
  })
  return { guard, applications }
}

const to = (section, path = '/' + section) => ({ meta: { section }, path, params: {} })

const toApp = applicationID => ({
  meta: { section: 'app' },
  path: '/app/' + applicationID,
  params: { applicationID },
})

describe('section access guard', () => {
  it('admits a section the user may access', async () => {
    const { guard } = guardWith(['admin/'])
    expect(await guard(to('admin'))).toBe(true)
  })

  it('refuses one the user may not, and says which', async () => {
    const { guard } = guardWith(['compose/'])
    expect(await guard(to('admin'))).toEqual({ name: 'home', query: { denied: 'admin' } })
  })

  it('admits the ungated section without consulting the list at all', async () => {
    const { guard, applications } = guardWith([])
    expect(await guard(to('home', '/'))).toBe(true)
    // Home is where a refusal lands, so it must not depend on a list that may
    // have failed to load — that would be a redirect loop.
    expect(applications.ready).not.toHaveBeenCalled()
  })

  it('refuses a section that declares no application', async () => {
    const { guard } = guardWith(['admin/', 'compose/'])
    expect(await guard(to('rogue'))).toEqual({ name: 'home', query: { denied: 'rogue' } })
  })

  it('refuses a route belonging to no section', async () => {
    const { guard } = guardWith(['admin/'])
    expect(await guard(to(undefined, '/stray'))).toEqual({
      name: 'home',
      query: { denied: '/stray' },
    })
  })

  it('admits a custom application the user may access', async () => {
    const { guard, applications } = guardWith([])
    expect(await guard(toApp('mine'))).toBe(true)
    // The section names no application of its own, so its url is never asked.
    expect(applications.canAccessApp).not.toHaveBeenCalled()
  })

  it('refuses a custom application the user may not access', async () => {
    const { guard } = guardWith([])
    expect(await guard(toApp('theirs'))).toEqual({ name: 'home', query: { denied: 'app' } })
  })

  it('refuses an application that is not a custom one', async () => {
    // A registry entry pointing at a section has no source to render.
    const { guard } = guardWith([])
    expect(await guard(toApp('section'))).toEqual({ name: 'home', query: { denied: 'app' } })
  })

  // Deleting one takes its tile out of the menu; the link has to stop working
  // too, or it stays open to whoever kept it.
  it('refuses a deleted application', async () => {
    const { guard } = guardWith([])
    expect(await guard(toApp('binned'))).toEqual({ name: 'home', query: { denied: 'app' } })
  })

  it('refuses an application the user cannot even read', async () => {
    const { guard } = guardWith([])
    expect(await guard(toApp('unknown'))).toEqual({ name: 'home', query: { denied: 'app' } })
  })

  it('waits for the application list before deciding', async () => {
    const order = []
    const ready = vi.fn(async () => {
      order.push('ready')
    })
    const { guard, applications } = guardWith(['admin/'], { ready })
    applications.canAccessApp = vi.fn(() => {
      order.push('check')
      return true
    })
    await guard(to('admin'))
    expect(order).toEqual(['ready', 'check'])
  })
})

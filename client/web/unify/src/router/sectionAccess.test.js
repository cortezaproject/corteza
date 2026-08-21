import { describe, expect, it, vi } from 'vitest'
import { makeSectionAccessGuard } from './sectionAccess'

const SECTIONS = {
  home: { id: 'home', app: null },
  admin: { id: 'admin', app: 'admin/' },
  compose: { id: 'compose', app: 'compose/' },
  // A section that forgot to declare one.
  rogue: { id: 'rogue' },
}

function guardWith(allowed, { ready = vi.fn().mockResolvedValue(undefined) } = {}) {
  const canAccessApp = vi.fn(url => allowed.includes(url))
  const applications = { ready, canAccessApp }
  const guard = makeSectionAccessGuard({
    useApplications: () => applications,
    sectionById: id => SECTIONS[id] || null,
  })
  return { guard, applications }
}

const to = (section, path = '/' + section) => ({ meta: { section }, path })

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

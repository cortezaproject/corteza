import { START_LOCATION } from 'vue-router'
import { describe, expect, it, vi } from 'vitest'
import { makeHomeEntryGuard, resolveHome } from './homeEntry'

const APPS = [
  { applicationID: '1', name: 'Projects', enabled: true, unify: { url: 'project/', home: true } },
  { applicationID: '2', name: 'TAQ', enabled: true, unify: { url: 'taq/' } },
  { applicationID: '3', name: 'Off', enabled: false, unify: { url: 'workflow/' } },
  { applicationID: '4', name: 'Google', enabled: true, unify: { url: 'www.google.com' } },
]

const openable = app => app.enabled

function guardWith({ ownID = '', apps = APPS } = {}) {
  const leave = vi.fn()
  const guard = makeHomeEntryGuard({
    useApplications: () => ({ ready: async () => {}, apps }),
    ownHomeID: () => ownID,
    openable,
    hrefOf: app => {
      const local = !app.unify.url.includes('.')
      return { href: local ? '/' + app.unify.url : 'https://' + app.unify.url, local }
    },
    leave,
  })
  return { guard, leave }
}

const home = (query = {}) => ({ name: 'home', path: '/', query })

describe('resolveHome', () => {
  it('prefers the own pick over the instance home', () => {
    expect(resolveHome({ ownID: '2', apps: APPS, openable }).app.name).toBe('TAQ')
  })

  it('falls back to the instance home when there is no own pick', () => {
    expect(resolveHome({ ownID: '', apps: APPS, openable }).app.name).toBe('Projects')
    expect(resolveHome({ ownID: '0', apps: APPS, openable }).app.name).toBe('Projects')
  })

  it('reports an own pick that cannot be opened rather than using the instance home', () => {
    expect(resolveHome({ ownID: '3', apps: APPS, openable })).toEqual({ unavailable: 'Off' })
    expect(resolveHome({ ownID: '99', apps: APPS, openable })).toEqual({ unavailable: '' })
  })

  it('is empty when nothing is set', () => {
    expect(resolveHome({ ownID: '', apps: APPS.slice(1), openable })).toEqual({})
  })
})

describe('makeHomeEntryGuard', () => {
  it('redirects the first navigation to `/` into the home application', async () => {
    const { guard } = guardWith({ ownID: '2' })
    expect(await guard(home(), START_LOCATION)).toBe('/taq/')
  })

  it('leaves an in-app trip to `/` on the launcher', async () => {
    const { guard } = guardWith({ ownID: '2' })
    expect(await guard(home(), { name: 'taq' })).toBe(true)
  })

  it('leaves a bounce carrying a query on the launcher', async () => {
    const { guard } = guardWith({ ownID: '2' })
    expect(await guard(home({ denied: 'admin' }), START_LOCATION)).toBe(true)
  })

  it('ignores the first navigation to any other route', async () => {
    const { guard } = guardWith({ ownID: '2' })
    expect(await guard({ name: 'taq', path: '/taq', query: {} }, START_LOCATION)).toBe(true)
  })

  it('shows the launcher with a warning when the home application cannot be opened', async () => {
    const { guard } = guardWith({ ownID: '3' })
    expect(await guard(home(), START_LOCATION)).toEqual({
      name: 'home',
      query: { homeUnavailable: 'Off' },
    })
  })

  it('stays on the launcher when no home application is set', async () => {
    const { guard } = guardWith({ apps: APPS.slice(1) })
    expect(await guard(home(), START_LOCATION)).toBe(true)
  })

  it('leaves the shell for a home application elsewhere', async () => {
    const { guard, leave } = guardWith({ ownID: '4' })
    expect(await guard(home(), START_LOCATION)).toBe(false)
    expect(leave).toHaveBeenCalledWith('https://www.google.com')
  })
})

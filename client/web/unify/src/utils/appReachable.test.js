import { describe, expect, it, vi } from 'vitest'
import { sectionAllows } from '@/router/sectionAccess'

// useAppReachable is the router + store wiring around sectionAllows; the rule
// itself is what matters, and what has to stay identical to the gate's.
const SECTIONS = {
  home: { id: 'home', app: null },
  admin: { id: 'admin', app: 'admin/' },
}

// Mirrors useAppReachable with the router and store stood in for.
function reachable(app, { allowed, resolve }) {
  const applications = { canAccessApp: url => allowed.includes(url) }
  const url = String(app?.unify?.url || '').replace(/^\/|\/$/g, '')
  if (!url) return false
  return sectionAllows(SECTIONS[resolve('/' + url)] || null, applications)
}

// /admin/system/labels resolves to the admin section, like the real router.
const resolve = path => (path === '/' ? 'home' : path.startsWith('/admin') ? 'admin' : undefined)

describe('app menu reachability', () => {
  it('offers an application whose section the user may enter', () => {
    const app = { unify: { url: 'admin/' } }
    expect(reachable(app, { allowed: ['admin/'], resolve })).toBe(true)
  })

  it('withholds one whose section the user may not enter', () => {
    const app = { unify: { url: 'admin/' } }
    expect(reachable(app, { allowed: [], resolve })).toBe(false)
  })

  it('judges a sub-path entry by the section it lands in, not its own access', () => {
    // The trap: an application registered inside another section. Granting it
    // must not offer a tile that the gate then bounces.
    const labels = { unify: { url: 'admin/system/labels' } }
    expect(reachable(labels, { allowed: [], resolve })).toBe(false)
    expect(reachable(labels, { allowed: ['admin/'], resolve })).toBe(true)
  })

  it('withholds an application pointing at no section at all', () => {
    expect(reachable({ unify: { url: 'one/' } }, { allowed: ['admin/'], resolve })).toBe(false)
  })

  it('withholds an application with no url', () => {
    expect(reachable({ unify: {} }, { allowed: ['admin/'], resolve })).toBe(false)
  })

  it('always offers the ungated section', () => {
    expect(reachable({ unify: { url: 'home/' } }, { allowed: [], resolve: () => 'home' })).toBe(true)
  })
})

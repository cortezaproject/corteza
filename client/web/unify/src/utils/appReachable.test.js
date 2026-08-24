import { describe, expect, it } from 'vitest'
import { appReachable } from './appReachable'

// The rule itself, with the router and store stood in for. It has to stay
// identical to the section gate's for a url a section serves, and answer from
// the application's own grant for one no section does.
const SECTIONS = {
  home: { id: 'home', app: null },
  admin: { id: 'admin', app: 'admin/' },
}

// /admin/system/labels resolves to the admin section, like the real router.
const sectionFor = path =>
  path === '/' ? SECTIONS.home : path.startsWith('/admin') ? SECTIONS.admin : null

const reachable = (app, allowed = []) =>
  appReachable(app, {
    sectionFor,
    applications: { canAccessApp: url => allowed.includes(url) },
  })

describe('app menu reachability', () => {
  it('offers an application whose section the user may enter', () => {
    expect(reachable({ unify: { url: 'admin/' } }, ['admin/'])).toBe(true)
  })

  it('withholds one whose section the user may not enter', () => {
    expect(reachable({ unify: { url: 'admin/' }, canAccessApplication: true })).toBe(false)
  })

  it('judges a sub-path entry by the section it lands in, not its own access', () => {
    // The trap: an application registered inside another section. Granting it
    // must not offer a tile that the gate then bounces.
    const labels = { unify: { url: 'admin/system/labels' }, canAccessApplication: true }
    expect(reachable(labels)).toBe(false)
    expect(reachable(labels, ['admin/'])).toBe(true)
  })

  it('judges an application no section serves by its own access', () => {
    // A custom entry pointing somewhere else entirely: no section can speak for
    // it, and withholding it outright greys out an application that works.
    for (const url of ['www.google.com', 'https://www.google.com', '//www.google.com']) {
      expect(reachable({ unify: { url }, canAccessApplication: true })).toBe(true)
      expect(reachable({ unify: { url }, canAccessApplication: false })).toBe(false)
    }
  })

  it('judges a path this shell does not serve by its own access', () => {
    // An application served outside the shell, still reached by a full load.
    const other = { unify: { url: 'unregistered/' } }
    expect(reachable({ ...other, canAccessApplication: true }, ['admin/'])).toBe(true)
    expect(reachable({ ...other, canAccessApplication: false }, ['admin/'])).toBe(false)
  })

  it('withholds an application with no url', () => {
    expect(reachable({ unify: {} }, ['admin/'])).toBe(false)
    expect(reachable({ unify: { url: '/' }, canAccessApplication: true }, ['admin/'])).toBe(false)
  })

  it('always offers the ungated section', () => {
    expect(
      appReachable(
        { unify: { url: 'home/' } },
        { sectionFor: () => SECTIONS.home, applications: { canAccessApp: () => false } },
      ),
    ).toBe(true)
  })
})

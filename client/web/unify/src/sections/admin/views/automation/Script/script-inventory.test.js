import { describe, it, expect } from 'vitest'
import {
  corredorBanner,
  groupScripts,
  matchesKindFilter,
  relativeTime,
  scriptBundle,
  scriptExtension,
  scriptKind,
  triggerRows,
} from './script-inventory'

// The names the fixture extension is listed under
const serverScript = { name: '/server-scripts/agent-sandbox/ContactActivate.js:default' }
const looseServerScript = { name: '/server-scripts/SinkHello.js:default' }
const clientScript = { name: '/client-scripts/compose/agent-sandbox/ContactGreet.js:default' }
const looseClientScript = { name: '/client-scripts/admin/DashboardHello.js:default' }

describe('scriptKind', () => {
  it('reads the kind off the path a script is listed under', () => {
    expect(scriptKind(serverScript)).toBe('server')
    expect(scriptKind(clientScript)).toBe('client')
  })

  it('calls anything it cannot place a server script', () => {
    expect(scriptKind({})).toBe('server')
    expect(scriptKind()).toBe('server')
  })
})

describe('scriptBundle', () => {
  it('reads the bundle off a client script path', () => {
    expect(scriptBundle(clientScript)).toBe('compose')
    expect(scriptBundle(looseClientScript)).toBe('admin')
  })

  it('prefers the bundle the server sent', () => {
    expect(scriptBundle({ ...clientScript, bundle: 'unify' })).toBe('unify')
  })

  it('leaves a server script without one', () => {
    expect(scriptBundle(serverScript)).toBe('')
  })
})

describe('scriptExtension', () => {
  it('takes the folder under the kind for a server script', () => {
    expect(scriptExtension(serverScript)).toBe('agent-sandbox')
  })

  it('takes the folder under the bundle for a client script', () => {
    expect(scriptExtension(clientScript)).toBe('agent-sandbox')
  })

  it('leaves a script sitting directly in the kind or bundle without one', () => {
    expect(scriptExtension(looseServerScript)).toBe('')
    expect(scriptExtension(looseClientScript)).toBe('')
  })
})

describe('triggerRows', () => {
  it('names what fires a trigger and what narrows it', () => {
    expect(
      triggerRows({
        triggers: [
          {
            eventTypes: ['beforeCreate'],
            resourceTypes: ['compose:record'],
            constraints: [
              { name: 'module', value: ['agent-contact'] },
              { name: 'namespace', value: ['agent-sandbox'] },
            ],
          },
        ],
      }),
    ).toEqual([
      {
        label: 'beforeCreate · compose:record',
        constraints: ['module = agent-contact', 'namespace = agent-sandbox'],
      },
    ])
  })

  it('has no rows for a script with no trigger', () => {
    expect(triggerRows({ triggers: null })).toEqual([])
    expect(triggerRows()).toEqual([])
  })
})

describe('groupScripts', () => {
  it('groups by extension, alphabetically, loose scripts last', () => {
    const grouped = groupScripts([
      looseServerScript,
      clientScript,
      serverScript,
      { name: '/server-scripts/other-extension/Thing.js:default' },
    ])

    expect(grouped.map(g => g.extension)).toEqual(['agent-sandbox', 'other-extension', ''])
    expect(grouped[0].items.map(s => s.name)).toEqual([clientScript.name, serverScript.name])
  })

  it('has no group at all for no scripts', () => {
    expect(groupScripts([])).toEqual([])
  })
})

describe('matchesKindFilter', () => {
  it('shows every kind while neither or both toggles are on', () => {
    expect(matchesKindFilter('server', { server: false, client: false })).toBe(true)
    expect(matchesKindFilter('client', { server: false, client: false })).toBe(true)
    expect(matchesKindFilter('server', { server: true, client: true })).toBe(true)
    expect(matchesKindFilter('client', { server: true, client: true })).toBe(true)
  })

  it('narrows to the one kind that is on', () => {
    expect(matchesKindFilter('server', { server: true, client: false })).toBe(true)
    expect(matchesKindFilter('client', { server: true, client: false })).toBe(false)
    expect(matchesKindFilter('client', { server: false, client: true })).toBe(true)
    expect(matchesKindFilter('server', { server: false, client: true })).toBe(false)
  })
})

describe('corredorBanner', () => {
  it('says nothing when the server reported nothing', () => {
    expect(corredorBanner({})).toBeNull()
    expect(corredorBanner()).toBeNull()
    expect(corredorBanner({ connected: true })).toBeNull()
  })

  it('reports Corredor turned off as information, not as a fault', () => {
    expect(corredorBanner({ enabled: false, connected: false })).toEqual({
      state: 'disabled',
      severity: 'info',
    })
  })

  it('warns when Corredor is on but out of reach', () => {
    expect(corredorBanner({ enabled: true, connected: false })).toEqual({
      state: 'unreachable',
      severity: 'warn',
    })
  })

  it('confirms a working connection', () => {
    expect(corredorBanner({ enabled: true, connected: true })).toEqual({
      state: 'connected',
      severity: 'success',
    })
  })

  it('says nothing while enabled is known and connected is not', () => {
    expect(corredorBanner({ enabled: true })).toBeNull()
  })
})

describe('relativeTime', () => {
  it('counts back in the coarsest unit that fits', () => {
    expect(relativeTime(new Date(Date.now() - 10 * 1000), 'en')).toMatch(/second/)
    expect(relativeTime(new Date(Date.now() - 5 * 60 * 1000), 'en')).toMatch(/5 minutes ago/)
    expect(relativeTime(new Date(Date.now() - 3 * 3600 * 1000), 'en')).toMatch(/3 hours ago/)
    expect(relativeTime(new Date(Date.now() - 2 * 86400 * 1000), 'en')).toMatch(/2 days ago/)
  })

  it('has nothing to say about a timestamp it cannot read', () => {
    expect(relativeTime(null, 'en')).toBe('')
    expect(relativeTime('not a date', 'en')).toBe('')
  })
})

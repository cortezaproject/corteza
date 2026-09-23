import { describe, it, expect } from 'vitest'
import {
  ANY_EXTENSION,
  corredorBanner,
  groupScripts,
  matchesExtensionFilter,
  matchesKindFilter,
  relativeTime,
  scriptBundle,
  scriptExtension,
  scriptKind,
  sortScripts,
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

describe('matchesExtensionFilter', () => {
  it('shows every script while no extension is chosen', () => {
    expect(matchesExtensionFilter(serverScript, ANY_EXTENSION)).toBe(true)
    expect(matchesExtensionFilter(looseClientScript, ANY_EXTENSION)).toBe(true)
  })

  it('narrows to the scripts of the chosen extension', () => {
    expect(matchesExtensionFilter(serverScript, 'agent-sandbox')).toBe(true)
    expect(matchesExtensionFilter(clientScript, 'agent-sandbox')).toBe(true)
    expect(matchesExtensionFilter(looseServerScript, 'agent-sandbox')).toBe(false)
  })

  it('treats the empty name as the scripts outside any extension', () => {
    expect(matchesExtensionFilter(looseServerScript, '')).toBe(true)
    expect(matchesExtensionFilter(looseClientScript, '')).toBe(true)
    expect(matchesExtensionFilter(serverScript, '')).toBe(false)
  })
})

describe('sortScripts', () => {
  const labelled = [
    { ...serverScript, label: 'Activate contact (server)' },
    { ...clientScript, label: 'Greet contact (client)' },
    { ...looseServerScript, label: 'Sink: hello (server)' },
  ]

  it('orders by the label a row shows', () => {
    expect(sortScripts(labelled, 'name', false).map(s => s.label)).toEqual([
      'Activate contact (server)',
      'Greet contact (client)',
      'Sink: hello (server)',
    ])
  })

  it('falls back to the machine name for a script with no label', () => {
    const unlabelled = [
      { name: '/server-scripts/Zulu.js:default' },
      { name: '/server-scripts/Alpha.js:default' },
    ]

    expect(sortScripts(unlabelled, 'name', false).map(s => s.name)).toEqual([
      '/server-scripts/Alpha.js:default',
      '/server-scripts/Zulu.js:default',
    ])
  })

  it('reverses on a descending sort', () => {
    expect(sortScripts(labelled, 'name', true).map(s => s.label)).toEqual([
      'Sink: hello (server)',
      'Greet contact (client)',
      'Activate contact (server)',
    ])
  })

  it('orders by extension and by kind', () => {
    // Two scripts of the same extension keep the order they came in.
    expect(sortScripts(labelled, 'extension', false).map(s => s.name)).toEqual([
      looseServerScript.name,
      serverScript.name,
      clientScript.name,
    ])

    expect(sortScripts(labelled, 'kind', false).map(s => s.name)).toEqual([
      clientScript.name,
      serverScript.name,
      looseServerScript.name,
    ])
  })

  it('orders by the timestamp the last-change column shows', () => {
    const stamped = [
      { name: 'b', updatedAt: '2026-02-01T00:00:00Z' },
      { name: 'a', updatedAt: '2026-01-01T00:00:00Z' },
      { name: 'c' },
    ]

    expect(sortScripts(stamped, 'changedAt', false).map(s => s.name)).toEqual(['c', 'a', 'b'])
  })

  it('leaves a column that carries no order alone, and never sorts in place', () => {
    const given = [...labelled]

    expect(sortScripts(given, 'triggers', false).map(s => s.label)).toEqual(
      labelled.map(s => s.label),
    )
    expect(sortScripts(given, 'name', true)).not.toBe(given)
    expect(given.map(s => s.label)).toEqual(labelled.map(s => s.label))
  })

  it('has nothing to sort for no scripts', () => {
    expect(sortScripts(null, 'name', false)).toEqual([])
    expect(sortScripts([], 'name', false)).toEqual([])
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

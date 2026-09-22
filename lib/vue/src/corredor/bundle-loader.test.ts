import { describe, it, expect, vi, beforeEach } from 'vitest'
import { eventbus } from '@planetcrust/human-js'
import { loadClientScripts, registerServerScripts, type ScriptEvent } from './bundle-loader'
import { UIHooks } from './ui-hooks'

const bundleSource = `
  window.testClientScripts = {
    scripts: [
      {
        name: '/client-scripts/test/Greet.js:default',
        label: 'Greet',
        triggers: [
          {
            eventTypes: ['onManual'],
            resourceTypes: ['compose:record'],
            uiProps: [{ name: 'app', value: 'compose' }],
          },
        ],
        exec: function (args, ctx) {
          window.execCalls.push({ args: args, ctx: ctx })
        },
      },
      {
        name: '/client-scripts/test/Refuse.js:default',
        label: 'Refuse',
        triggers: [
          { eventTypes: ['beforeFormSubmit'], resourceTypes: ['ui:compose:record-page'] },
        ],
        exec: function () {
          return false
        },
      },
    ],
  }
`

interface ExecCall {
  args: { $record?: unknown }
  ctx: unknown
}

declare global {
  interface Window {
    execCalls: ExecCall[]
    testClientScripts?: unknown
  }
}

function systemAPI(data: unknown) {
  return {
    automationBundleEndpoint: ({ bundle, type, ext }: Record<string, string>) =>
      `/automation/${bundle}-${type}.${ext}`,
    api: () => ({ get: vi.fn(() => Promise.resolve({ data })) }),
  }
}

const ctx = { withArgs: (args: unknown) => ({ withArgs: args }) as never }

beforeEach(() => {
  window.execCalls = []
  delete window.testClientScripts
})

describe('loadClientScripts', () => {
  it('evals the bundle and registers every trigger it declares', async () => {
    const scriptBus = new eventbus.EventBus()
    const scripts = await loadClientScripts({
      systemAPI: systemAPI(bundleSource),
      scriptBus,
      bundle: 'test',
      ctx,
    })

    expect(scripts.map(s => s.name)).toEqual([
      '/client-scripts/test/Greet.js:default',
      '/client-scripts/test/Refuse.js:default',
    ])

    // The script name is what the bus matches an explicit dispatch on
    expect(scripts[0].triggers?.[0].scriptName).toBe('/client-scripts/test/Greet.js:default')
  })

  it('runs the script with the event arguments and the context', async () => {
    const scriptBus = new eventbus.EventBus()
    const record = { recordID: '1' }

    await loadClientScripts({ systemAPI: systemAPI(bundleSource), scriptBus, bundle: 'test', ctx })
    await scriptBus.Dispatch(
      { resourceType: 'compose:record', eventType: 'onManual', args: { record } },
      '/client-scripts/test/Greet.js:default',
    )

    expect(window.execCalls).toHaveLength(1)
    // ArgsProxy hands the very object the event carries to the script, so what
    // the script writes to it reaches whoever dispatched the event
    expect(window.execCalls[0].args.$record).toBe(record)
    expect(window.execCalls[0].ctx).toEqual({ withArgs: window.execCalls[0].args })
  })

  it('aborts the dispatch when the script returns false', async () => {
    const scriptBus = new eventbus.EventBus()

    await loadClientScripts({ systemAPI: systemAPI(bundleSource), scriptBus, bundle: 'test', ctx })

    await expect(
      scriptBus.Dispatch({ resourceType: 'ui:compose:record-page', eventType: 'beforeFormSubmit' }),
    ).rejects.toThrow('Aborted')
  })

  it('resolves empty on an empty bundle', async () => {
    const scriptBus = new eventbus.EventBus()

    expect(
      await loadClientScripts({ systemAPI: systemAPI(''), scriptBus, bundle: 'test', ctx }),
    ).toEqual([])
  })

  it('resolves empty when the bundle defines no global', async () => {
    const scriptBus = new eventbus.EventBus()

    expect(
      await loadClientScripts({
        systemAPI: systemAPI('var somethingElse = 1'),
        scriptBus,
        bundle: 'test',
        ctx,
      }),
    ).toEqual([])
  })

  it('resolves empty rather than throwing when the bundle cannot be fetched', async () => {
    const scriptBus = new eventbus.EventBus()
    const api = {
      automationBundleEndpoint: () => '/automation/test-client-scripts.js',
      api: () => ({ get: () => Promise.reject(new Error('404')) }),
    }

    expect(await loadClientScripts({ systemAPI: api, scriptBus, bundle: 'test', ctx })).toEqual([])
  })
})

describe('registerServerScripts', () => {
  const scripts = [
    {
      name: '/server-scripts/Activate.js:default',
      label: 'Activate',
      triggers: [
        {
          eventTypes: ['onManual'],
          resourceTypes: ['compose:record'],
          uiProps: [{ name: 'app', value: 'compose' }],
        },
      ],
    },
    {
      name: '/server-scripts/BeforeCreate.js:default',
      label: 'Before create',
      triggers: [{ eventTypes: ['beforeCreate'], resourceTypes: ['compose:record'] }],
    },
    {
      name: '/client-scripts/compose/Greet.js:default',
      label: 'Greet',
      triggers: [
        {
          eventTypes: ['onManual'],
          resourceTypes: ['compose:record'],
          uiProps: [{ name: 'app', value: 'compose' }],
        },
      ],
    },
  ]

  it('routes a manual server script to the handler with its name', async () => {
    const scriptBus = new eventbus.EventBus()
    const uiHooks = new UIHooks({ apps: ['compose'] })
    const handler = vi.fn(() => Promise.resolve(undefined))

    registerServerScripts({ scriptBus, uiHooks, scripts, handler })

    const ev: ScriptEvent = {
      resourceType: 'compose:record',
      eventType: 'onManual',
      args: { record: { recordID: '1' } },
    }
    await scriptBus.Dispatch(ev, '/server-scripts/Activate.js:default')

    expect(handler).toHaveBeenCalledTimes(1)
    expect(handler.mock.calls[0]).toEqual([ev, '/server-scripts/Activate.js:default'])
  })

  it('leaves client scripts and implicit triggers to the bundle and the API', async () => {
    const scriptBus = new eventbus.EventBus()
    const uiHooks = new UIHooks({ apps: ['compose'] })
    const handler = vi.fn(() => Promise.resolve(undefined))

    registerServerScripts({ scriptBus, uiHooks, scripts, handler })

    await scriptBus.Dispatch(
      { resourceType: 'compose:record', eventType: 'onManual' },
      '/client-scripts/compose/Greet.js:default',
    )
    await scriptBus.Dispatch({ resourceType: 'compose:record', eventType: 'beforeCreate' })

    expect(handler).not.toHaveBeenCalled()
  })

  it('offers every manual script as a button, client ones included', () => {
    const scriptBus = new eventbus.EventBus()
    const uiHooks = new UIHooks({ apps: ['compose'] })

    registerServerScripts({
      scriptBus,
      uiHooks,
      scripts,
      handler: () => Promise.resolve(undefined),
    })

    expect(
      uiHooks.Find('compose:record', undefined, undefined, 'compose').map(b => b.script),
    ).toEqual(['/client-scripts/compose/Greet.js:default', '/server-scripts/Activate.js:default'])
  })

  it('does nothing with no scripts', () => {
    const uiHooks = new UIHooks({ apps: ['compose'] })
    const register = vi.fn()

    registerServerScripts({
      scriptBus: { Register: register },
      uiHooks,
      scripts: [],
      handler: () => Promise.resolve(undefined),
    })

    expect(register).not.toHaveBeenCalled()
  })
})

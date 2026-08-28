import Vue from 'vue'
import { expect } from 'chai'
import { C311Runtime, createC311Provider, default as c311Plugin } from './c311'

type TestProvider = {
  getSession: () => Promise<any>
}

describe('C311 provider selection and runtime', () => {
  let originalWindow: any

  beforeEach(() => {
    originalWindow = (globalThis as any).window
    ;(globalThis as any).window = {}
  })

  afterEach(() => {
    ;(globalThis as any).window = originalWindow
    delete (Vue.prototype as any).$C311
  })

  it('prefers an explicitly configured provider object or constructor', () => {
    const objectProvider: TestProvider = { getSession: async () => ({ authenticated: true }) }
    expect(createC311Provider({})).to.equal(undefined)

    ;(globalThis as any).window.C311Provider = objectProvider
    expect(createC311Provider({})).to.equal(objectProvider)

    class ConfiguredProvider implements TestProvider {
      async getSession (): Promise<any> { return { authenticated: true } }
    }
    ;(globalThis as any).window.C311Provider = ConfiguredProvider
    expect(createC311Provider({})).to.be.instanceOf(ConfiguredProvider)
  })

  it('selects mock mode before the HTTP provider and passes fixture options', () => {
    const options: any[] = []
    class MockProvider implements TestProvider {
      constructor (value: any) { options.push(value) }
      async getSession (): Promise<any> { return { authenticated: false } }
    }
    class HttpProvider implements TestProvider {
      constructor (_transport: unknown) {}
      async getSession (): Promise<any> { return { authenticated: true } }
    }
    class FetchTransport {
      constructor (value: any) { expect(value.baseURL).to.equal('/api') }
    }

    const browserWindow = (globalThis as any).window
    browserWindow.C311Mode = 'mock'
    browserWindow.C311MockRole = 'service_agent'
    browserWindow.C311MockScenario = 'empty'
    browserWindow.C311MockSession = 'expired'
    const mock = createC311Provider({ MockC311Provider: MockProvider as any, C311HttpProvider: HttpProvider as any, C311FetchTransport: FetchTransport as any })
    expect(mock).to.be.instanceOf(MockProvider)
    expect(options).to.deep.equal([{ role: 'service_agent', scenario: 'empty', sessionVariant: 'expired' }])

    browserWindow.C311Mode = 'http'
    browserWindow.CortezaAPI = '/api'
    const http = createC311Provider({ C311HttpProvider: HttpProvider as any, C311FetchTransport: FetchTransport as any })
    expect(http).to.be.instanceOf(HttpProvider)
  })

  it('supports cached, forced and expired sessions while exposing capabilities and scopes', async () => {
    let calls = 0
    const provider: TestProvider = {
      getSession: async () => {
        calls += 1
        return { authenticated: true, expires_at: '2999-01-01T00:00:00.000Z', actor: { capabilities: ['read'], scopes: ['scope.read'] } }
      },
    }
    const runtime = new C311Runtime(provider)
    expect(await runtime.loadSession()).to.have.property('authenticated', true)
    expect(await runtime.loadSession()).to.have.property('authenticated', true)
    expect(calls).to.equal(1)
    expect(runtime.can('read')).to.equal(true)
    expect(runtime.can('write')).to.equal(false)
    expect(runtime.hasScope('scope.read')).to.equal(true)
    expect(runtime.hasScope('scope.write')).to.equal(false)
    await runtime.loadSession(true)
    expect(calls).to.equal(2)

    runtime.session = { authenticated: true, expires_at: '2000-01-01T00:00:00.000Z' }
    expect(runtime.isExpired()).to.equal(true)
    runtime.clearSession()
    expect(runtime.session).to.equal(null)
    expect(runtime.error).to.equal(null)
  })

  it('returns the current session while loading and records provider failures', async () => {
    let resolveSession: ((value: any) => void) | undefined
    const pending = new Promise(resolve => { resolveSession = resolve })
    const provider: TestProvider = { getSession: async () => pending }
    const runtime = new C311Runtime(provider)
    const first = runtime.loadSession()
    expect(await runtime.loadSession()).to.equal(null)
    resolveSession?.({ authenticated: true })
    expect(await first).to.deep.equal({ authenticated: true })

    const failure = new Error('session unavailable')
    const failed = new C311Runtime({ getSession: async () => { throw failure } })
    expect(await failed.loadSession()).to.equal(null)
    expect(failed.error).to.equal(failure)
    expect(failed.session).to.equal(null)
  })

  it('installs a reactive runtime on Vue instances', async () => {
    const provider: TestProvider = { getSession: async () => ({ authenticated: false }) }
    c311Plugin({ provider })(Vue)
    expect((Vue.prototype as any).$C311).to.be.instanceOf(C311Runtime)
    expect(await (Vue.prototype as any).$C311.loadSession()).to.deep.equal({ authenticated: false })
  })
})

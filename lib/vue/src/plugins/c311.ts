import Vue, { PluginFunction } from 'vue'
import type { C311Provider, C311Session } from '../libs/c311'

export interface C311PluginOptions {
  provider?: C311Provider
  mockProvider?: C311Provider
  mode?: 'mock' | 'http'
  baseURL?: string
  providerFactory?: (options: C311PluginOptions) => C311Provider
}

export class C311Runtime {
  provider: C311Provider
  session: C311Session | null = null
  loading = false
  error: unknown = null

  constructor (provider: C311Provider) {
    this.provider = provider
  }

  async loadSession (force = false): Promise<C311Session | null> {
    if (this.loading) return this.session
    if (this.session && !force && !this.isExpired()) return this.session

    this.loading = true
    this.error = null
    try {
      this.session = await this.provider.getSession()
      return this.session
    } catch (error) {
      this.error = error
      this.session = null
      return null
    } finally {
      this.loading = false
    }
  }

  clearSession (): void {
    this.session = null
    this.error = null
  }

  isExpired (): boolean {
    const expiresAt = this.session?.expires_at
    return !!expiresAt && Date.parse(expiresAt) <= Date.now()
  }

  can (capability: string): boolean {
    return !!this.session?.actor?.capabilities?.includes(capability)
  }

  hasScope (scope: string): boolean {
    return !!this.session?.actor?.scopes?.includes(scope)
  }
}

function fallbackProvider (): C311Provider {
  return {
    async getSession () {
      return { authenticated: false, actor: null, expires_at: null }
    },
  }
}

function configuredProvider (options: C311PluginOptions): C311Provider {
  if (options.provider) return options.provider
  if (options.providerFactory) return options.providerFactory(options)
  if (options.mockProvider) return options.mockProvider
  if (typeof window !== 'undefined' && window.C311Provider) return window.C311Provider
  return fallbackProvider()
}

export default function c311Plugin (options: C311PluginOptions = {}): PluginFunction<C311PluginOptions> {
  return function install (VueInstance): void {
    const provider = configuredProvider(options)
    const runtime = Vue.observable(new C311Runtime(provider))
    VueInstance.prototype.$C311 = runtime
  }
}

declare module 'vue/types/vue' {
  interface Vue {
    $C311: C311Runtime
  }
}

declare global {
  interface Window {
    C311Provider?: C311Provider
    C311Mode?: 'mock' | 'http'
    C311MockRole?: string
    C311MockScenario?: string
    C311MockSession?: 'current' | 'expired'
    CortezaAPI?: string
  }
}

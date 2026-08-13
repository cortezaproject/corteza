import { describe, expect, it, vi, beforeEach } from 'vitest'

/**
 * The bootstrap's two documented contracts (src.intent.md, plugins.intent.md):
 * "nothing renders before auth resolves", and on `Unauthenticated` the auth flow
 * starts and "setup aborts silently".
 *
 * It did neither. The catch started the flow and then RESOLVED, so main.js
 * mounted an app on which none of the plugins in the success path had been
 * installed, and App.vue's first useI18n() threw "Need to install with
 * `app.use` function". A cold-load trace put that crash at 617ms and the
 * navigation to /auth/callback at 820ms — the broken mount was only ever hidden
 * by losing a race to the redirect.
 */

const stubs = vi.hoisted(() => ({
  handle: vi.fn(),
  startAuthenticationFlow: vi.fn(),
  settingsInit: vi.fn(() => Promise.resolve()),
  localeGet: vi.fn(() => Promise.resolve({})),
  // Reset per test; the signed-in user's language drives both the i18n bundle
  // and PrimeVue's date format.
  userMeta: {},
}))

vi.mock('@planetcrust/human-vue', () => {
  const named = name => ({ pluginName: name, install: () => {} })

  return {
    AuthPlugin: {
      pluginName: 'AuthPlugin',
      install: app => {
        app.config.globalProperties.$Auth = {
          handle: stubs.handle,
          startAuthenticationFlow: stubs.startAuthenticationFlow,
          user: { meta: stubs.userMeta },
        }
      },
    },
    SystemAPIPlugin: {
      pluginName: 'SystemAPIPlugin',
      install: app => {
        app.config.globalProperties.$SystemAPI = { localeGet: stubs.localeGet }
      },
    },
    SettingsPlugin: {
      pluginName: 'SettingsPlugin',
      install: app => {
        app.config.globalProperties.$Settings = {
          init: stubs.settingsInit,
          get: () => undefined,
        }
      },
    },
    ComposeAPIPlugin: named('ComposeAPIPlugin'),
    DiscoveryAPIPlugin: named('DiscoveryAPIPlugin'),
    AutomationAPIPlugin: named('AutomationAPIPlugin'),
    FederationAPIPlugin: named('FederationAPIPlugin'),
    EventBusPlugin: named('EventBusPlugin'),
    I18nPlugin: named('I18nPlugin'),
    PrimeVueComponentsPlugin: named('PrimeVueComponentsPlugin'),
    ToastPlugin: named('ToastPlugin'),
    getTheme: () => ({}),
    setThemes: () => {},
    // Echoes the language it was handed so the assertion below can prove the
    // user's locale actually reaches the PrimeVue config rather than a default.
    primeVueLocale: lang => ({ dateFormat: `format-for-${lang}` }),
  }
})

// The router pulls in every section; none of it is under test here.
vi.mock('../router', () => ({ default: { pluginName: 'router', install: () => {} } }))
vi.mock('pinia', () => ({ createPinia: () => ({ pluginName: 'pinia', install: () => {} }) }))
vi.mock('primevue/config', () => ({ default: { pluginName: 'PrimeVue', install: () => {} } }))
vi.mock('primevue/toastservice', () => ({
  default: { pluginName: 'ToastService', install: () => {} },
}))
vi.mock('primevue/confirmationservice', () => ({
  default: { pluginName: 'ConfirmationService', install: () => {} },
}))
vi.mock('primevue/dialogservice', () => ({
  default: { pluginName: 'DialogService', install: () => {} },
}))
vi.mock('primevue/ripple', () => ({ default: {} }))

const { setupAndAuthenticate } = await import('./index')

// A stand-in for the Vue app that records which plugins were installed, which
// is the thing that decides whether mounting is safe.
function fakeApp() {
  const installed = []

  // Keyed by plugin name so a test can inspect what a plugin was configured
  // with, not only that it was installed.
  const optionsFor = {}

  const app = {
    installed,
    optionsFor,
    config: { globalProperties: {} },
    directive: () => app,
    use(plugin, options) {
      const name = plugin?.pluginName ?? 'anonymous'
      installed.push(name)
      optionsFor[name] = options
      plugin?.install?.(app, options)
      return app
    },
  }

  return app
}

describe('setupAndAuthenticate', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    stubs.settingsInit.mockResolvedValue()
    stubs.localeGet.mockResolvedValue({})
    Object.keys(stubs.userMeta).forEach(k => delete stubs.userMeta[k])
  })

  it("configures PrimeVue with the user's own date format", async () => {
    // Dates used to render US-formatted for everyone because PrimeVue's locale
    // was never configured at all.
    stubs.userMeta.preferredLanguage = 'sl'
    stubs.handle.mockResolvedValue()

    const app = fakeApp()
    await setupAndAuthenticate(app)

    expect(app.optionsFor.PrimeVue.locale).toEqual({ dateFormat: 'format-for-sl' })
  })

  it('leaves dates to the browser when the user has no language set', async () => {
    // The UI has to fall back to English — it is the only bundle there is — but
    // dates must not: an English interface does not mean the reader wants
    // American date order. Passing undefined lets Intl use the browser locale.
    stubs.handle.mockResolvedValue()

    const app = fakeApp()
    await setupAndAuthenticate(app)

    expect(app.optionsFor.PrimeVue.locale).toEqual({ dateFormat: 'format-for-undefined' })
    expect(stubs.localeGet).toHaveBeenCalledWith({ lang: 'en', application: 'human-webapp' })
  })

  it('resolves false and starts the auth flow when unauthenticated', async () => {
    stubs.handle.mockRejectedValue(new Error('Unauthenticated'))

    const app = fakeApp()
    const ready = await setupAndAuthenticate(app)

    expect(ready).toBe(false)
    expect(stubs.startAuthenticationFlow).toHaveBeenCalledOnce()
  })

  it('installs no renderable plugins when it aborts for auth', async () => {
    stubs.handle.mockRejectedValue(new Error('Unauthenticated'))

    const app = fakeApp()
    await setupAndAuthenticate(app)

    // i18n above all: App.vue's setup() calls useI18n() immediately, so mounting
    // without it is the crash this guards against.
    expect(app.installed).not.toContain('I18nPlugin')
    expect(app.installed).not.toContain('router')
    expect(app.installed).not.toContain('PrimeVue')
  })

  it('resolves true once the app is fully configured', async () => {
    stubs.handle.mockResolvedValue()

    const app = fakeApp()
    const ready = await setupAndAuthenticate(app)

    expect(ready).toBe(true)
    expect(app.installed).toContain('I18nPlugin')
    expect(app.installed).toContain('router')
    expect(app.installed).toContain('PrimeVue')
  })

  it('still boots when the locale bundle cannot be fetched', async () => {
    stubs.handle.mockResolvedValue()
    stubs.localeGet.mockRejectedValue(new Error('network'))

    const app = fakeApp()

    // Missing translations fall back to {} and must never block boot.
    await expect(setupAndAuthenticate(app)).resolves.toBe(true)
    expect(app.installed).toContain('I18nPlugin')
  })

  it('propagates any error that is not an auth challenge', async () => {
    stubs.handle.mockRejectedValue(new Error('boom'))

    await expect(setupAndAuthenticate(fakeApp())).rejects.toThrow('boom')
    expect(stubs.startAuthenticationFlow).not.toHaveBeenCalled()
  })
})

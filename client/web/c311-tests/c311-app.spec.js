import Vue from 'vue'

const mockStore = { dispatch: jest.fn() }
const mockRouter = { beforeEach: jest.fn() }
const mockI18n = () => ({
  i18next: {
    isInitialized: true,
    on: jest.fn(),
    changeLanguage: jest.fn(),
  },
})
const mockC311I18n = {
  installC311Translations: jest.fn(),
  readC311Locale: jest.fn(() => 'es'),
  persistC311Locale: jest.fn(),
}
const mockCorredor = {
  ComposeCtx: class ComposeCtx {},
  WebappCtx: class WebappCtx {},
}
const mockWebsocket = { init: jest.fn(), endpoint: jest.fn(() => '/socket'), config: {} }
const mockPlugin = () => ({ install: jest.fn() })
const mockPlugins = {
  Auth: jest.fn(mockPlugin),
  C311: jest.fn(mockPlugin),
  createC311Provider: jest.fn(() => ({ getSession: async () => ({ authenticated: false }) })),
  CortezaAPI: jest.fn(mockPlugin),
  DiscoveryAPI: jest.fn(mockPlugin),
  EventBus: jest.fn(mockPlugin),
  UIHooks: jest.fn(mockPlugin),
  Settings: mockPlugin(),
  Reminder: mockPlugin(),
}

jest.mock('../compose/src/config-check', () => ({}))
jest.mock('../compose/src/console-splash', () => ({}))
jest.mock('../compose/src/filters', () => ({}))
jest.mock('../compose/src/plugins', () => ({}))
jest.mock('../compose/src/mixins', () => ({}))
jest.mock('../compose/src/components', () => ({}))
jest.mock('../compose/src/store', () => ({ __esModule: true, default: mockStore }))
jest.mock('../compose/src/router', () => ({ __esModule: true, default: mockRouter }))
jest.mock('../admin/src/config-check', () => ({}))
jest.mock('../admin/src/console-splash', () => ({}))
jest.mock('../admin/src/filters', () => ({}))
jest.mock('../admin/src/plugins', () => ({}))
jest.mock('../admin/src/mixins', () => ({}))
jest.mock('../admin/src/components', () => ({}))
jest.mock('../admin/src/store', () => ({ __esModule: true, default: mockStore }))
jest.mock('../admin/src/router', () => ({ __esModule: true, default: mockRouter }))
jest.mock('vue', () => {
  function MockVue (options) {
    Object.assign(this, options)
    this.options = options
  }
  MockVue.use = jest.fn()
  MockVue.prototype = {}
  return { __esModule: true, default: MockVue }
})
jest.mock('bootstrap-vue', () => ({ install: jest.fn() }))
jest.mock('vue-router', () => function MockRouter () {})
jest.mock('vue-native-websocket', () => ({ install: jest.fn() }))
jest.mock('@cortezaproject/corteza-vue', () => ({
  plugins: mockPlugins,
  mixins: { corredor: {} },
  corredor: mockCorredor,
  websocket: mockWebsocket,
  i18n: jest.fn(mockI18n),
  c311I18n: mockC311I18n,
}))
jest.mock('@cortezaproject/corteza-js', () => ({
  compose: {
    TriggerComposeServerScriptOnManual: jest.fn(() => jest.fn()),
  },
  system: {
    TriggerSystemServerScriptOnManual: jest.fn(() => jest.fn()),
  },
}))
jest.mock('vuex', () => ({ __esModule: true, default: { install: jest.fn() }, mapGetters: () => ({ isRbacLoaded: () => true }) }))

function makeRuntime (route) {
  const handlers = {}
  const i18next = {
    isInitialized: true,
    on: jest.fn((event, handler) => { handlers[event] = handler }),
    changeLanguage: jest.fn(),
  }
  const user = { userID: 'actor-c311', meta: { preferredLanguage: 'en', theme: '' } }
  const runtime = {
    $route: route,
    $i18n: { i18next },
    $auth: {
      accessToken: undefined,
      user,
      vue: () => ({ handle: async () => ({ user }) }),
      startAuthenticationFlow: jest.fn(),
    },
    $Settings: {
      init: async () => undefined,
      attachment: () => null,
      get: () => false,
    },
    $SystemAPI: {
      setHeader: jest.fn().mockReturnThis(),
      automationList: async () => [],
    },
    $ComposeAPI: {
      setHeader: jest.fn().mockReturnThis(),
      automationList: async () => [],
    },
    $AutomationAPI: {},
    $FederationAPI: {},
    $store: mockStore,
    $on: jest.fn(),
    websocket: jest.fn(),
    textDirectionality: () => 'ltr',
    loadBundle: async () => undefined,
    makeAutomationScriptsRegistrator: () => () => undefined,
    i18nLoaded: false,
    loaded: false,
  }
  return { runtime, handlers, user, i18next }
}

describe('C311 application shell initialization', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    window.C311Mode = undefined
  })

  it('lets Compose public C311 routes load without the OAuth flow', async () => {
    const appFactory = require('../compose/src/app').default
    const definition = appFactory({ router: mockRouter })
    const { runtime, i18next } = makeRuntime({ name: 'c311.submit', meta: { c311: { public: true } } })
    runtime.$i18n.i18next = i18next
    await definition.created.call(runtime)

    expect(runtime.loaded).toBe(true)
    expect(runtime.websocket).not.toHaveBeenCalled()
    expect(runtime.$auth.startAuthenticationFlow).not.toHaveBeenCalled()
    expect(mockC311I18n.installC311Translations).toHaveBeenCalled()
  })

  it('initializes Compose authenticated C311 locale persistence', async () => {
    const appFactory = require('../compose/src/app').default
    const definition = appFactory({ router: mockRouter })
    const { runtime, handlers } = makeRuntime({ name: 'root', meta: {} })
    await definition.created.call(runtime)

    expect(runtime.loaded).toBe(true)
    expect(runtime.websocket).toHaveBeenCalled()
    expect(runtime.$i18n.i18next.changeLanguage).toHaveBeenCalledWith('es')
    handlers.languageChanged('es-MX')
    expect(mockC311I18n.persistC311Locale).toHaveBeenCalledWith('es', 'actor-c311')
  })

  it('allows Admin mock routes and exposes the mock route computed flag', async () => {
    const appFactory = require('../admin/src/app').default
    const definition = appFactory({ router: mockRouter })
    const { runtime } = makeRuntime({ name: 'c311.staff', meta: { c311: { public: true } } })
    window.C311Mode = 'mock'
    expect(definition.computed.isC311MockRoute.call(runtime)).toBe(true)
    await definition.created.call(runtime)

    expect(runtime.loaded).toBe(true)
    expect(runtime.websocket).not.toHaveBeenCalled()
  })

  it('initializes Admin authenticated C311 locale persistence', async () => {
    const appFactory = require('../admin/src/app').default
    const definition = appFactory({ router: mockRouter })
    const { runtime, handlers } = makeRuntime({ name: 'root', meta: {} })
    await definition.created.call(runtime)

    expect(runtime.loaded).toBe(true)
    expect(runtime.websocket).toHaveBeenCalled()
    handlers.languageChanged('vi-VN')
    expect(mockC311I18n.persistC311Locale).toHaveBeenCalledWith('vi', 'actor-c311')
  })

  it('initializes the shared C311 plugin from both application plugin entrypoints', () => {
    window.C311Mode = 'mock'
    jest.requireActual('../compose/src/plugins/index')
    jest.requireActual('../admin/src/plugins/index')

    expect(mockPlugins.createC311Provider).toHaveBeenCalledTimes(2)
    expect(mockPlugins.C311).toHaveBeenCalledTimes(2)
    expect(mockPlugins.C311.mock.calls[0][0]).toEqual(expect.objectContaining({ provider: expect.any(Object) }))
  })
})

import { apiClients } from '@planetcrust/human-js'
import type { App } from 'vue'

interface Options {
  baseURL?: string
  accessTokenFn?: () => string | undefined
}

const getBaseURL = (service: string, opt: Options = {}) => {
  if (opt.baseURL) {
    return opt.baseURL
  } else {
    // @ts-ignore
    if (!window.HumanAPI) {
      throw new Error('config.js missing or window.HumanAPI not set')
    }

    // @ts-ignore
    return `${window.HumanAPI}/${service}`
  }
}

const getAccessTokenFn = (app: App, opt: Options = {}) => {
  if (opt.accessTokenFn) {
    return opt.accessTokenFn
  }

  return app.config.globalProperties.$Auth.accessTokenFn
}

const getOptions = (service: string, app: App, opt: Options = {}) => {
  return {
    baseURL: getBaseURL(service, opt),
    accessTokenFn: getAccessTokenFn(app, opt),
  }
}

export const SystemAPIPlugin = {
  install(app: App, opt: Options) {
    const SystemAPI = new apiClients.System(getOptions('system', app, opt))
    app.config.globalProperties.$SystemAPI = SystemAPI
    app.provide('$SystemAPI', SystemAPI)
  },
}

export const ComposeAPIPlugin = {
  install(app: App, opt: Options) {
    const ComposeAPI = new apiClients.Compose(getOptions('compose', app, opt))
    app.config.globalProperties.$ComposeAPI = ComposeAPI
    app.provide('$ComposeAPI', ComposeAPI)
  },
}

export const AutomationAPIPlugin = {
  install(app: App, opt: Options) {
    const AutomationAPI = new apiClients.Automation(getOptions('automation', app, opt))
    app.config.globalProperties.$AutomationAPI = AutomationAPI
    app.provide('$AutomationAPI', AutomationAPI)
  },
}

export const FederationAPIPlugin = {
  install(app: App, opt: Options) {
    const FederationAPI = new apiClients.Federation(getOptions('federation', app, opt))
    app.config.globalProperties.$FederationAPI = FederationAPI
    app.provide('$FederationAPI', FederationAPI)
  },
}

export const DiscoveryAPIPlugin = {
  install(app: App, opt: Options) {
    const baseURL =
      opt?.baseURL ??
      (window as unknown as Record<string, string>).HumanDiscoveryAPI ??
      'http://localhost:3200/'
    const DiscoveryAPI = new apiClients.Discovery({
      baseURL,
      accessTokenFn: getAccessTokenFn(app, opt),
    })
    app.config.globalProperties.$DiscoveryAPI = DiscoveryAPI
    app.provide('$DiscoveryAPI', DiscoveryAPI)
  },
}

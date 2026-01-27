import { apiClients } from '@cortezaproject/corteza-js-next'
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
    if (!window.CortezaAPI) {
      throw new Error('config.js missing or window.CortezaAPI not set')
    }

    // @ts-ignore
    return `${window.CortezaAPI}/${service}`
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

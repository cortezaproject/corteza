/* eslint-disable @typescript-eslint/ban-ts-ignore */
import { apiClients } from '@cortezaproject/corteza-js'
import { PluginFunction } from 'vue'

interface JWTFetcher {
  (): string|null;
}

interface Options {
  baseURL?: string;
  accessTokenFn?: () => string | undefined;
}

/**
 * Generic Corteza API plugin
 *
 * Install a specific plugin:
 * Vue.use(plugins.CortezaAPI('compose'))
 *
 * @constructor
 */
export default function (service: string, opt: Options = {}): PluginFunction<Options> {
  if (!opt.baseURL) {
    // @ts-ignore
    if (!window.CortezaAPI) {
      throw new Error('config.js missing or window.CortezaAPI not set')
    }

    // @ts-ignore
    opt.baseURL = `${window.CortezaAPI}/${service}`
  }

  return function (Vue): void {
    service = service.substring(0, 1).toUpperCase() + service.substring(1)

    if (!opt.accessTokenFn) {
      /**
       * Checking if auth plugin was initialized before and
       * hooking on to it's accessTokenFn
       */
      opt.accessTokenFn = Vue.prototype.$auth.accessTokenFn
    }

    // @ts-ignore
    // makes Vue.$<service>API (Vue.$SystemAPI, Vue.$ComposeAPI, Vue.$FederationAPI, Vue.$AutomationAPI) available
    Vue.prototype[`$${service}API`] = new apiClients[service](opt)
  }
}

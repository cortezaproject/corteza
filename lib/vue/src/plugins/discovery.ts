import { apiClients } from '@cortezaproject/corteza-js'
import { PluginFunction } from 'vue'

interface Options {
  baseURL?: string;
  accessTokenFn?: () => string | undefined;
}

/**
 * Corteza Discovery API plugin
 *
 * Install:
 * Vue.use(plugins.DiscoveryAPI())
 *
 * @constructor
 */
export default function (opt: Options = {}): PluginFunction<Options> {
  if (!opt.baseURL) {
    // @ts-ignore
    if (window.CortezaDiscoveryAPI) {
      // @ts-ignore
      opt.baseURL = `${window.CortezaDiscoveryAPI}/`
    } else {
      // Fallback for development environments
      opt.baseURL = 'http://localhost:3200/'
    }
  }

  return function (Vue): void {
    if (!opt.accessTokenFn) {
      /**
       * Checking if auth plugin was initialized before and
       * hooking on to it's accessTokenFn
       */
      opt.accessTokenFn = Vue.prototype.$auth.accessTokenFn
    }

    if (opt.baseURL) {
      // @ts-ignore
      Vue.prototype.$DiscoveryAPI = new apiClients.Discovery(opt)
    } else {
      console.warn('window.CortezaDiscoveryAPI not set, $DiscoveryAPI not initialized')
    }
  }
}

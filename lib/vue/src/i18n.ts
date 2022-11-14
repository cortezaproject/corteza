/* eslint-disable @typescript-eslint/ban-ts-ignore */
import { VueConstructor } from 'vue'
import i18next, { InitOptions } from 'i18next'
import http from 'i18next-http-backend'
import detector from 'i18next-browser-languagedetector'
import multiload from 'i18next-multiload-backend-adapter'
import VueI18Next from '@panter/vue-i18next'
import moment from 'moment'
import Pseudo from 'i18next-pseudo'

interface Options {
  app: string;

  resources: object;

  lng: string;
  fallbackLng: string | false;

  /**
   * Namespace(s) to preload
   */
  ns: string | Array<string>;

  /**
   * Where too look for keys that are not found in the default namespace,
   */
  fallbackNS: string | false;

  /**
   * What namespace to use when not explicitly defined
   * When empty, default is set to the value of (the first item in) ns
   */
  defaultNS: string;

  baseURL: string;

  pseudo: boolean;
}

/**
 * Initializes i18n options, registers plugin on a given Vue instance and returns the options
 *
 * To be used as:
 * import { i18n } from '@cortezaproject/corteza-vue'
 * new Vue({
 *   i18n: i18n(Vue, {
 *     app: 'corteza-webapp-....'
 *     namespaces: [ .... ]
 *    }),
*   })
 *
 * The most convenient way to use it:
 * i18n(Vue, 'app name', 'namespace...', 'additional namespace...')
 */
export default (Vue: VueConstructor, app: string | Partial<Options>, ...namespaces: Array<string>): VueI18Next => {
  const devMode = process.env.NODE_ENV !== 'production'
  const defNS = 'translation'

  let opt: Partial<Options> = {}
  if (typeof app === 'string') {
    opt = {
      app,
    }
  } else {
    opt = app
  }

  const {
    // keeping lng without a value
    // lang-auto-detect plugin
    lng,
    fallbackLng = 'en',
    fallbackNS = false,
  } = opt

  let ns: Array<string> = []
  if (!Array.isArray(opt.ns)) {
    ns = [opt.ns || defNS]
  } else {
    ns = opt.ns
  }

  ns.push(...namespaces)

  const defaultNS = opt.defaultNS || ns[0]

  if (!opt.baseURL) {
    // @ts-ignore
    if (!window.CortezaAPI) {
      throw new Error('config.js missing or window.CortezaAPI not set')
    }

    // @ts-ignore
    opt.baseURL = `${window.CortezaAPI}/system`
  }

  const pseudo = devMode && (
    !!opt.pseudo ||
    // @ts-ignore
    !!window.i18nPseudoModeEnabled ||
    // @ts-ignore
    window.location.search.indexOf('i18nPseudoModeEnabled') > -1
  )

  let postProcess: Array<string> = []

  if (pseudo) {
    postProcess = ['pseudo']
  }

  const options: InitOptions = {
    debug: devMode,

    lng,

    fallbackLng,
    ns,
    fallbackNS,
    defaultNS,

    postProcess,

    initImmediate: false,

    detection: {
      // to overwrite, to use user defined, to guess user's lang
      order: ['querystring', 'localStorage', 'cookie', 'navigator'],
      caches: devMode ? [] : ['localStorage', 'cookie'],
    },

    backend: {
      // @ts-ignore
      backend: http,
      backendOption: {
        loadPath: `${opt.baseURL}/locale/{{lng}}/${opt.app}`,
      },
    },
  }

  i18next
    .use(detector)
    .use(multiload)
    .use(new Pseudo({
      enabled: pseudo,
    }))
    .init(options)

  Vue.use(VueI18Next)

  // Set locales for other libs we use
  // @todo this needs to be set after language is detected
  moment.locale(lng)

  return new VueI18Next(i18next)
}

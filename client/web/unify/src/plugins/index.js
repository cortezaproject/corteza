import '@/assets/styles.css'
import 'primeicons/primeicons.css'

import {
    AuthPlugin,
    AutomationAPIPlugin,
    ComposeAPIPlugin,
    DiscoveryAPIPlugin,
    EventBusPlugin,
    FederationAPIPlugin,
    I18nPlugin,
    PrimeVueComponentsPlugin,
    SettingsPlugin,
    SystemAPIPlugin,
    ToastPlugin,
    getTheme,
    setThemes,
} from '@planetcrust/human-vue'
import { isPlainObject, mergeWith } from 'lodash-es'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'
import DialogService from 'primevue/dialogservice'
import Ripple from 'primevue/ripple'
import ToastService from 'primevue/toastservice'
import { localeApplications } from '../sections'
import router from '../router'

/**
 * Sets up PrimeVue with theming and services
 */
function setupPrimeVue(app, theme) {
  app.use(PrimeVue, {
    theme: {
      preset: getTheme(theme),
      options: {
        darkModeSelector: '.dark',
        cssLayer: {
          name: 'primevue',
          order: 'tailwind-base, primevue, tailwind-utilities',
        },
      },
    },
    ripple: true,
  })

  app.directive('ripple', Ripple)

  // PrimeVue services
  app.use(ToastService)
  app.use(ConfirmationService)
  app.use(DialogService)

  // Human toast wrapper
  app.use(ToastPlugin)

  // Register common PrimeVue components globally
  app.use(PrimeVueComponentsPlugin)
}

// True if a value carries actual translation text (a string anywhere inside).
function hasTranslationContent(value) {
  if (typeof value === 'string') return true
  if (isPlainObject(value)) return Object.values(value).some(hasTranslationContent)
  return false
}

/**
 * Deep-merges locale bundles. Apps sometimes define the same key with a
 * different shape — flat string in one app, nested object in another
 * (e.g. `navigation.automation` is "Automation" in compose but
 * `{ group, items }` in admin). A plain deep-merge lets whichever bundle is
 * merged last clobber the other shape and silently drop subkeys. On such a
 * string↔object collision we keep whichever side actually contains
 * translations (so admin's `{ group, items }` wins over a flat label, but a
 * real "General" string wins over an unused garbage object).
 */
function localeMerge(target, source) {
  return mergeWith(target, source, (objValue, srcValue) => {
    const objIsObject = isPlainObject(objValue)
    const srcIsObject = isPlainObject(srcValue)
    if (objIsObject === srcIsObject) return undefined // same shape → default merge
    const objectSide = objIsObject ? objValue : srcValue
    const primitiveSide = objIsObject ? srcValue : objValue
    return hasTranslationContent(objectSide) ? objectSide : primitiveSide
  })
}

/**
 * Loads and merges the locale bundles for every active section so a single
 * i18n instance resolves chrome + every section's namespaces.
 */
async function loadMergedTranslations(api, locale) {
  const bundles = await Promise.all(
    localeApplications.map(application =>
      api.localeGet({ lang: locale, application }).catch(() => ({})),
    ),
  )
  return bundles.reduce((acc, bundle) => localeMerge(acc, bundle), {})
}

/**
 * Main app setup and authentication flow
 */
export function setupAndAuthenticate(app) {
  app.use(AuthPlugin, { app: import.meta.env.VITE_APP_ID, rootApp: true })

  const $Auth = app.config.globalProperties.$Auth

  return $Auth
    .handle()
    .then(async () => {
      // API plugins
      app.use(SystemAPIPlugin)
      app.use(ComposeAPIPlugin)
      app.use(DiscoveryAPIPlugin)
      app.use(AutomationAPIPlugin)
      app.use(FederationAPIPlugin)

      // Settings
      app.use(SettingsPlugin, {
        api: app.config.globalProperties.$SystemAPI,
      })

      // State management & routing
      app.use(createPinia())
      app.use(EventBusPlugin)
      app.use(router)

      // i18n — merge every active section's locale bundle into one instance
      const locale = $Auth.user.meta.preferredLanguage || 'en'
      const translations = await loadMergedTranslations(
        app.config.globalProperties.$SystemAPI,
        locale,
      )

      app.use(I18nPlugin, {
        locale: locale,
        translations: translations,
      })

      // UI setup (after settings loaded for theme)
      const $Settings = app.config.globalProperties.$Settings

      return $Settings.init().then(() => {
        setThemes($Settings.get('ui.studio.themes'))
        setupPrimeVue(app, $Auth.user.meta.theme)
      })
    })
    .catch(err => {
      if (err instanceof Error && err.message === 'Unauthenticated') {
        $Auth.startAuthenticationFlow()
        return
      }
      throw err
    })
}

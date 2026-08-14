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
  primeVueLocale,
  setThemes,
} from '@planetcrust/human-vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'
import DialogService from 'primevue/dialogservice'
import Ripple from 'primevue/ripple'
import ToastService from 'primevue/toastservice'
import router from '../router'
import TrimOnBlurPlugin from './trimOnBlur'

/**
 * Sets up PrimeVue with theming and services
 */
function setupPrimeVue(app, theme, locale) {
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
    // Date ordering, day/month names and week start come from Intl rather than
    // PrimeVue's US default; merged over its built-in locale, so every string
    // this does not name keeps its default. `locale` is the user's explicit
    // preference or undefined — undefined means the browser's own, which is
    // what someone reading an English UI still expects their dates in. Reaches
    // every date input in the app through CInputDateTime.
    locale: primeVueLocale(locale),
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

/**
 * Main app setup and authentication flow
 *
 * Resolves true when the app is configured and safe to mount, false when setup
 * deliberately abandoned it to start the auth flow. The caller must not mount on
 * false: none of the plugins below were installed, so App.vue's first useI18n()
 * throws. That is invisible in practice only because the auth redirect usually
 * wins the race — a trace of a cold load showed the mount error at 617ms and the
 * navigation to /auth/callback at 820ms.
 */
export function setupAndAuthenticate(app) {
  app.use(AuthPlugin, { app: import.meta.env.VITE_APP_ID, rootApp: true })

  // Global trim-on-blur for text inputs/textareas across the app.
  app.use(TrimOnBlurPlugin)

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
      const pinia = createPinia()
      app.use(pinia)
      app.use(EventBusPlugin)
      app.use(router)

      exposeForDebugging(app, { pinia, router })

      // i18n — single consolidated locale bundle for chrome + every section.
      // The UI falls back to English because that is the only bundle there is;
      // dates must NOT, so they take the unset value and let Intl use the
      // browser's locale instead (see setupPrimeVue).
      const preferredLanguage = $Auth.user.meta.preferredLanguage
      const locale = preferredLanguage || 'en'
      const translations = await app.config.globalProperties.$SystemAPI
        .localeGet({ lang: locale, application: 'human-webapp' })
        .catch(() => ({}))

      app.use(I18nPlugin, {
        locale: locale,
        translations: translations,
      })

      // UI setup (after settings loaded for theme)
      const $Settings = app.config.globalProperties.$Settings

      return $Settings.init().then(() => {
        setThemes($Settings.get('ui.studio.themes'))
        setupPrimeVue(app, $Auth.user.meta.theme, preferredLanguage)

        return true
      })
    })
    .catch(err => {
      if (err instanceof Error && err.message === 'Unauthenticated') {
        $Auth.startAuthenticationFlow()

        // The browser is on its way to the auth server; there is nothing
        // configured to render in the meantime.
        return false
      }
      throw err
    })
}

// Hands a dev build's router, stores and global properties to whatever is
// driving the browser — a devtools console, or an automated UI check.
//
// Without it, reading component or store state from outside the app means
// editing the component to park something on `window`, running, then
// remembering to take it out again: a source edit as a debugging step, on code
// that is not the one under investigation. `import.meta.env.DEV` is statically
// false in a production build, so this whole function is dropped at build time.
function exposeForDebugging(app, { pinia, router }) {
  if (!import.meta.env.DEV || typeof window === 'undefined') return

  window.__human = {
    app,
    router,
    pinia,

    // Route table as name → path, so a check can look a route up instead of
    // grepping the section's index.js for it.
    routes: () => router.getRoutes().map(r => ({ name: r.name, path: r.path, meta: r.meta })),

    // Every initialised Pinia store, by id, with its state readable.
    stores: () => Object.fromEntries(pinia._s),

    // The `$Foo` globals (`$SystemAPI`, `$Settings`, `$Auth`, …).
    globals: () => app.config.globalProperties,
  }
}

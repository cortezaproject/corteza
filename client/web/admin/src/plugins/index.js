import '@/assets/styles.css'
import 'primeicons/primeicons.css'

import {
  AuthPlugin,
  AutomationAPIPlugin,
  ComposeAPIPlugin,
  EventBusPlugin,
  FederationAPIPlugin,
  I18nPlugin,
  PrimeVueComponentsPlugin,
  SettingsPlugin,
  SystemAPIPlugin,
  ToastPlugin,
  getTheme,
  setThemes,
} from '@cortezaproject/corteza-vue-next'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'
import DialogService from 'primevue/dialogservice'
import Ripple from 'primevue/ripple'
import ToastService from 'primevue/toastservice'
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

  // Corteza toast wrapper
  app.use(ToastPlugin)

  // Register common PrimeVue components globally
  app.use(PrimeVueComponentsPlugin)
}

/**
 * Main app setup and authentication flow
 */
export function setupAndAuthenticate(app) {
  app.use(AuthPlugin, { app: import.meta.env.VITE_APP_ID })

  const $Auth = app.config.globalProperties.$Auth

  return $Auth
    .handle()
    .then(async () => {
      // API plugins
      app.use(SystemAPIPlugin)
      app.use(ComposeAPIPlugin)
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

      // i18n
      const locale = $Auth.user.meta.preferredLanguage || 'en'
      const translations = await app.config.globalProperties.$SystemAPI.localeGet({
        lang: locale,
        application: import.meta.env.VITE_APP_NAME,
      })

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

import {
  AuthPlugin,
  AutomationAPIPlugin,
  ComposeAPIPlugin,
  I18nPlugin,
  SettingsPlugin,
  SystemAPIPlugin,
  setThemes,
} from '@cortezaproject/corteza-vue-next'
import { createPinia } from 'pinia'
import { UIPlugin } from '../components'
import router from '../router'

export function setupAndAuthenticate(app) {
  app.use(AuthPlugin, { app: import.meta.env.VITE_APP_ID, rootApp: true })

  const $Auth = app.config.globalProperties.$Auth

  return $Auth
    .handle()
    .then(async () => {
      app.use(SystemAPIPlugin)
      app.use(ComposeAPIPlugin)
      app.use(AutomationAPIPlugin)

      app.use(SettingsPlugin, {
        api: app.config.globalProperties.$SystemAPI,
      })

      app.use(createPinia())
      app.use(router)

    const locale = $Auth.user.meta.preferredLanguage || 'en'
      const translations = await app.config.globalProperties.$SystemAPI.localeGet({
        lang: locale,
        application: import.meta.env.VITE_APP_NAME,
      })

      app.use(I18nPlugin, {
        locale: locale,
        translations: translations,
      })

      const $Settings = app.config.globalProperties.$Settings

      return $Settings.init().then(() => {
        setThemes($Settings.get('ui.studio.themes'))
        app.use(UIPlugin, { theme: $Auth.user.meta.theme })
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

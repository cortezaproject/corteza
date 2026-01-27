import { App } from 'vue'
import { createI18n } from 'vue-i18n'

interface Options {
  appName: string
  locale: string
  translations: any
}

export const I18nPlugin = {
  async install(app: App, options: Options) {
    try {
      app.use(
        createI18n({
          legacy: false, // Composition API
          locale: options.locale,
          fallbackLocale: 'en',
          messages: options.translations,
        }),
      )
    } catch (error) {
      console.error('Failed to load translations:', error)
    }
  },
}

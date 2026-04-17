import { config } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { beforeEach, vi } from 'vitest'

const i18n = createI18n({ legacy: false, locale: 'en', messages: { en: {} } })

config.global.plugins = [i18n]

beforeEach(() => {
  vi.spyOn(console, 'error').mockImplementation(() => {})
})

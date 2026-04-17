import { config } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import PrimeVue from 'primevue/config'
import { PrimeVueComponentsPlugin } from '../plugins/primevue-components'
import { beforeEach, vi } from 'vitest'

const i18n = createI18n({ legacy: false, locale: 'en', messages: { en: {} } })

config.global.plugins = [
  i18n,
  [PrimeVue, { unstyled: true }],
  PrimeVueComponentsPlugin,
]

beforeEach(() => {
  vi.spyOn(console, 'error').mockImplementation(() => {})
  vi.spyOn(console, 'warn').mockImplementation(() => {})
})

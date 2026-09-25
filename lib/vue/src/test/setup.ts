import { config } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import PrimeVue from 'primevue/config'
import { PrimeVueComponentsPlugin } from '../plugins/primevue-components'
import { beforeEach, vi } from 'vitest'

const i18n = createI18n({ legacy: false, locale: 'en', messages: { en: {} } })

// jsdom has no layout, so a Range reports no geometry at all. The editor asks
// for it whenever it scrolls the caret into view.
const emptyRect = { top: 0, left: 0, bottom: 0, right: 0, width: 0, height: 0, x: 0, y: 0 }
Range.prototype.getClientRects = () => Object.assign([], { item: () => null })
Range.prototype.getBoundingClientRect = () => ({ ...emptyRect, toJSON: () => emptyRect })

// jsdom has no ResizeObserver; PrimeVue's Textarea watches its own size with one.
class ResizeObserverStub {
  observe() {}
  unobserve() {}
  disconnect() {}
}
globalThis.ResizeObserver ??= ResizeObserverStub as unknown as typeof ResizeObserver

config.global.plugins = [i18n, [PrimeVue, { unstyled: true }], PrimeVueComponentsPlugin]

beforeEach(() => {
  vi.spyOn(console, 'error').mockImplementation(() => {})
  vi.spyOn(console, 'warn').mockImplementation(() => {})
})

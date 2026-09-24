import { App as McpApp, applyDocumentTheme } from '@modelcontextprotocol/ext-apps'
import type { McpUiHostContext } from '@modelcontextprotocol/ext-apps'
import { usePreset } from '@primeuix/themes'
import { getTheme } from '@planetcrust/human-vue/src/composables/useTheme'
import PrimeVue from 'primevue/config'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Menu from 'primevue/menu'
import messages from 'virtual:human-locale'
import { type Component, createApp, provide, h, ref } from 'vue'
import { createI18n } from 'vue-i18n'
import './styles.css'

export const mcpAppKey = Symbol('mcpApp')

type Theme = 'light' | 'dark'

// The host's current theme, for what reads it outside CSS (echarts).
export const hostTheme = ref<Theme>('light')

function applyTheme(ctx?: McpUiHostContext) {
  const theme: Theme = ctx?.theme === 'dark' ? 'dark' : 'light'
  hostTheme.value = theme
  applyDocumentTheme(theme)
  return getTheme(theme)
}

// Mounts a view with Human's PrimeVue theme and translations, following the
// host's light or dark, and connects it to the host. The connected App is
// provided under mcpAppKey; results arrive through its ontoolresult.
//
// Only the PrimeVue components lib/vue's reused components expect to find
// registered are registered: registering them all quadruples every view.
export async function mount(root: Component, name: string, version = '0.1.0') {
  const bridge = new McpApp({ name, version })

  const app = createApp({
    setup() {
      provide(mcpAppKey, bridge)
      return () => h(root)
    },
  })

  app.use(PrimeVue, {
    theme: {
      preset: applyTheme(),
      options: {
        darkModeSelector: '.dark',
        cssLayer: { name: 'primevue', order: 'tailwind-base, primevue, tailwind-utilities' },
      },
    },
  })
  app.component('DataTable', DataTable).component('Column', Column).component('Menu', Menu)

  const i18n = createI18n({ legacy: false, locale: 'en', fallbackLocale: 'en', messages })
  app.use(i18n)

  bridge.onhostcontextchanged = () => usePreset(applyTheme(bridge.getHostContext()))

  app.mount('#app')
  await bridge.connect()
  usePreset(applyTheme(bridge.getHostContext()))

  return bridge
}

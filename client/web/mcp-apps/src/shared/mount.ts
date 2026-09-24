import { App as McpApp, applyDocumentTheme } from '@modelcontextprotocol/ext-apps'
import type { McpUiHostContext } from '@modelcontextprotocol/ext-apps'
import { usePreset } from '@primeuix/themes'
import { getTheme } from '@planetcrust/human-vue/src/composables/useTheme'
import PrimeVue from 'primevue/config'
import { type Component, createApp, provide, h } from 'vue'
import './styles.css'

export const mcpAppKey = Symbol('mcpApp')

type Theme = 'light' | 'dark'

function applyTheme(ctx?: McpUiHostContext) {
  const theme: Theme = ctx?.theme === 'dark' ? 'dark' : 'light'
  applyDocumentTheme(theme)
  return getTheme(theme)
}

// Mounts a view with Human's PrimeVue theme, following the host's light or
// dark, and connects it to the host. The connected App is provided under
// mcpAppKey; results arrive through its ontoolresult.
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

  bridge.onhostcontextchanged = () => usePreset(applyTheme(bridge.getHostContext()))

  app.mount('#app')
  await bridge.connect()
  usePreset(applyTheme(bridge.getHostContext()))

  return bridge
}

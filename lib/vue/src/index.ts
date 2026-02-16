// No CSS imports - client's Tailwind handles all styling

// Export Vue plugins
export { default as AuthPlugin } from './plugins/auth'
export {
  AutomationAPIPlugin,
  ComposeAPIPlugin,
  FederationAPIPlugin,
  SystemAPIPlugin,
} from './plugins/corteza-api'
export { I18nPlugin } from './plugins/i18n'
export { SettingsPlugin } from './plugins/settings'
export { PrimeVueComponentsPlugin } from './plugins/primevue-components'
export { ToastPlugin } from './plugins/toast'

// Export stores
export { useRBACStore } from './composables/useRBAC'
export { useComposeResourceStore } from './stores/useComposeResourceStore'
export { useConfirmDelete } from './composables/useConfirmDelete'
export { useResourceList } from './composables/useResourceList'
export { getTheme, setThemes, useTheme } from './composables/useTheme'
export { useMinDuration, withMinDuration } from './composables/useMinDuration'

// Export filters
export * as filters from './filters'

// Export components
export * as components from './components'

// Global rich text content styles — apply .rt-content class to any container showing Tiptap HTML
import './assets/css/rt-content.css'

// Export Vue plugins
export { default as AuthPlugin } from './plugins/auth'
export {
  AutomationAPIPlugin,
  ComposeAPIPlugin,
  DiscoveryAPIPlugin,
  FederationAPIPlugin,
  SystemAPIPlugin,
} from './plugins/human-api'
export { I18nPlugin } from './plugins/i18n'
export { SettingsPlugin } from './plugins/settings'
export { PrimeVueComponentsPlugin } from './plugins/primevue-components'
export { ToastPlugin } from './plugins/toast'
export { EventBusPlugin } from './plugins/event-bus'

// Export stores
export { useRBACStore } from './composables/useRBAC'
export { useApplicationsStore } from './stores/useApplicationsStore'
export { useComposeResourceStore } from './stores/useComposeResourceStore'
export { useNotificationsStore } from './stores/useNotificationsStore'
export { useAgentChatStore } from './stores/useAgentChatStore'
export { useWorkflowPromptsStore } from './stores/useWorkflowPromptsStore'
export { useRightSidebarStore } from './stores/useRightSidebarStore'
export { useConfirmDelete } from './composables/useConfirmDelete'
export { useResourceList } from './composables/useResourceList'
export { getTheme, setThemes, useTheme } from './composables/useTheme'
export { useMinDuration, withMinDuration } from './composables/useMinDuration'
export { useUserResolver } from './composables/useUserResolver'
export { useFileUpload } from './composables/useFileUpload'
export { providePermissions, usePermissions, PermissionsKey } from './composables/usePermissions'
export { useRightSidebarResize } from './composables/useRightSidebarResize'
export { useAgentRouteContextProvider } from './composables/useAgentRouteContextProvider'
export { useUnsavedGuard } from './composables/useUnsavedGuard'
export { resolveAppLogoUrl, appIconMap, defaultAppIcon } from './utils/appIcons'
export * as websocket from './libs/websocket'

// Export filters
export * as filters from './filters'

// Export components
export * as components from './components'

// Direct named exports for specific components
export { CEmojiPicker } from './components'
export { emojiData } from './components'
export { CChatbotInbox, makeChatbotInboxTranslations } from './components'
export { CAgentChat, CAgentSidebar, CAgentSidebarButton, CChatMessages, makeAgentChatTranslations } from './components'
export type { AgentChatTranslations } from './components/agent/translations'

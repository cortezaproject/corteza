/* eslint-disable no-unused-vars */
/**
 * Vue 3 Type Declarations for Corteza
 *
 * This file provides TypeScript type augmentations for Vue's global properties
 * and component definitions used throughout the Corteza applications.
 */

// Vue SFC module declaration
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<object, object, unknown>
  export default component
}

// Import types for augmentation
import type { Auth } from './plugins/auth'
import type { Settings } from './plugins/settings'

// API Client types - import from declaration files
import type AutomationAPI from '@cortezaproject/corteza-js-next/src/api-clients/automation'
import type ComposeAPI from '@cortezaproject/corteza-js-next/src/api-clients/compose'
import type FederationAPI from '@cortezaproject/corteza-js-next/src/api-clients/federation'
import type SystemAPI from '@cortezaproject/corteza-js-next/src/api-clients/system'

// Augment Vue's ComponentCustomProperties for global properties
declare module 'vue' {
  interface ComponentCustomProperties {
    /**
     * Authentication plugin instance
     * Provides OAuth2 authentication flow, user info, and token management
     */
    $Auth: Auth

    /**
     * System API client
     * Handles users, roles, settings, applications, etc.
     */
    $SystemAPI: SystemAPI

    /**
     * Compose API client
     * Handles namespaces, modules, records, pages, charts, etc.
     */
    $ComposeAPI: ComposeAPI

    /**
     * Automation API client
     * Handles workflows, triggers, sessions, etc.
     */
    $AutomationAPI: AutomationAPI

    /**
     * Federation API client
     * Handles federated nodes, sync, etc.
     */
    $FederationAPI: FederationAPI

    /**
     * Settings manager
     * Provides access to system settings with reactive updates
     */
    $Settings: Settings

    /**
     * Translation function from vue-i18n
     * @param key - Translation key
     * @param values - Interpolation values
     */
    $t: (key: string, ...args: unknown[]) => string
  }
}

// Augment Vue's GlobalComponents for globally registered PrimeVue components
declare module 'vue' {
  interface GlobalComponents {
    // PrimeVue Components (registered globally via UIPlugin)
    Button: (typeof import('primevue/button'))['default']
    Card: (typeof import('primevue/card'))['default']
    Checkbox: (typeof import('primevue/checkbox'))['default']
    InputText: (typeof import('primevue/inputtext'))['default']
    Select: (typeof import('primevue/select'))['default']
    Menu: (typeof import('primevue/menu'))['default']
    Sidebar: (typeof import('primevue/sidebar'))['default']
    TieredMenu: (typeof import('primevue/tieredmenu'))['default']
    IconField: (typeof import('primevue/iconfield'))['default']
    InputIcon: (typeof import('primevue/inputicon'))['default']
    Avatar: (typeof import('primevue/avatar'))['default']
    Toast: (typeof import('primevue/toast'))['default']
    Tree: (typeof import('primevue/tree'))['default']
  }
}

// Extend Window interface for Corteza configuration
declare global {
  interface Window {
    /**
     * Corteza API base URL
     * Set in public/config.js
     * @example 'https://corteza.example.com/api'
     */
    CortezaAPI?: string

    /**
     * Corteza Auth URL (optional, auto-derived from CortezaAPI if not set)
     * @example 'https://corteza.example.com/auth'
     */
    CortezaAuth?: string

    /**
     * Corteza Webapp base URL (optional)
     * Used for constructing callback URLs
     */
    CortezaWebapp?: string

    /**
     * Enable i18n pseudo mode for translation testing
     */
    i18nPseudoModeEnabled?: boolean
  }
}

// Vite environment variables
interface ImportMetaEnv {
  /**
   * Application display title
   */
  readonly VITE_APP_TITLE: string

  /**
   * Application identifier (used for auth callbacks)
   */
  readonly VITE_APP_ID: string

  /**
   * Application name (package name)
   */
  readonly VITE_APP_NAME: string

  /**
   * Base URL for the application
   */
  readonly BASE_URL: string

  /**
   * Current mode (development, production, etc.)
   */
  readonly MODE: string

  /**
   * Is production build
   */
  readonly PROD: boolean

  /**
   * Is development build
   */
  readonly DEV: boolean
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

// Global constants defined by Vite
declare const VERSION: string
declare const BUILD_TIME: string

export {}

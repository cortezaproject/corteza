<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <!-- Section whose registry application is disabled (enabled=false) -->
  <CAppDisabled v-if="!loading && !appEnabled" />

  <div v-else class="h-screen flex">
    <!-- Per-section left sidebar (only sections that declare one) -->
    <CSidebar v-if="hasSidebar" v-model="expanded">
      <template #body>
        <component :is="sidebarComponent" />
      </template>
    </CSidebar>

    <!-- Main content area -->
    <div
      class="flex-1 flex flex-col transition-[margin] duration-300 max-w-full"
      :style="{ marginLeft: contentMargin }"
    >
      <header>
        <CTopbar
          v-model:sidebar-expanded="expanded"
          :sidebar-disabled="sidebarDisabled"
          :settings="topbarSettings"
          :labels="{
            appMenu: $t('navigation.appMenu'),
            home: $t('navigation.home'),
            helpBuyHuman: $t('navigation.help.buyHuman'),
            helpManageSubscription: $t('navigation.help.manageSubscription'),
            helpPackageDetails: $t('navigation.help.packageDetails'),
            helpVersion: $t('navigation.help.version'),
            userSettingsProfile: $t('navigation.userSettings.profile'),
            userSettingsChangePassword: $t('navigation.userSettings.changePassword'),
            userSettingsLogout: $t('navigation.userSettings.logout'),
            userSettingsTheme: $t('navigation.userSettings.theme'),
            lightTheme: $t('navigation.themes.labels.light'),
            darkTheme: $t('navigation.themes.labels.dark'),
          }"
          :custom-profile-items="customProfileItems"
          @app-menu-click="rightSidebarStore.open('app-list')"
        >
          <template
            v-if="
              $Settings.get('discovery.enabled', false) && $Settings.get('ui.topbar.showSearch', true)
            "
            #right-tools
          >
            <Button
              v-tooltip.bottom="$t('navigation.search.label')"
              icon="pi pi-search"
              severity="secondary"
              variant="text"
              rounded
              @click="searchRef.open()"
            />
          </template>
        </CTopbar>
      </header>

      <main class="flex-1 overflow-hidden">
        <RouterView />
      </main>

      <Toast
        position="top-right"
        :pt="{
          root: {
            style: {
              top: 'var(--topbar-height)',
            },
          },
        }"
      />

      <ConfirmDialog />
    </div>

    <CAppListSidebar
      :labels="{
        title: $t('navigation.appList.title'),
        search: $t('navigation.appList.search'),
        noResults: $t('navigation.appList.noResults'),
        noApps: $t('navigation.appList.noApps'),
      }"
    />

    <CPrompts />
    <CNotificationSidebar />
    <CAgentSidebar :context-provider="agentContextProvider" />
    <CPermissionsDialog />
    <CTopbarSearch
      ref="searchRef"
      :labels="{
        placeholder: $t('navigation.search.placeholder'),
        noResults: () => $t('navigation.search.noResults'),
        notFoundNamespace: $t('navigation.search.notFoundNamespace'),
        notFoundPage: $t('navigation.search.notFoundPage'),
        recordRedirectError: $t('navigation.search.recordRedirectError'),
        recentSearches: $t('navigation.search.recentSearches'),
        clearHistory: $t('navigation.search.clearHistory'),
        openInNewTab: $t('navigation.search.openInNewTab'),
        numberOfResults: count => $t('navigation.search.numberOfResults', { count }),
      }"
    />
  </div>
</template>

<script setup>
import { sectionById } from '@/sections'
import {
  components,
  providePermissions,
  useAgentRouteContextProvider,
  useApplicationsStore,
  useNotificationsStore,
  useRBACStore,
  useRightSidebarStore,
  useWorkflowPromptsStore,
  websocket,
} from '@planetcrust/human-vue'
import { computed, inject, onBeforeUnmount, onMounted, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterView, useRoute } from 'vue-router'

// Shared resource stores (defined in compose, used app-wide). Provided here so
// every section + lib input reads ONE cache; core reference data is preloaded
// once at boot, the long tail (records, per-namespace modules) is on-demand.
import { useModuleStore } from '@/sections/compose/stores/module'
import { useNamespaceStore } from '@/sections/compose/stores/namespace'
import { useRecordStore } from '@/sections/compose/stores/record'
import { useUserStore } from '@/sections/compose/stores/user'

const {
  CTopbar,
  CLoaderLogo,
  CSidebar,
  CAppListSidebar,
  CPrompts,
  CNotificationSidebar,
  CPermissionsDialog,
  CAgentSidebar,
  CTopbarSearch,
  CAppDisabled,
} = components

providePermissions()

const APP_KIND = 'unify'

const { t } = useI18n()

const $Auth = inject('$Auth')
const $Settings = inject('$Settings')
const $SystemAPI = inject('$SystemAPI')
const $AutomationAPI = inject('$AutomationAPI')
const $ComposeAPI = inject('$ComposeAPI')
const $eventBus = inject('$eventBus', null)

const applicationsStore = useApplicationsStore()
const notificationsStore = useNotificationsStore()
const rightSidebarStore = useRightSidebarStore()
const workflowPromptsStore = useWorkflowPromptsStore()
const rbacStore = useRBACStore()

// Shared resource stores, provided app-wide (lib field inputs/viewers inject
// these; compose views import them directly — same singletons either way).
const usersStore = useUserStore()
const namespaceStore = useNamespaceStore()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()
provide('$userStore', usersStore)
provide('$recordStore', recordStore)
provide('$namespaceStore', namespaceStore)
provide('$moduleStore', moduleStore)

const route = useRoute()

// --- Section awareness -------------------------------------------------------
const activeSection = computed(() => sectionById(route.meta.section))

// Usability gate: a section whose registry application is disabled (enabled=
// false) shows the "disabled" screen, like the standalone apps did. Reactive
// to route.path so it re-evaluates on client-side navigation.
const appEnabled = computed(() => applicationsStore.isPathEnabled(route.path))

// Identity used for per-webapp filtering (e.g. workflow prompt `meta.webapps`).
// In the unified app the "current webapp" is the active section, so prompts
// keep the same visibility they had in the standalone apps.
const currentWebapp = computed(() => route.meta.section || APP_KIND)
const sidebarComponent = computed(() => activeSection.value?.sidebar || null)
const hasSidebar = computed(() => !!sidebarComponent.value)

const topbarSettings = computed(() => ({
  ...$Settings.get('ui.topbar', {}),
  ...(activeSection.value?.topbar || {}),
}))

// Agent context carries the active section so the agent sees where the user is;
// a section may enrich it (e.g. compose resolves namespace/page/record).
const baseAgentContext = useAgentRouteContextProvider(APP_KIND)
const agentContextProvider = () => {
  const ctx = { ...baseAgentContext(), section: route.meta.section || '' }
  const enrich = activeSection.value?.agentContext
  return enrich ? enrich(ctx) : ctx
}

// Topbar profile-menu items contributed by the active section (e.g. compose
// reminders). `t` is passed in so section modules don't call useI18n() outside setup.
const customProfileItems = computed(() => activeSection.value?.profileItems?.(t) || [])

// --- Sidebar expand/collapse -------------------------------------------------
const expanded = ref(false)
const isMobile = ref(window.innerWidth < 1024)

const sidebarDisabled = computed(() => {
  if (!hasSidebar.value) return true
  const disabled = activeSection.value?.sidebarDisabledRoutes || []
  return disabled.includes(route.name?.toString() || '')
})

const contentMargin = computed(() =>
  !isMobile.value && expanded.value && !sidebarDisabled.value ? 'var(--sidebar-width)' : '0',
)

const handleResize = () => {
  isMobile.value = window.innerWidth < 1024
}

watch(
  sidebarDisabled,
  disabled => {
    if (disabled) expanded.value = false
    else if (activeSection.value?.autoExpandSidebar && !isMobile.value) expanded.value = true
  },
  { immediate: true },
)

watch(
  () => route.meta.section,
  (newSection, oldSection) => {
    if (newSection !== oldSection) rightSidebarStore.closeSectionPanels()
  },
)

// --- Shell state -------------------------------------------------------------
const searchRef = ref(null)
const loading = ref(true)
let realtimeClient

const logoUrl = computed(() => $Settings.attachment('ui.mainLogo'))

const baseTitle = import.meta.env.VITE_APP_TITLE || 'Human'
watch(
  () => notificationsStore.unreadCount,
  count => {
    document.title = count > 0 ? `(${count}) ${baseTitle}` : baseTitle
  },
)

onMounted(async () => {
  window.addEventListener('resize', handleResize)

  const fetchPromise = Promise.all([
    applicationsStore.fetchApplications(),
    notificationsStore.fetchNotifications($SystemAPI),
    workflowPromptsStore.update($AutomationAPI, currentWebapp.value),
    rbacStore.load([$SystemAPI, $AutomationAPI, $ComposeAPI]),
    // Preload bounded reference data once; both are cache-guarded so they
    // won't refetch on later navigation. Records / per-namespace modules stay
    // on-demand and populate their shared store lazily.
    namespaceStore.load(),
    usersStore.load({ limit: 500 }),
  ])
  const delayPromise = new Promise(resolve => setTimeout(resolve, 2000))

  Promise.all([fetchPromise, delayPromise]).finally(() => {
    loading.value = false
  })

  realtimeClient = websocket.createRealtimeClient({
    auth: $Auth,
    onMessage: ({ data }) => {
      const msg = JSON.parse(data)
      // Re-broadcast every realtime message so active sections can react to
      // their own types (e.g. compose handles 'reminder').
      $eventBus?.emit('realtime', msg)
      switch (msg['@type']) {
        case 'workflowSessionPrompt':
          workflowPromptsStore.newPrompt(msg['@value'], currentWebapp.value)
          break
        case 'workflowSessionResumed':
          workflowPromptsStore.clear(msg['@value'])
          break
        case 'notification':
          notificationsStore.addNotification(msg['@value'])
          break
        case 'notification.read':
          notificationsStore.updateReadNotification(msg['@value'])
          break
        case 'notification.unread':
          notificationsStore.updateUnreadNotification(msg['@value'])
          break
        case 'notification.read.all':
          notificationsStore.updateAllReadNotifications(msg['@value'])
          break
        case 'notification.unread.all':
          notificationsStore.updateAllUnreadNotifications(msg['@value'])
          break
        case 'notification.delete':
          notificationsStore.removeNotification(msg['@value'])
          break
      }
    },
  })

  realtimeClient.connect()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  realtimeClient?.disconnect?.()
})
</script>

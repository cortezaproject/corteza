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

    <!-- Main content area. min-w-0 is load-bearing: without it the column's
         automatic minimum size is its content's min-content width, so a wide
         block (a record list with many columns) makes it wider than the space
         beside the sidebar and the margin carries the topbar off screen. -->
    <div
      class="flex-1 min-w-0 flex flex-col transition-[margin] duration-300"
      :style="{ marginLeft: contentMargin }"
    >
      <header>
        <CTopbar
          v-model:sidebar-expanded="expanded"
          :sidebar-disabled="sidebarDisabled"
          :settings="topbarSettings"
          :home-url="topbarHomeURL"
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
              $Settings.get('discovery.enabled', false) &&
              $Settings.get('ui.topbar.showSearch', true)
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
        position="bottom-right"
        :pt="{
          root: {
            style: {
              bottom: '4rem',
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
        tileMenu: $t('navigation.appList.tileMenu'),
        setHome: $t('navigation.appList.setHome'),
        clearHome: $t('navigation.appList.clearHome'),
        homeSet: $t('navigation.appList.homeSet', { app: '{app}' }),
        homeSetTitle: $t('navigation.appList.homeSetTitle'),
        homeClearedTitle: $t('navigation.appList.homeClearedTitle'),
        homeErrorTitle: $t('navigation.appList.homeErrorTitle'),
        homeCleared: $t('navigation.appList.homeCleared', { app: '{app}' }),
        homeOwn: $t('navigation.appList.homeOwn'),
        homeGlobal: $t('navigation.appList.homeGlobal'),
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
import { homeHref } from '@/router'
import { sectionById } from '@/sections'
import { previewsSwitchedOffApp } from '@/sections/app/preview'
import {
  components,
  providePermissions,
  useAgentRouteContextProvider,
  useApplicationsStore,
  useNotificationsStore,
  useRBACStore,
  useRightSidebarStore,
  useSystemNotifications,
  useWorkflowPromptsStore,
  websocket,
} from '@planetcrust/human-vue'
import { computed, inject, onBeforeUnmount, onMounted, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterView, useRoute } from 'vue-router'

import { useModuleStore, useNamespaceStore, useUserStore } from '@planetcrust/human-vue'
import { appIconMap } from '@/utils/appIcons'
import { useAppReachable } from '@/utils/appReachable'
import { useDocumentTitle } from '@/utils/documentTitle'
import { appIconImage, useFavicon } from '@/utils/favicon'

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
const $eventBus = inject('$eventBus', null)

const applicationsStore = useApplicationsStore()
const notificationsStore = useNotificationsStore()
// The configured app icon (a custom upload, or the default).
const iconUrl = computed(() => $Settings.attachment('ui.iconLogo'))
const systemNotifications = useSystemNotifications({ icon: () => appIconImage(iconUrl.value) })
const rightSidebarStore = useRightSidebarStore()
const workflowPromptsStore = useWorkflowPromptsStore()
const rbacStore = useRBACStore()

const usersStore = useUserStore()
const namespaceStore = useNamespaceStore()
const moduleStore = useModuleStore()
provide('$appIconMap', appIconMap)
// The app menu offers exactly what the router's section gate will admit.
provide('$appReachable', useAppReachable())

const route = useRoute()

// --- Section awareness -------------------------------------------------------
const activeSection = computed(() => sectionById(route.meta.section))

// Usability gate: a section whose registry application is disabled (enabled=
// false) shows the "disabled" screen, like the standalone apps did. Reactive
// to route.path so it re-evaluates on client-side navigation. A switched-off
// custom app stays open, as a preview, to whoever may change its page.
const appEnabled = computed(
  () =>
    applicationsStore.isPathEnabled(route.path) ||
    previewsSwitchedOffApp(applicationsStore.appForPath(route.path)),
)

// Identity used for per-webapp filtering (e.g. workflow prompt `meta.webapps`).
// In the unified app the "current webapp" is the active section, so prompts
// keep the same visibility they had in the standalone apps.
const currentWebapp = computed(() => route.meta.section || APP_KIND)
const sidebarComponent = computed(() => activeSection.value?.sidebar || null)
const hasSidebar = computed(() => !!sidebarComponent.value)

// The home button leads to the home application, like opening `/` does.
const topbarHomeURL = computed(() => homeHref())

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
  // Routes opt out of the sidebar via `meta: { hideSidebar: true }` (refactor-
  // safe — the flag travels with the route, so renames can't silently break it).
  return !!route.meta.hideSidebar
})

const contentMargin = computed(() =>
  !isMobile.value && expanded.value && !sidebarDisabled.value ? 'var(--sidebar-width)' : '0',
)

const handleResize = () => {
  isMobile.value = window.innerWidth < 1024
}

// Sidebar expand/collapse is remembered per webapp in localStorage. First visit
// to a section falls back to its `sidebarExpandedByDefault` flag; after that the
// user's last manual choice wins.
const SIDEBAR_PREFS_KEY = 'ui.sidebar.expanded'

const loadSidebarPrefs = () => {
  try {
    return JSON.parse(localStorage.getItem(SIDEBAR_PREFS_KEY)) || {}
  } catch {
    return {}
  }
}

const saveSidebarPref = (webapp, value) => {
  const prefs = loadSidebarPrefs()
  prefs[webapp] = value
  localStorage.setItem(SIDEBAR_PREFS_KEY, JSON.stringify(prefs))
}

const resolveExpanded = () => {
  if (sidebarDisabled.value || isMobile.value) return false
  const prefs = loadSidebarPrefs()
  const key = currentWebapp.value
  if (key in prefs) return prefs[key]
  return !!activeSection.value?.sidebarExpandedByDefault
}

// Guards the persistence watcher against programmatic (resolve-driven) writes so
// only genuine user toggles are stored. Relies on flush:'sync' below.
let applyingResolved = false
const applyResolved = () => {
  applyingResolved = true
  expanded.value = resolveExpanded()
  applyingResolved = false
}

watch([sidebarDisabled, currentWebapp, isMobile], applyResolved, { immediate: true })

watch(
  expanded,
  value => {
    if (applyingResolved || sidebarDisabled.value || isMobile.value) return
    saveSidebarPref(currentWebapp.value, value)
  },
  { flush: 'sync' },
)

watch(
  () => route.meta.section,
  (newSection, oldSection) => {
    if (newSection !== oldSection) rightSidebarStore.closeSectionPanels()
  },
)

// Sections may declare a `preload` the shell runs when they become active, so
// section-wide data (e.g. the projects list feeding the sidebar) is loaded on
// app load / section entry rather than lazily when a drawer is first opened.
watch(activeSection, section => section?.preload?.(), { immediate: true })

// --- Shell state -------------------------------------------------------------
const searchRef = ref(null)
const loading = ref(true)
let realtimeClient

const logoUrl = computed(() => $Settings.attachment('ui.mainLogo'))

// Tab title = the heading the active view teleports into the topbar, with the
// unread count kept in front of it unless notifications are muted.
useDocumentTitle(() => notificationsStore.badgeCount)
useFavicon(iconUrl, () => notificationsStore.badgeCount > 0)

onMounted(async () => {
  window.addEventListener('resize', handleResize)

  // allSettled, not all: every one of these is an optional cache, and being
  // refused one is an ordinary state now that permissions are deny-by-default —
  // a user with no roles may list neither users nor namespaces. Promise.all
  // rejects on the first refusal and the chain below has no catch, so the shell
  // booted with an uncaught error for exactly the users it is meant to serve.
  const fetchPromise = Promise.allSettled([
    // ready(), not fetchApplications(): the router gate already asked for the
    // list before this mounted, and re-fetching would double the request.
    applicationsStore.ready(),
    notificationsStore.fetchNotifications(),
    workflowPromptsStore.update(currentWebapp.value),
    rbacStore.load(),
    // Preload bounded reference data once; both are cache-guarded so they
    // won't refetch on later navigation. Records / per-namespace modules stay
    // on-demand and populate their shared store lazily.
    namespaceStore.load(),
    usersStore.load({ limit: 500 }),
  ]).then(results => {
    // Refused is expected; failed for any other reason is not, and silence
    // would make a broken preload look like empty data.
    results
      .filter(r => r.status === 'rejected')
      .forEach(r => console.warn('shell preload skipped:', r.reason?.message || r.reason))
  })

  // Brief splash floor so the logo doesn't flash-and-vanish on fast loads,
  // without forcing a long wait when data is ready sooner.
  const delayPromise = new Promise(resolve => setTimeout(resolve, 1000))

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
      // Notification.* types are owned by the notifications store; a new one
      // is also offered to the OS.
      if (notificationsStore.handleRealtime(msg)) {
        if (msg['@type'] === 'notification') systemNotifications.notify(msg['@value'])
        return
      }
      switch (msg['@type']) {
        case 'workflowSessionPrompt':
          workflowPromptsStore.newPrompt(msg['@value'], currentWebapp.value)
          break
        case 'workflowSessionResumed':
          workflowPromptsStore.clear(msg['@value'])
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

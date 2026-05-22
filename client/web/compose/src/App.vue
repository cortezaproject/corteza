<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <CAppDisabled v-if="!loading && !appEnabled" />

  <div v-else class="h-screen flex">
    <!-- Sidebar: fixed position via Drawer -->
    <CSidebar v-model="expanded">
      <template v-if="!sidebarDisabled" #header>
        <CSidebarNamespaceSwitcher />
      </template>
      <template v-if="!sidebarDisabled" #body>
        <CSidebarNavigation />
      </template>
    </CSidebar>

    <!-- Main content area: pushed by sidebar margin on desktop -->
    <div
      class="flex-1 flex flex-col min-w-0 transition-[margin] duration-300"
      :style="{ marginLeft: contentMargin }"
    >
      <header>
        <CTopbar
          v-model:sidebar-expanded="expanded"
          :sidebar-disabled="sidebarDisabled"
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
            lightTheme: $t('general.themes.labels.light'),
            darkTheme: $t('general.themes.labels.dark'),
          }"
          :custom-profile-items="profileReminderItems"
          @app-menu-click="appListVisible = true"
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
        position="top-right"
        :pt="{
          root: {
            style: {
              top: 'calc(var(--topbar-height) + 1rem)',
            },
          },
        }"
      />

      <ConfirmDialog />
    </div>

    <CAppListSidebar
      v-model:visible="appListVisible"
      :labels="{
        title: $t('navigation.appList.title'),
        search: $t('navigation.appList.search'),
        noResults: $t('navigation.appList.noResults'),
        noApps: $t('navigation.appList.noApps'),
      }"
    />

    <CPrompts />
    <CNotificationSidebar />
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
    <CAgentSidebar :context-provider="agentContextProvider" />
    <ReminderSidebar />
    <ReminderToastHost />
    <CPermissionsDialog />
    <CTranslatorDialog />
  </div>
</template>

<script setup>
import CSidebarNamespaceSwitcher from '@/components/CSidebarNamespaceSwitcher.vue'
import CSidebarNavigation from '@/components/CSidebarNavigation.vue'
import ReminderSidebar from '@/components/Reminders/ReminderSidebar.vue'
import ReminderToastHost from '@/components/Reminders/ReminderToastHost.vue'
import CTranslatorDialog from '@/components/Translator/CTranslatorDialog.vue'
import { useModuleStore } from '@/stores/module'
import { useNamespaceStore } from '@/stores/namespace'
import { usePageStore } from '@/stores/page'
import { useRecordStore } from '@/stores/record'
import { useReminderStore } from '@/stores/reminder'
import { useUserStore } from '@/stores/user'
import {
  components,
  providePermissions,
  useAgentRouteContextProvider,
  useApplicationsStore,
  useNotificationsStore,
  useRBACStore,
  useWorkflowPromptsStore,
  websocket,
} from '@planetcrust/human-vue'
import { computed, inject, onBeforeUnmount, onMounted, provide, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterView, useRoute } from 'vue-router'
const {
  CTopbar,
  CLoaderLogo,
  CSidebar,
  CAppListSidebar,
  CPrompts,
  CNotificationSidebar,
  CPermissionsDialog,
  CAgentSidebar,
  CAppDisabled,
  CTopbarSearch,
} = components

// Provide permissions dialog context for the entire app
providePermissions()

const { t } = useI18n()
const $Auth = inject('$Auth')
const $Settings = inject('$Settings')
const $ComposeAPI = inject('$ComposeAPI')
const $SystemAPI = inject('$SystemAPI')
const $AutomationAPI = inject('$AutomationAPI')

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

const namespaceStore = useNamespaceStore()
const usersStore = useUserStore()
const recordStore = useRecordStore()
const moduleStore = useModuleStore()
const pageStore = usePageStore()

// Base route-only context — every webapp uses this; compose layers
// resolved entities on top.
const baseContextProvider = useAgentRouteContextProvider('compose')

function agentContextProvider() {
  const ctx = baseContextProvider()
  const params = ctx.routeParams || {}

  // Resolve namespace by URL part (slug or numeric ID, depending on the route).
  if (params.slug) {
    const ns = namespaceStore.getByUrlPart?.(params.slug)
    if (ns) {
      ctx.namespace = {
        namespaceID: String(ns.namespaceID),
        slug: ns.slug,
        name: ns.name,
      }
    }
  }

  if (params.pageID) {
    const page = pageStore.getByID?.(params.pageID)
    if (page) {
      ctx.page = {
        pageID: String(page.pageID),
        handle: page.handle,
        title: page.title,
        moduleID: page.moduleID ? String(page.moduleID) : undefined,
      }
    }
  }

  if (params.moduleID) {
    const mod = moduleStore.getByID?.(params.moduleID)
    if (mod) {
      ctx.module = {
        moduleID: String(mod.moduleID),
        handle: mod.handle,
        name: mod.name,
      }
    }
  }

  if (params.recordID) {
    const record =
      recordStore.records?.get?.(params.recordID) ||
      recordStore.labelCache?.get?.(params.recordID)
    if (record) {
      ctx.record = {
        recordID: String(record.recordID),
        moduleID: record.moduleID ? String(record.moduleID) : undefined,
      }
      if (record.values) {
        if (Array.isArray(record.values)) {
          ctx.record.values = record.values.reduce((acc, v) => {
            if (v && v.name) acc[v.name] = v.value
            return acc
          }, {})
        } else if (typeof record.values === 'object') {
          ctx.record.values = { ...record.values }
        }
      }
    }
  }

  return ctx
}
const rbacStore = useRBACStore()
const applicationsStore = useApplicationsStore()
const notificationsStore = useNotificationsStore()
const workflowPromptsStore = useWorkflowPromptsStore()
const reminderStore = useReminderStore()

const appListVisible = ref(false)
const searchRef = ref(null)
const appEnabled = computed(() => applicationsStore.isCurrentAppEnabled())
let realtimeClient

const profileReminderItems = computed(() => {
  const label =
    reminderStore.activeCount > 0
      ? `${t('navigation.userSettings.reminders')} (${reminderStore.activeCount})`
      : t('navigation.userSettings.reminders')

  return [
    {
      label,
      icon: 'pi pi-clock',
      command: () => {
        reminderStore.toggleVisibility()
      },
    },
  ]
})

// Provide stores to field editor/viewer components in lib/vue
provide('$userStore', usersStore)
provide('$recordStore', recordStore)

onMounted(() => {
  const fetchPromises = [
    namespaceStore.load({ force: true }),
    usersStore.load({ limit: 500 }),
    rbacStore.load([$ComposeAPI, $SystemAPI]),
    applicationsStore.fetchApplications(),
    notificationsStore.fetchNotifications($SystemAPI),
    reminderStore.fetchReminders(),
    workflowPromptsStore.update($AutomationAPI, 'compose'),
  ]
  const delayPromise = new Promise(resolve => setTimeout(resolve, 1000))

  Promise.all([...fetchPromises, delayPromise]).finally(() => {
    loading.value = false
  })

  realtimeClient = websocket.createRealtimeClient({
    auth: $Auth,
    onMessage: ({ data }) => {
      const msg = JSON.parse(data)
      switch (msg['@type']) {
        case 'workflowSessionPrompt':
          workflowPromptsStore.newPrompt(msg['@value'], 'compose')
          break
        case 'workflowSessionResumed':
          workflowPromptsStore.clear(msg['@value'])
          break
        case 'notification':
          notificationsStore.addNotification(msg['@value'])
          break
        case 'reminder':
          reminderStore.handleRealtimeReminder(msg['@value'])
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

const expanded = ref(false)
const isMobile = ref(window.innerWidth < 1024)

const contentMargin = computed(() => {
  // Push on desktop only when sidebar is expanded and not on a disabled route
  return !isMobile.value && expanded.value && !sidebarDisabled.value ? 'var(--sidebar-width)' : '0'
})

// Handle resize to detect mobile
const handleResize = () => {
  isMobile.value = window.innerWidth < 1024
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  realtimeClient?.disconnect?.()
  reminderStore.dispose()
})

// Route-based sidebar control
const route = useRoute()

// Routes where sidebar should be disabled
const disabledRoutes = [
  'namespaces',
  'namespace.list',
  'namespace.edit',
  'namespace.create',
  'namespace.clone',
]

const sidebarDisabled = computed(() => {
  return disabledRoutes.includes(route.name?.toString() || '')
})

// Control sidebar state based on route
watch(
  sidebarDisabled,
  disabled => {
    if (disabled) {
      // Close sidebar on disabled routes
      expanded.value = false
    } else if (!isMobile.value) {
      // Auto-expand on desktop for enabled routes
      expanded.value = true
    }
  },
  { immediate: true },
)
</script>

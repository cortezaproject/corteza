<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <CAppDisabled v-if="!loading && !appEnabled" />

  <div v-else class="h-screen flex">
    <!-- Sidebar: fixed position via Drawer -->
    <CSidebar v-model="expanded">
      <template #body>
        <CSidebarNavigation />
      </template>
    </CSidebar>

    <!-- Main content area: pushed by sidebar margin on desktop -->
    <div
      class="flex-1 flex flex-col transition-[margin] duration-300 max-w-full"
      :style="{ marginLeft: contentMargin }"
    >
      <header>
        <CTopbar
          v-model:sidebar-expanded="expanded"
          :sidebar-disabled="sidebarDisabled"
          :settings="$Settings.get('ui.topbar', {})"
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
          @app-menu-click="appListVisible = true"
        >
          <template
            v-if="$Settings.get('discovery.enabled', false) && $Settings.get('ui.topbar.showSearch', true)"
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
          messageIcon: { class: 'hidden' },
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
    <CAgentSidebar />
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
import CSidebarNavigation from '@/components/CSidebarNavigation.vue'
import {
  components,
  providePermissions,
  withMinDuration,
  useApplicationsStore,
  useNotificationsStore,
  useRBACStore,
  useWorkflowPromptsStore,
  websocket,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import { useAutomationStore } from '@/stores/automation'
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

providePermissions()

const $Auth = inject('$Auth')
const $Settings = inject('$Settings')
const $AutomationAPI = inject('$AutomationAPI')
const $SystemAPI = inject('$SystemAPI')
const store = useAutomationStore()
const applicationsStore = useApplicationsStore()
const notificationsStore = useNotificationsStore()
const workflowPromptsStore = useWorkflowPromptsStore()
const rbacStore = useRBACStore()

const appListVisible = ref(false)
const searchRef = ref(null)
const appEnabled = computed(() => applicationsStore.isCurrentAppEnabled())
let realtimeClient

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

onMounted(async () => {
  window.addEventListener('resize', handleResize)

  try {
    await withMinDuration(
      Promise.all([
        store.loadCatalog($AutomationAPI),
        store.fetchList($AutomationAPI),
        applicationsStore.fetchApplications(),
        notificationsStore.fetchNotifications($SystemAPI),
        workflowPromptsStore.update($AutomationAPI, 'taq'),
        rbacStore.load([$SystemAPI, $AutomationAPI]),
      ]),
    )
  } catch (e) {
    console.error('Failed to load catalog:', e)
  } finally {
    loading.value = false
  }

  realtimeClient = websocket.createRealtimeClient({
    auth: $Auth,
    onMessage: ({ data }) => {
      const msg = JSON.parse(data)
      switch (msg['@type']) {
        case 'workflowSessionPrompt':
          workflowPromptsStore.newPrompt(msg['@value'], 'taq')
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

const expanded = ref(false)
const isMobile = ref(window.innerWidth < 1024)

// Route-based sidebar control
const route = useRoute()
const disabledRoutes = ['list'] // "disable it on taq list"

const sidebarDisabled = computed(() => {
  return disabledRoutes.includes(route.name?.toString() || '')
})

const contentMargin = computed(() => {
  return !isMobile.value && expanded.value && !sidebarDisabled.value ? 'var(--sidebar-width)' : '0'
})

const handleResize = () => {
  isMobile.value = window.innerWidth < 1024
}

watch(
  sidebarDisabled,
  disabled => {
    if (disabled) {
      expanded.value = false
    }
  },
  { immediate: true },
)
</script>

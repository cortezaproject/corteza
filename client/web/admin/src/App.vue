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
      class="flex-1 flex flex-col transition-[margin] duration-300"
      :style="{ marginLeft: contentMargin }"
    >
      <header>
        <CTopbar
          v-model:sidebar-expanded="expanded"
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
        />
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
  </div>
</template>

<script setup>
import CSidebarNavigation from '@/components/CSidebarNavigation.vue'
import {
  components,
  providePermissions,
  useApplicationsStore,
  useNotificationsStore,
  useRBACStore,
  useWorkflowPromptsStore,
  websocket,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
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
} = components

// Provide permissions dialog context for the entire app
providePermissions()

const $Auth = inject('$Auth')
const $Settings = inject('$Settings')
const $SystemAPI = inject('$SystemAPI')
const $AutomationAPI = inject('$AutomationAPI')

const applicationsStore = useApplicationsStore()
const notificationsStore = useNotificationsStore()
const workflowPromptsStore = useWorkflowPromptsStore()
const rbacStore = useRBACStore()
const appListVisible = ref(false)
const appEnabled = computed(() => applicationsStore.isCurrentAppEnabled())
let realtimeClient

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

onMounted(() => {
  const delayPromise = new Promise(resolve => setTimeout(resolve, 1500))

  Promise.all([
    delayPromise,
    applicationsStore.fetchApplications(),
    notificationsStore.fetchNotifications($SystemAPI),
    workflowPromptsStore.update($AutomationAPI, 'admin'),
    rbacStore.load([$SystemAPI, $AutomationAPI]),
  ]).finally(() => {
    loading.value = false
  })

  realtimeClient = websocket.createRealtimeClient({
    auth: $Auth,
    onMessage: ({ data }) => {
      const msg = JSON.parse(data)
      switch (msg['@type']) {
        case 'workflowSessionPrompt':
          workflowPromptsStore.newPrompt(msg['@value'], 'admin')
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

const expanded = ref(false)
const isMobile = ref(window.innerWidth < 1024)

const contentMargin = computed(() => {
  return !isMobile.value && expanded.value ? 'var(--sidebar-width)' : '0'
})

const handleResize = () => {
  isMobile.value = window.innerWidth < 1024
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  // Auto-expand on desktop
  if (!isMobile.value) {
    expanded.value = true
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  realtimeClient?.disconnect?.()
})
</script>

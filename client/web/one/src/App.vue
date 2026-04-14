<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <div class="h-screen flex flex-col">
    <header>
      <CTopbar
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          appMenu: $t('navigation.appMenu'),
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
        :hide-app-selector="true"
        :sidebar-disabled="true"
        :hide-logo="true"
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

    <CPrompts />
    <CAgentSidebar />
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
  </div>
</template>

<script setup>
import {
  components,
  useApplicationsStore,
  useNotificationsStore,
  useWorkflowPromptsStore,
  websocket,
} from '@cortezaproject/corteza-vue-next'
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
const { CTopbar, CLoaderLogo, CPrompts, CNotificationSidebar, CAgentSidebar, CTopbarSearch } = components

const $Auth = inject('$Auth')
const $Settings = inject('$Settings')
const $SystemAPI = inject('$SystemAPI')
const $AutomationAPI = inject('$AutomationAPI')

const applicationsStore = useApplicationsStore()
const notificationsStore = useNotificationsStore()
const workflowPromptsStore = useWorkflowPromptsStore()
const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)
const searchRef = ref(null)
let realtimeClient

onMounted(() => {
  const fetchPromise = Promise.all([
    applicationsStore.fetchApplications(),
    notificationsStore.fetchNotifications($SystemAPI),
    workflowPromptsStore.update($AutomationAPI, 'one'),
  ])
  const delayPromise = new Promise(resolve => setTimeout(resolve, 2000))

  Promise.all([fetchPromise, delayPromise]).finally(() => {
    loading.value = false
  })

  realtimeClient = websocket.createRealtimeClient({
    auth: $Auth,
    onMessage: ({ data }) => {
      const msg = JSON.parse(data)
      switch (msg['@type']) {
        case 'workflowSessionPrompt':
          workflowPromptsStore.newPrompt(msg['@value'], 'one')
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
  realtimeClient?.disconnect?.()
})
</script>

<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <div class="h-screen flex flex-col overflow-hidden">
    <header>
      <CTopbar
        :settings="{
          ...$Settings.get('ui.topbar', {}),
          hideAppSelector: true,
          hideAgentSidebar: true,
          hideNotifications: true,
        }"
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
        :sidebar-disabled="true"
        :hide-logo="true"
      />
    </header>

    <Teleport to="#topbar-title" defer>
      <img
        v-if="logoUrl"
        :src="logoUrl"
        class="h-8 max-w-[12rem] object-contain object-left"
        alt="Logo"
      />
    </Teleport>

    <main class="flex-1 overflow-hidden min-h-0">
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
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterView } from 'vue-router'
const { CTopbar, CLoaderLogo, CPrompts } = components

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
let realtimeClient

const baseTitle = import.meta.env.VITE_APP_TITLE || 'Two'
watch(
  () => notificationsStore.unreadCount,
  count => {
    document.title = count > 0 ? `(${count}) ${baseTitle}` : baseTitle
  },
)

onMounted(() => {
  const fetchPromise = Promise.all([
    applicationsStore.fetchApplications(),
    notificationsStore.fetchNotifications($SystemAPI),
    workflowPromptsStore.update($AutomationAPI, 'two'),
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
          workflowPromptsStore.newPrompt(msg['@value'], 'two')
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

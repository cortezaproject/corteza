<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <div class="h-screen flex flex-col">
    <header>
      <CTopbar
        :sidebar-expanded="expanded"
        :labels="{
          appMenu: $t('navigation.appMenu'),
          helpForum: $t('navigation.help.forum'),
          helpDocumentation: $t('navigation.help.documentation'),
          helpFeedback: $t('navigation.help.feedback'),
          helpVersion: $t('navigation.help.version'),
          userSettingsProfile: $t('navigation.userSettings.profile'),
          userSettingsChangePassword: $t('navigation.userSettings.changePassword'),
          userSettingsLogout: $t('navigation.userSettings.logout'),
          userSettingsTheme: $t('navigation.userSettings.theme'),
          lightTheme: $t('general.themes.labels.light'),
          darkTheme: $t('general.themes.labels.dark'),
        }"
      />
    </header>

    <CSidebar v-model="expanded" :expand-on-click="true" :disabled-routes="disabledRoutes" />

    <main class="flex-1 overflow-hidden">
      <RouterView />
    </main>

    <Toast
      :pt="{
        root: {
          style: {
            top: 'calc(var(--topbar-height) + 20px)',
            right: '17px',
          },
        },
        messageIcon: {
          style: {
            display: 'none',
          },
        },
      }"
    />
  </div>
</template>

<script setup>
import { useNamespaceStore } from '@/stores/namespace'
import { useUserStore } from '@/stores/user'
import { components } from '@cortezaproject/corteza-vue-next'
import Toast from 'primevue/toast'
import { computed, inject, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
const { CTopbar, CLoaderLogo, CSidebar } = components

const $Settings = inject('$Settings')

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

const namespaceStore = useNamespaceStore()
const usersStore = useUserStore()

onMounted(() => {
  const fetchPromises = [namespaceStore.load({ force: true }), usersStore.load({ limit: 500 })]
  const delayPromise = new Promise(resolve => setTimeout(resolve, 2000))

  Promise.all([...fetchPromises, delayPromise]).finally(() => {
    loading.value = false
  })
})

const expanded = ref(false)

const disabledRoutes = [
  'namespaces',
  'namespace.list',
  'namespace.edit',
  'namespace.create',
  'namespace.clone',
  'namespace.manage',
]
</script>

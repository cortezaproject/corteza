<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <div class="h-screen flex">
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
            helpForum: $t('navigation.help.forum'),
            helpDocumentation: $t('navigation.help.documentation'),
            helpFeedback: $t('navigation.help.feedback'),
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
  </div>
</template>

<script setup>
import CSidebarNavigation from '@/components/CSidebarNavigation.vue'
import { components, withMinDuration, useApplicationsStore } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import { useAutomationStore } from '@/stores/automation'
const { CTopbar, CLoaderLogo, CSidebar, CAppListSidebar } = components

const $Settings = inject('$Settings')
const $AutomationAPI = inject('$AutomationAPI')
const store = useAutomationStore()
const applicationsStore = useApplicationsStore()

const appListVisible = ref(false)

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
      ]),
    )
  } catch (e) {
    console.error('Failed to load catalog:', e)
  } finally {
    loading.value = false
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
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

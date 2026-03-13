<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <div class="h-screen flex">
    <!-- Sidebar -->
    <CSidebar v-model="expanded">
      <template #body>
        <CSidebarNavigation />
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
          :hide-app-selector="true"
        />
      </header>

      <main class="flex-1 overflow-hidden">
        <RouterView />
      </main>

      <Toast
        position="bottom-left"
        :pt="{
          root: {
            style: {
              bottom: 'var(--topbar-height)',
              left: '17px',
            },
          },
          messageIcon: {
            style: {
              display: 'none',
            },
          },
        }"
      />

      <ConfirmDialog />
    </div>
  </div>
</template>

<script setup>
import CSidebarNavigation from '@/components/CSidebarNavigation.vue'
import { useAgentStore } from '@/stores/agent'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { RouterView, useRoute } from 'vue-router'
const { CTopbar, CLoaderLogo, CSidebar } = components

const $Settings = inject('$Settings')
const $SystemAPI = inject('$SystemAPI')
const agentStore = useAgentStore()

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

onMounted(async () => {
  window.addEventListener('resize', handleResize)

  try {
    await agentStore.fetchList($SystemAPI)
  } catch (e) {
    console.error('Failed to load agents:', e)
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
const disabledRoutes = ['root']

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

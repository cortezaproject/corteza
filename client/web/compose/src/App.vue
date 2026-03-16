<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <div class="h-screen flex">
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
      class="flex-1 flex flex-col transition-[margin] duration-300"
      :style="{ marginLeft: contentMargin }"
    >
      <header>
        <CTopbar
          v-model:sidebar-expanded="expanded"
          :sidebar-disabled="sidebarDisabled"
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

      <main class="flex-1 overflow-hidden">
        <RouterView />
      </main>

      <Toast
        position="top-center"
        :pt="{
          root: {
            style: {
              top: 'var(--topbar-height)',
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
import CSidebarNamespaceSwitcher from '@/components/CSidebarNamespaceSwitcher.vue'
import CSidebarNavigation from '@/components/CSidebarNavigation.vue'
import { useNamespaceStore } from '@/stores/namespace'
import { useRecordStore } from '@/stores/record'
import { useUserStore } from '@/stores/user'
import { components, useRBACStore } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onBeforeUnmount, onMounted, provide, ref, watch } from 'vue'
import { RouterView, useRoute } from 'vue-router'
const { CTopbar, CLoaderLogo, CSidebar } = components

const $Settings = inject('$Settings')
const $ComposeAPI = inject('$ComposeAPI')
const $SystemAPI = inject('$SystemAPI')

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

const namespaceStore = useNamespaceStore()
const usersStore = useUserStore()
const recordStore = useRecordStore()
const rbacStore = useRBACStore()

// Provide stores to field editor/viewer components in lib/vue
provide('$userStore', usersStore)
provide('$recordStore', recordStore)

onMounted(() => {
  const fetchPromises = [
    namespaceStore.load({ force: true }),
    usersStore.load({ limit: 500 }),
    rbacStore.load([$ComposeAPI, $SystemAPI]),
  ]
  const delayPromise = new Promise(resolve => setTimeout(resolve, 1000))

  Promise.all([...fetchPromises, delayPromise]).finally(() => {
    loading.value = false
  })
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
  'namespace.manage',
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

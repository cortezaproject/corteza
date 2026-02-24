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
      class="flex-1 flex flex-col transition-[margin] duration-300"
      :style="{ marginLeft: contentMargin }"
    >
      <header>
        <CTopbar
          v-model:sidebar-expanded="expanded"
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
import CSidebarNavigation from '@/components/CSidebarNavigation.vue'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
const { CTopbar, CLoaderLogo, CSidebar } = components

const $Settings = inject('$Settings')

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

onMounted(() => {
  const delayPromise = new Promise(resolve => setTimeout(resolve, 1500))

  delayPromise.finally(() => {
    loading.value = false
  })
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
})
</script>

<template>
  <CLoaderLogo :show="loading" :logo-url="logoUrl" />

  <div class="h-screen flex flex-col">
    <header>
      <CTopbar
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
      position="top-center"
      :pt="{
        messageIcon: {
          style: {
            display: 'none',
          },
        },
      }"
    />

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { components, withMinDuration } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
import { useAutomationStore } from '@/stores/automation'
const { CTopbar, CLoaderLogo } = components

const $Settings = inject('$Settings')
const $AutomationAPI = inject('$AutomationAPI')
const store = useAutomationStore()

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

onMounted(async () => {
  try {
    await withMinDuration(store.loadCatalog($AutomationAPI))
  } catch (e) {
    console.error('Failed to load catalog:', e)
  } finally {
    loading.value = false
  }
})
</script>

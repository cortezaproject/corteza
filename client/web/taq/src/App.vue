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

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, ref } from 'vue'
import { RouterView } from 'vue-router'
const { CTopbar, CLoaderLogo } = components

const $Settings = inject('$Settings')

const logoUrl = computed(() => {
  return $Settings.attachment('ui.mainLogo')
})

const loading = ref(true)

onMounted(() => {
  // Simple loading delay - no store fetching needed yet
  setTimeout(() => {
    loading.value = false
  }, 1000)
})
</script>

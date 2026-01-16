<template>
  <div class="d-flex flex-column w-100 h-viewport overflow-hidden">
    <header
      v-show="loaded"
    >
      <c-topbar
        hide-app-selector
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          appMenu: $t('appMenu'),
          helpForum: $t('help.forum'),
          helpDocumentation: $t('help.documentation'),
          helpFeedback: $t('help.feedback'),
          helpVersion: $t('help.version'),
          userSettingsLoggedInAs: $t('userSettings.loggedInAs', { user }),
          userSettingsProfile: $t('userSettings.profile'),
          userSettingsChangePassword: $t('userSettings.changePassword'),
          userSettingsLogout: $t('userSettings.logout'),
          userSettingsTheme: $t('userSettings.theme'),
          lightTheme: $t('themes.labels.light'),
          darkTheme: $t('themes.labels.dark'),
        }"
      />
    </header>

    <main
      v-show="loaded"
      class="flex-fill overflow-hidden"
    >
      <c-app-selector
        :logo="logo"
      />
    </main>

    <c-loader-logo
      v-if="!loaded"
      :logo="logo"
    />

    <c-prompts />

    <c-extend-session
      v-if="isAutoLogoutEnabled"
      :timeout="$Settings.get('auth.autoLogout.timeout')"
      :labels="{
        extend: $t('extendSession.labels.extend'),
        warning: (countdownTime) => $t('extendSession.labels.warning', { countdownTime }),
      }"
    />

    <c-notification-sidebar v-if="!$Settings.get('ui.topbar', {}).hideNotifications" />
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import CAppSelector from '../components/CAppSelector'
import { components } from '@cortezaproject/corteza-vue'
const { CTopbar, CLoaderLogo, CPrompts, CExtendSession, CNotificationSidebar } = components

export default {
  i18nOptions: {
    namespaces: 'navigation',
  },

  components: {
    CAppSelector,
    CTopbar,
    CLoaderLogo,
    CPrompts,
    CExtendSession,
    CNotificationSidebar,
  },

  data () {
    return {
      loaded: false,
    }
  },

  computed: {
    icon () {
      return this.$Settings.attachment('ui.iconLogo')
    },

    logo () {
      return this.$Settings.attachment('ui.mainLogo')
    },

    user () {
      const { user } = this.$auth
      return user.name || user.handle || user.email || ''
    },

    isAutoLogoutEnabled () {
      return this.$Settings.get('auth.autoLogout.enabled')
    },
  },

  created () {
    this.preloadApplications()
      .then(() => {
        setTimeout(() => {
          this.loaded = true
        }, 2000)
      })
  },

  methods: {
    ...mapActions({
      preloadApplications: 'applications/load',
    }),
  },
}
</script>
<style scoped>
.h-viewport {
  height: 100vh;
  height: 100dvh;
}
</style>

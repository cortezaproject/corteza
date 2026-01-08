<template>
  <div class="d-flex flex-column w-100 h-viewport overflow-hidden">
    <header>
      <c-topbar
        :expanded="expanded"
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
          lightTheme: $t('general:themes.labels.light'),
          darkTheme: $t('general:themes.labels.dark'),
        }"
      >
        <template #title>
          {{ $t("discovery") }}
        </template>
      </c-topbar>
    </header>

    <aside>
      <c-sidebar
        :expanded.sync="expanded"
        :icon="icon"
        :logo="logo"
        expand-on-click
      >
        <template #body-expanded>
          <filters />
        </template>
      </c-sidebar>
    </aside>

    <main class="d-inline-flex overflow-hidden">
      <!--
        Content spacer
        Large and xl screens should push in content when the nav is expanded
      -->
      <template>
        <div
          class="sidebar-spacer d-print-none"
          :class="{
            expanded: expanded,
          }"
        />
      </template>

      <div class="flex-grow-1 overflow-hidden">
        <search />
      </div>
    </main>

    <c-extend-session
      v-if="isAutoLogoutEnabled"
      :timeout="$Settings.get('auth.autoLogout.timeout')"
      :labels="{
        extend: $t('general:extendSession.labels.extend'),
        warning: (countdownTime) =>
          $t('general:extendSession.labels.warning', { countdownTime }),
      }"
    />
    <c-notification-sidebar
      v-if="!$Settings.get('ui.topbar', {}).hideNotifications"
    />

    <rag :show-rag="showRag" />
    <b-button
      variant="primary"
      pill
      class="shadow position-fixed chat-trigger-btn"
      style="
        bottom: 1.5rem;
        right: 1.5rem;
        width: 3.5rem;
        height: 3.5rem;
        z-index: 1000;
      "
      aria-label="Toggle chat"
      @click="showRag = !showRag"
    >
      <font-awesome-icon :icon="['fas', 'robot']" class="h4 mb-0" />
    </b-button>
  </div>
</template>

<script>
import Search from '../components/Search.vue'
import Filters from '../components/Filters.vue'
import Rag from '../components/Rag.vue'
import { components } from '@cortezaproject/corteza-vue'
const { CTopbar, CSidebar, CExtendSession, CNotificationSidebar } = components

export default {
  i18nOptions: {
    namespaces: 'navigation',
  },

  components: {
    CTopbar,
    CSidebar,
    Search,
    Filters,
    Rag,
    CExtendSession,
    CNotificationSidebar,
  },

  data () {
    return {
      expanded: false,
      showRag: false,
    }
  },

  computed: {
    user () {
      const { user } = this.$auth
      return user.name || user.handle || user.email || ''
    },

    icon () {
      return this.$Settings.attachment('ui.iconLogo')
    },

    logo () {
      return this.$Settings.attachment('ui.mainLogo')
    },

    isAutoLogoutEnabled () {
      return this.$Settings.get('auth.autoLogout.enabled')
    },
  },
}
</script>
<style scoped>
.h-viewport {
  height: 100vh;
  height: 100dvh;
}
</style>

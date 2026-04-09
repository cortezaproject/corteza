<template>
  <div class="d-flex flex-column w-100 h-viewport overflow-hidden">
    <header>
      <c-topbar
        :expanded="expanded"
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          appMenu: $t('navigation:appMenu'),
          helpForum: $t('navigation:help.forum'),
          helpDocumentation: $t('navigation:help.documentation'),
          helpFeedback: $t('navigation:help.feedback'),
          helpVersion: $t('navigation:help.version'),
          userSettingsLoggedInAs: $t('navigation:userSettings.loggedInAs', { user }),
          userSettingsProfile: $t('navigation:userSettings.profile'),
          userSettingsChangePassword: $t('navigation:userSettings.changePassword'),
          userSettingsLogout: $t('navigation:userSettings.logout'),
          userSettingsTheme: $t('navigation:userSettings.theme'),
          lightTheme: $t('general:themes.labels.light'),
          darkTheme: $t('general:themes.labels.dark'),
        }"
      >
        <template #title>
          <portal-target name="topbar-title" />
        </template>

        <template #tools>
          <portal-target name="topbar-tools" />
        </template>
      </c-topbar>
    </header>

    <aside>
      <c-sidebar
        :icon="icon"
        :logo="logo"
        :disabled-routes="['root', 'workflow.list']"
      >
        <template #header-expanded>
          <portal-target name="sidebar-header-expanded" />
        </template>

        <template #body-expanded>
          <portal-target name="sidebar-body-expanded" />
        </template>

        <template #footer-expanded>
          <portal-target name="sidebar-footer-expanded" />
        </template>
      </c-sidebar>
    </aside>

    <main class="d-inline-flex h-100 overflow-auto">
      <!--
        Content spacer
        Large and xl screens should push in content when the nav is expanded
      -->
      <template>
        <div
          class="sidebar-spacer d-print-none"
          :class="{
            'expanded': expanded,
          }"
        />
      </template>
      <router-view />
    </main>

    <c-permissions-modal
      :labels="{
        save: $t('permissions:ui.save'),
        cancel: $t('permissions:ui.cancel'),
        loading: $t('permissions:ui.loading'),
        edit: {
          label: $t('permissions:ui.edit.label'),
          description: $t('permissions:ui.edit.description'),
        },
        evaluate: {
          title: $t('permissions:ui.evaluate.title'),
          description: $t('permissions:ui.evaluate.description'),
        },
        add: {
          label: $t('permissions:ui.add.label'),
          title: $t('permissions:ui.add.title'),
          save: $t('permissions:ui.add.save'),
          role: {
            label: $t('permissions:ui.add.role.label'),
            placeholder: $t('permissions:ui.add.role.placeholder'),
          },
          user: {
            label: $t('permissions:ui.add.user.label'),
            placeholder: $t('permissions:ui.add.user.placeholder'),
          },
        },
      }"
    />

    <c-extend-session
      v-if="isAutoLogoutEnabled"
      :timeout="$Settings.get('auth.autoLogout.timeout')"
      :labels="{
        extend: $t('general:extendSession.labels.extend'),
        warning: (countdownTime) => $t('general:extendSession.labels.warning', { countdownTime }),
      }"
    />

    <c-notification-sidebar v-if="!$Settings.get('ui.topbar', {}).hideNotifications" />
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'
const { CPermissionsModal, CTopbar, CSidebar, CExtendSession, CNotificationSidebar } = components

export default {
  components: {
    CPermissionsModal,
    CTopbar,
    CSidebar,
    CExtendSession,
    CNotificationSidebar,
  },

  data () {
    return {
      // Sidebar and Topbar
      expanded: false,
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

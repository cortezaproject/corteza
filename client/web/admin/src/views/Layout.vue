<template>
  <div class="d-flex flex-column w-100 h-viewport overflow-hidden">
    <header>
      <c-topbar
        :expanded="expanded"
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          appMenu: $t('navigation.appMenu'),
          helpForum: $t('navigation.help.forum'),
          helpDocumentation: $t('navigation.help.documentation'),
          helpFeedback: $t('navigation.help.feedback'),
          helpVersion: $t('navigation.help.version'),
          userSettingsLoggedInAs: $t('navigation.userSettings.loggedInAs', { user }),
          userSettingsProfile: $t('navigation.userSettings.profile'),
          userSettingsChangePassword: $t('navigation.userSettings.changePassword'),
          userSettingsLogout: $t('navigation.userSettings.logout'),
          userSettingsTheme: $t('navigation.userSettings.theme'),
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

    <aside
      v-if="allowed"
    >
      <c-sidebar
        :expanded.sync="expanded"
        :icon="icon"
        :logo="logo"
        expand-on-click
        :right="textDirectionality() === 'rtl'"
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

      <portal to="sidebar-body-expanded">
        <c-the-main-nav />
      </portal>
    </aside>

    <main
      v-if="allowed"
      class="d-inline-flex h-100 overflow-auto"
    >
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
      <div class="d-flex flex-column w-100 flex-fill">
        <router-view />
      </div>
    </main>

    <c-prompts />

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
import { components, mixins } from '@cortezaproject/corteza-vue'
import CTheMainNav from 'corteza-webapp-admin/src/components/CTheMainNav'
import { mapGetters } from 'vuex'

const { CExtendSession, CPermissionsModal, CPrompts, CTopbar, CSidebar, CNotificationSidebar } = components

export default {
  i18nOptions: {
    namespaces: 'admin',
  },

  components: {
    CPermissionsModal,
    CPrompts,
    CTopbar,
    CSidebar,
    CTheMainNav,
    CExtendSession,
    CNotificationSidebar,
  },

  mixins: [
    mixins.corredor,
  ],

  data () {
    return {
      expanded: false,

      allowed: false,
      error: null,
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

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

  created () {
    this.$root.$on('alert', this.displayToast)

    const rulesToCheck = [
      // Grant
      { resource: 'system/', operation: 'grant' },
      { resource: 'compose/', operation: 'grant' },
      { resource: 'federation/', operation: 'grant' },
      { resource: 'automation/', operation: 'grant' },
      // System - Search
      { resource: 'system/', operation: 'users.search' },
      { resource: 'system/', operation: 'roles.search' },
      { resource: 'system/', operation: 'applications.search' },
      { resource: 'system/', operation: 'templates.search' },
      { resource: 'system/', operation: 'auth-clients.search' },
      { resource: 'system/', operation: 'queues.search' },
      { resource: 'system/', operation: 'apigw-routes.search' },
      { resource: 'system/', operation: 'dal-connections.search' },
      // System - Create
      { resource: 'system/', operation: 'auth-client.create' },
      { resource: 'system/', operation: 'role.create' },
      { resource: 'system/', operation: 'user.create' },
      { resource: 'system/', operation: 'application.create' },
      { resource: 'system/', operation: 'queue.create' },
      { resource: 'system/', operation: 'apigw-route.create' },
      { resource: 'system/', operation: 'dal-connection.create' },
      // System - Manage/Read
      { resource: 'system/', operation: 'settings.read' },
      { resource: 'system/', operation: 'settings.manage' },
      { resource: 'system/', operation: 'action-log.read' },
      { resource: 'system/', operation: 'dal-sensitivity-level.manage' },
      // Compose
      { resource: 'compose/', operation: 'settings.read' },
      { resource: 'compose/', operation: 'settings.manage' },
      // Automation
      { resource: 'automation/', operation: 'workflows.search' },
      { resource: 'automation/', operation: 'workflow.create' },
      { resource: 'automation/', operation: 'sessions.search' },
      // Federation
      { resource: 'federation/', operation: 'pair' },
    ]

    this.allowed = rulesToCheck.some(({ resource, operation }) => this.can(resource, operation))

    // If not allowed to access, show error prompt and redirect after a delay
    if (!this.allowed) {
      this.toastDanger(this.$t('notification:notAllowed'))

      setTimeout(() => {
        window.location = '/..'
      }, 5000)
    }
  },

  methods: {
    displayToast ({ title, message, variant, countdown }) {
      this.$bvToast.toast(message, {
        title,
        variant,
        solid: true,
        autoHideDelay: countdown,
        toaster: 'b-toaster-bottom-right',
      })
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

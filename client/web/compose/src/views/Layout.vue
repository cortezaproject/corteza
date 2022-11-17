<template>
  <div class="d-flex flex-column w-100 vh-100 overflow-hidden">
    <header>
      <c-topbar
        :sidebar-pinned="pinned"
        :settings="$Settings.get('ui.topbar', {})"
        :labels="{
          helpForum: $t('help.forum'),
          helpDocumentation: $t('help.documentation'),
          helpFeedback: $t('help.feedback'),
          helpVersion: $t('help.version'),
          userSettingsLoggedInAs: $t('userSettings.loggedInAs', { user }),
          userSettingsProfile: $t('userSettings.profile'),
          userSettingsChangePassword: $t('userSettings.changePassword'),
          userSettingsLogout: $t('userSettings.logout'),
        }"
      >
        <template #title>
          <portal-target name="topbar-title" />
        </template>

        <template #tools>
          <portal-target name="topbar-tools" />
        </template>

        <template #avatar-dropdown>
          <portal-target name="topbar-avatar-dropdown" />
        </template>
      </c-topbar>
    </header>

    <aside>
      <c-sidebar
        :expanded.sync="expanded"
        :pinned.sync="pinned"
        :icon="icon"
        :logo="logo"
        :disabled-routes="disabledRoutes"
        expand-on-hover
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
    </aside>

    <main class="d-inline-flex h-100 overflow-auto d-print-flex">
      <!--
        Content spacer
        Large and xl screens should push in content when the nav is expanded
      -->
      <template>
        <div
          class="spacer d-print-none"
          :class="{
            'expanded': expanded && pinned,
          }"
        />
      </template>
      <router-view />
    </main>

    <c-prompts />

    <c-toaster
      :toasts="toasts"
    />

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
          label: $t('permissions:ui.evaluate.label'),
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

    <c-translation-modal />
  </div>
</template>
<script>
import moment from 'moment'
import CTranslationModal from '../components/Translator/CTranslatorModal'
import { components } from '@cortezaproject/corteza-vue'
const { CToaster, CPrompts, CPermissionsModal, CTopbar, CSidebar } = components

export default {
  i18nOptions: {
    namespaces: 'navigation',
  },

  components: {
    CPrompts,
    CTranslationModal,
    CPermissionsModal,
    CTopbar,
    CSidebar,
    CToaster, // Only used for reminders
  },

  data () {
    return {
      // Sidebar and Topbar
      expanded: false,
      pinned: false,

      toasts: [],

      disabledRoutes: [
        'namespaces',
        'namespace.list',
        'namespace.edit',
        'namespace.create',
        'namespace.clone',
        'namespace.manage',
      ],
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
  },

  watch: {
    icon: {
      immediate: true,
      handler (icon) {
        if (icon) {
          const favicon = document.getElementById('favicon')
          favicon.href = icon
        }
      },
    },
  },

  created () {
    /**
     * Listen for incoming warnings, alerts and other messages
     * from the (mostly) Corredor scripts and display them using toasts
     */
    this.$root.$on('alert', ({ message, ...params }) => this.toast(message, params))
    this.$root.$on('reminder.show', this.showReminder)

    this.$root.$on('check-namespace-sidebar', this.checkNamespaceSidebar)
  },

  beforeDestroy () {
    this.$root.$off('alert')
    this.$root.$off('reminder.show', this.showReminder)
    this.$root.$off('check-namespace-sidebar', this.checkNamespaceSidebar)
  },

  methods: {
    checkNamespaceSidebar (showSidebar) {
      const defaultDisabledRoutes = [
        'namespaces',
        'namespace.list',
        'namespace.edit',
        'namespace.create',
        'namespace.clone',
        'namespace.manage',
      ]
      const namespaceRoutes = ['page', 'pages', 'page.record', 'page.record.create', 'page.record.edit']

      this.disabledRoutes = [...defaultDisabledRoutes, ...(showSidebar ? [] : namespaceRoutes)]
    },

    removeToast (reminderID) {
      const i = this.toasts.findIndex(r => r.reminderID === reminderID)
      if (i > -1) {
        this.toasts.splice(i, 1)
      }
    },

    onReminderHide ({ reminderID }) {
      // Dismiss
      this.$SystemAPI.reminderDismiss({ reminderID })
        .then(() => {
          this.removeToast(reminderID)
          this.$root.$emit('reminder.updated', reminderID)
        })
    },

    // Duration is in ms
    onReminderSnooze ({ reminderID }, { duration }) {
      const remindAt = moment().add(duration, 'ms').toISOString()
      this.$SystemAPI.reminderSnooze({ reminderID, remindAt })
        .then(() => {
          this.removeToast(reminderID)
          this.$root.$emit('reminder.updated', reminderID)
        })
    },

    showReminder (r) {
      const i = this.toasts.findIndex(({ reminderID }) => reminderID === r.reminderID)
      if (i > -1 && (!r.editedAt || r.editedAt === this.toasts[i].editedAt)) {
        // Same reminder; no need to push it again
        return
      }

      r.options = {
        variant: 'warning',
        'no-auto-hide': true,
        solid: true,
        ...r.options,
      }

      r.actions.dismiss = {
        cb: this.onReminderHide,
        kind: 'Button',
        label: `<b>${this.$t('general:reminder.dismiss')}</b>`,
        options: {
          variant: 'warning',
          class: ['float-right'],
        },
      }

      r.actions.snooze = {
        cb: this.onReminderSnooze,
        label: `<b>${this.$t('general:reminder.snooze')}</b>`,
        kind: 'Select',
        options: {
          variant: 'outline-warning',
          class: ['float-left'],
          items: [
            { kind: 'item-button', label: this.$t('general:label.timeMinute', { t: 5 }), value: { duration: 1000 * 60 * 5 } },
            { kind: 'item-button', label: this.$t('general:label.timeMinute', { t: 15 }), value: { duration: 1000 * 60 * 15 } },
            { kind: 'item-button', label: this.$t('general:label.timeMinute', { t: 30 }), value: { duration: 1000 * 60 * 30 } },
            { kind: 'item-button', label: this.$t('general:label.timeHour', { t: 1 }), value: { duration: 1000 * 60 * 60 * 1 } },
            { kind: 'item-button', label: this.$t('general:label.timeHour', { t: 2 }), value: { duration: 1000 * 60 * 60 * 2 } },
            { kind: 'item-button', label: this.$t('general:label.timeHour', { t: 24 }), value: { duration: 1000 * 60 * 60 * 24 } },
          ],
        },
      }

      if (i > -1) {
        this.toasts.splice(i, 1, r)
      } else {
        this.toasts.push(r)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.spacer {
  min-width: 0;
  -webkit-transition: min-width 0.2s ease-in-out;
  -moz-transition: min-width 0.2s ease-in-out;
  -o-transition: min-width 0.2s ease-in-out;
  transition: min-width 0.2s ease-in-out;

  &.expanded {
    min-width: $sidebar-width;
    -webkit-transition: min-width 0.2s ease-in-out;
    -moz-transition: min-width 0.2s ease-in-out;
    -o-transition: min-width 0.2s ease-in-out;
    transition: min-width 0.2s ease-in-out;
  }
}
</style>

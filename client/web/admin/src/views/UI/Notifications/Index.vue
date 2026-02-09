<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-ui-notification-settings
      v-if="settings"
      :settings="settings"
      :processing="notifications.processing"
      :success="notifications.success"
      :can-manage="canManage"
      @submit="onSubmit($event, 'notifications')"
    />
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CUINotificationSettings from 'corteza-webapp-admin/src/components/Settings/UI/CUINotificationSettings'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'ui.settings',
    keyPrefix: 'editor.notifications',
  },

  components: {
    'c-ui-notification-settings': CUINotificationSettings,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      settings: undefined,

      notifications: {
        processing: false,
        success: false,
      },
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canManage () {
      return this.can('system/', 'settings.manage')
    },
  },

  created () {
    this.fetchSettings()
  },

  methods: {
    fetchSettings () {
      this.incLoader()

      this.$Settings.fetch()
      return this.$SystemAPI.settingsList({ prefix: 'ui' })
        .then(settings => {
          this.settings = {}

          settings.forEach(({ name, value }) => {
            this.$set(this.settings, name, value)
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.notifications.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onSubmit (settings, type) {
      this[type].processing = true

      const values = Object.entries(settings).map(([name, value]) => {
        return { name, value }
      })

      this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          return this.fetchSettings()
        })
        .then(() => {
          this.animateSuccess(type)
          this.toastSuccess(this.$t('notification:settings.notifications.update.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.notifications.update.error')))
        .finally(() => {
          this[type].processing = false
        })
    },
  },
}
</script>

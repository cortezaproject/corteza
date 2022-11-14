<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-ui-logo-editor
      :settings="settings"
      :can-manage="canManage"
    />

    <c-ui-topbar-settings
      :settings="settings"
      :processing="topbar.processing"
      :success="topbar.success"
      :can-manage="canManage"
      class="mt-3"
      @submit="onTopbarSubmit"
    />
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CUILogoEditor from 'corteza-webapp-admin/src/components/Settings/UI/CUILogoEditor'
import CUITopbarSettings from 'corteza-webapp-admin/src/components/Settings/UI/CUITopbarSettings'
import { mapGetters } from 'vuex'

const prefix = 'ui.'

export default {
  i18nOptions: {
    namespaces: [ 'ui.settings' ],
    keyPrefix: 'editor',
  },

  components: {
    'c-ui-logo-editor': CUILogoEditor,
    'c-ui-topbar-settings': CUITopbarSettings,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      settings: {},

      topbar: {
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
      this.$SystemAPI.settingsList({ prefix: prefix })
        .then(settings => {
          settings.forEach(({ name, value }) => {
            this.$set(this.settings, name, value)
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.ui.fetch.error')))
        .finally(() => {
          this.decLoader()
        })
    },

    onTopbarSubmit (settings) {
      this.topbar.processing = true

      const values = Object.entries(settings).map(([name, value]) => {
        return { name, value }
      })

      this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.animateSuccess('topbar')
          this.toastSuccess(this.$t('notification:settings.compose.update.success'))
          this.$Settings.fetch()
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.compose.update.error')))
        .finally(() => {
          this.topbar.processing = false
        })
    },
  },
}
</script>

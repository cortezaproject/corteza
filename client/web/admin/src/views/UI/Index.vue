<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-ui-logo-editor
      v-if="settings"
      :settings="settings"
      :can-manage="canManage"
    />

    <c-ui-branding-editor
      v-if="settings"
      :settings="settings"
      :processing="branding.processing"
      :success="branding.success"
      :can-manage="canManage"
      @submit="onSubmit($event, 'branding')"
    />

    <c-ui-custom-css
      v-if="settings"
      :settings="settings"
      :processing="customCSS.processing"
      :success="customCSS.success"
      :can-manage="canManage"
      class="mt-3"
      @submit="onSubmit($event, 'customCSS')"
    />

    <c-ui-topbar-settings
      v-if="settings"
      :settings="settings"
      :processing="topbar.processing"
      :success="topbar.success"
      :can-manage="canManage"
      class="mt-3"
      @submit="onSubmit($event, 'topbar')"
    />
  </b-container>
</template>

<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CUILogoEditor from 'corteza-webapp-admin/src/components/Settings/UI/CUILogoEditor'
import CUIBrandingEditor from '../../components/Settings/UI/CUIBrandingEditor.vue'
import CUICustomCSS from 'corteza-webapp-admin/src/components/Settings/UI/CUICustomCSS.vue'
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
    'c-ui-custom-css': CUICustomCSS,
    'c-ui-branding-editor': CUIBrandingEditor,
    'c-ui-topbar-settings': CUITopbarSettings,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      settings: undefined,

      topbar: {
        processing: false,
        success: false,
      },

      customCSS: {
        processing: false,
        success: false,
      },

      branding: {
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
          this.settings = {}

          settings.forEach(({ name, value }) => {
            this.$set(this.settings, name, value)
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.ui.fetch.error')))
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
          this.animateSuccess(type)
          this.toastSuccess(this.$t('notification:settings.ui.update.success'))
          this.$Settings.fetch()
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.ui.update.error')))
        .finally(() => {
          this[type].processing = false

          // Refresh the page if branding or custom CSS was updated
          if (type === 'branding' || type === 'customCSS') {
            window.location.reload()
          }
        })
    },
  },
}
</script>

<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-system-editor-auth
      v-if="Object.keys(getAuth).length"
      :settings="getAuth"
      :processing="auth.processing"
      :success="auth.success"
      :can-manage="canManage"
      @submit="onAuthSubmit"
    />

    <c-system-editor-external
      v-model="settings"
      class="mt-3"
      :processing="external.processing"
      :success="external.success"
      :can-manage="canManage"
      @submit="onExternalSubmit"
    />

    <c-system-editor-auth-bg-screen
      :settings="getAuthBackground"
      :can-manage="canManage"
      :processing="authBackground.processing"
      :success="authBackground.success"
      class="mt-3"
      @onUpload="onBackgroundImageUpload"
      @resetAttachment="onResetBackgroundImage"
      @submit="onAuthBackgroundSubmit"
    />
  </b-container>
</template>
<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CSystemEditorAuth from 'corteza-webapp-admin/src/components/Settings/System/CSystemEditorAuth'
import CSystemEditorExternal from 'corteza-webapp-admin/src/components/Settings/System/CSystemEditorExternal'
import CSystemEditorAuthBgScreen from 'corteza-webapp-admin/src/components/Settings/System/CSystemEditorAuthBgScreen'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor',
  },

  components: {
    CSystemEditorAuth,
    CSystemEditorExternal,
    CSystemEditorAuthBgScreen,
  },

  mixins: [
    editorHelpers,
  ],

  data () {
    return {
      settings: [],

      auth: {
        processing: false,
        success: false,
      },

      external: {
        processing: false,
        success: false,
      },

      authBackground: {
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

    getAuth () {
      return this.filterSettings('auth')
    },

    getAuthBackground () {
      return this.filterSettings('auth.ui')
    },

  },

  created () {
    this.fetchSettings()
  },

  methods: {
    fetchSettings () {
      this.incLoader()

      this.$SystemAPI.settingsList()
        .then(settings => {
          this.settings = settings
        })
        .catch(e => {
          this.toastErrorHandler(this.$t('notification:settings.system.fetch.error'))(e)
          this.$router.push({ name: 'dashboard' })
        })
        .finally(() => {
          this.decLoader()
        })
    },

    filterSettings (prefix) {
      if (this.settings.length > 0) {
        return this.settings.reduce((map, obj) => {
          const { name, value } = obj
          if (name.startsWith(prefix)) {
            map[name] = value
          }
          return map
        }, {})
      }
      return {}
    },

    onAuthSubmit (auth) {
      this.auth.processing = true

      const values = Object.entries(auth).map(([name, value]) => {
        return { name, value }
      })

      this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.animateSuccess('auth')
          this.toastSuccess(this.$t('notification:settings.system.auth.success'))
          this.$Settings.fetch()
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.system.auth.error')))
        .finally(() => {
          this.auth.processing = false
        })
    },

    onExternalSubmit (external) {
      this.external.processing = true

      this.$SystemAPI.settingsUpdate({ values: external })
        .then(() => {
          this.animateSuccess('external')
          this.toastSuccess(this.$t('notification:settings.system.external.success'))
          this.$Settings.fetch()
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.system.external.error')))
        .finally(() => {
          this.external.processing = false
          this.fetchSettings()
        })
    },

    onBackgroundImageUpload () {
      this.fetchSettings()
    },

    onResetBackgroundImage (name) {
      this.$SystemAPI.settingsUpdate({ values: [{ name, value: undefined }], upload: {} })
        .then(() => {
          this.fetchSettings()
        })
    },

    onAuthBackgroundSubmit (authBackground) {
      this.authBackground.processing = true

      const values = ([
        {
          name: 'auth.ui.styles',
          value: authBackground,
        },
      ])

      this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.animateSuccess('authBackground')
          this.toastSuccess(this.$t('notification:settings.system.bgScreen.style.success'))
          this.$Settings.fetch()
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.system.bgScreen.style.error')))
        .finally(() => {
          this.authBackground.processing = false
        })
    },
  },
}
</script>

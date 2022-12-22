<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-system-editor-auth
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
  </b-container>
</template>
<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CSystemEditorAuth from 'corteza-webapp-admin/src/components/Settings/System/CSystemEditorAuth'
import CSystemEditorExternal from 'corteza-webapp-admin/src/components/Settings/System/CSystemEditorExternal'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'system.settings',
    keyPrefix: 'editor',
  },

  components: {
    CSystemEditorAuth,
    CSystemEditorExternal,
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
      if (this.settings.length > 0) {
        return this.settings.reduce((map, obj) => {
          const { name, value } = obj
          const split = name.split('.')
          if (split[0] === 'auth' && split[1] !== 'external') {
            map[name] = value
          }
          return map
        }, {})
      }
      return {}
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
  },
}
</script>

<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-system-email-server
      :key="JSON.stringify(server)"
      v-model="server"
      :processing="auth.processing"
      :success="auth.success"
      :disabled="!canManage"
      @submit="onEmailServerSubmit($event)"
    />
  </b-container>
</template>
<script>
import editorHelpers from 'corteza-webapp-admin/src/mixins/editorHelpers'
import CSystemEmailServer from 'corteza-webapp-admin/src/components/Settings/Email/CSystemEmailServer'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'system.email',
    keyPrefix: 'editor',
  },

  components: {
    CSystemEmailServer,
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

    server () {
      let s = this.settings.find(({ name }) => name === 'smtp.servers')

      if (!s || !s.value || s.value.length === 0 || (typeof s.value[0] !== 'object')) {
        return {}
      }

      return s.value[0]
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

    /**
     * Handles data change, supporting only one server for now
     * @param server
     */
    onEmailServerSubmit (server) {
      const values = [
        { name: 'smtp.servers', value: [server] },
      ]

      this.external.processing = true

      this.$SystemAPI.settingsUpdate({ values })
        .then(() => {
          this.animateSuccess('external')
          this.toastSuccess(this.$t('notification:settings.system.external.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.system.external.error')))
        .finally(() => {
          this.external.processing = false
        })
    },
  },
}
</script>

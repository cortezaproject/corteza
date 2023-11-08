<template>
  <b-container
    class="pt-2 pb-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <c-system-email-server
      :key="JSON.stringify(server)"
      v-model="server"
      :processing="external.processing"
      :processing-smtp-test="auth.processing"
      :success="external.success"
      :success-smtp-tes="auth.success"
      :disabled="!canManage"
      @submit="onEmailServerSubmit($event)"
      @smtpConnectionCheck="onSmtpConnectionCheck($event)"
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

    onSmtpConnectionCheck (server) {
      this.external.processing = true

      // Append the list of recepient's email  Addresses
      const recepients = []
      recepients.push(server.from)

      this.$SystemAPI.smtpConfigurationCheckerCheck({
        host: server.host,
        port: parseInt(server.port),
        recipients: recepients,
        username: server.user,
        password: server.pass,
        tlsInsecure: server.tlsInsecure,
        tlsServerName: server.tlsServerName,
      })
        .then(response => {
          if (Object.values(response).every(resp => resp === '')) {
            this.animateSuccess('external')
            this.toastSuccess(this.$t('notification:settings.system.smtpCheck.success'))
          }

          Object.keys(response).forEach(key => {
            if (response[key]) {
              this.toastWarning(`${key}: ${response[key]}`)
            }
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:settings.system.smtpCheck.error')))
        .finally(() => {
          this.external.processing = false
        })
    },
  },

}
</script>

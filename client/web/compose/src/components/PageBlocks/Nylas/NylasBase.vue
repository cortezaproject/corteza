<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="hideComponent"
      class="d-flex flex-column align-items-center justify-content-center h-100 overflow-hidden"
    >
      <b-spinner
        v-if="processing"
      />

      <span
        v-else-if="!isExternalConfigured"
      >
        {{ $t('noExternal') }}
      </span>

      <b-button
        v-else-if="!accessToken"
        :href="`${$auth.cortezaAuthURL}/external/nylas`"
        target="_blank"
      >
        {{ $t('connect') }}
      </b-button>

      <span
        v-else-if="!options.componentID"
      >
        {{ $t('noComponentID') }}
      </span>
    </div>

    <component
      :is="nylasComponent"
      v-else
      :component-i-d="options.componentID"
      :access-token="accessToken"
      :prefill-values="prefillValues"
    />
  </wrap>
</template>

<script>
import base from '../base'
import * as Components from './Components/loader'

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'nylas.viewer',
  },

  extends: base,

  data () {
    return {
      processing: false,

      isExternalConfigured: false,

      accessToken: '',

      tokenCheckInterval: undefined,

      prefillValues: {
        to: [],
        subject: '',
        body: '',
        queryString: '',
      },
    }
  },

  computed: {
    nylasComponent () {
      return Components[this.options.kind]
    },

    hideComponent () {
      return this.processing || !this.isExternalConfigured || !this.accessToken || !this.options.componentID
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
      async handler () {
        this.processing = true
        // Check if nylas is configured as a provider
        const { enabled: externalEnabled = false, providers = [] } = this.$Settings.get('auth.external', {})
        const { enabled: nylasEnabled = false, usage = [] } = providers.find(({ handle }) => handle === 'nylas') || {}

        if (!externalEnabled || !nylasEnabled || !usage || !usage.includes('api')) {
          this.isExternalConfigured = false
          this.processing = false
          return
        }

        this.isExternalConfigured = true

        // Check if token exists, if yes get it, otherwise check it every 3 seconds
        await this.checkNylasAccessToken()
          .finally(() => {
            this.processing = false
          })

        this.processPrefillValues()
      },
    },
  },

  methods: {
    processPrefillValues () {
      if (this.module) {
        this.prefillValues = {
          to: this.mapFieldToValue('to').map(v => ({ email: v })),
          subject: this.mapFieldToValue('subject').join(','),
          body: this.mapFieldToValue('body').join('<br />'),
          queryString: this.mapFieldToValue('queryString').join(','),
        }
      }
    },

    mapFieldToValue (property) {
      const ID = this.options.prefill[property]

      if (!ID) {
        return []
      }

      const { name, isMulti } = this.module.fields.find(f => f.fieldID === ID) || {}
      const value = this.record.values[name]

      if (!value) {
        return []
      }

      return isMulti ? this.record.values[name] : [this.record.values[name]]
    },

    checkNylasAccessToken () {
      return this.$SystemAPI.userListCredentials({ userID: this.$auth.user.userID })
        .then(credentials => {
          if (!credentials || !credentials.length) {
            return
          }

          this.accessToken = (credentials.find(({ kind }) => kind === 'access-nylas') || {}).credentials

          if (this.accessToken) {
            if (this.tokenCheckInterval) {
              clearInterval(this.tokenCheckInterval)
            }
          } else if (!this.tokenCheckInterval) {
            this.tokenCheckInterval = setInterval(() => {
              this.checkNylasAccessToken()
            }, 3000)
          }
        })
    },
  },
}
</script>

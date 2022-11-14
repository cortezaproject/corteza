<template>
  <b-container
    v-if="isDC !== null"
    fluid
    class="d-flex flex-column p-3"
  >
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <div
      class="flex-shrink-1"
    >
      <p>
        {{ $t('description.first') }}<br>
        {{ $t('description.second') }}
      </p>

      <b-row
        align-v="stretch"
      >
        <b-col
          v-for="option in options"
          :key="option.title"
          cols="12"
          md="6"
          xl="3"
          class="mb-3"
        >
          <b-card
            :title="option.title"
            class="card-hover-popup shadow-sm h-100"
            body-class="d-flex flex-column"
          >
            <b-card-text class="flex-grow-1">
              {{ option.description }}
            </b-card-text>

            <b-button
              :variant="option.button.variant || 'light'"
              :to="option.button.to"
            >
              {{ option.button.label }}
            </b-button>
          </b-card>
        </b-col>
      </b-row>
    </div>

    <div
      class="d-flex flex-column h-100"
    >
      <h6
        class="text-primary"
      >
        {{ $t('connection-location') }}
      </h6>

      <connection-map
        :connections="connections"
        class="align-self-center justify-self-center flex-grow-1 rounded-lg shadow-lg"
      />
    </div>
  </b-container>
</template>

<script>
import ConnectionMap from 'corteza-webapp-privacy/src/components/ConnectionMap'

export default {
  name: 'Dashboard',

  i18nOptions: {
    namespaces: 'dashboard',
  },

  components: {
    ConnectionMap,
  },

  data () {
    return {
      processing: false,

      isDC: null,

      connections: [],

      userOptions: [
        {
          title: this.$t('user-options.data-overview.title'),
          description: this.$t('user-options.data-overview.description'),
          button: { label: this.$t('user-options.data-overview.button-label'), to: { name: 'data-overview' } },
        },
        {
          title: this.$t('user-options.privacy-requests.title'),
          description: this.$t('user-options.privacy-requests.description'),
          button: { label: this.$t('user-options.privacy-requests.button-label'), to: { name: 'request.list' } },
        },
        {
          title: this.$t('user-options.export.title'),
          description: this.$t('user-options.export.description'),
          button: { label: this.$t('user-options.export.button-label'), to: { name: 'request.create', params: { kind: 'export' } } },
        },

        {
          title: this.$t('user-options.delete.title'),
          description: this.$t('user-options.delete.description'),
          button: { label: this.$t('user-options.delete.button-label'), variant: 'danger', to: { name: 'request.create', params: { kind: 'delete' } } },
        },
      ],

      dcOptions: [
        {
          title: this.$t('dc-options.privacy-requests.title'),
          description: this.$t('dc-options.privacy-requests.description'),
          button: { label: this.$t('dc-options.privacy-requests.button-label'), to: { name: 'request.list' } },
        },
        {
          title: this.$t('dc-options.sensitive-data.title'),
          description: this.$t('dc-options.sensitive-data.description'),
          button: { label: this.$t('dc-options.sensitive-data.button-label'), to: { name: 'sensitive-data' } },
        },
      ],
    }
  },

  computed: {
    options () {
      return this.isDC ? this.dcOptions : this.userOptions
    },
  },

  created () {
    this.fetchConnections()
    this.checkIsDC()
  },

  methods: {
    fetchConnections () {
      this.processing = true

      this.$SystemAPI.dataPrivacyConnectionList()
        .then(({ set = [] }) => {
          this.connections = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:connection-load-failed')))
        .finally(() => {
          this.processing = false
        })
    },

    checkIsDC () {
      this.$SystemAPI.roleList({ query: 'data-privacy-officer', memberID: this.$auth.user.userID })
        .then(({ set = [] }) => {
          this.isDC = !!set.length
        })
    },
  },
}
</script>

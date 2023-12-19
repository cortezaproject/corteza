<template>
  <b-container
    class="d-flex flex-column p-3"
  >
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <b-card
      class="shadow-sm mb-4"
    >
      <b-form-group
        :label="$t('connection.label')"
        label-class="text-primary"
        class="mb-0"
      >
        <c-input-select
          v-model="connectionID"
          :disabled="processing.connections"
          :options="connections"
          :clearable="false"
          :reduce="o => o.connectionID"
          :placeholder="$t('connection.placeholder')"
          :get-option-label="({ handle, meta }) => meta.name || handle"
          :get-option-key="getOptionKey"
        />
      </b-form-group>
    </b-card>

    <div
      v-if="processing.sensitiveData"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <h5
      v-else-if="!(connectionID && modules[connectionID])"
      class="text-center mt-5"
    >
      {{ $t('no-data-available') }}
    </h5>

    <module-records
      v-else
      v-slot="{ value }"
      :modules="modules[connectionID]"
    >
      <p
        v-for="(v, vi) in value.value"
        :key="vi"
        class="mb-0"
        :class="{ 'mt-1': vi > 0 }"
      >
        {{ v }}
      </p>
    </module-records>

    <portal to="editor-toolbar">
      <editor-toolbar
        :processing="processing.connections || processing.sensitiveData"
        :back-link="{ name: 'data-overview' }"
      >
        <b-button
          data-test-id="button-request-deletion"
          :disabled="processing.connections || processing.sensitiveData"
          variant="light"
          size="lg"
          class="ml-1"
          @click="$router.push({ name: 'request.create', params: { kind: 'delete', connection } })"
        >
          {{ $t('request-deletion') }}
        </b-button>

        <b-button
          data-test-id="button-request-correction"
          :disabled="processing.connections || processing.sensitiveData"
          variant="primary"
          size="lg"
          class="ml-1"
          @click="$router.push({ name: 'request.create', params: { kind: 'correct', connection } })"
        >
          {{ $t('request-correction') }}
        </b-button>
      </editor-toolbar>
    </portal>
  </b-container>
</template>

<script>
import EditorToolbar from 'corteza-webapp-privacy/src/components/Common/EditorToolbar'
import ModuleRecords from 'corteza-webapp-privacy/src/components/Common/ModuleRecords'

export default {
  name: 'ApplicationDataOverview',

  i18nOptions: {
    namespaces: 'data-overview',
    keyPrefix: 'application',
  },

  components: {
    EditorToolbar,
    ModuleRecords,
  },

  data () {
    return {
      processing: {
        connections: true,
        sensitiveData: true,
      },

      connectionID: undefined,

      connections: [],

      modules: {},
    }
  },

  computed: {
    connection () {
      return this.connections.find(({ connectionID }) => connectionID === this.connectionID) || {}
    },
  },

  watch: {
    connectionID: {
      handler (connectionID = '') {
        this.fetchSensitiveData(connectionID)
      },
    },
  },

  created () {
    this.fetchConnections()
  },

  methods: {
    fetchConnections () {
      this.processing.connections = true

      this.$SystemAPI.dataPrivacyConnectionList()
        .then(({ set = [] }) => {
          this.connections = set
          const { connectionID } = set[0] || {}
          this.connectionID = connectionID
        })
        .catch(this.toastErrorHandler(this.$t('notification:connection-load-failed')))
        .finally(() => {
          this.processing.connections = false
        })
    },

    fetchSensitiveData (connectionID) {
      if (connectionID) {
        this.processing.sensitiveData = true

        this.$ComposeAPI.dataPrivacyRecordList({ connectionID: [connectionID] })
          .then(({ set = [] }) => {
            if (set.length) {
              this.$set(this.modules, connectionID, set)
            }
          })
          .catch(this.toastErrorHandler(this.$t('notification:sensitive-data-fetch-failed')))
          .finally(() => {
            this.processing.sensitiveData = false
          })
      }
    },

    getOptionKey ({ connectionID }) {
      return connectionID
    },
  },
}
</script>

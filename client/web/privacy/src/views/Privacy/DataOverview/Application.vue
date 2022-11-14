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
        <vue-select
          v-model="connection"
          :disabled="processing.connections"
          :options="connections"
          :clearable="false"
          :placeholder="$t('connection.placeholder')"
          :get-option-label="({ handle, meta }) => meta.name || handle"
          class="h-100 bg-white"
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
      v-else-if="!(connection && modules[connection.connectionID])"
      class="text-center mt-5"
    >
      {{ $t('no-data-available') }}
    </h5>

    <module-records
      v-else
      v-slot="{ value }"
      :modules="modules[connection.connectionID]"
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
          :disabled="processing.connections || processing.sensitiveData"
          variant="light"
          size="lg"
          class="ml-1"
          @click="$router.push({ name: 'request.create', params: { kind: 'delete', connection } })"
        >
          {{ $t('request-deletion') }}
        </b-button>

        <b-button
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
import VueSelect from 'vue-select'

export default {
  name: 'ApplicationDataOverview',

  i18nOptions: {
    namespaces: 'data-overview',
    keyPrefix: 'application',
  },

  components: {
    VueSelect,
    EditorToolbar,
    ModuleRecords,
  },

  data () {
    return {
      processing: {
        connections: true,
        sensitiveData: true,
      },

      connection: undefined,

      connections: [],

      modules: {},
    }
  },

  watch: {
    connection: {
      handler ({ connectionID } = {}) {
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
          this.connection = set[0]
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
  },
}
</script>

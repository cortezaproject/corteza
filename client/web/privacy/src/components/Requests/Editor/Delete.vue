<template>
  <div
    class="d-flex flex-column h-100"
  >
    <b-card
      class="shadow-sm mb-3"
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
          :get-option-label="({ handle, meta }) => meta.name || handle"
          :placeholder="$t('connection.placeholder')"
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
      v-slot="{ namespace, module, record, value }"
      :modules="modules[connection.connectionID]"
    >
      <b-form-checkbox
        v-if="value.value.length"
        v-model="value.selected"
        @change="updateValue({ namespace, module, recordID: record.recordID, field: value.name, value: $event, orgValue: value.value })"
      >
        <span
          v-for="(v, vi) in value.value"
          :key="vi"
          class="py-2 mb-0"
        >
          <del
            v-if="value.selected"
          >
            {{ v }}
          </del>
          <span
            v-else
          >
            {{ v }}
          </span>
        </span>
      </b-form-checkbox>
    </module-records>

    <portal to="editor-toolbar">
      <editor-toolbar
        :processing="processing.connections || processing.sensitiveData"
        :back-link="{ name: 'data-overview.application' }"
        submit-show
        :submit-label="$t('submit')"
        :submit-disabled="!valid || !connection"
        @submit="$emit('submit', { kind: 'delete', payload })"
      />
    </portal>
  </div>
</template>

<script>
import EditorToolbar from 'corteza-webapp-privacy/src/components/Common/EditorToolbar'
import ModuleRecords from 'corteza-webapp-privacy/src/components/Common/ModuleRecords'
import VueSelect from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'request',
    keyPrefix: 'edit.delete',
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

      payload: {},
      valid: false,
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
          if (!this.$route.params.connection) {
            this.connection = set[0]
          } else {
            this.connection = this.$route.params.connection
          }
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

            // Reset payload
            this.payload = {
              connectionID,
              modules: {},
            }
          })
          .catch(this.toastErrorHandler(this.$t('notification:sensitive-data-fetch-failed')))
          .finally(() => {
            this.processing.sensitiveData = false
          })
      }
    },

    updateValue ({ namespace, module, recordID, field, value, orgValue = [] }) {
      const { connectionID } = this.connection
      const { namespaceID, name: namespaceName } = namespace
      const { moduleID, name: moduleName } = module

      if (!value) {
        delete this.payload.modules[moduleID].records[recordID].values[field]

        if (Object.keys(this.payload.modules[moduleID].records[recordID].values).length === 0) {
          delete this.payload.modules[moduleID].records[recordID]

          if (Object.keys(this.payload.modules[moduleID].records).length === 0) {
            delete this.payload.modules[moduleID]
          }
        }

        this.valid = Object.keys(this.payload.modules).length > 0

        return
      }

      if (connectionID) {
        if (!this.payload.modules[moduleID]) {
          this.payload.modules[moduleID] = {
            namespace: namespaceName,
            namespaceID,
            module: moduleName,
            moduleID,
            records: {},
          }
        }

        if (!this.payload.modules[moduleID].records[recordID]) {
          this.payload.modules[moduleID].records[recordID] = {
            values: {},
          }
        }

        this.payload.modules[moduleID].records[recordID].values[field] = orgValue
        this.valid = true
      }
    },
  },
}
</script>

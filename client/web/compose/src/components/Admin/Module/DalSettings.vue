<template>
  <div
    v-if="module"
  >
    <b-form-group
      :label="$t('connection.label')"
      :description="$t('connection.description')"
      label-class="text-primary"
    >
      <vue-select
        v-model="module.config.dal.connectionID"
        :options="connections"
        :disabled="processing"
        :clearable="false"
        :reduce="s => s.connectionID"
        :placeholder="$t('connection.placeholder')"
        :get-option-label="getConnectionLabel"
        class="bg-white"
      />
    </b-form-group>

    <b-form-group
      :label="$t('ident.label')"
      :description="$t('ident.description', { interpolation: { prefix: '{{{', suffix: '}}}' } })"
      label-class="text-primary"
    >
      <b-form-input
        v-model="module.config.dal.ident"
        :placeholder="$t('ident.placeholder')"
      />
    </b-form-group>

    <b-form-group
      :label="$t('module-fields.label')"
      :description="$t('module-fields.description')"
    >
      <dal-field-store-encoding
        v-for="({ field, storeIdent, label }) in moduleFields"
        :key="field"
        :config="moduleFieldEncoding[field] || {}"
        :field="field"
        :label="label"
        :default-strategy="moduleFieldDefaultEncodingStrategy"
        :store-ident="storeIdent"
        @change="applyModuleFieldStrategyConfig(field, $event)"
      />
    </b-form-group>

    <b-form-group
      :label="$t('system-fields.label')"
      :description="$t('system-fields.description')"
    >
      <dal-field-store-encoding
        v-for="({ field, storeIdent, label }) in systemFields"
        :key="field"
        :config="systemFieldEncoding[field] || {}"
        :field="field"
        :label="label"
        :store-ident="storeIdent"
        :allow-omit-strategy="true"
        @change="applySystemFieldStrategyConfig(field, $event)"
      />
    </b-form-group>
  </div>
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { moduleFieldStrategyConfig, systemFieldStrategyConfig, types } from './encoding-strategy'
import VueSelect from 'vue-select'
import DalFieldStoreEncoding from 'corteza-webapp-compose/src/components/Admin/Module/DalFieldStoreEncoding'

const PrimaryConnType = 'corteza::system:primary-dal-connection'

export default {
  i18nOptions: {
    namespaces: 'module',
    keyPrefix: 'edit.config.dal',
  },

  components: {
    VueSelect,
    DalFieldStoreEncoding,
  },

  props: {
    module: {
      type: compose.Module,
      required: true,
    },
  },

  data () {
    const systemFieldEncoding = this.module.config.dal.systemFieldEncoding || {}
    const systemFields = [
      { field: 'id', storeIdent: 'id' },
      { field: 'namespaceID', storeIdent: 'rel_namespace' },
      { field: 'moduleID', storeIdent: 'rel_module' },
      { field: 'revision', storeIdent: 'revision' },
      { field: 'meta', storeIdent: 'meta' },
      { field: 'ownedBy', storeIdent: 'owned_by' },
      { field: 'createdAt', storeIdent: 'created_at' },
      { field: 'createdBy', storeIdent: 'created_by' },
      { field: 'updatedAt', storeIdent: 'updated_at' },
      { field: 'updatedBy', storeIdent: 'updated_by' },
      { field: 'deletedAt', storeIdent: 'deleted_at' },
      { field: 'deletedBy', storeIdent: 'deleted_by' },
    ].map(sf => ({ ...sf, label: this.$t(`field:system.${sf.field}`) }))

    return {
      processing: false,
      connections: [],

      moduleFields: [],
      moduleFieldEncoding: [],

      systemFields,
      systemFieldEncoding: systemFields.reduce((enc, { field }) => {
        enc[field] = systemFieldEncoding[field] || {}
        return enc
      }, {}),
    }
  },

  computed: {
    moduleFieldDefaultEncodingStrategy () {
      return types.JSON
    },
  },

  watch: {
    'module.fields': {
      handler (m) {
        this.moduleFields = []

        for (const f of this.module.fields) {
          const a = {
            field: f.name,
            label: f.label || f.name,
            storeIdent: f.name,
          }

          // In case of a JSON encoding strategy, default to values
          const strat = f.config.dal.encodingStrategy
          if (!strat || strat[types.JSON]) {
            a.storeIdent = 'values'
          }

          this.moduleFields.push(a)
        }

        this.moduleFieldEncoding = this.moduleFields.reduce((enc, { field }) => {
          const f = this.module.findField(field)
          if (f) {
            enc[field] = f.config.dal.encodingStrategy || {}
          }

          return enc
        }, {})
      },
      deep: true,
      immediate: true,
    },
  },

  mounted () {
    this.fetchConnections()
  },

  methods: {
    async fetchConnections () {
      this.processing = true
      return this.$SystemAPI.dalConnectionList()
        .then(({ set = [] }) => {
          this.connections = set

          const { connectionID } = this.module.config.dal || {}
          if (!connectionID || connectionID === NoID) {
            const primaryConnectionID = (this.connections.find(c => c.type === PrimaryConnType) || { connectionID: NoID }).connectionID
            this.module.config.dal.connectionID = primaryConnectionID
          }
        })
        .catch(this.toastErrorHandler(this.$t('connections.fetch-failed')))
        .finally(() => {
          this.processing = false
        })
    },

    getConnectionLabel ({ connectionID, handle, meta = {} }) {
      return meta.name || handle || connectionID
    },

    applyModuleFieldStrategyConfig (field, { strategy, config }) {
      const value = moduleFieldStrategyConfig(strategy, config)

      // merge new config into existing
      this.moduleFieldEncoding = { ...this.moduleFieldEncoding, [field]: value }

      // filter out empty configs and update the original config
      const moduleField = this.module.findField(field)
      if (moduleField) {
        moduleField.config.dal.encodingStrategy = value
      }
    },

    applySystemFieldStrategyConfig (field, { strategy, config }) {
      const value = systemFieldStrategyConfig(strategy, config)

      // merge new config into existing
      this.systemFieldEncoding = { ...this.systemFieldEncoding, [field]: value }

      // filter out empty configs and update the original config
      this.module.config.dal.systemFieldEncoding = Object.entries(this.systemFieldEncoding)
        .reduce((enc, [f, c]) => {
          if (c === null || Object.keys(c).length) {
            enc[f] = c
          }
          return enc
        }, {})
    },
  },
}
</script>

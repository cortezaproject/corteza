<template>
  <div
    v-if="module"
  >
    <b-form-group
      :label="$t('connection.label')"
      :description="$t('connection.description')"
      label-class="text-primary"
    >
      <c-input-select
        v-model="module.config.dal.connectionID"
        :options="connections"
        :get-option-key="getOptionKey"
        :disabled="processing"
        :clearable="false"
        :reduce="s => s.connectionID"
        :placeholder="$t('connection.placeholder')"
        :get-option-label="getConnectionLabel"
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
      label-class="text-primary"
    >
      <dal-field-store-encoding
        v-for="({ field, storeIdent, label, isMulti }, i) in moduleFields"
        :key="i"
        :config="moduleFieldEncoding[field] || {}"
        :field="field"
        :label="label"
        :is-multi="isMulti"
        :default-strategy="moduleFieldDefaultEncodingStrategy"
        :store-ident="storeIdent"
        @change="applyModuleFieldStrategyConfig(field, $event)"
      />
    </b-form-group>

    <b-form-group
      :description="$t('system-fields.description')"
    >
      <div class="my-4 d-flex justify-content-between align-items-center flex-wrap">
        <label>
          {{ $t('system-fields.label') }}
        </label>
        <b-form-radio-group
          v-model="selectedGroup"
          buttons
          button-variant="outline-secondary"
          size="sm"
          name="buttons"
          :options="optionsGroups"
          class="mb-3"
          @change="applySelectedSystemFields"
        />
      </div>

      <dal-field-store-encoding
        v-for="({ field, storeIdent, label, disabled }, i) in systemFields"
        :key="i"
        :config="systemFieldEncoding[field] || {}"
        :field="field"
        :label="label"
        :store-ident="storeIdent"
        :allow-omit-strategy="true"
        :disabled="disabled"
        @change="applySystemFieldStrategyConfig(field, $event)"
      />
    </b-form-group>
  </div>
</template>

<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { moduleFieldStrategyConfig, systemFieldStrategyConfig, types } from './encoding-strategy'
import DalFieldStoreEncoding from 'corteza-webapp-compose/src/components/Admin/Module/DalFieldStoreEncoding'

const PrimaryConnType = 'corteza::system:primary-dal-connection'

export default {
  i18nOptions: {
    namespaces: 'module',
    keyPrefix: 'edit.config.dal',
  },

  components: {
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
      { field: 'id', storeIdent: 'id', disabled: true },
      { field: 'namespaceID', storeIdent: 'rel_namespace', group: 'partition' },
      { field: 'moduleID', storeIdent: 'rel_module', group: 'partition' },
      { field: 'revision', storeIdent: 'revision', group: 'extras' },
      { field: 'meta', storeIdent: 'meta', group: 'extras' },
      { field: 'ownedBy', storeIdent: 'owned_by', group: 'user_reference' },
      { field: 'createdAt', storeIdent: 'created_at', group: 'timestamps' },
      { field: 'createdBy', storeIdent: 'created_by', group: 'user_reference' },
      { field: 'updatedAt', storeIdent: 'updated_at', group: 'timestamps' },
      { field: 'updatedBy', storeIdent: 'updated_by', group: 'user_reference' },
      { field: 'deletedAt', storeIdent: 'deleted_at', group: 'timestamps' },
      { field: 'deletedBy', storeIdent: 'deleted_by', group: 'user_reference' },
    ].map(sf => ({ ...sf, label: this.$t(`field:system.${sf.field}`) }))

    return {
      processing: false,
      connections: [],

      moduleFields: [],
      moduleFieldEncoding: [],
      selectedGroup: '',

      systemFields,
      systemFieldEncoding: systemFields.reduce((enc, { field }) => {
        enc[field] = systemFieldEncoding[field] || {}
        return enc
      }, {}),
      optionsGroups: [
        { text: this.$t('system-fields.grouptypes.all'), value: 'all' },
        { text: this.$t('system-fields.grouptypes.partition'), value: 'partition' },
        { text: this.$t('system-fields.grouptypes.userReference'), value: 'user_reference' },
        { text: this.$t('system-fields.grouptypes.timestamps'), value: 'timestamps' },
        { text: this.$t('system-fields.grouptypes.extras'), value: 'extras' },
      ],
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
            isMulti: f.isMulti,
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

  beforeDestroy () {
    this.setDefaultValues()
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

    applySelectedSystemFields (selectedOption) {
      this.systemFieldEncoding = this.systemFields.reduce((enc, { field, group }) => {
        if (field !== 'id') {
          if (selectedOption === 'all') {
            enc[field] = {}
          } else {
            enc[field] = group === selectedOption ? {} : { omit: true }
          }
        }
        return enc
      }, {})
    },

    getOptionKey ({ connectionID }) {
      return connectionID
    },

    setDefaultValues () {
      this.processing = false
      this.connections = []
      this.moduleFields = []
      this.moduleFieldEncoding = []
      this.selectedGroup = ''
      this.systemFields = []
      this.systemFieldEncoding = []
      this.optionsGroups = []
    },
  },
}
</script>

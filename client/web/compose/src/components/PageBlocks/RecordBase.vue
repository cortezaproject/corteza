<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <div
      v-else-if="fieldModule"
      class="mt-3"
    >
      <div
        v-for="(field, index) in fields"
        :key="index"
        :class="{ 'd-flex flex-column mb-3 px-3': canDisplay(field) }"
      >
        <template
          v-if="canDisplay(field)"
        >
          <label
            class="text-primary mb-0"
            :class="{ 'mb-0': !!(field.options.description || {}).view || false }"
          >
            {{ field.label || field.name }}
            <hint
              :id="field.fieldID"
              :text="((field.options.hint || {}).view || '')"
              class="d-inline-block"
            />
          </label>

          <small
            class="text-muted"
          >
            {{ (field.options.description || {}).view }}
          </small>
          <div
            v-if="field.canReadRecordValue"
            class="value mt-2"
          >
            <field-viewer
              v-bind="{ ...$props, field }"
              :extra-options="options"
              :record="fieldRecord"
            />
          </div>
          <i
            v-else
            class="text-primary"
          >
            {{ $t('field.noPermission') }}
          </i>
        </template>
      </div>
    </div>
  </wrap>
</template>
<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { mapActions } from 'vuex'
import base from './base'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import Hint from 'corteza-webapp-compose/src/components/Common/Hint.vue'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import conditionalFields from 'corteza-webapp-compose/src/mixins/conditionalFields'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
    Hint,
  },

  extends: base,

  mixins: [
    users,
    records,
    conditionalFields,
  ],

  data () {
    return {
      referenceRecord: undefined,
      referenceModule: undefined,
    }
  },

  computed: {
    fields () {
      if (!this.fieldModule) {
        // No module, no fields
        return []
      }

      if (!this.options.fields || this.options.fields.length === 0) {
        // No fields defined in the options, show all (buy system)
        return this.fieldModule.fields
      }

      // Show filtered & ordered list of fields
      return this.fieldModule.filterFields(this.options.fields).map(f => {
        f.label = f.isSystem ? this.$t(`field:system.${f.name}`) : f.label || f.name
        return f
      })
    },

    fieldModule () {
      return this.options.referenceField ? this.referenceModule : this.module
    },

    fieldRecord () {
      return this.options.referenceField ? this.referenceRecord : this.record
    },

    processing () {
      return !this.fieldRecord || this.evaluating
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler (recordID) {
        if (!recordID) return

        let resolutions = []

        if (recordID !== NoID) {
          resolutions = [
            this.fetchUsers(this.fields, [this.record]),
            this.fetchRecords(this.namespace.namespaceID, this.fields, [this.record]),
          ]
        }

        this.evaluating = true

        Promise.all([
          ...resolutions,
          this.evaluateExpressions(),
        ]).finally(() => {
          this.evaluating = false
        })

        if (this.options.referenceModuleID) {
          this.fetchReferenceModule(this.options.referenceModuleID)
        }
      },
    },

    options: {
      deep: true,
      handler (options) {
        if (options.referenceModuleID) {
          this.fetchReferenceModule(options.referenceModuleID)
        }
      },
    },
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
    }),

    fetchReferenceModule (moduleID) {
      if (!moduleID) {
        this.referenceModule = undefined
        return
      }

      this.findModuleByID({ namespace: this.namespace.namespaceID, moduleID: this.options.referenceModuleID })
        .then(module => {
          this.referenceModule = new compose.Module({ ...module })

          if (this.options.referenceField) {
            this.loadRecord(this.referenceModule)
          }
        })
    },

    loadRecord (module) {
      if (!module) return

      const { namespaceID, moduleID } = module
      const { referenceField } = this.options
      const field = this.module.fields.find(({ fieldID }) => fieldID === referenceField)

      const recordID = this.record.values[field.name]

      if (!recordID || !field) {
        this.referenceRecord = new compose.Record(this.fieldModule, {})
        return
      }

      if (field.isMulti) {
        this.referenceRecord = new compose.Record(this.fieldModule, {})
        return
      }

      this.$ComposeAPI.recordRead({ namespaceID, moduleID, recordID })
        .then(record => {
          this.referenceRecord = new compose.Record(this.fieldModule, { ...record })
        })
        .catch((e) => {
          this.referenceRecord = new compose.Record(this.fieldModule, {})
          this.toastErrorHandler(this.$t('notification:record.loadFailed'))(e)
        })
    },
  },
}
</script>
<style lang="scss">
.value {
  min-height: 1.2rem;
}
</style>

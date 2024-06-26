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
      ref="fieldContainer"
      class="mt-3 px-3"
      :class="fieldLayoutClass"
    >
      <template v-for="field in fields">
        <b-form-group
          v-if="canDisplay(field)"
          :key="field.id"
          :data-test-id="getFieldCypressId(field.label || field.name)"
          :label-cols-md="options.horizontalFieldLayoutEnabled && '6'"
          :label-cols-xl="options.horizontalFieldLayoutEnabled && '5'"
          :content-cols-md="options.horizontalFieldLayoutEnabled && '6'"
          :content-cols-xl="options.horizontalFieldLayoutEnabled && '7'"
          :class="columnWrapClass"
          :style="fieldWidth"
        >
          <template #label>
            <div
              class="d-flex text-primary mb-0"
            >
              <span class="d-flex">
                {{ field.label || field.name }}
              </span>

              <c-hint :tooltip="((field.options.hint || {}).view || '')" />

              <div
                v-if="options.inlineRecordEditEnabled && isFieldEditable(field)"
                class="inline-actions ml-2"
              >
                <b-button
                  v-b-tooltip.noninteractive.hover="{ title: $t('field.inlineEdit.button.title'), container: '#body' }"
                  variant="outline-light"
                  size="sm"
                  :disabled="editable"
                  class="text-secondary border-0"
                  @click="editInlineField(fieldRecord, field)"
                >
                  <font-awesome-icon
                    :icon="['fas', 'pen']"
                  />
                </b-button>
              </div>
            </div>

            <div
              class="small text-muted"
              :class="{ 'mb-1': !!(field.options.description || {}).view }"
            >
              {{ (field.options.description || {}).view }}
            </div>
          </template>

          <div
            v-if="field.canReadRecordValue"
            class="value align-self-center"
          >
            <field-viewer
              v-bind="{ ...$props, field }"
              :extra-options="options"
              :record="fieldRecord"
            />
          </div>
          <i
            v-else
            class="text-muted"
          >
            {{ $t('field.noPermission') }}
          </i>
        </b-form-group>
      </template>
    </div>

    <!-- Modal for inline editing -->
    <bulk-edit-modal
      v-if="options.inlineRecordEditEnabled"
      :namespace="namespace"
      :module="fieldModule"
      :selected-records="inlineEdit.recordIDs"
      :selected-fields="inlineEdit.fields"
      :initial-record="inlineEdit.record"
      :query="inlineEdit.query"
      :modal-title="$t('field.inlineEdit.modal.title')"
      open-on-select
      @save="onInlineEdit()"
      @close="onInlineEditClose()"
    />
  </wrap>
</template>
<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { mapActions } from 'vuex'
import axios from 'axios'
import base from './base'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import BulkEditModal from 'corteza-webapp-compose/src/components/Public/Record/BulkEdit'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import conditionalFields from 'corteza-webapp-compose/src/mixins/conditionalFields'
import recordLayout from 'corteza-webapp-compose/src/mixins/recordLayout'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
    BulkEditModal,
  },

  extends: base,

  mixins: [
    users,
    records,
    conditionalFields,
    recordLayout,
  ],

  data () {
    return {
      referenceRecord: undefined,
      referenceModule: undefined,
      inlineEdit: {
        fields: [],
        recordIDs: [],
        initialRecord: {},
      },

      abortableRequests: [],
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

    fieldLayoutClass () {
      const classes = {
        default: 'd-flex flex-column',
        noWrap: 'd-flex gap-2',
        wrap: 'row no-gutters',
      }

      return classes[this.options.recordFieldLayoutOption]
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

    fieldWidth () {
      if (this.options.recordFieldLayoutOption !== 'noWrap') {
        return {}
      }

      return { 'min-width': '13rem' }
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
          setTimeout(() => {
            this.evaluating = false
          }, 300)
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

    processing: {
      handler (newVal) {
        if (this.options.recordFieldLayoutOption !== 'wrap') return

        if (!newVal && this.fieldModule) {
          this.$nextTick(() => {
            this.initializeResizeObserver(this.$refs.fieldContainer, this.options.recordFieldLayoutOption)
          })
        }
      },
    },

    'options.recordFieldLayoutOption': {
      handler (newVal) {
        if (newVal === 'wrap' && this.fieldModule) {
          this.initializeResizeObserver(this.$refs.fieldContainer, this.options.recordFieldLayoutOption)
        } else if (this.resizeObserver) {
          this.resizeObserver.unobserve(this.$refs.fieldContainer)
          this.columnWrapClass = ''
        }
      },
    },
  },

  beforeDestroy () {
    this.abortRequests()
    this.setDefaultValues()
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

      const { response, cancel } = this.$ComposeAPI
        .recordReadCancellable({ namespaceID, moduleID, recordID })

      this.abortableRequests.push(cancel)

      response()
        .then(record => {
          this.referenceRecord = new compose.Record(this.fieldModule, { ...record })
        })
        .catch(e => {
          if (!axios.isCancel(e)) {
            this.referenceRecord = new compose.Record(this.fieldModule, {})
            this.toastErrorHandler(this.$t('notification:record.loadFailed'))(e)
          }
        })
    },

    editInlineField (record, field) {
      this.inlineEdit.fields = [field.name]
      this.inlineEdit.record = record.clone()
      this.inlineEdit.recordIDs = [record.recordID]
      this.inlineEdit.query = `recordID = ${record.recordID}`
    },

    onInlineEdit () {
      this.inlineEdit.fields = []
      this.inlineEdit.recordIDs = []
      this.inlineEdit.record = {}
      this.inlineEdit.query = ''

      this.$root.$emit('refetch-record-blocks')
    },

    onInlineEditClose () {
      this.inlineEdit.fields = []
      this.inlineEdit.record = {}
      this.inlineEdit.query = ''
    },

    setDefaultValues () {
      this.referenceRecord = undefined
      this.referenceModule = undefined
      this.inlineEdit = {}
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    getFieldCypressId (field) {
      return `field-${field.toLowerCase().split(' ').join('-')}`
    },
  },
}
</script>

<style scoped>
.field-col > * {
  margin-left: 1rem;
  margin-right: 1rem;
}
</style>

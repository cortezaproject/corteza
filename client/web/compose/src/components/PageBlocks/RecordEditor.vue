<template>
  <wrap
    v-bind="$props"
    card-class="position-static"
    v-on="$listeners"
  >
    <div
      v-if="isProcessing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <div
      v-else-if="module"
      ref="fieldContainer"
      class="mt-3"
      :class="fieldLayoutClass"
    >
      <template v-for="field in fields">
        <div
          v-if="canDisplay(field)"
          :key="`${field.fieldID || field.originalName}-${field.name}`"
          :class="`field-container ${columnWrapClass}`"
          :style="fieldWidth"
        >
          <template v-if="field.isParentField">
            <b-form-group
              :label-cols-md="horizontal && '5'"
              :label-cols-xl="horizontal && '4'"
              :content-cols-md="horizontal && '7'"
              :content-cols-xl="horizontal && '8'"
            >
              <template #label>
                <div
                  class="d-flex align-items-center text-primary mb-0"
                >
                  <span
                    class="d-flex"
                  >
                    {{ field.label || field.name }}
                  </span>

                  <c-hint :tooltip="((field.options.hint || {}).view || '')" />
                </div>

                <div
                  class="small text-muted"
                  :class="{ 'mb-1': !!(field.options.description || {}).view }"
                >
                  {{ (field.options.description || {}).view }}
                </div>
              </template>

              <div class="value align-self-center">
                <template v-if="resolvedParentRecord">
                  <field-viewer
                    :field="field"
                    :record="resolvedParentRecord"
                    :module="getModuleByID(field.parentModuleID)"
                    :namespace="namespace"
                    :extra-options="options"
                    value-only
                  />
                </template>
                <span
                  v-else
                  class="text-muted"
                >&mdash;</span>
              </div>
            </b-form-group>
          </template>

          <template v-else>
            <field-editor
              v-if="isFieldEditable(field)"
              v-bind="{ ...$props, errors: fieldErrors(field.name) }"
              :horizontal="horizontal"
              :field="field"
              :extra-options="options"
              @change="onFieldChange(field)"
            />

            <b-form-group
              v-else
              :label-cols-md="horizontal && '5'"
              :label-cols-xl="horizontal && '4'"
              :content-cols-md="horizontal && '7'"
              :content-cols-xl="horizontal && '8'"
            >
              <template #label>
                <div
                  class="d-flex align-items-center text-primary mb-0"
                >
                  <span
                    class="d-flex"
                  >
                    {{ field.label || field.name }}
                  </span>

                  <c-hint :tooltip="((field.options.hint || {}).view || '')" />
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
                  :field="field"
                  v-bind="{ ...$props, errors: fieldErrors(field.name) }"
                  value-only
                />
              </div>

              <div
                v-else
              >
                <i
                  class="text-muted"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </div>
            </b-form-group>
          </template>
        </div>
      </template>
    </div>
  </wrap>
</template>
<script>
import { NoID, compose } from '@cortezaproject/corteza-js'
import base from './base'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import conditionalFields from 'corteza-webapp-compose/src/mixins/conditionalFields'
import recordLayout from 'corteza-webapp-compose/src/mixins/recordLayout'
import { debounce } from 'lodash'
import { mapGetters } from 'vuex'
import axios from 'axios'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldEditor,
    FieldViewer,
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
      resolvedParentRecord: null,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    fields () {
      if (!this.module) {
        return []
      }

      if (!this.options.fields || this.options.fields.length === 0) {
        return this.module.fields
      }

      const childFieldConfigs = (this.options.fields || []).filter(f => !f.isParentField)
      const parentFieldConfigs = (this.options.fields || []).filter(f => f.isParentField)

      let moduleFields = []
      if (childFieldConfigs.length > 0) {
        moduleFields = this.module.filterFields(childFieldConfigs)
      } else {
        moduleFields = this.module.fields
      }

      const configured = moduleFields.map(f => {
        f.label = f.isSystem ? this.$t(`field:system.${f.name}`) : f.label || f.name
        f.isParentField = false
        return f
      })

      const parentModule = this.parentModule
      const parentConfigured = parentFieldConfigs
        .map(pf => {
          const actualField = parentModule ? parentModule.fields.find(f => f.name === pf.originalName) : null
          if (!actualField) return null
          return {
            ...actualField,
            isParentField: true,
            parentModuleID: pf.parentModuleID,
            originalName: pf.originalName,
            label: actualField.label || actualField.name,
          }
        })
        .filter(Boolean)

      if (parentFieldConfigs.length > 0 && this.options.fields.length > 0) {
        const allConfigured = [...configured, ...parentConfigured]
        const orderedFields = []

        this.options.fields.forEach(configField => {
          const match = allConfigured.find(c => {
            if (configField.isParentField) {
              return c.isParentField && c.originalName === configField.originalName && c.parentModuleID === configField.parentModuleID
            }
            return !c.isParentField && (c.name === configField.name || c.fieldID === configField.name)
          })
          if (match) orderedFields.push(match)
        })

        allConfigured.forEach(c => {
          if (!orderedFields.find(o => o === c)) {
            orderedFields.push(c)
          }
        })

        return orderedFields
      }

      return [...configured, ...parentConfigured]
    },

    fieldLayoutClass () {
      const classes = {
        default: 'd-flex flex-column px-3',
        noWrap: 'd-flex gap-2 pl-3',
        wrap: 'row no-gutters',
      }

      return classes[this.options.recordFieldLayoutOption]
    },

    fieldWidth () {
      if (this.options.recordFieldLayoutOption !== 'noWrap') {
        return {}
      }

      return { 'min-width': '20rem' }
    },

    isProcessing () {
      return this.loadingRecord || !this.record || this.evaluating
    },

    horizontal () {
      return this.block.options.horizontalFieldLayoutEnabled
    },

    isNew () {
      return this.record && this.record.recordID === NoID
    },

    parentFieldConfigs () {
      if (!this.options.fields || !this.options.fields.length) {
        return []
      }
      return this.options.fields.filter(f => f.isParentField && f.parentModuleID)
    },

    parentModule () {
      if (!this.options.includeParentFields || !this.options.parentField) {
        return null
      }

      const linkField = this.module.fields.find(f => f.name === this.options.parentField)
      if (!linkField || linkField.kind !== 'Record' || !linkField.options || !linkField.options.moduleID) {
        return null
      }

      return this.getModuleByID(linkField.options.moduleID) || null
    },

    parentLinkFieldName () {
      return this.options.parentField || null
    },
  },

  watch: {
    loadingRecord: {
      immediate: true,
      handler () {
        const { recordID } = this.record || {}

        if (!recordID || this.loadingRecord) return

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

        if (this.options.includeParentFields) {
          this.fetchParentRecord()
        }
      },
    },

    isProcessing: {
      handler (newVal) {
        if (this.options.recordFieldLayoutOption !== 'wrap') return

        if (!newVal && this.module) {
          this.$nextTick(() => {
            this.initializeResizeObserver(this.$refs.fieldContainer)
          })
        } else if (this.resizeObserver) {
          this.resizeObserver.unobserve(this.$refs.fieldContainer)
          this.columnWrapClass = ''
        }
      },
    },
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.destroyEvents()
  },

  methods: {
    async fetchParentRecord () {
      if (!this.options.includeParentFields || !this.parentLinkFieldName || !this.parentFieldConfigs.length) {
        return
      }

      if (!this.record || !this.record.recordID || this.record.recordID === NoID) {
        return
      }

      const linkField = this.module.fields.find(f => f.name === this.parentLinkFieldName)
      if (!linkField || linkField.kind !== 'Record' || !linkField.options || !linkField.options.moduleID) {
        return
      }

      const parentModuleID = linkField.options.moduleID
      const parentRecordID = this.record.values[this.parentLinkFieldName]

      if (!parentRecordID || parentRecordID === '0') {
        this.resolvedParentRecord = null
        return
      }

      const parentModule = this.getModuleByID(parentModuleID)
      if (!parentModule) return

      try {
        const { response } = this.$ComposeAPI.recordReadCancellable({
          namespaceID: this.namespace.namespaceID,
          moduleID: parentModuleID,
          recordID: parentRecordID,
        })

        const record = await response()
        this.resolvedParentRecord = new compose.Record(parentModule, { ...record })

        const selectedParentFields = this.parentFieldConfigs
          .map(pf => parentModule.fields.find(f => f.name === pf.originalName))
          .filter(Boolean)

        await Promise.all([
          this.fetchRecords(this.namespace.namespaceID, selectedParentFields, [this.resolvedParentRecord]),
          this.fetchUsers(selectedParentFields, [this.resolvedParentRecord]),
        ])
      } catch (e) {
        if (!axios.isCancel(e)) {
          console.warn('[RecordEditor] Failed to fetch parent record:', e)
          this.resolvedParentRecord = null
        }
      }
    },

    onFieldChange: debounce(function (field) {
      this.$root.$emit('record-field-change', {
        fieldName: field.name,
      })
    }, 500),

    createEvents () {
      this.$root.$on('record-field-change', this.evaluateExpressions)
    },

    destroyEvents () {
      this.$root.$off('record-field-change', this.evaluateExpressions)
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

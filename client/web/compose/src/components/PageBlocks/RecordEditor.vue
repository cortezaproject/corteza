<template>
  <wrap
    v-bind="$props"
    card-class="position-static"
    v-on="$listeners"
  >
    <div
      v-if="processing"
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
          :key="field.id"
          :class="columnWrapClass"
          :style="fieldWidth"
        >
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
        </div>
      </template>
    </div>
  </wrap>
</template>
<script>
import { validator, NoID } from '@cortezaproject/corteza-js'
import base from './base'
import users from 'corteza-webapp-compose/src/mixins/users'
import records from 'corteza-webapp-compose/src/mixins/records'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import conditionalFields from 'corteza-webapp-compose/src/mixins/conditionalFields'
import recordLayout from 'corteza-webapp-compose/src/mixins/recordLayout'
import { debounce } from 'lodash'

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

  props: {
    errors: {
      type: validator.Validated,
      required: false,
      default: undefined,
    },
  },

  computed: {
    fields () {
      if (!this.module) {
        // No module, no fields
        return []
      }

      if (!this.options.fields || this.options.fields.length === 0) {
        // No fields defined in the options, show all (but system)
        return this.module.fields
      }

      // Show filtered & ordered list of fields
      return this.module.filterFields(this.options.fields).map(f => {
        f.label = f.isSystem ? this.$t(`field:system.${f.name}`) : f.label || f.name
        return f
      })
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

    errorID () {
      const { recordID = NoID } = this.record || {}
      return recordID === NoID ? 'parent:0' : recordID
    },

    processing () {
      return !this.record || this.evaluating
    },

    horizontal () {
      return this.block.options.horizontalFieldLayoutEnabled
    },

    isNew () {
      return this.record && this.record.recordID === NoID
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
      },
    },

    processing: {
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

  beforeDestroy () {
    this.destroyEvents(this.$refs.fieldContainer)
  },

  methods: {
    /**
     * Returns errors, filtered for a specific field
     * @param name
     * @returns {validator.Validated} filtered validation results
     */
    fieldErrors (name) {
      if (!this.errors) {
        return new validator.Validated()
      }

      return this.errors
        .filterByMeta('field', name)
        .filterByMeta('id', this.errorID)
    },

    onFieldChange: debounce(function (field) {
      this.evaluateExpressions()

      this.$root.$emit('record-field-change', {
        fieldName: field.name,
      })
    }, 500),
  },
}
</script>

<style scoped>
.field-col > * {
  margin-left: 1rem;
  margin-right: 1rem;
}
</style>

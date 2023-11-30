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
      class="mt-3 px-3"
    >
      <template v-for="field in fields">
        <div
          v-if="canDisplay(field)"
          :key="field.id"
          class="mb-3"
        >
          <field-editor
            v-if="isFieldEditable(field)"
            v-bind="{ ...$props, errors: fieldErrors(field.name) }"
            :horizontal="horizontal"
            :field="field"
            @change="onFieldChange()"
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
                <span class="d-inline-block mw-100 py-1">
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
              class="value mt-1 align-self-center"
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
                class="text-primary"
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
      },
    },
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

    isFieldEditable (field) {
      if (!field) return false

      const { canCreateOwnedRecord } = this.module || {}
      const { createdAt, canManageOwnerOnRecord } = this.record || {}
      const { name, canUpdateRecordValue, isSystem, expressions = {} } = field || {}

      if (!canUpdateRecordValue) return false

      if (isSystem) {
        // Make ownedBy field editable if correct permissions
        if (name === 'ownedBy') {
          // If not created we check module permissions, otherwise the canManageOwnerOnRecord
          return createdAt ? canManageOwnerOnRecord : canCreateOwnedRecord
        }

        return false
      }

      return !expressions.value
    },

    onFieldChange: debounce(function () {
      this.evaluateExpressions()
    }, 500),
  },
}
</script>

<style lang="scss">
.value {
  min-height: 1.2rem;
}
</style>

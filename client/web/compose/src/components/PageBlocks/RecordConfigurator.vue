<template>
  <b-tab
    :title="$t('record.label')"
    no-body
  >
    <div class="px-3 pt-3">
      <h5 class="mb-3">
        {{ $t('recordList.record.generalLabel') }}
      </h5>

      <b-row>
        <b-col cols="12">
          <b-form-group
            :label="$t('general.module')"
            label-class="text-primary"
          >
            <b-form-input
              v-if="module"
              v-model="module.name"
              type="text"
              readonly
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.inlineEdit.enabled')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="inlineRecordEditEnabled"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.inlineEdit.allowAddField')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.inlineRecordEditAllowAddField"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.horizontalFormLayout')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.horizontalFieldLayoutEnabled"
              switch
              :disabled="options.recordFieldLayoutOption === 'noWrap'"
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.fieldsLayoutMode.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.recordFieldLayoutOption"
              :options="recordFieldLayoutOptions"
              :reduce="option => option.value"
              :get-option-key="option => option.label"
              @input="handleRecordFieldLayout"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.referenceRecordField')"
            :description="$t('record.referenceRecordFieldDescription')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.referenceField"
              :options="recordSelectorFields"
              :get-option-label="getFieldLabel"
              :get-option-key="getOptionKey"
              :placeholder="$t('record.referenceRecordFieldPlaceholder')"
              :reduce="getOptionKey"
              @input="updateReferenceModule($event, [])"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </div>

    <hr v-if="module">

    <div
      v-if="module"
      class="px-3"
    >
      <h5 class="mb-3">
        {{ $t('module:general.fields') }}
      </h5>

      <b-row>
        <b-col cols="12">
          <field-picker
            :module="fieldModule"
            :extra-module="parentModuleConfig.extraModule"
            :extra-module-fields="parentModuleConfig.extraModuleFields"
            :fields.sync="options.fields"
            style="height: 52vh;"
          />
        </b-col>
      </b-row>

      <b-row class="mt-2">
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.parentFields.enable')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.includeParentFields"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          v-if="options.includeParentFields"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.parentFields.field')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.parentField"
              :options="availableParentFields"
              label="label"
              :reduce="f => f.name"
              :placeholder="$t('general.label.none')"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <b-row
        v-if="isRecordFieldUsedConfigured"
        class="mt-3"
      >
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.recordSelectorDisplayOptions')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="options.recordSelectorDisplayOption"
              :options="recordDisplayOptions"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.recordSelectorCanAddRecord')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.recordSelectorShowAddRecordButton"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('record.recordSelectorAddRecordDisplayOption')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="options.recordSelectorAddRecordDisplayOption"
              :options="recordDisplayOptions"
              :disabled="!options.recordSelectorShowAddRecordButton"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </div>

    <hr>

    <div class="px-3">
      <h5 class="d-flex align-items-center mb-2">
        {{ $t('record.fieldConditions.label') }}

        <c-hint
          :tooltip="$t('record.fieldConditions.tooltip.performance')"
          icon-class="text-warning"
        />

        <b-button
          variant="link"
          :href="visibilityDocumentationURL"
          target="_blank"
          class="p-0 ml-auto"
        >
          {{ $t('general:label.examples') }}
        </b-button>
      </h5>

      <b-row class="mt-3">
        <b-col cols="12">
          <b-form-group
            :label="$t('record.fieldConditions.clearAllOnHide')"
            :description="$t('record.fieldConditions.clearAllOnHideDescription')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.clearConditionalFieldsOnHide"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <c-form-table-wrapper
        :labels="{
          addButton: $t('general:label.add')
        }"
        class="my-3"
        @add-item="addRule"
      >
        <b-table-simple
          v-if="block.options.fieldConditions.length > 0"
          borderless
          small
          responsive
        >
          <b-thead>
            <b-tr>
              <b-th
                class="text-primary"
              >
                {{ $t('record.fieldConditions.field') }}
              </b-th>
              <b-th
                class="text-primary"
              >
                {{ $t('record.fieldConditions.condition') }}
              </b-th>
              <b-th
                class="text-primary text-center"
                style="width: 150px;"
              >
                <div class="d-flex align-items-center justify-content-center">
                  {{ $t('record.fieldConditions.clearOnHide') }}
                  <c-hint
                    :tooltip="$t('record.fieldConditions.clearOnHideTooltip')"
                    class="ml-1"
                  />
                </div>
              </b-th>
              <b-th />
            </b-tr>
          </b-thead>

          <b-tbody>
            <b-tr
              v-for="(condition, i) in block.options.fieldConditions"
              :key="i"
            >
              <b-td style="width: 33%; min-width: 250px;">
                <c-input-select
                  v-model="condition.field"
                  :options="block.options.fields"
                  :placeholder="$t('record.fieldConditions.selectPlaceholder')"
                  :selectable="option => isSelectable(option)"
                  :get-option-label="getOptionLabel"
                  :get-option-key="getOptionKey"
                  :reduce="option => option.isSystem ? option.name : option.fieldID"
                />
              </b-td>

              <b-td
                class="align-middle"
                style="min-width: 300px;"
              >
                <b-input-group>
                  <b-input-group-prepend>
                    <b-input-group-text variant="extra-light">
                      ƒ
                    </b-input-group-text>
                  </b-input-group-prepend>

                  <c-input-expression
                    v-model="condition.condition"
                    auto-complete
                    :placeholder="$t('record.fieldConditions.placeholder')"
                    :suggestion-params="visibilityAutoCompleteParams"
                    class="flex-grow-1"
                  />
                </b-input-group>
              </b-td>

              <b-td style="width: 20px;">
                <div class="d-flex align-items-center justify-content-center">
                  <c-input-checkbox
                    v-model="condition.clearOnHide"
                    switch
                  />
                </div>
              </b-td>

              <b-td
                class="text-right"
                style="width: 4rem;"
              >
                <c-input-confirm
                  show-icon
                  @confirmed="deleteRule(i)"
                />
              </b-td>
            </b-tr>
          </b-tbody>
        </b-table-simple>
      </c-form-table-wrapper>

      <i18next
        path="general.visibility.condition.description.record-page"
        tag="small"
        class="text-muted"
      >
        <code>record.values.fieldName</code>
        <code>user.(userID/email...)</code>
        <code>user.userID == record.createdBy</code>
        <code>record.values.fieldName == "value"</code>
        <code>record.ownedBy == user.userID</code>
        <code>screen.width &lt; 1024</code>
      </i18next>
    </div>
  </b-tab>
</template>
<script>
import base from './base'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'
import { mapActions, mapGetters } from 'vuex'
import { NoID, compose } from '@cortezaproject/corteza-js'
import autocomplete from 'corteza-webapp-compose/src/mixins/autocomplete.js'
import { components } from '@cortezaproject/corteza-vue'

const { CInputExpression } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Record',

  components: {
    FieldPicker,
    CInputExpression,
  },

  extends: base,

  mixins: [autocomplete],

  data () {
    return {
      referenceModule: undefined,
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    visibilityDocumentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/compose-configuration/page-layouts.html#visibility-condition`
    },

    recordDisplayOptions () {
      return [
        { value: 'sameTab', text: this.$t('record.openInSameTab') },
        { value: 'newTab', text: this.$t('record.openInNewTab') },
        { value: 'modal', text: this.$t('record.openInModal') },
      ]
    },

    recordFieldLayoutOptions () {
      return [
        { value: 'default', label: this.$t('record.fieldsLayoutMode.default') },
        { value: 'noWrap', label: this.$t('record.fieldsLayoutMode.noWrap') },
        { value: 'wrap', label: this.$t('record.fieldsLayoutMode.wrap') },
      ]
    },

    recordSelectorFields () {
      return this.module.fields.filter(f => f.kind === 'Record' && !f.isMulti)
    },

    fieldModule () {
      return (this.options.referenceField && this.referenceModule) ? this.referenceModule : this.module
    },

    inlineRecordEditEnabled: {
      get () {
        return !!this.options.inlineRecordEditEnabled
      },
      set (v) {
        this.options.inlineRecordEditEnabled = v
      },
    },

    isRecordFieldUsedConfigured () {
      if (this.options.fields.length === 0) {
        return this.module.fields.some(f => f.kind === 'Record')
      } else {
        return this.options.fields.some(f => f.kind === 'Record')
      }
    },

    availableParentFields () {
      if (!this.module) {
        return []
      }

      return this.module.fields
        .filter(field => {
          if (field.kind !== 'Record') return false
          if (field.isMulti) return false
          if (!field.options || !field.options.moduleID) return false
          return true
        })
        .map(field => {
          const linkedModule = this.getModuleByID(field.options.moduleID)
          return {
            ...field,
            label: linkedModule
              ? `${field.label || field.name} → ${linkedModule.name}`
              : field.label || field.name,
          }
        })
    },

    parentModuleConfig () {
      if (!this.options.includeParentFields || !this.options.parentField) {
        return {
          extraModule: null,
          extraModuleFields: [],
        }
      }

      if (!this.module) {
        return {
          extraModule: null,
          extraModuleFields: [],
        }
      }

      const linkField = this.module.fields.find(f => f.name === this.options.parentField)
      if (!linkField || linkField.kind !== 'Record' || !linkField.options || !linkField.options.moduleID) {
        return {
          extraModule: null,
          extraModuleFields: [],
        }
      }

      const parentModule = this.getModuleByID(linkField.options.moduleID)
      if (!parentModule) {
        return {
          extraModule: null,
          extraModuleFields: [],
        }
      }

      return {
        extraModule: {
          moduleID: parentModule.moduleID,
          name: parentModule.name,
        },
        extraModuleFields: parentModule.fields.map(f => ({
          ...f,
          originalName: f.name,
        })),
      }
    },
  },

  watch: {
    'options.parentField' (newVal, oldVal) {
      if (newVal !== oldVal) {
        if (this.options.fields && this.options.fields.length > 0) {
          this.options.fields = this.options.fields.filter(f => !f.isParentField)
        }
      }
    },

    'options.includeParentFields' (newVal) {
      if (!newVal) {
        if (this.options.fields && this.options.fields.length) {
          this.options.fields = this.options.fields.filter(f => !f.isParentField)
        }
        this.options.parentField = null
      }
    },
  },

  created () {
    if (this.options.referenceField) {
      this.updateReferenceModule(this.options.referenceField, this.options.fields)
    }
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
    }),

    addRule () {
      this.options.fieldConditions.push({
        field: undefined,
        condition: '',
        clearOnHide: false,
      })
    },

    deleteRule (i) {
      this.options.fieldConditions.splice(i, 1)
    },

    isSelectable (option) {
      return !this.block.options.fieldConditions.find(({ field }) => field === option.fieldID || field === option.name) && !option.isRequired
    },

    getOptionLabel (option) {
      return option.label || option.name
    },

    getFieldLabel ({ name, label }) {
      return label || name
    },

    getOptionKey ({ fieldID, name }) {
      return fieldID !== NoID ? fieldID : name
    },

    updateReferenceModule (fieldID, fields) {
      if (!fieldID) {
        this.block.options.fields = []
        this.block.options.referenceModuleID = undefined
        return
      }

      const field = this.recordSelectorFields.find(f => f.fieldID === fieldID)
      const moduleID = field && field.options && field.options.moduleID

      if (moduleID) {
        this.findModuleByID({ namespace: this.namespace.namespaceID, moduleID })
          .then(module => {
            this.block.options.fields = fields
            this.block.options.referenceModuleID = module.moduleID
            this.referenceModule = new compose.Module({ ...module })
          })
      }
    },

    handleRecordFieldLayout (v) {
      if (v !== 'noWrap') return

      this.block.options.horizontalFieldLayoutEnabled = false
    },
  },
}
</script>

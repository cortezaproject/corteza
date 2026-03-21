<template>
  <b-tab
    :title="$t('collapsibleRecord.label')"
    no-body
  >
    <div class="px-3 pt-3">
      <h5 class="mb-3">
        {{ $t('collapsibleRecord.headerSection.label') }}
      </h5>

      <b-row>
        <!-- Title Expression -->
        <b-col cols="12">
          <b-form-group
            :label="$t('collapsibleRecord.headerSection.titleSection.label')"
            :description="$t('collapsibleRecord.headerSection.expressionHelp')"
            label-class="text-primary"
          >
            <c-input-expression
              v-model="options.titleExpression"
              auto-complete
              :placeholder="$t('collapsibleRecord.headerSection.titleSection.placeholder')"
              :suggestion-params="expressionParams"
            />
          </b-form-group>
        </b-col>

        <!-- Subtitle Expression -->
        <b-col cols="12">
          <b-form-group
            :label="$t('collapsibleRecord.headerSection.subtitleSection.label')"
            :description="$t('collapsibleRecord.headerSection.expressionHelp')"
            label-class="text-primary"
          >
            <c-input-expression
              v-model="options.subtitleExpression"
              auto-complete
              :placeholder="$t('collapsibleRecord.headerSection.subtitleSection.placeholder')"
              :suggestion-params="expressionParams"
            />
          </b-form-group>
        </b-col>

        <!-- Subtitle Show When Collapsed -->
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('collapsibleRecord.headerSection.subtitleSection.showWhenCollapsed')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.subtitleShowWhenCollapsed"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </div>

    <hr>

    <div
      v-if="module"
      class="px-3"
    >
      <h5 class="mb-3">
        {{ $t('collapsibleRecord.bodySection.label') }}
      </h5>

      <b-row>
        <!-- Body Field Selector -->
        <b-col cols="12">
          <b-form-group
            :label="$t('collapsibleRecord.bodySection.bodyField.label')"
            :description="$t('collapsibleRecord.bodySection.bodyField.description')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.bodyField"
              :options="bodyFieldOptions"
              :get-option-label="getFieldLabel"
              :get-option-key="getOptionKey"
              :placeholder="$t('collapsibleRecord.bodySection.bodyField.placeholder')"
              :reduce="getOptionKey"
            />
          </b-form-group>
        </b-col>

        <!-- Inline Edit Enabled -->
        <b-col
          v-if="options.bodyField"
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('collapsibleRecord.bodySection.bodyField.inlineEditEnabled')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.inlineEditEnabled"
              switch
              :labels="checkboxLabel"
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
        {{ $t('collapsibleRecord.otherFields.label') }}
      </h5>

      <b-row>
        <!-- Other Fields Position -->
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('collapsibleRecord.otherFields.position.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.otherFieldsPosition"
              :options="positionOptions"
              :reduce="option => option.value"
            />
          </b-form-group>
        </b-col>

        <!-- Horizontal Field Layout -->
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

        <!-- Fields Layout Mode -->
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
              :options="fieldsLayoutOptions"
              :reduce="option => option.value"
              @input="handleFieldsLayout"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <!-- Field Picker -->
      <b-row class="mt-3">
        <b-col cols="12">
          <field-picker
            :module="module"
            :fields.sync="options.fields"
            :exclude-fields="[options.bodyField]"
            style="height: 40vh;"
          />
        </b-col>
      </b-row>
    </div>

    <hr>

    <div class="px-3">
      <h5 class="mb-3">
        {{ $t('collapsibleRecord.collapseSection.label') }}
      </h5>

      <b-row>
        <!-- Default Collapsed -->
        <b-col cols="12">
          <b-form-group
            :label="$t('collapsibleRecord.collapseSection.defaultCollapsed')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="options.defaultCollapsed"
              switch
              :labels="checkboxLabel"
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
                class="text-center"
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
import { NoID } from '@cortezaproject/corteza-js'
import autocomplete from 'corteza-webapp-compose/src/mixins/autocomplete.js'
import { components } from '@cortezaproject/corteza-vue'

const { CInputExpression } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'CollapsibleRecord',

  components: {
    FieldPicker,
    CInputExpression,
  },

  extends: base,

  mixins: [autocomplete],

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },

  computed: {
    visibilityDocumentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/compose-configuration/page-layouts.html#visibility-condition`
    },

    positionOptions () {
      return [
        { value: 'default', label: this.$t('collapsibleRecord.otherFields.position.default') },
        { value: 'above', label: this.$t('collapsibleRecord.otherFields.position.above') },
        { value: 'below', label: this.$t('collapsibleRecord.otherFields.position.below') },
      ]
    },

    fieldsLayoutOptions () {
      return [
        { value: 'default', label: this.$t('record.fieldsLayoutMode.default') },
        { value: 'noWrap', label: this.$t('record.fieldsLayoutMode.noWrap') },
        { value: 'wrap', label: this.$t('record.fieldsLayoutMode.wrap') },
      ]
    },

    bodyFieldOptions () {
      if (!this.module) {
        return []
      }
      return this.module.fields
    },

    expressionParams () {
      return [
        { label: 'record.values.fieldName', text: 'record.values.fieldName' },
        { label: 'user.userID', text: 'user.userID' },
        { label: 'record.createdBy', text: 'record.createdBy' },
        { label: 'record.ownedBy', text: 'record.ownedBy' },
      ]
    },

    visibilityAutoCompleteParams () {
      return [
        { label: 'record.values.fieldName', text: 'record.values.fieldName' },
        { label: 'user.userID', text: 'user.userID' },
        { label: 'user.email', text: 'user.email' },
        { label: 'record.createdBy', text: 'record.createdBy' },
        { label: 'record.ownedBy', text: 'record.ownedBy' },
        { label: 'recordID', text: 'recordID' },
        { label: 'ownerID', text: 'ownerID' },
        { label: 'userID', text: 'userID' },
      ]
    },
  },

  methods: {
    getFieldLabel ({ name, label }) {
      return label || name
    },

    getOptionKey ({ fieldID, name }) {
      return fieldID !== NoID ? fieldID : name
    },

    handleFieldsLayout (v) {
      if (v !== 'noWrap') return

      this.block.options.horizontalFieldLayoutEnabled = false
    },

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
  },
}
</script>

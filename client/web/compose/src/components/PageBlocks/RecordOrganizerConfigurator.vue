<template>
  <b-tab :title="$t('recordOrganizer.label')">
    <b-form-group
      :label="$t('general.module')"
      label-class="text-primary"
    >
      <c-input-select
        v-model="options.moduleID"
        :options="modules"
        label="name"
        :reduce="m => m.moduleID"
        :placeholder="$t('recordOrganizer.module.placeholder')"
        default-value="0"
        required
      />
    </b-form-group>

    <div v-if="selectedModule">
      <b-form-group
        :label="$t('field.selector.available')"
        label-class="text-primary"
      >
        <div class="d-flex">
          <div class="border fields w-100 p-2">
            <div
              v-for="field in allFields"
              :key="field.name"
              class="field"
            >
              <span v-if="field.label">{{ field.label }} ({{ field.name }})</span>

              <span v-else>{{ field.name }}</span>

              <span class="small float-right">
                <span v-if="field.isSystem">{{ $t('field.selector.systemField') }}</span>

                <span v-else>{{ field.kind }}</span>
              </span>
            </div>
          </div>
        </div>
      </b-form-group>

      <b-form-group
        :label="$t('recordList.record.prefilterLabel')"
        label-class="text-primary"
      >
        <c-input-expression
          v-model.trim="options.filter"
          height="3.688rem"
          lang="javascript"
          :suggestion-params="recordAutoCompleteParams"
          :placeholder="$t('recordList.record.prefilterPlaceholder')"
        />
        <i18next
          path="recordList.record.prefilterFootnote"
          tag="small"
          class="text-muted"
        >
          <code>${record.values.fieldName}</code>
          <code>${recordID}</code>
          <code>${ownerID}</code>
          <span><code>${userID}</code>, <code>${user.name}</code></span>
        </i18next>
      </b-form-group>

      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('recordOrganizer.labelField.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.labelField"
              :options="selectedModuleFields"
              :reduce="o => o.name"
              :get-option-label="fieldLabel"
              :placeholder="$t('general.label.none')"
            />
            <b-form-text>{{ $t('recordOrganizer.labelField.footnote') }}</b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('recordOrganizer.descriptionField.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.descriptionField"
              :options="selectedModuleFields"
              :reduce="o => o.name"
              :get-option-label="descriptionLabel"
              :placeholder="$t('general.label.none')"
            />

            <b-form-text class="text-secondary small">
              {{ $t('recordOrganizer.descriptionField.footnote') }}
            </b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('recordOrganizer.groupField.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.groupField"
              :options="groupFields"
              :reduce="o => o.name"
              :get-option-label="groupFieldLabel"
              :placeholder="$t('general.label.none')"
            />

            <b-form-text class="text-secondary small">
              {{ $t('recordOrganizer.groupField.footnote') }}
            </b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('recordOrganizer.group.label')"
            label-class="text-primary"
          >
            <field-editor
              v-if="options.groupField"
              v-bind="mock"
              value-only
              class="mb-0"
            />

            <b-form-input
              v-else
              disabled
            />

            <b-form-text class="text-secondary small">
              {{ $t('recordOrganizer.group.footnote') }}
            </b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('recordOrganizer.positionField.label')"
            label-class="text-primary"
          >
            <c-input-select
              v-model="options.positionField"
              :placeholder="$t('recordOrganizer.positionField.placeholder')"
              :reduce="f => f.name"
              label="label"
            />

            <b-form-text class="text-secondary small">
              {{ $t('recordOrganizer.positionField.footnote') }}
            </b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('recordOrganizer.onRecordClick')"
            label-class="text-primary"
            class="mb-0"
          >
            <b-form-select
              v-model="options.displayOption"
              :options="displayOptions"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </div>
  </b-tab>
</template>
<script>
import FieldEditor from '../ModuleFields/Editor'
import { mapGetters } from 'vuex'
import { compose, validator } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import autocomplete from 'corteza-webapp-compose/src/mixins/autocomplete.js'

import base from './base'

const { CInputExpression } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'RecordOrganizer',

  components: {
    FieldEditor,
    CInputExpression,
  },

  extends: base,

  mixins: [autocomplete],

  data () {
    return {
      /*
        This are mocks that allow us to use the field editor component.
        Since we want all the field kinds to work properly out of the box, the field editor component is best for this case.
      */
      mock: {
        namespace: undefined,
        module: undefined,
        field: undefined,
        record: undefined,
        errors: new validator.Validated(),
      },
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
    }),

    selectedModule () {
      return this.modules.find(m => m.moduleID === this.options.moduleID)
    },

    selectedModuleFields () {
      if (this.selectedModule) {
        return [...this.selectedModule.fields].sort((a, b) => a.label.localeCompare(b.label))
      }
      return []
    },

    allFields () {
      if (this.options.moduleID) {
        return [
          ...this.selectedModuleFields,
          ...this.selectedModule.systemFields().map(sf => {
            sf.label = this.$t(`field:system.${sf.name}`)
            return sf
          }),
        ]
      }
      return []
    },

    positionFields () {
      return this.selectedModuleFields.filter(({ kind, isMulti }) => kind === 'Number' && !isMulti)
    },

    groupFields () {
      return this.selectedModuleFields.filter(({ isMulti }) => !isMulti)
    },

    group () {
      return this.allFields.find(f => f.name === this.options.groupField)
    },

    displayOptions () {
      return [
        { value: 'sameTab', text: this.$t('recordOrganizer.openInSameTab') },
        { value: 'newTab', text: this.$t('recordOrganizer.openInNewTab') },
        { value: 'modal', text: this.$t('recordOrganizer.openInModal') },
      ]
    },

    recordAutoCompleteParams () {
      return this.processRecordAutoCompleteParams({ module: this.selectedModule })
    },
  },

  watch: {
    'options.moduleID': {
      handler () {
        this.options.labelField = ''
        this.options.descriptionField = ''
        this.options.positionField = ''
        this.options.groupField = ''
      },
    },

    'options.groupField': {
      immediate: true,
      handler (newGroupField, oldGroupField) {
        // If this is not the immediate call
        if (oldGroupField) {
          this.options.group = undefined
        }

        if (newGroupField) {
          newGroupField = this.groupFields.find(f => f.name === newGroupField)
          this.mock.namespace = this.namespace
          this.mock.field = compose.ModuleFieldMaker(newGroupField)
          this.mock.field.apply({ name: 'group' })
          this.mock.module = new compose.Module({ fields: [this.mock.field] }, this.namespace)
          this.mock.record = new compose.Record(this.mock.module, { group: this.options.group })
        }
      },
    },

    'mock.record.values.group': {
      handler (group) {
        this.options.group = group
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    setDefaultValues () {
      this.mock = []
    },

    fieldLabel (option) {
      return `${option.label || option.name} (${option.kind})`
    },

    descriptionLabel (option) {
      return `${option.label || option.name} (${option.kind})`
    },

    groupFieldLabel (option) {
      return `${option.label || option.name}`
    },
  },
}
</script>
<style lang="scss" scoped>
.fields {
  height: 150px;
  overflow-y: auto;
  cursor: default;
}
</style>

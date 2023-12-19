<template>
  <b-tab :title="$t('recordOrganizer.label')">
    <b-form-group
      :label="$t('general.module')"
    >
      <c-input-select
        v-model="options.moduleID"
        :options="moduleOptions"
        label="name"
        :reduce="m => m.moduleID"
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
        :label-cols="3"
        :label="$t('recordList.record.prefilterLabel')"
        horizontal
        breakpoint="md"
        label-class="text-primary"
      >
        <b-form-textarea
          v-model.trim="options.filter"
          :value="true"
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
          <code>${userID}</code>
        </i18next>
      </b-form-group>

      <b-form-group
        :label-cols="3"
        :label="$t('recordOrganizer.labelField.label')"
        horizontal
        breakpoint="md"
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

      <b-form-group
        :label-cols="3"
        :label="$t('recordOrganizer.descriptionField.label')"
        horizontal
        breakpoint="md"
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

      <b-form-group
        :label-cols="3"
        :label="$t('recordOrganizer.positionField.label')"
        horizontal
        breakpoint="md"
        label-class="text-primary"
      >
        <b-form-select v-model="options.positionField">
          <option value="">
            {{ $t('general.label.none') }}
          </option>
          <option
            v-for="(field, index) in positionFields"
            :key="index"
            :value="field.name"
          >
            {{ field.label || field.name }}
          </option>
        </b-form-select>
        <b-form-text class="text-secondary small">
          {{ $t('recordOrganizer.positionField.footnote') }}
        </b-form-text>
      </b-form-group>

      <b-form-group
        :label-cols="3"
        :label="$t('recordOrganizer.groupField.label')"
        horizontal
        breakpoint="md"
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

      <b-form-group
        v-if="options.groupField"
        :label="$t('recordOrganizer.group.label')"
        :label-cols="3"
        breakpoint="md"
        horizontal
        label-class="text-primary"
      >
        <field-editor
          class="mb-0"
          value-only
          v-bind="mock"
        />
        <b-form-text class="text-secondary small">
          {{ $t('recordOrganizer.group.footnote') }}
        </b-form-text>
      </b-form-group>

      <b-form-group
        v-if="options.groupField"
        :label="$t('recordOrganizer.onRecordClick')"
        :label-cols="3"
        breakpoint="md"
        horizontal
        label-class="text-primary"
        class="mb-0"
      >
        <b-form-select
          v-model="options.displayOption"
          :options="displayOptions"
        />
      </b-form-group>
    </div>
  </b-tab>
</template>
<script>
import FieldEditor from '../ModuleFields/Editor'
import { mapGetters } from 'vuex'
import { compose, validator, NoID } from '@cortezaproject/corteza-js'
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'RecordOrganizer',

  components: {
    FieldEditor,
  },

  extends: base,

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

    moduleOptions () {
      return [
        { moduleID: NoID, name: this.$t('general.label.none') },
        ...this.modules,
      ]
    },

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

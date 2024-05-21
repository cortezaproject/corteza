<template>
  <div>
    <b-form-group
      :label="$t('kind.record.moduleLabel')"
      label-class="text-primary"
    >
      <c-input-select
        v-model="f.options.moduleID"
        :options="moduleOptions"
        label="name"
        :placeholder="$t('kind.record.modulePlaceholder')"
        default-value="0"
        :reduce="module => module.moduleID"
      />
    </b-form-group>

    <template
      v-if="selectedModule"
    >
      <b-form-group label-class="d-flex align-items-center text-primary p-0">
        <template #label>
          {{ $t('kind.record.moduleField') }}
          <c-hint
            :tooltip="$t('kind.record.tooltip.moduleField')"
            icon-class="text-warning"
          />
        </template>

        <c-input-select
          v-model="f.options.labelField"
          :options="fieldOptions"
          label="text"
          :placeholder="$t('kind.record.pickField')"
          :reduce="field => field.value"
        />
      </b-form-group>

      <div
        v-if="labelField && labelField.kind === 'Record'"
      >
        <b-form-group
          :label="$t('kind.record.fieldFromModuleField')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="f.options.recordLabelField"
            :options="labelFieldOptions"
            :disabled="!labelFieldModule"
            label="text"
            :placeholder="$t('kind.record.pickField')"
            :reduce="field => field.value"
          />
        </b-form-group>
      </div>

      <b-form-group
        :label="$t('kind.record.queryFieldsLabel')"
        label-class="text-primary"
      >
        <c-input-select
          v-model="f.options.queryFields"
          :options="queryFieldOptions"
          label="text"
          :reduce="field => field.value"
          :placeholder="$t('kind.record.queryFieldsPlaceholder')"
          multiple
        />
      </b-form-group>

      <b-form-group
        :label="$t('kind.record.prefilterLabel')"
        label-class="text-primary"
      >
        <b-form-textarea
          v-model="f.options.prefilter"
          :placeholder="$t('kind.record.prefilterPlaceholder')"
        />
        <i18next
          path="kind.record.prefilterFootnote"
          tag="small"
          class="text-muted"
        >
          <code>${record.values.fieldName}</code>
          <code>${recordID}</code>
          <code>${ownerID}</code>
          <span><code>${userID}</code>, <code>${user.name}</code></span>
        </i18next>
      </b-form-group>
    </template>

    <template v-if="field.isMulti">
      <b-form-group
        :label="$t('kind.select.optionType.label')"
        label-class="text-primary"
      >
        <b-form-radio-group
          v-model="f.options.selectType"
          :options="selectOptions"
          stacked
          @change="updateIsUniqueMultiValue"
        />
      </b-form-group>

      <b-form-group
        v-if="shouldAllowDuplicates"
      >
        <b-form-checkbox
          v-model="f.options.isUniqueMultiValue"
          :value="false"
          :unchecked-value="true"
        >
          {{ $t('kind.select.allow-duplicates') }}
        </b-form-checkbox>
      </b-form-group>
    </template>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import { nonQueryableFieldKinds } from 'corteza-webapp-compose/src/lib/record-filter'
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  data () {
    return {
      selectOptions: [
        { text: this.$t('kind.select.optionType.default'), value: 'default', allowDuplicates: true },
        { text: this.$t('kind.select.optionType.multiple'), value: 'multiple' },
        { text: this.$t('kind.select.optionType.each'), value: 'each', allowDuplicates: true },
      ],
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
    }),

    moduleOptions () {
      let modules = this.modules

      // If current module hasn't been created add it to modules
      if (this.module.moduleID === NoID) {
        modules = [
          ({ moduleID: '-1', name: this.module.name || this.$t('kind.record.currentUnnamedModule') }),
          ...modules,
        ]
      }

      return modules
    },

    selectedModule () {
      if (this.field.options.moduleID === '-1') {
        return this.module
      } else if (this.field.options.moduleID !== NoID) {
        return this.$store.getters['module/getByID'](this.field.options.moduleID)
      }

      return undefined
    },

    fieldOptions () {
      const fields = this.selectedModule
        ? this.selectedModule.fields
          .map(({ label, name, kind }) => {
            return { value: name, text: label || name, kind }
          })
        : []
      return [
        ...fields.sort((a, b) => a.text.localeCompare(b.text)),
      ]
    },

    queryFieldOptions () {
      return this.fieldOptions.filter(({ kind }) => !nonQueryableFieldKinds.includes(kind))
    },

    labelField () {
      if (this.field.options.labelField) {
        return this.selectedModule.fields.find(({ name }) => name === this.field.options.labelField)
      }

      return undefined
    },

    labelFieldModule () {
      if (this.labelField) {
        return this.$store.getters['module/getByID'](this.labelField.options.moduleID)
      }

      return undefined
    },

    labelFieldOptions () {
      let fields = []

      if (this.labelField && this.labelFieldModule) {
        fields = this.labelFieldModule.fields.map(({ label, name }) => { return { value: name, text: label || name } })

        return [
          ...fields.sort((a, b) => a.text.localeCompare(b.text)),
        ]
      }

      return fields
    },

    labelFieldQueryOptions () {
      return this.labelFieldOptions.filter(({ name }) => name !== this.field.options.recordLabelField)
    },

    shouldAllowDuplicates () {
      if (!this.f.isMulti) return false

      const { allowDuplicates } = this.selectOptions.find(({ value }) => value === this.f.options.selectType) || {}
      return !!allowDuplicates
    },
  },

  watch: {
    'field.options.moduleID' () {
      this.f.options.labelField = ''
      this.f.options.queryFields = []
      this.f.options.selectType = 'default'
    },
  },

  methods: {
    updateIsUniqueMultiValue (value) {
      const { allowDuplicates = false } = this.selectOptions.find(({ value: v }) => v === value) || {}
      if (!allowDuplicates) {
        this.f.options.isUniqueMultiValue = true
      }
    },
  },
}
</script>

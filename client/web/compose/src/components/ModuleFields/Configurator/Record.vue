<template>
  <div>
    <b-form-group
      :label="$t('kind.record.moduleLabel')"
    >
      <b-form-select
        v-model="f.options.moduleID"
        :options="moduleOptions"
        text-field="name"
        value-field="moduleID"
      />
    </b-form-group>

    <template
      v-if="selectedModule"
    >
      <b-form-group
        :label="$t('kind.record.moduleField')"
      >
        <b-form-select
          v-model="f.options.labelField"
          :options="fieldOptions"
        />
      </b-form-group>

      <div
        v-if="labelField && labelField.kind === 'Record'"
      >
        <b-form-group
          :label="$t('kind.record.fieldFromModuleField')"
        >
          <b-form-select
            v-model="f.options.recordLabelField"
            :options="labelFieldOptions"
            :disabled="!labelFieldModule"
          />
        </b-form-group>
      </div>

      <b-form-group
        :label="$t('kind.record.queryFieldsLabel')"
      >
        <b-form-select
          v-model="f.options.queryFields"
          class="form-control"
          :options="queryFieldOptions"
          multiple
        />
      </b-form-group>

      <b-form-group
        :label="$t('kind.record.prefilterLabel')"
      >
        <b-form-textarea
          v-model="f.options.prefilter"
          class="form-control"
          :placeholder="$t('kind.record.prefilterPlaceholder')"
        />
        <b-form-text>
          <i18next
            path="kind.record.prefilterFootnote"
            tag="label"
          >
            <code>${recordID}</code>
            <code>${ownerID}</code>
            <code>${userID}</code>
          </i18next>
        </b-form-text>
      </b-form-group>
    </template>

    <b-form-group
      v-if="field.isMulti"
      :label="$t('kind.select.optionType.label')"
    >
      <b-form-radio-group
        v-model="f.options.selectType"
        :options="selectOptions"
        stacked
      />
    </b-form-group>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  data () {
    return {
      selected: null,
      selectOptions: [
        { text: this.$t('kind.select.optionType.default'), value: 'default' },
        { text: this.$t('kind.select.optionType.multiple'), value: 'multiple' },
        { text: this.$t('kind.select.optionType.each'), value: 'each' },
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

      return [
        { moduleID: NoID, name: this.$t('kind.record.modulePlaceholder'), disabled: true },
        ...modules,
      ]
    },

    selectedModule () {
      if (this.field.options.moduleID === '-1') {
        return this.module
      } else if (this.field.options.moduleID !== NoID) {
        return this.$store.getters['module/getByID'](this.field.options.moduleID)
      } else {
        return undefined
      }
    },

    fieldOptions () {
      const fields = this.selectedModule
        ? this.selectedModule.fields
          .map(({ label, name }) => { return { value: name, text: label || name } })
        : []
      return [
        {
          value: undefined,
          text: this.$t('kind.record.pickField'),
          disabled: true,
        },
        ...fields.sort((a, b) => a.text.localeCompare(b.text)),
      ]
    },

    queryFieldOptions () {
      return this.fieldOptions.slice(1)
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
          {
            value: undefined,
            text: this.$t('kind.record.pickField'),
            disabled: true,
          },
          ...fields.sort((a, b) => a.text.localeCompare(b.text)),
        ]
      }

      return fields
    },

    labelFieldQueryOptions () {
      return this.labelFieldOptions.filter(({ name }) => name !== this.field.options.recordLabelField)
    },
  },

  watch: {
    'field.options.moduleID' () {
      this.f.options.labelField = undefined
      this.f.options.queryFields = []
      this.f.options.selectType = 'default'
    },
  },
}
</script>

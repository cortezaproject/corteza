<template>
  <b-tab :title="$t('label')">
    <b-form-group>
      <b-form-checkbox
        v-model="options.preload"
        :value="true"
        :unchecked-value="false"
      >
        {{ $t('preload') }}
      </b-form-checkbox>
    </b-form-group>
    <b-form-group
      v-if="module"
      :label="$t('fields.label')"
    >
      <b-form-checkbox
        v-model="displayAllFields"
        :value="true"
        :unchecked-value="false"
        class="mb-2"
      >
        {{ $t('fields.show-all.label') }}
      </b-form-checkbox>
      <b-table
        :items="module.fields"
        :fields="columns"
      >
        <template #cell(kind)="{ item: field }">
          {{ $t(`field:kind.${field.kind.toLowerCase()}.label`) }}
          <span
            v-if="isRecord(field)"
          >
            ({{ refModuleName(field) }})
          </span>
        </template>
        <template #cell(displayedFields)="{ item: field }">
          <b-form-checkbox
            v-if="displayAllFields"
            checked="true"
            switch
            disabled
            inline
          />
          <b-form-checkbox-group
            v-else
            v-model="options.displayedFields"
            switches
          >
            <b-form-checkbox
              :value="field.name"
              inline
            />
          </b-form-checkbox-group>
        </template>
        <template #cell(expRefFields)="{ item: field }">
          <b-form-checkbox-group
            v-if="isRef(field)"
            v-model="options.expRefFields"
            switches
          >
            <b-form-checkbox
              :value="field.name"
              inline
            />
          </b-form-checkbox-group>
        </template>
      </b-table>
    </b-form-group>
  </b-tab>
</template>
<script>
import base from './base'
import { NoID } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'recordRevisions.configurator',
  },

  name: 'RecordRevisions',

  components: {},

  extends: base,

  data () {
    const { displayedFields = [] } = this.block.options

    return {
      displayAllFields: !displayedFields || displayedFields.length === 0,

      /**
       * store displayed fields in block options
       * while user madly toggles display-all switch
       */
      displayedFieldsBackup: [...displayedFields],

      columns: [
        {
          key: 'label',
          label: this.$t('fields.columns.label.label'),
          formatter: (label, { name }) => label || name,
        },
        {
          key: 'kind',
          label: this.$t('fields.columns.kind.label'),
        },
        {
          key: 'displayedFields',
          label: this.$t('fields.columns.displayed-fields.label'),
        },
        // supporting reference expansion opens too many questions and
        // makes the current UI a bit confusing
        //
        // disabled for now
        //
        // {
        //   key: 'expRefFields',
        //   label: this.$t('fields.columns.expanded-references.label'),
        // },
      ],
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),
  },

  watch: {
    /**
     * Copy changes back to options
     * when user toggles display-all switch
     */
    displayAllFields (all) {
      this.$set(this.options, 'displayedFields', all ? [] : this.displayedFieldsBackup)
    },

    /**
     * Ensure fields are backed-up on each change
     * expect when user toggles display-all switch
     */
    'options.displayedFields' (fields) {
      if (!this.displayAllFields && fields.length > 0) {
        this.displayedFieldsBackup = fields
      }
    },
  },

  methods: {
    isRef (f) {
      return this.isRecord(f) || this.isUser(f)
    },

    isRecord (f) {
      return f.kind === 'Record'
    },

    isUser (f) {
      return f.kind === 'User'
    },

    refModuleName ({ options: { moduleID = NoID } }) {
      const m = moduleID === NoID ? null : this.getModuleByID(moduleID)
      return m ? m.name || m.handle : this.$t('errors.invalid-module-id')
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

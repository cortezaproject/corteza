<template>
  <b-tab :title="$t('label')">
    <b-row>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('sortDirection.label')"
          label-class="text-primary"
          :description="$t('sortDirection.footnote')"
        >
          <c-input-select
            v-model="options.sortDirection"
            :options="sortDirections"
            label="label"
            :clearable="false"
            :reduce="o => o.value"
          />
        </b-form-group>
      </b-col>
      <b-col
        cols="12"
        lg="6"
      >
        <b-form-group
          :label="$t('preload')"
          label-class="text-primary"
        >
          <c-input-checkbox
            v-model="options.preload"
            switch
            :labels="checkboxLabel"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <b-form-group
      v-if="module"
      :label="$t('fields.label')"
      label-class="text-primary"
    >
      <field-picker
        :module="module"
        :fields.sync="displayedFieldsArray"
        style="height: 50vh;"
      />
    </b-form-group>
  </b-tab>
</template>
<script>
import base from './base'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'

export default {
  i18nOptions: {
    namespaces: 'block',
    keyPrefix: 'recordRevisions.configurator',
  },

  name: 'RecordRevisions',

  components: {
    FieldPicker,
  },

  extends: base,

  data () {
    return {
      sortDirections: [
        { label: this.$t('sortDirection.desc'), value: 'desc' },
        { label: this.$t('sortDirection.asc'), value: 'asc' },
      ],

      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
  },

  computed: {
    displayedFieldsArray: {
      get () {
        return this.module.filterFields(this.options.displayedFields)
      },

      set (fields) {
        this.options.displayedFields = fields.map(f => f.name)
      },
    },
  },

  created () {
    if (!this.options.sortDirection) {
      this.options.sortDirection = 'desc'
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    setDefaultValues () {
      this.sortDirections = []
    },
  },
}
</script>

<template>
  <div
    v-if="module"
  >
    <h5>{{ $t('record-duplication.title') }}</h5>

    <b-form-group
      :label="$t('record-duplication.strict-fields.label')"
      label-class="text-primary pb-0"
      class="my-4"
    >
      <small>{{ $t('record-duplication.strict-fields.description') }}</small>
      <field-picker
        :module="module"
        :fields.sync="strictFields"
        :field-subset="module.fields"
        disable-system-fields
        style="max-height: 35vh;"
        class="mt-3"
      />
    </b-form-group>

    <hr>

    <b-form-group
      :label="$t('record-duplication.non-strict-fields.label')"
      label-class="text-primary pb-0"
      class="mt-4"
    >
      <small>{{ $t('record-duplication.non-strict-fields.description') }}</small>
      <field-picker
        :module="module"
        :fields.sync="nonStrictFields"
        :field-subset="module.fields"
        disable-system-fields
        style="max-height: 35vh;"
        class="mt-3"
      />
    </b-form-group>
  </div>
</template>

<script>
import { compose } from '@cortezaproject/corteza-js'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'

export default {
  i18nOptions: {
    namespaces: 'module',
    keyPrefix: 'edit.config.validation',
  },

  components: {
    FieldPicker,
  },

  props: {
    module: {
      type: compose.Module,
      required: true,
    },
  },

  computed: {
    strictFields: {
      // Get list of strict fields that are inside rules
      get () {
        return this.getRuleFields(true).map(({ attributes }) => {
          return { name: attributes[0] }
        })
      },

      // Merge current non strict fields with the new strict
      set (fields = []) {
        const fieldNames = fields.map(({ name }) => name)

        this.module.config.recordDeDup.rules = [
          ...this.getRuleFields(false).filter(({ attributes }) => !fieldNames.includes(attributes[0])),
          ...fieldNames.map(name => {
            return {
              name: 'case-sensitive',
              strict: true,
              attributes: [name],
            }
          }),
        ]
      },
    },

    nonStrictFields: {
      // Get list of non-strict fields that are inside rules
      get () {
        return this.getRuleFields(false).map(({ attributes }) => {
          return { name: attributes[0] }
        })
      },

      // Merge current strict fields with the new non-strict
      set (fields = []) {
        const fieldNames = fields.map(({ name }) => name)

        this.module.config.recordDeDup.rules = [
          ...this.getRuleFields(true).filter(({ attributes }) => !fieldNames.includes(attributes[0])),
          ...fieldNames.map(name => {
            return {
              name: 'case-sensitive',
              strict: false,
              attributes: [name],
            }
          }),
        ]
      },
    },
  },

  methods: {
    getRuleFields (strictValue) {
      return this.module.config.recordDeDup.rules.filter(({ name, strict, attributes = [] }) => {
        return strict === strictValue && name === 'case-sensitive' && attributes.length
      })
    },
  },
}
</script>

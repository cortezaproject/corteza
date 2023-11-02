<template>
  <div>
    <b-form-group
      :label="$t('kind.bool.checkedValueLabel')"
      label-class="text-primary"
    >
      <b-input-group>
        <b-form-input
          v-model="f.options.trueLabel"
          :placeholder="$t('kind.bool.checkedValuePlaceholder')"
        />

        <b-input-group-append>
          <field-bool-translator
            v-if="field"
            :field="field"
            :module="module"
            :highlight-key="`meta.bool.true.label`"
            size="sm"
            :disabled="isNew"
          />
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <b-form-group
      :label="$t('kind.bool.uncheckedValueLabel')"
      label-class="text-primary"
    >
      <b-input-group>
        <b-form-input
          v-model="f.options.falseLabel"
          :placeholder="$t('kind.bool.uncheckedValuePlaceholder')"
        />

        <b-input-group-append>
          <field-bool-translator
            v-if="field"
            :field="field"
            :module="module"
            :highlight-key="`meta.bool.fals.label`"
            size="sm"
            :disabled="isNew"
          />
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <b-form-group
      :label="$t('kind.bool.toggleTypeLabel')"
      label-class="text-primary"
    >
      <b-form-checkbox v-model="f.options.switch">
        {{ $t('kind.bool.toggleType') }}
      </b-form-checkbox>
    </b-form-group>
  </div>
</template>

<script>
import base from './base'
import { NoID } from '@cortezaproject/corteza-js'
import FieldBoolTranslator from 'corteza-webapp-compose/src/components/Admin/Module/FieldBoolTranslator'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  components: {
    FieldBoolTranslator,
  },

  extends: base,

  computed: {
    isNew () {
      return this.module.moduleID === NoID || this.field.fieldID === NoID
    },
  },
}
</script>

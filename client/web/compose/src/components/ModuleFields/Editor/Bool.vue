<template>
  <b-form-group
    label-class="d-flex align-items-center text-primary p-0"
  >
    <template
      v-if="field.options.switch"
      #label
    >
      {{ label }}

      <hint
        :id="field.fieldID"
        :text="hint"
      />
    </template>

    <div
      v-if="field.options.switch"
      class="small text-muted"
      :class="{ 'mb-1': description }"
    >
      {{ description }}
    </div>

    <c-input-checkbox
      v-model="value"
      :switch="field.options.switch"
      :labels="field.options.switch ? checkboxLabel : {}"
    >
      <div
        v-if="!field.options.switch"
        class="d-flex align-items-center text-primary"
      >
        {{ label }}

        <hint
          :id="field.fieldID"
          :text="hint"
        />
      </div>
    </c-input-checkbox>

    <div
      v-if="!field.options.switch"
      class="small text-muted"
    >
      {{ description }}
    </div>
    <errors :errors="errors" />
  </b-form-group>
</template>
<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  computed: {
    value: {
      get () {
        return this.record.values[this.field.name] === '1'
      },

      set (value) {
        this.record.values[this.field.name] = value ? '1' : '0'
        this.$emit('change', value)
      },
    },

    checkboxLabel () {
      return {
        on: this.field.options.trueLabel || this.$t('general:label.yes'),
        off: this.field.options.falseLabel || this.$t('general:label.no'),
      }
    },
  },
}
</script>

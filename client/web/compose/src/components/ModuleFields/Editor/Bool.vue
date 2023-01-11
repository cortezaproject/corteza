<template>
  <b-form-group>
    <b-form-checkbox v-model="value">
      {{ field.label || field.name }}
      <hint
        :id="field.fieldID"
        :text="hint"
        class="d-inline-block text-primary"
      />
    </b-form-checkbox>

    <small
      class="form-text font-weight-light text-muted"
    >
      {{ description }}
    </small>
    <errors :errors="errors" />
  </b-form-group>
</template>
<script>
import base from './base'

export default {
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
  },
}
</script>

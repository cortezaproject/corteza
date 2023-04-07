<template>
  <b-form-group
    :label="label"
    :label-cols="labelCols"
  >
    <b-form-input
      ref="colorPickerInput"
      data-test-id="input-color-picker"
      v-model="value"
      type="color"
      :value="value"
      :title="tooltip"
      :disabled="disabled"
      :debounce="debounce"
      :style="{'cursor': disabled ? 'not-allowed' : 'default'}"
      @input="value => $emit('input', value)"
    />
  </b-form-group>
</template>

<script>
export default {
  name: 'CInputColorPicker',

  props: {
    value: {
      type: String,
      default: '',
    },

    label: {
      type: String,
      default: '',
    },

    tooltip: {
      type: String,
      default: '',
    },

    labelCols: {
      type: Number,
      default: 4,
    },

    disabled: {
      type: Boolean,
    },

    debounce: {
      type: Number,
      default: 300,
    },
  },
}
</script>
<style>
input:disabled {
  background: none !important;
}
</style>
<template>
  <chrome
    data-test-id="input-color-picker"
    :value="value || 'rgba(0,0,0,0)'"
    @input="color => updateColor(color)"
  />
</template>

<script>
import { debounce } from 'lodash'
import { Chrome } from 'vue-color'

export default {
  name: 'CInputColorPicker',

  components: {
    Chrome,
  },

  props: {
    value: {
      type: String,
      default: '',
    },
  },

  methods: {
    updateColor: debounce(function (color = '') {
      this.$emit('input', color.hex8)
    }, 300),
  },
}
</script>

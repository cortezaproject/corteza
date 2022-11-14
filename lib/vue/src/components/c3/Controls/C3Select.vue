<template>
  <b-form-group
    :label="label"
  >
    <b-form-select
      :value="value"
      :options="options"
      @input="$emit('update', $event)"
    />
  </b-form-group>
</template>

<script>
export default {
  props: {
    label: {
      required: true,
      type: String,
    },
    value: {
      type: String,
      default: () => undefined,
    },
    options: {
      default: () => [],
      type: Array,
    },
  },
}
</script>

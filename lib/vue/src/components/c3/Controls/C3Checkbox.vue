<template>
  <b-form-group>
    <b-form-checkbox
      :checked="value"
      @change="$emit('update', $event)"
    >
      {{ label }}
    </b-form-checkbox>
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
      type: Boolean,
      default: () => false,
    },
  },
}
</script>

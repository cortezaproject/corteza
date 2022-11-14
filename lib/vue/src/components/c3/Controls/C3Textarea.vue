<template>
  <b-form-group
    :label="label"
  >
    <b-form-textarea
      :value="value"
      @update="$emit('update', $event)"
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
      default: () => '',
    },
  },
}
</script>

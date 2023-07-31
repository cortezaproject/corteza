<template>
  <c-input-select
    key="name"
    :value="value"
    :options="columns"
    :get-option-label="getColumnLabel"
    :get-option-key="getOptionKey"
    :placeholder="$t('general:label.none')"
    :reduce="r => r.name"
    append-to-body
    v-on="$listeners"
  />
</template>

<script>
export default {

  props: {
    columns: {
      type: Array,
      required: true,
    },

    value: {
      type: String,
      default: '',
    },
  },

  methods: {
    getColumnLabel ({ name, label }) {
      return `${label} (${name})`
    },

    getOptionKey ({ name }) {
      return name
    },
  },
}
</script>

<style lang="scss">
.column-selector {
  min-width: 15vw !important;
  max-width: 150px !important;
}
</style>

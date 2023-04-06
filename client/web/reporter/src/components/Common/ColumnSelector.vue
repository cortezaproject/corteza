<template>
  <vue-select
    key="name"
    :value="value"
    :options="columns"
    :get-option-label="getColumnLabel"
    :get-option-key="getOptionKey"
    :placeholder="$t('general:label.none')"
    :reduce="r => r.name"
    append-to-body
    :calculate-position="calculateDropdownPosition"
    class="column-selector bg-white"
    v-on="$listeners"
  />
</template>

<script>
import { VueSelect } from 'vue-select'

export default {
  components: {
    VueSelect,
  },

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

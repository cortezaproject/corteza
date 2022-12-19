<template>
  <vue-select
    key="name"
    :value="value"
    :options="columns"
    :get-option-label="getColumnLabel"
    :placeholder="$t('general:label.none')"
    :reduce="r => r.name"
    append-to-body
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
  },
}
</script>

<style lang="scss">
.column-selector {
  min-width: 15vw !important;
  max-width: 150px !important;
}
</style>

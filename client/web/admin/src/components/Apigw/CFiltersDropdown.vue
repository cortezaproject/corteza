<template>
  <b-dropdown
    :text="$t('filters.addFilter')"
    class="min-content"
    variant="primary"
  >
    <template v-if="filterList.length">
      <b-dropdown-item
        v-for="(filter, index) in filterList"
        :key="index"
        :disabled="filter.disabled"
        href="#"
        @click="onAddFilter(filter)"
      >
        {{ filter.label }}
      </b-dropdown-item>
    </template>
    <b-dropdown-item
      v-else
      disabled
      href="#"
    >
      <span class="text-danger">
        {{ $t('filters.filterListEmpty') }}
      </span>
    </b-dropdown-item>
  </b-dropdown>
</template>

<script>
export default {
  props: {
    availableFilters: {
      type: Array,
      required: true,
    },
    filters: {
      type: Array,
      required: true,
    },
  },

  computed: {
    filterList () {
      return this.availableFilters.map(f => {
        return { ...f, disabled: !!(this.filters || []).some(filter => filter.ref === f.ref) }
      })
    },
  },

  methods: {
    onAddFilter (filter) {
      const add = { ...filter, created: true, params: [] }
      const { params = [] } = filter

      for (const p of params) {
        add.params.push({ ...p, options: { ...p.options } })
      }

      this.$emit('addFilter', add)
    },
  },
}
</script>

<style scoped>
.min-content{
    max-width: min-content;
}
</style>

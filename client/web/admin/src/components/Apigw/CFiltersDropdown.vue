<template>
  <b-dropdown
    data-test-id="dropdown-add-filter"
    :text="$t('filters.addFilter')"
    variant="primary"
    size="lg"
    class="d-inline-flex mt-3"
  >
    <template v-if="filterList.length">
      <b-dropdown-item
        v-for="(filter, index) in filterList"
        :key="index"
        :data-test-id="filterDropdownCypressId(filter.label)"
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
      <span
        data-test-id="filter-list-empty"
        class="text-danger"
      >
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

    filterDropdownCypressId (filter) {
      return filter.toLowerCase().split(' ').join('-')
    },
  },
}
</script>

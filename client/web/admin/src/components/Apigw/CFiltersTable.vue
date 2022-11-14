<template>
  <div>
    <b-table-simple
      class="filter-table"
      hover
    >
      <b-thead>
        <b-tr>
          <b-th>{{ $t('filters.list.filters') }}</b-th>
          <b-th>{{ $t('filters.list.status') }}</b-th>
          <b-th />
        </b-tr>
      </b-thead>

      <draggable
        v-model="sortableFilters"
        tag="b-tbody"
      >
        <b-tr
          v-for="(filter, index) in sortableFilters"
          :key="index"
          class="pointer"
          :class="[selectedRow===index ? 'row-selected' : 'row-not-selected']"
          @click.stop="onRowClick(filter,index)"
        >
          <b-td class="align-baseline">
            {{ filter.label }}
          </b-td>
          <b-td class="align-baseline">
            {{ $t(`filters.${filter.enabled ? 'enabled' : 'disabled'}`) }}
          </b-td>
          <b-td class="text-right align-baseline">
            <b-button
              variant="danger"
              class="my-1"
              size="sm"
              @click.stop="onRemoveFilter(filter)"
            >
              {{ $t('filters.list.remove') }}
            </b-button>
          </b-td>
        </b-tr>
      </draggable>
    </b-table-simple>
    <h6
      v-if="!sortableFilters.length"
      class="d-flex justify-content-center align-items-center mb-3"
    >
      {{ $t('filters.list.noFilters') }}
    </h6>
  </div>
</template>

<script>
import draggable from 'vuedraggable'
export default {
  components: {
    draggable,
  },
  props: {
    filters: {
      type: Array,
      required: true,
    },
    step: {
      type: Number,
      default: () => 0,
    },
  },

  data () {
    return {
      selectedRow: 0,
      selectedFilter: {},
    }
  },
  computed: {
    sortableFilters: {
      get () {
        return this.filters
      },

      set (v) {
        this.$emit('sortFilters', v)
      },
    },
  },
  methods: {
    onAddFilter (filter) {
      if (!this.filters.find(f => f.ref === filter.ref)) {
        this.filters.push({ ...filter })
      }

      if (this.filters.length === 1) {
        this.$emit('filterSelect', filter)
      }
    },

    onRemoveFilter (filter) {
      this.$emit('removeFilter', filter)
    },

    onRowClick (filter, index) {
      this.selectedRow = index
      this.selectedFilter = filter
      this.$emit('filterSelect', filter)
    },
  },
}
</script>

<style lang="scss">
.filter-table .row-selected{
  background: #F3F3F5;
}
.cursor-default{
  cursor: default;
}
</style>

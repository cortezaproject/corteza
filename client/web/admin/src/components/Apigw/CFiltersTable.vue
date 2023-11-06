<template>
  <div>
    <b-table-simple
      hover
      class="mb-0"
    >
      <b-thead head-variant="light">
        <b-tr>
          <b-th />
          <b-th>{{ $t('filters.list.filters') }}</b-th>
          <b-th>{{ $t('filters.list.status') }}</b-th>
          <b-th />
        </b-tr>
      </b-thead>

      <draggable
        v-if="!fetching"
        v-model="sortableFilters"
        :options="{ handle: '.handle' }"
        tag="b-tbody"
      >
        <b-tr
          v-for="(filter, index) in sortableFilters"
          :key="index"
        >
          <td
            class="handle align-middle grab"
            style="width: 1%"
          >
            <font-awesome-icon
              :icon="['fas', 'bars']"
              class="text-light"
            />
          </td>
          <b-td class="align-middle">
            {{ filter.label }}
          </b-td>
          <b-td class="align-middle">
            {{ $t(`filters.${filter.enabled ? 'enabled' : 'disabled'}`) }}
          </b-td>
          <b-td class="text-right align-middle">
            <b-button
              size="sm"
              variant="link"
              @click.stop="onRowClick(filter, index)"
            >
              <font-awesome-icon
                :icon="['far', 'edit']"
              />
            </b-button>
            <c-input-confirm
              show-icon
              class="ml-1"
              @confirmed="onRemoveFilter(filter)"
              @click.stop
            />
          </b-td>
        </b-tr>
      </draggable>
    </b-table-simple>

    <div
      class="d-flex flex-column align-items-center justify-content-center h-100 overflow-hidden"
    >
      <b-spinner
        v-if="fetching"
        class="my-4"
      />

      <p
        v-else-if="!sortableFilters.length"
        data-test-id="no-filters"
        class="my-4"
      >
        {{ $t('filters.list.noFilters') }}
      </p>
    </div>
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
    fetching: {
      type: Boolean,
      value: false,
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

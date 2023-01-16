<template>
  <b-card
    data-test-id="card-filter-list"
    header-bg-variant="white"
    footer-bg-variant="white"
    body-class="p-0"
    class="apigw shadow-sm mt-3"
  >
    <c-filter-modal
      :visible="!!selectedFilter"
      :filter="selectedFilter"
      @submit="onSubmit"
      @reset="onReset"
    />

    <b-tabs
      data-test-id="filter-steps"
      active-nav-item-class="active-tab bg-white"
      class="border-0"
      content-class="border-bottom"
    >
      <b-tab
        v-for="(step, index) in steps"
        :key="index"
        class="border-0"
        :title="$t(`filters.step_title.${step}`)"
        @click="onActivateTab(index)"
      >
        <c-filters-table
          :filters="getSelectedFiltersByStep"
          :step="index"
          :fetching="fetching"
          @filterSelect="onFilterSelect"
          @removeFilter="onRemoveFilter"
          @sortFilters="onSortFilters"
        />
      </b-tab>
    </b-tabs>

    <template #header>
      <h3 class="m-0">
        {{ $t('filters.title') }}
      </h3>
    </template>

    <template #footer>
      <div
        class="d-flex justify-content-between"
      >
        <c-filters-dropdown
          :available-filters="getAvailableFiltersByStep"
          :filters="getSelectedFiltersByStep"
          @addFilter="onAddFilter"
        />
        <c-submit-button
          :processing="processing"
          :success="success"
          :disabled="disabled"
          @submit="$emit('submit')"
        />
      </div>
    </template>
  </b-card>
</template>
<script>
import CFilterModal from 'corteza-webapp-admin/src/components/Apigw/CFilterModal'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import CFiltersTable from 'corteza-webapp-admin/src/components/Apigw/CFiltersTable'
import CFiltersDropdown from 'corteza-webapp-admin/src/components/Apigw/CFiltersDropdown'

const mapKindToStep = {
  prefilter: 0,
  processer: 1,
  postfilter: 2,
}

export default {
  components: {
    CFilterModal,
    CSubmitButton,
    CFiltersTable,
    CFiltersDropdown,
  },
  props: {
    fetching: {
      type: Boolean,
      value: false,
    },
    processing: {
      type: Boolean,
      value: false,
    },
    success: {
      type: Boolean,
      value: false,
    },
    filters: {
      type: Array,
      required: true,
    },
    availableFilters: {
      type: Array,
      required: true,
    },
    steps: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
      selectedFilter: null,
      selectedTab: 0,
    }
  },

  computed: {
    disabled () {
      return !(this.filters.some(({ updated, created, deleted }) => updated || created || deleted))
    },

    getSelectedFilter () {
      return this.selectedFilter ? this.selectedFilter : null
    },

    getAvailableFiltersByStep () {
      return (this.availableFilters || []).filter(({ kind }) => {
        return mapKindToStep[kind] === this.selectedTab
      })
    },

    getSelectedFiltersByStep () {
      return (this.filters || []).filter(({ kind, deleted }) => {
        return mapKindToStep[kind] === this.selectedTab && !deleted
      }).sort((a, b) => a.weight - b.weight)
    },

    disabledRemoveButton () {
      return !this.filters.some(({ options }) => (options || { checked: false }).checked)
    },
  },

  methods: {
    onAddFilter (filter) {
      const i = this.filters.findIndex(({ ref, deleted }) => ref === filter.ref && !deleted)

      if (i < 0) {
        this.selectedFilter = filter
      } else {
        this.selectedFilter = this.filters[i]
      }
    },

    onSubmit (filter) {
      const i = this.filters.findIndex(({ ref, deleted }) => ref === filter.ref && !deleted)

      if (i < 0) {
        filter.weight = this.getSelectedFiltersByStep.length
        this.filters.push(filter)
      } else {
        this.filters.splice(this.filters.findIndex(({ ref, deleted }) => ref === filter.ref && !deleted), 1, filter)
      }

      this.$emit('update:filters', this.filters)
    },

    onReset () {
      this.selectedFilter = undefined
    },

    onSortFilters (sortedFilters) {
      this.filters.forEach(filter => {
        const i = sortedFilters.findIndex(({ ref, deleted }) => ref === filter.ref && !deleted)

        if (i >= 0) {
          filter.weight = sortedFilters.findIndex(({ ref, deleted }) => ref === filter.ref && !deleted)
          filter.updated = true
        }
      })
      this.filters.sort((a, b) => a.weight - b.weight)
    },

    onRemoveFilter (filter) {
      if (filter.filterID) {
        this.filters.splice(this.filters.findIndex(({ filterID, deleted }) => filterID === filter.filterID && !deleted), 1, { ...filter, deleted: true })
      } else {
        this.filters.splice(this.filters.findIndex(({ ref, deleted }) => ref === filter.ref && !deleted), 1)
      }

      this.$emit('update:filters', this.filters)
    },

    onFilterSelect (filter = {}) {
      this.selectedFilter = { ...filter }
    },

    onActivateTab (index) {
      this.selectedTab = index
    },
  },
}
</script>
<style lang="scss" >

.apigw {
  .nav-link {
    color: $primary;
    border-width: 0 0 3px 0 !important;
    border-color: transparent !important;
  }

  .active-tab {
    color: $primary !important;
    border-color: $primary !important;
  }
}
</style>

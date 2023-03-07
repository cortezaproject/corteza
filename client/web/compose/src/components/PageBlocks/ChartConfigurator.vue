<template>
  <b-tab :title="$t('chart.label')">
    <b-form-group
      :label="$t('chart.display')"
    >
      <b-input-group class="d-flex w-100">
        <vue-select
          v-model="block.options.chartID"
          :options="charts"
          :placeholder="$t('chart.pick')"
          :reduce="option => option.chartID"
          label="name"
          append-to-body
          class="chart-selector bg-white"
          @input="chartSelected"
        />
        <b-input-group-append>
          <b-button
            :title="$t('chart.openInBuilder')"
            :disabled="!selectedChart || (!selectedChart.canUpdateChart && !selectedChart.canDeleteChart)"
            variant="light"
            class="d-flex align-items-center"
            :to="{ name: 'admin.charts.edit', params: { chartID: (selectedChart || {}).chartID }, query: null }"
          >
            <font-awesome-icon :icon="['fas', 'external-link-alt']" />
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <template v-if="isDrillDownAvailable">
      <b-form-group
        :description="$t('chart.drillDown.description')"
        label-class="d-flex align-items-center"
        class="mb-1"
      >
        <template #label>
          {{ $t('chart.drillDown.label') }}
          <b-form-checkbox
            v-model="options.drillDown.enabled"
            switch
            class="ml-1"
          />
        </template>

        <vue-select
          v-model="options.drillDown.blockID"
          :options="drillDownOptions"
          :disabled="!options.drillDown.enabled"
          :get-option-label="o => o.title || o.kind"
          :reduce="option => option.blockID"
          :clearable="true"
          :placeholder="$t('chart.drillDown.openInModal')"
          append-to-body
          class="block-selector bg-white"
        />
      </b-form-group>
    </template>
  </b-tab>
</template>
<script>
import base from './base'
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import { VueSelect } from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Chart',

  components: {
    VueSelect,
  },

  extends: base,

  computed: {
    ...mapGetters({
      charts: 'chart/set',
    }),

    selectedChart () {
      if (!this.options.chartID || this.options.chartID === NoID) {
        return
      }

      return this.charts.find(({ chartID }) => chartID === this.options.chartID)
    },

    selectedChartModuleID () {
      if (!this.selectedChart) return

      const { moduleID } = (this.selectedChart.config.reports[0] || {})

      return moduleID
    },

    isDrillDownAvailable () {
      if (!this.selectedChart) return

      const { metrics = [] } = (this.selectedChart.config.reports[0] || {})

      return !metrics.some(({ type }) => type === 'gauge')
    },

    drillDownOptions () {
      return this.page.blocks.filter(({ blockID, kind, options = {} }) => kind === 'RecordList' && blockID !== NoID && options.moduleID === this.selectedChartModuleID)
    },
  },

  methods: {
    chartSelected () {
      this.options.drillDown = {
        enabled: false,
        blockID: '',
      }
    },
  },
}
</script>

<style lang="scss">
.chart-selector {
  position: relative;
  -ms-flex: 1 1 auto;
  flex: 1 1 auto;
  width: 1%;
  margin-bottom: 0;
}

.block-selector, .chart-selector {
  &:not(.vs--open) .vs__selected + .vs__search {
    // force this to not use any space
    // we still need it to be rendered for the focus
    width: 0;
    padding: 0;
    margin: 0;
    border: none;
    height: 0;
  }

  .vs__selected-options {
    // do not allow growing
    width: 0;
  }

  .vs__selected {
    display: block;
    white-space: nowrap;
    text-overflow: ellipsis;
    max-width: 100%;
    overflow: hidden;
  }
}

.vs__dropdown-menu .vs__dropdown-option {
  text-overflow: ellipsis;
  overflow: hidden !important;
}
</style>

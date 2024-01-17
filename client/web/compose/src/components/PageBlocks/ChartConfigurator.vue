<template>
  <b-tab :title="$t('chart.label')">
    <b-form-group
      :label="$t('chart.display')"
      label-class="text-primary"
    >
      <b-input-group class="d-flex w-100">
        <c-input-select
          v-model="block.options.chartID"
          :options="charts"
          :get-option-key="getOptionKey"
          :placeholder="$t('chart.pick')"
          :reduce="option => option.chartID"
          label="name"
          :selectable="c => !c.deletedAt"
          @input="chartSelected"
        />

        <b-input-group-append>
          <b-button
            v-b-tooltip.hover="{ title: $t(chartSelectorTooltip), container: '#body' }"
            :disabled="selectedChart && (!selectedChart.canUpdateChart && !selectedChart.canDeleteChart)"
            variant="extra-light"
            class="d-flex align-items-center"
            :to="{ name: chartExternalLink, params: { chartID: (selectedChart || {}).chartID }, query: null }"
          >
            <font-awesome-icon :icon="['fas', 'external-link-alt']" />
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </b-form-group>

    <template v-if="isDrillDownAvailable">
      <b-form-group
        :label="$t('chart.drillDown.label')"
        :description="$t('chart.drillDown.description')"
        label-class="d-flex align-items-center text-primary"
        class="mb-1"
      >
        <template #label>
          {{ $t('chart.drillDown.label') }}

          <c-input-checkbox
            v-model="options.drillDown.enabled"
            switch
            class="ml-1"
          />
        </template>

        <b-input-group>
          <c-input-select
            v-model="options.drillDown.blockID"
            :options="drillDownOptions"
            :get-option-key="getOptionKey"
            :disabled="!options.drillDown.enabled"
            :get-option-label="o => o.title || o.kind"
            :reduce="option => option.blockID"
            :clearable="true"
            :placeholder="$t('chart.drillDown.openInModal')"
          />

          <b-input-group-append>
            <column-picker
              ref="columnPicker"
              :module="selectedChartModule"
              :fields="selectedDrilldownFields"
              :disabled="!!options.drillDown.blockID || !options.drillDown.enabled"
              variant="extra-light"
              size="md"
              @updateFields="onUpdateFields"
            >
              <font-awesome-icon :icon="['fas', 'wrench']" />
            </column-picker>
          </b-input-group-append>
        </b-input-group>
      </b-form-group>
    </template>
  </b-tab>
</template>
<script>
import base from './base'
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'
import ColumnPicker from 'corteza-webapp-compose/src/components/Admin/Module/Records/ColumnPicker'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Chart',

  components: {
    ColumnPicker,
  },

  extends: base,

  computed: {
    ...mapGetters({
      charts: 'chart/set',
      getModuleByID: 'module/getByID',
    }),

    selectedChart () {
      if (!this.options.chartID || this.options.chartID === NoID) {
        return
      }

      return this.charts.find(({ chartID }) => chartID === this.options.chartID)
    },

    chartExternalLink () {
      return !this.selectedChart ? 'admin.charts' : 'admin.charts.edit'
    },

    chartSelectorTooltip () {
      return !this.selectedChart ? 'chart.openChartList' : 'chart.openInBuilder'
    },

    selectedChartModuleID () {
      if (!this.selectedChart) return

      const { moduleID } = (this.selectedChart.config.reports[0] || {})

      return moduleID
    },

    selectedChartModule () {
      if (!this.selectedChartModuleID) return

      return this.getModuleByID(this.selectedChartModuleID)
    },

    selectedDrilldownFields () {
      if (!this.selectedChart) return []

      return this.options.drillDown.recordListOptions.fields
    },

    isDrillDownAvailable () {
      if (!this.selectedChart) return

      const { metrics = [] } = (this.selectedChart.config.reports[0] || {})

      return !metrics.some(({ type }) => type === 'gauge' || type === 'radar')
    },

    drillDownOptions () {
      return this.blocks.filter(({ blockID, kind, options = {} }) => kind === 'RecordList' && blockID !== NoID && options.moduleID === this.selectedChartModuleID)
    },
  },

  methods: {
    chartSelected () {
      this.block.resetDrillDown()
    },

    getOptionKey ({ chartID }) {
      return chartID
    },

    onUpdateFields (fields) {
      this.options.drillDown.recordListOptions.fields = fields.map(({ fieldID }) => fieldID)
    },
  },
}
</script>

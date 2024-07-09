<template>
  <wrap
    v-bind="$props"
    body-class="position-relative"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <chart-component
      v-if="chart"
      :key="key"
      :chart="chart"
      :record="record"
      :reporter="reporter"
      @drill-down="drillDown"
    />

    <template v-if="options.liveFilterEnabled">
      <b-button
        variant="outline-light"
        class="chart__livefilter-btn position-absolute d-flex d-print-none border-0 px-1"
        :class="[!!liveFilterValue && !liveFilterModal.show ? 'text-primary' : 'text-secondary']"
        @click="showFilterModal"
      >
        <font-awesome-icon
          :icon="['fas', 'filter']"
        />
      </b-button>

      <b-modal
        v-model="liveFilterModal.show"
        :title="$t('chart.filter.modal.title')"
        :ok-title="$t('general:label.saveAndClose')"
        centered
        size="lg"
        cancel-variant="light"
        no-fade
        @ok="updateLiveFilter"
        @hide="liveFilterModal.show = undefined"
      >
        <b-form-group
          v-if="originalFilter"
          :label="$t('chart.filter.modal.originalFilter.label')"
          label-class="text-primary"
        >
          <b-form-textarea
            :value="originalFilter"
            readonly
          />
        </b-form-group>

        <b-form-group
          :label="$t('chart.filter.modal.liveFilter.label')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="liveFilterModal.value"
            :options="predefinedFilters"
            label="text"
            :reduce="filter => filter.value"
            :placeholder="$t('chart.filter.modal.liveFilter.placeholder')"
          />
        </b-form-group>

        <div v-if="originalFilter && liveFilterModal.value">
          <hr class="my-3">

          <b-form-group
            :label="$t('chart.filter.modal.filterPreview.label')"
            label-class="text-primary"
          >
            <b-form-textarea
              :value="liveFilterPreview"
              readonly
            />
          </b-form-group>

          <b-form-group
            :label="$t('chart.filter.modal.options.label')"
            label-class="text-primary pb-0"
          >
            <b-form-radio
              v-for="(option, ci) in filterOptions"
              :key="ci"
              v-model="liveFilterModal.option"
              :value="option.value"
            >
              {{ $t(`chart.filter.modal.options.${option.label}`) }}
            </b-form-radio>
          </b-form-group>
        </div>

        <template #modal-footer="{ ok }">
          <div
            class="d-flex justify-content-between align-items-center w-100"
          >
            <b-button
              type="button"
              variant="light"
              @click="resetLiveFilter()"
            >
              {{ $t('general:label.reset') }}
            </b-button>

            <div class="d-flex gap-1">
              <b-button
                variant="light"
                type="button"
                rounded
                @click="cancelLiveFilter()"
              >
                {{ $t('general:label.cancel') }}
              </b-button>
              <b-button
                variant="primary"
                @click="ok()"
              >
                {{ $t('general:label.save') }}
              </b-button>
            </div>
          </div>
        </template>
      </b-modal>
    </template>
  </wrap>
</template>

<script>
import { mapActions } from 'vuex'
import base from './base'
import ChartComponent from '../Chart'
import { NoID, compose } from '@cortezaproject/corteza-js'
import { evaluatePrefilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  components: {
    ChartComponent,
  },

  extends: base,

  data () {
    return {
      chart: null,

      originalFilter: undefined,
      filter: undefined,

      drillDownFilter: undefined,

      liveFilterModal: {
        show: false,
        value: undefined,
        option: 'AND',
      },

      predefinedFilters: [
        ...compose.chartUtil.predefinedFilters.map(pf => ({ ...pf, text: this.$t(`chart:edit.filter.${pf.text}`) })),
      ],

      selectedFilter: undefined,

      customDate: {
        start: undefined,
        end: undefined,
      },

      liveFilterValue: undefined,
      liveFilterOption: 'AND',

      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },

      filterOptions: [
        { label: 'and', value: 'AND' },
        { label: 'or', value: 'OR' },
        { label: 'overwrite', value: '' },
      ],
    }
  },

  computed: {
    isDrillDownEnabled () {
      if (!this.options) return false

      return this.options.drillDown && this.options.drillDown.enabled
    },

    liveFilterPreview () {
      return this.getFilter(this.liveFilterModal.value, this.liveFilterModal.option)
    },
  },

  watch: {
    options: {
      deep: true,
      handler () {
        this.refresh()
      },
    },
  },

  mounted () {
    this.fetchChart()
    this.refreshBlock(this.refresh)
    this.createEvents()
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      findChartByID: 'chart/findByID',
      findModuleByID: 'module/findByID',
    }),

    createEvents () {
      this.$root.$on('drill-down-chart', this.drillDown)
      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { filter } = this.filter

      if (isFieldInFilter(fieldName, filter)) {
        this.refresh()
      }
    },

    refreshOnRelatedRecordsUpdate ({ moduleID, notPageID }) {
      if (this.filter.moduleID === moduleID && this.page.pageID !== notPageID) {
        this.refresh()
      }
    },

    async fetchChart (params = {}) {
      const { chartID } = this.options

      if (!chartID) {
        return
      }

      const { namespaceID } = this.namespace

      return this.findChartByID({ chartID, namespaceID, ...params }).then((chart) => {
        this.chart = chart

        if (this.isDrillDownEnabled) {
          const { moduleID, dimensions = [] } = this.chart.config.reports[0] || {}

          this.findModuleByID({ namespace: this.namespace, moduleID }).then(chartModule => {
            if (!chartModule) {
              return
            }

            const { field } = dimensions[0] || {}
            const { name, label } = chartModule.fields.find(({ name }) => name === field) || {}
            this.filter.field = { name, label }
          })
        }
      }).catch(this.toastErrorHandler(this.$t('notification:chart.loadFailed')))
    },

    reporter (r = {}) {
      if (!this.originalFilter) {
        this.originalFilter = r.filter
        this.filter = r
      }

      let filter = this.getFilter()

      if (filter) {
        // If we use ${record} or ${ownerID} and there is no record, resolve empty
        /* eslint-disable no-template-curly-in-string */
        if (!this.record && (filter.includes('${record') || filter.includes('${ownerID}'))) {
          return new Promise((resolve) => resolve([]))
        }

        filter = evaluatePrefilter(filter, {
          record: this.record,
          user: this.$auth.user || {},
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })

        this.filter.filter = filter
      }

      const { namespaceID } = this.namespace

      return this.$ComposeAPI.recordReport({ namespaceID, ...r, filter })
    },

    getFilter (liveFilter = this.liveFilterValue, option = this.liveFilterOption) {
      if (liveFilter) {
        return this.originalFilter && option ? `(${this.originalFilter}) ${option} (${liveFilter})` : liveFilter
      }

      return this.originalFilter
    },

    updateLiveFilter () {
      this.liveFilterValue = this.liveFilterModal.value
      this.liveFilterOption = this.liveFilterModal.option

      this.liveFilterModal.show = false
      this.refresh()
    },

    resetLiveFilter () {
      this.liveFilterModal.value = undefined
      this.liveFilterModal.option = 'AND'
    },

    cancelLiveFilter () {
      this.liveFilterModal.show = undefined
    },

    refresh () {
      this.fetchChart({ force: true }).then(() => {
        this.chart.config.noAnimation = true
        this.key++
      })
    },

    /**
     *
     * Based on drill down configuration, either changes the linked block on the page
     * or opens it in a modal wit the filter and dimensions from the chart and the clicked value
     */
    drillDown ({ trueName, value }) {
      const { chartID, drillDown } = this.options

      if (!drillDown.enabled) {
        return
      }

      const report = this.chart.config.reports[0] || {}
      const { yAxis = {} } = report

      // If trueName exists we use it as value, otherwise we need to look at the actual value based on if it is horizontal or vertical
      let drillDownValue = trueName
      if (!trueName) {
        drillDownValue = yAxis.horizontal ? value[1] : value[0]
      }

      // Get recordListID that is linked
      let { moduleID, dimensions, filter, field } = this.filter
      const { name, label } = field || {}

      // Construct filter
      const dimensionFilter = dimensions ? `(${dimensions} = '${drillDownValue}')` : ''

      if (drillDown.blockID) {
        // Use linked record list to display drill down data
        const { pageID = NoID } = this.page
        const { recordID = NoID } = this.record || {}
        // Construct its uniqueID to identify it
        const recordListUniqueID = [pageID, recordID, drillDown.blockID, false].map(v => v || NoID).join('-')

        this.$root.$emit(`drill-down-recordList:${recordListUniqueID}`, {
          prefilter: dimensionFilter,
          name: name || label || dimensions,
          value: drillDownValue,
        })
      } else {
        filter = filter ? `(${filter})` : ''
        const prefilter = [dimensionFilter, filter].filter(f => f).join(' AND ')

        const { title } = this.block
        const { fields = [] } = this.options.drillDown.recordListOptions || {}

        // Open in modal
        const block = new compose.PageBlockRecordList({
          title: title ? `${title} - "${drillDownValue}"` : drillDownValue,
          blockID: `drillDown-${chartID}`,
          options: {
            moduleID,
            fields,
            prefilter,
            presort: 'createdAt DESC',
            hideRecordReminderButton: true,
            hideRecordViewButton: false,
            hideConfigureFieldsButton: false,
            hideImportButton: true,
            enableRecordPageNavigation: true,
            selectable: true,
            allowExport: true,
            perPage: 14,
            showTotalCount: true,
            recordDisplayOption: 'modal',
          },
        })

        this.$root.$emit('magnify-page-block', { block })
      }
    },

    showFilterModal () {
      this.liveFilterModal.value = this.liveFilterValue
      this.liveFilterModal.option = this.liveFilterOption

      this.liveFilterModal.show = true
    },

    setDefaultValues () {
      this.chart = null
      this.filter = undefined
      this.drillDownFilter = undefined
    },

    destroyEvents () {
      this.$root.$off('drill-down-chart', this.drillDown)
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
    },
  },
}
</script>

<style lang="scss" scoped>
.chart__livefilter-btn {
  right: 0.5rem;
  top: 1.5rem;
  z-index: 1;
}
</style>

<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <chart-component
      v-if="chart"
      :key="key"
      :chart="chart"
      :record="record"
      :reporter="reporter"
      @click="drillDown"
    />
  </wrap>
</template>
<script>
import { mapActions } from 'vuex'
import base from './base'
import ChartComponent from '../Chart'
import { NoID, compose } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  i18nOptions: {
    namespaces: 'notification',
  },

  components: {
    ChartComponent,
  },
  extends: base,

  data () {
    return {
      chart: null,

      filter: undefined,

      drillDownFilter: undefined,
    }
  },

  mounted () {
    this.fetchChart()
    this.refreshBlock(this.refresh)

    this.$root.$on('drill-down-chart', this.drillDown)
  },

  methods: {
    ...mapActions({
      findChartByID: 'chart/findByID',
    }),

    async fetchChart (params = {}) {
      const { chartID } = this.options

      if (!chartID) {
        return
      }

      const { namespaceID } = this.namespace

      return this.findChartByID({ chartID, namespaceID, ...params }).then((chart) => {
        this.chart = chart
      }).catch(this.toastErrorHandler(this.$t('chart.loadFailed')))
    },

    reporter (r) {
      this.filter = r

      let filter = r.filter

      if (filter) {
        // If we use ${record} or ${ownerID} and there is no record, resolve empty
        /* eslint-disable no-template-curly-in-string */
        if (!this.record && (filter.includes('${record') || filter.includes('${ownerID}'))) {
          return new Promise((resolve) => resolve([]))
        }

        filter = evaluatePrefilter(filter, {
          record: this.record,
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })
      }

      const { namespaceID } = this.namespace

      return this.$ComposeAPI.recordReport({ namespaceID, ...r, filter })
    },

    refresh () {
      this.fetchChart({ force: true }).then(() => {
        this.chart.config.noAnimation = true
        this.key++
      })
    },

    /**
     *
     * @param {*} name
     * Based on drill down configuration, either changes the linked block on the page
     * or opens it in a modal wit the filter and dimensions from the chart and the clicked value
     */
    drillDown ({ name }) {
      const { chartID, drillDown } = this.options

      if (!drillDown.enabled) {
        return
      }

      // Get recordListID that is linked
      let { moduleID, dimensions, filter } = this.filter

      // Construct filter
      const dimensionFilter = dimensions ? `(${dimensions} = '${name}')` : ''
      filter = filter ? `(${filter})` : ''
      const prefilter = [dimensionFilter, filter].filter(f => f).join(' AND ')

      if (drillDown.blockID) {
        // Use linked record list to display drill down data
        const { pageID = NoID } = this.page
        const { recordID = NoID } = this.record || {}
        // Construct its uniqueID to identify it
        const recordListUniqueID = [pageID, recordID, drillDown.blockID].map(v => v || NoID).join('-')

        this.$root.$emit(`drill-down-recordList:${recordListUniqueID}`, prefilter)
      } else {
        const { title } = this.block

        // Open in modal
        const block = new compose.PageBlockRecordList({
          title: title ? `${title} - "${name}"` : name,
          blockID: `drillDown-${chartID}`,
          options: {
            moduleID,
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
            magnifyOption: 'modal',
            recordDisplayOption: 'modal',
          },
        })

        this.$root.$emit('magnify-page-block', { block })
      }
    },
  },

}
</script>

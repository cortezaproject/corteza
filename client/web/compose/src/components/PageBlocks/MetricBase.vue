<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <template v-else>
      <div
        v-for="(m, mi) in options.metrics"
        :key="mi"
        class="d-flex align-items-center justify-content-center overflow-hidden h-100"
      >
        <div
          v-for="(v, i) in formatResponse(m, mi)"
          :key="i"
          class="w-100 h-100 px-2 py-1"
          :class="m.drillDown.enabled ? 'pointer' : ''"
          @click="drillDown(m, mi)"
        >
          <metric-item
            :metric="m"
            :value="v"
          />
        </div>
      </div>
    </template>
  </wrap>
</template>

<script>
import base from './base'
import numeral from 'numeral'
import moment from 'moment'
import { debounce } from 'lodash'
import MetricItem from './Metric/Item'
import { NoID, compose } from '@cortezaproject/corteza-js'
import { evaluatePrefilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    MetricItem,
  },

  extends: base,

  props: {
    // Preview mode automatically generates data instead of fetching from the API
    previewMode: {
      type: Boolean,
      required: false,
      default: false,
    },
  },

  data () {
    return {
      processing: false,
      reports: [],

      abortableRequests: [],
    }
  },

  watch: {
    'record.updatedAt': {
      immediate: true,
      handler () {
        this.refresh()
      },
    },

    options: {
      deep: true,
      handler: debounce(function () {
        this.refresh()
      }, 300),
    },
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.abortRequests()
    this.destroyEvents()
    this.setDefaultValues()
  },

  created () {
    this.refreshBlock(this.refresh)
  },

  methods: {
    createEvents () {
      this.$root.$on('metric.update', this.refresh)
      this.$root.$on(`refetch-non-record-blocks:${this.page.pageID}`, this.refresh)
      this.$root.$on('drill-down-chart', this.drillDown)
      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { metrics } = this.options

      if (metrics.some(({ filter }) => isFieldInFilter(fieldName, filter))) {
        this.refresh()
      }
    },

    /**
     * Performs some post processing on the provided data
     */
    formatResponse (m, i) {
      const vals = this.reports[i]
      if (!vals) {
        return []
      }

      return vals.map(({ label, value }) => {
        if (m.numberFormat) {
          value = numeral(value).format(m.numberFormat)
        }

        if (m.dateFormat) {
          label = moment(label).format(m.dateFormat)
        }

        return {
          label,
          value,
        }
      })
    },

    /**
     * Pulls fresh data from the API
     */
    async refresh () {
      this.processing = true

      try {
        const rtr = []
        const namespaceID = this.namespace.namespaceID
        const reporter = r => {
          const { response, cancel } = this.$ComposeAPI
            .recordReportCancellable({ ...r, namespaceID })

          this.abortableRequests.push(cancel)

          return response()
        }

        for (const m of this.options.metrics) {
          if (m.moduleID) {
            // prepare a fresh metric with an evaluated prefilter
            const auxM = { ...m }
            if (auxM.filter) {
              auxM.filter = evaluatePrefilter(auxM.filter, {
                record: this.record,
                user: this.$auth.user || {},
                recordID: (this.record || {}).recordID || NoID,
                ownerID: (this.record || {}).ownedBy || NoID,
                userID: (this.$auth.user || {}).userID || NoID,
              })
            }

            if (auxM.transformFx) {
              auxM.transformFx = evaluatePrefilter(auxM.transformFx, {
                record: this.record,
                user: this.$auth.user || {},
                recordID: (this.record || {}).recordID || NoID,
                ownerID: (this.record || {}).ownedBy || NoID,
                userID: (this.$auth.user || {}).userID || NoID,
              })
            }

            const vals = await this.block.fetch({ m: auxM }, reporter)
            rtr.push(vals)
          }
        }

        this.reports = rtr
        setTimeout(() => {
          this.processing = false
        }, 300)
      } catch {
        setTimeout(() => {
          this.processing = false
        }, 300)
      }
    },
    /**
     *
     * @param {*} name
     * Based on drill down configuration, either changes the linked block on the page
     * or opens it in a modal wit the filter and dimensions from the chart and the clicked value
     */
    drillDown ({ label: name = '', filter, moduleID, drillDown }, metricIndex) {
      if (!drillDown.enabled) {
        return
      }

      if (drillDown.blockID) {
        // Use linked record list to display drill down data
        const { pageID = NoID } = this.page
        const { recordID = NoID } = this.record || {}
        // Construct its uniqueID to identify it
        const recordListUniqueID = [pageID, recordID, drillDown.blockID, false].map(v => v || NoID).join('-')
        this.$root.$emit(`drill-down-recordList:${recordListUniqueID}`, filter)
      } else {
        // Open in modal
        const metricID = `${this.block.blockID}-${name.replace(/\s+/g, '-').toLowerCase()}-${moduleID}-${metricIndex}`
        const { fields = [] } = this.options.metrics[metricIndex].drillDown.recordListOptions || {}

        const block = new compose.PageBlockRecordList({
          title: name || this.$t('metric.metricDrillDown'),
          blockID: `drillDown-${metricID}`,
          options: {
            moduleID,
            fields,
            prefilter: filter,
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

    setDefaultValues () {
      this.processing = false
      this.reports = []
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    refreshOnRelatedRecordsUpdate ({ moduleID, notPageID }) {
      if (this.page.pageID === notPageID) {
        return
      }

      const metrics = this.options.metrics

      const hasMatchingModule = metrics.some((m) => {
        return m.moduleID === moduleID
      })

      if (hasMatchingModule) {
        this.refresh()
      }
    },

    destroyEvents () {
      this.$root.$off('metric.update', this.refresh)
      this.$root.$off(`refetch-non-record-blocks:${this.page.pageID}`, this.refresh)
      this.$root.$off('drill-down-chart', this.drillDown)
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
    },
  },
}
</script>
<style scoped lang="scss">
h3 {
  line-height: 1;
}
</style>

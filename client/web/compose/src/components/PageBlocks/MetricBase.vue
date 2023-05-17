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
          <!-- <h3 :style="genStyle(m.labelStyle)">
            {{ v.label }}
          </h3> -->
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
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
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
    }
  },

  watch: {
    'record.recordID': {
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
    this.$root.$on('metric.update', this.refresh)
    this.$root.$on(`refetch-non-record-blocks:${this.page.pageID}`, this.refresh)
    this.$root.$on('drill-down-chart', this.drillDown)
  },

  beforeDestroy () {
    this.$root.$off('metric.update', this.refresh)
    this.$root.$off(`refetch-non-record-blocks:${this.page.pageID}`)
  },

  created () {
    this.refreshBlock(this.refresh)
  },

  methods: {
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
        const reporter = r => this.$ComposeAPI.recordReport({ ...r, namespaceID })

        for (const m of this.options.metrics) {
          if (m.moduleID) {
            // prepare a fresh metric with an evaluated prefilter
            const auxM = { ...m }
            if (auxM.filter) {
              auxM.filter = evaluatePrefilter(auxM.filter, {
                record: this.record,
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
        this.processing = false
      } catch {
        this.processing = false
      }
    },
    /**
     *
     * @param {*} name
     * Based on drill down configuration, either changes the linked block on the page
     * or opens it in a modal wit the filter and dimensions from the chart and the clicked value
     */
    drillDown ({ label: name, filter, moduleID, drillDown }, metricIndex) {
      if (!drillDown.enabled) {
        return
      }

      if (drillDown.blockID) {
        // Use linked record list to display drill down data
        const { pageID = NoID } = this.page
        const { recordID = NoID } = this.record || {}
        // Construct its uniqueID to identify it
        const recordListUniqueID = [pageID, recordID, drillDown.blockID].map(v => v || NoID).join('-')
        this.$root.$emit(`drill-down-recordList:${recordListUniqueID}`, filter)
      } else {
        // Open in modal
        const metricID = `${this.block.blockID}-${name.replace(/\s+/g, '-').toLowerCase()}-${moduleID}-${metricIndex}`

        const block = new compose.PageBlockRecordList({
          title: name,
          blockID: `drillDown-${metricID}`,
          options: {
            moduleID,
            prefilter: filter,
            presort: '',
            hideRecordReminderButton: true,
            hideRecordViewButton: false,
            hideConfigureFieldsButton: false,
            hideImportButton: true,
            selectable: true,
            allowExport: true,
            perPage: 14,
            showTotalCount: true,
            magnifyOption: 'modal',
          },
        })
        this.$root.$emit('magnify-page-block', { block })
      }
    },
  },
}
</script>
<style scoped lang="scss">
h3 {
  line-height: 1;
}
</style>

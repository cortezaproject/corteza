<template>
  <b-tab :title="$t('report.label')">
    <b-row>
      <b-col>
        <b-form-group
          :label="$t('report.label')"
        >
          <b-form-select
            v-model="options.reportID"
            :options="reportOptions"
            text-field="handle"
            value-field="reportID"
          />
        </b-form-group>
      </b-col>
      <b-col
        v-if="selectedReport && scenarioOptions.length > 1"
      >
        <b-form-group
          :label="$t('report.scenario.label')"
        >
          <b-form-select
            v-model="options.scenarioID"
            :options="scenarioOptions"
            text-field="label"
            value-field="scenarioID"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <b-form-group
      v-if="selectedReport"
      :label="$t('report.element.label')"
      :description="$t('report.element.description')"
    >
      <b-form-select
        v-model="options.elementID"
        :options="elementOptions"
        text-field="name"
        value-field="elementID"
      />
    </b-form-group>
  </b-tab>
</template>
<script>
import base from '../base'
import { NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  extends: base,

  data () {
    return {
      reports: [],
    }
  },

  computed: {
    reportOptions () {
      return [
        { reportID: NoID, handle: this.$t('general:label.none') },
        ...this.reports,
      ]
    },

    selectedReport () {
      const { reportID = NoID } = this.options

      if (reportID !== NoID) {
        return this.reports.find(r => r.reportID === reportID)
      }

      return undefined
    },

    scenarioOptions () {
      const { scenarios = [] } = this.selectedReport || {}

      return [
        { scenarioID: NoID, label: this.$t('general:label.none') },
        ...scenarios,
      ]
    },

    elementOptions () {
      const elements = [{ elementID: NoID, name: this.$t('general:label.none') }]

      if (this.selectedReport) {
        this.selectedReport.blocks.forEach(b => {
          elements.push({
            label: b.title || `${this.$t('general:label.block')} ${b.key}`,
            options: b.elements.map(({ elementID, name, kind }) => ({ elementID, name: name || kind })),
          })
        })
      }

      return elements
    },
  },

  watch: {
    'options.blockID' () {
      this.options.elementID = NoID
    },
  },

  created () {
    this.fetchReports()
  },

  methods: {
    fetchReports () {
      this.$SystemAPI.reportList()
        .then(({ set = [] }) => {
          this.reports = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.listFetchFailed')))
    },
  },
}
</script>

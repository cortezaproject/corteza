<template>
  <div
    class="d-flex overflow-auto p-2 w-100"
  >
    <portal
      v-if="!fetchingReport"
      to="topbar-title"
    >
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <c-input-select
        v-if="scenarioOptions.length"
        v-model="scenarios.selected"
        :options="scenarioOptions"
        :get-option-key="getOptionKey"
        :placeholder="$t('pick-scenario')"
        :disabled="processing || fetchingReport"
        size="sm"
        style="max-width: 300px; min-width: 150px;"
        @input="refreshReport()"
      />

      <b-button-group
        v-if="canUpdate"
        size="sm"
      >
        <b-button
          variant="primary"
          class="d-flex align-items-center justify-content-center"
          :to="reportBuilder"
        >
          {{ $t('report.builder') }}
          <font-awesome-icon
            class="ml-2"
            :icon="['fas', 'tools']"
          />
        </b-button>
        <b-button
          v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.edit.report'), container: '#body' }"
          variant="primary"
          class="d-flex align-items-center justify-content-center"
          style="margin-left:2px;"
          :to="reportEditor"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <div
      v-if="fetchingReport"
      class="d-flex align-items-center justify-content-center w-100 h-100"
    >
      <b-spinner />
    </div>

    <grid
      v-if="report && canRead && showReport && !fetchingReport"
      :blocks="report.blocks"
    >
      <template
        slot-scope="{ block, index }"
      >
        <block
          :index="index"
          :block="block"
          :scenario="currentSelectedScenario"
          :report-i-d="reportID"
        />
      </template>
    </grid>
  </div>
</template>

<script>
import { system } from '@cortezaproject/corteza-js'
import Grid from 'corteza-webapp-reporter/src/components/Report/Grid'
import Block from 'corteza-webapp-reporter/src/components/Report/Blocks'

export default {
  name: 'ReportView',

  i18nOptions: {
    namespaces: 'view',
  },

  components: {
    Grid,
    Block,
  },

  data () {
    return {
      processing: false,
      showReport: true,
      fetchingReport: false,

      report: undefined,
      dataframes: [],

      scenarios: {
        selected: undefined,
      },

    }
  },

  computed: {
    reportID () {
      return this.$route.params.reportID
    },

    pageTitle () {
      const title = this.report ? (this.report.meta.name || this.report.handle) : ''
      return title || this.$t('report.view')
    },

    canRead () {
      return this.report ? this.report.canReadReport : false
    },

    canUpdate () {
      return this.isNew ? this.canCreate : (this.report && this.report.canUpdateReport) || false
    },

    reportBuilder () {
      return this.report ? { name: 'report.builder', params: { reportID: this.report.reportID } } : undefined
    },

    reportEditor () {
      return this.report ? { name: 'report.edit', params: { reportID: this.report.reportID } } : undefined
    },

    reportDatasources () {
      return this.report ? this.report.sources : []
    },

    reportScenarios () {
      return this.report ? this.report.scenarios : []
    },

    scenarioOptions () {
      return this.report ? this.reportScenarios.map(({ label }) => label) : []
    },

    currentSelectedScenario () {
      return this.scenarios.selected ? this.reportScenarios.find(({ label }) => label === this.scenarios.selected) : undefined
    },
  },

  watch: {
    reportID: {
      immediate: true,
      handler (reportID) {
        this.scenarios.selected = undefined
        this.report = undefined

        if (reportID) {
          this.fetchingReport = true
          this.processing = true
          this.fetchReport(reportID)
            .catch(() => {
              this.toastErrorHandler(this.$t('notification:report.loadFailed'))
            })
            .finally(() => {
              setTimeout(() => {
                this.fetchingReport = false
                this.processing = false
              }, 400)
            })
        }
      },
    },
  },

  methods: {
    refreshReport () {
      this.showReport = false
      return setTimeout(() => {
        this.showReport = true
      }, 50)
    },

    async fetchReport (reportID) {
      this.report = undefined

      return this.$SystemAPI.reportRead({ reportID })
        .then(report => {
          this.report = new system.Report(report)

          this.report.blocks = this.report.blocks.map(({ xywh, ...p }, i) => {
            const [x, y, w, h] = xywh
            return { ...p, x, y, w, h, i }
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.fetchFailed')))
    },

    getOptionKey (scenario) {
      return scenario
    },
  },
}
</script>

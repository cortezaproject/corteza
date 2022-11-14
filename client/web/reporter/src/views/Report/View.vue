<template>
  <div
    class="d-flex overflow-auto px-2 w-100"
  >
    <portal to="topbar-title">
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <div
        v-if="scenarioOptions.length"
        class="d-inline-block mr-2"
      >
        <vue-select
          v-model="scenarios.selected"
          :options="scenarioOptions"
          :placeholder="$t('pick-scenario')"
          class="bg-white rounded"
          @input="refreshReport()"
        />
      </div>

      <b-button-group
        v-if="canUpdate"
        size="sm"
        class="mr-1"
      >
        <b-button
          variant="primary"
          class="d-flex align-items-center justify-content-center"
          :to="reportBuilder"
        >
          {{ $t('report.builder') }}
          <font-awesome-icon
            class="ml-2"
            :icon="['fas', 'cogs']"
          />
        </b-button>
        <b-button
          variant="primary"
          style="margin-left:2px;"
          :title="$t('tooltip.edit.report')"
          :to="reportEditor"
        >
          <font-awesome-icon
            :icon="['fas', 'pen']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <grid
      v-if="report && canRead && showReport"
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
import VueSelect from 'vue-select'

export default {
  name: 'ReportView',

  i18nOptions: {
    namespaces: 'view',
  },

  components: {
    Grid,
    Block,
    VueSelect,
  },

  data () {
    return {
      processing: false,
      showReport: true,

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

        if (reportID) {
          this.fetchReport(reportID)
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
      this.processing = true

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
        .finally(() => {
          this.processing = false
        })
    },
  },
}
</script>

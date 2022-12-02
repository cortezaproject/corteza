<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <display-element
      v-else-if="displayElement"
      :key="key"
      :display-element="displayElement"
      :labels="{
        previous: $t('recordList.pagination.prev'),
        next: $t('recordList.pagination.next'),
      }"
      @update="getDataframes"
    />
  </wrap>
</template>
<script>
import base from '../base'
import { system, reporter, NoID } from '@cortezaproject/corteza-js'
import DisplayElement from './DisplayElements'

export default {
  name: 'ReportBase',

  components: {
    DisplayElement,
  },

  extends: base,

  data () {
    return {
      processing: false,
      report: undefined,
      displayElement: undefined,
    }
  },

  watch: {
    options: {
      deep: true,
      immediate: true,
      handler ({ reportID = NoID }) {
        if (reportID !== NoID) {
          this.fetchReport(reportID)
        }
      },
    },
  },

  created () {
    this.refreshBlock(this.refresh, true)
  },

  methods: {
    fetchReport (reportID) {
      this.processing = true

      this.$SystemAPI.reportRead({ reportID })
        .then(report => {
          this.report = new system.Report(report)

          this.getDataframes()
        }).catch(this.toastErrorHandler(this.$t('notification:report.fetchFailed')))
        .finally(() => {
          this.processing = false
        })
    },

    getDataframes (definition = {}) {
      const { elementID } = this.options

      if (elementID) {
        const block = this.report.blocks.find(({ elements }) => {
          return elements.some(e => e.elementID === elementID)
        })

        let element = (block.elements || []).find(e => e.elementID === elementID)

        if (element && element.kind !== 'Text') {
          element = reporter.DisplayElementMaker(element)

          const scenarioDefinition = this.getScenarioDefinition(element)

          Object.entries(definition).forEach(([key, value]) => {
            definition[key] = { ...value, ...scenarioDefinition[key] }
          })

          const { dataframes: frames = [] } = element.reportDefinitions({ ...definition, ...scenarioDefinition })

          if (frames.length) {
            this.$SystemAPI.reportRun({ frames, reportID: this.options.reportID })
              .then(({ frames: dataframes = [] }) => {
                this.displayElement = {
                  ...element,
                  dataframes,
                }
              }).catch((e) => {
                this.toastErrorHandler(this.$t('notification:report.run-failed'))(e)
              })
          }
        } else {
          this.displayElement = element
        }
      }
    },

    getScenarioDefinition (element) {
      const scenarioDefinition = {}
      const { scenarioID } = this.options

      const scenario = this.report.scenarios.find(({ label }) => scenarioID === label)

      // Generate filter for each load datasource
      if (scenario && scenario.filters) {
        element.options.datasources.forEach(({ name }) => {
          scenarioDefinition[name] = {
            ref: name,
            filter: scenario.filters[name] || {},
          }
        })
      }

      return scenarioDefinition
    },

    refresh () {
      this.fetchReport(this.options.reportID)
    },
  },
}
</script>

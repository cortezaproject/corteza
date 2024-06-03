<template>
  <div class="d-flex flex-column align-items-center justify-content-center h-100 position-relative">
    <div
      v-if="processing"
      class="d-flex flex-column align-items-center justify-content-center flex-fill"
    >
      <b-spinner />
    </div>

    <c-chart
      v-if="renderer"
      :chart="renderer"
      class="flex-fill p-1"
      @click="drillDown"
    />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { chartConstructor } from 'corteza-webapp-compose/src/lib/charts'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
const { CChart } = components

export default {
  i18nOptions: {
    namespaces: 'notification',
  },

  components: {
    CChart,
  },

  props: {
    chart: {
      type: Object,
      required: true,
    },
    reporter: {
      type: Function,
      required: true,
    },
    record: {
      type: compose.Record,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      processing: false,

      valueMap: new Map(),

      renderer: undefined,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      getUserByID: 'user/findByID',
    }),
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler () {
        const { pageID = NoID } = this.$route.params
        this.$root.$on(`refetch-non-record-blocks:${pageID}`, this.requestChartUpdate)
        this.updateChart()
      },
    },
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      resolveUsers: 'user/resolveUsers',
    }),

    async updateChart () {
      this.renderer = undefined

      const [report = {}] = this.chart.config.reports

      if (!report.moduleID) {
        return
      }

      this.processing = true

      const chart = chartConstructor(this.chart)

      try {
        chart.isValid()

        const data = await chart.fetchReports({ reporter: this.reporter })

        const module = this.getModuleByID(report.moduleID)
        const fields = [
          ...module.fields,
          ...module.systemFields(),
        ]

        if (!!data.labels && Array.isArray(data.labels)) {
          // Get dimension field kind
          const [dimension = {}] = report.dimensions
          let { field } = dimension

          if (!module) throw new Error('Module not found')

          field = fields.find(({ name }) => name === field)

          if (!field) throw new Error('Dimension field not found')

          if (field.kind === 'Bool') {
            const { trueLabel, falseLabel } = field.options

            data.labels = data.labels.map(value => {
              return value === '1' ? trueLabel || this.$t('general:label.yes') : falseLabel || this.$t('general:label.no')
            })
          } else if (field.kind === 'Select') {
            data.labels = data.labels.map(value => {
              const { text } = field.options.options.find(o => o.value === value) || {}
              const label = text || value
              this.valueMap[label] = value

              return label
            })
          } else if (field.kind === 'User') {
            // Fetch and map users to labels
            await this.resolveUsers(data.labels)
            data.labels = data.labels.map(userID => {
              const label = field.formatter(this.getUserByID(userID)) || userID
              this.valueMap[label] = userID
              return label
            })
          } else if (field.kind === 'Record') {
            // Fetch and map records and their values to labels
            const { namespaceID } = this.chart || {}
            const recordModule = this.getModuleByID(field.options.moduleID)
            if (recordModule && data.labels) {
              const isValidRecordID = (recordID) => recordID !== dimension.default && recordID !== 'undefined'
              await Promise.all(data.labels.map(recordID => {
                if (isValidRecordID(recordID)) {
                  return this.$ComposeAPI.recordRead({ namespaceID, moduleID: recordModule.moduleID, recordID }).then(record => {
                    record = new compose.Record(recordModule, record)

                    if (field.options.recordLabelField) {
                      // Get actual field
                      const relatedField = recordModule.fields.find(({ name }) => name === field.options.labelField)

                      return this.$ComposeAPI.recordRead({ namespaceID, moduleID: relatedField.options.moduleID, recordID: record.values[field.options.labelField] }).then(labelRecord => {
                        record.values[field.options.labelField] = (labelRecord.values.find(({ name }) => name === this.field.options.recordLabelField) || {}).value
                        return record
                      })
                    } else {
                      return record
                    }
                  })
                } else {
                  const record = { values: {} }
                  record.values[field.options.labelField] = recordID
                  return record
                }
              })).then(records => {
                data.labels = records.map(record => {
                  const value = field.options.labelField ? record.values[field.options.labelField] : record.recordID
                  const label = Array.isArray(value) ? value.join(', ') : value
                  this.valueMap[label] = record.recordID

                  return value
                })
              })
            }
          }
        }

        data.datasets = data.datasets.map((dataset = {}, i) => {
          const { label } = dataset

          if (label === 'count') {
            dataset.label = this.$t('chart:general.label.count')
          } else {
            const field = fields.find(({ name }) => name === label) || {}
            dataset.label = field.label || label
          }

          return dataset
        })

        data.labels = data.labels.map(l => l === 'undefined' ? this.$t('chart:undefined') : l)
        data.customColorSchemes = this.$Settings.get('ui.charts.colorSchemes', [])
        data.themeVariables = this.getThemeVariables()

        this.renderer = chart.makeOptions(data)
      } catch (e) {
        this.processing = false
        this.toastErrorHandler(this.$t('chart.optionsBuildFailed'))(e)
      }

      setTimeout(() => {
        this.processing = false
        this.$emit('updated')
      }, 300)
    },

    drillDown (e) {
      e.trueName = this.valueMap[e.name] || e.name

      return this.$emit('drill-down', e)
    },

    getThemeVariables () {
      const getCssVariable = (variableName) => {
        return getComputedStyle(document.documentElement).getPropertyValue(variableName).trim()
      }

      // Turn below into an object with key value pairs
      return ['white', 'black', 'primary', 'secondary', 'success', 'warning', 'danger', 'light', 'extra-light', 'dark', 'font-regular'].reduce((acc, variable) => {
        acc[variable] = getCssVariable(`--${variable}`)
        return acc
      }, {})
    },

    requestChartUpdate ({ name, handle } = {}) {
      if (!name && !handle) {
        this.updateChart()
      }

      if (name && this.chart && this.chart.name === name) {
        this.updateChart()
      }

      if (handle && this.chart && this.chart.handle === handle) {
        this.updateChart()
      }
    },

    error (msg) {
      /* eslint-disable no-console */
      console.error(msg)
    },

    setDefaultValues () {
      this.processing = false
      this.renderer = undefined
      this.valueMap = {}
    },

    destroyEvents () {
      const { pageID = NoID } = this.$route.params
      this.$root.$off(`refetch-non-record-blocks:${pageID}`)
    },
  },
}
</script>

<template>
  <div class="h-100">
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <div
      class="justify-content-center position-relative h-100 p-1 overflow-hidden"
      :class="{ 'd-none': processing, 'd-flex': !processing }"
    >
      <canvas
        ref="chartCanvas"
        class="mh-100 w-auto"
      />
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { chartConstructor } from 'corteza-webapp-compose/src/lib/charts'
import ChartJS from 'chart.js'
import Funnel from 'chartjs-plugin-funnel'
import Gauge from 'chartjs-gauge'
import csc from 'chartjs-plugin-colorschemes'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'notification',
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

      renderer: null,
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
        this.$nextTick(() => {
          this.updateChart()
        })
      },
    },
  },

  mounted () {
    this.$root.$on('chart.update', this.requestChartUpdate)
  },

  beforeDestroy () {
    if (this.renderer) {
      this.renderer.destroy()
    }
  },

  methods: {
    async updateChart () {
      const [report = {}] = this.chart.config.reports

      if (!report.moduleID) {
        return
      }

      this.processing = true

      const chart = chartConstructor(this.chart)

      try {
        chart.isValid()

        const data = await chart.fetchReports({ reporter: this.reporter })

        if (!!data.labels && Array.isArray(data.labels)) {
          // Get dimension field kind
          const [dimension = {}] = report.dimensions
          let { field, meta = {} } = dimension
          const module = this.getModuleByID(report.moduleID)

          if (module) {
            field = [
              ...module.fields,
              ...module.systemFields(),
            ].find(({ name }) => name === field)

            if (meta.fields) {
              data.labels = data.labels.map(value => {
                const { text } = field.options.options.find(o => o.value === value) || {}
                return text || value
              })
            }
          }

          if (field && ['User', 'Record'].includes(field.kind)) {
            if (field.kind === 'User') {
              // Fetch and map users to labels
              await this.$store.dispatch('user/fetchUsers', data.labels)
              data.labels = data.labels.map(label => {
                return field.formatter(this.getUserByID(label)) || label
              })
            } else {
              // Fetch and map records and their values to labels
              const { namespaceID } = this.chart || {}
              const recordModule = this.getModuleByID(field.options.moduleID)
              if (recordModule && data.labels) {
                await Promise.all(data.labels.map(recordID => {
                  if (recordID && recordID !== 'undefined') {
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
                    return Array.isArray(value) ? value.join(', ') : value
                  })
                })
              }
            }
          }

          data.labels = data.labels.map(l => l === 'undefined' ? this.$t('chart:undefined') : l)
          const options = chart.makeOptions(data)
          const plugins = chart.plugins()
          if (!options) {
            this.toastWarning(this.$t('chart.optionsBuildFailed'))
          }
          const type = chart.baseChartType(data.datasets)

          const canvas = this.$refs.chartCanvas || undefined

          if (canvas) {
            const newRenderer = () => new ChartJS(this.$refs.chartCanvas.getContext('2d'), { options, plugins: [...plugins, Funnel, Gauge, csc], data, type })
            if (this.renderer) {
              this.renderer.destroy()
            }
            this.renderer = newRenderer()
          }
        } else {
          data.labels = []
        }
      } catch (e) {
        this.processing = false
        this.toastErrorHandler(this.$t('chart.optionsBuildFailed'))(e)
      }
      this.processing = false
      this.$emit('updated')
    },

    requestChartUpdate ({ name, handle } = {}) {
      if (!name && !handle) {
        this.$nextTick(() => this.updateChart())
      }

      if (name && this.chart && this.chart.name === name) {
        this.$nextTick(() => this.updateChart())
      }

      if (handle && this.chart && this.chart.handle === handle) {
        this.$nextTick(() => this.updateChart())
      }
    },

    error (msg) {
      /* eslint-disable no-console */
      console.error(msg)
    },
  },
}
</script>

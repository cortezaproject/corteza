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
    />
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
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
    const { pageID = NoID } = this.$route.params
    this.$root.$off(`refetch-non-record-blocks:${pageID}`)
  },

  methods: {
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

          this.renderer = chart.makeOptions(data)
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
  },
}
</script>

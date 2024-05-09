import { system } from '@cortezaproject/corteza-js'

export default {
  data () {
    return {
      processing: false,
      processingSave: false,
      processingDelete: false,
      processingClone: false,
      report: undefined,
    }
  },

  methods: {
    async fetchReport (reportID) {
      this.processing = true

      return this.$SystemAPI.reportRead({ reportID })
        .then(report => {
          this.report = new system.Report(report)
          this.initialReportState = this.report.clone()
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.fetchFailed')))
        .finally(() => {
          this.processing = false
        })
    },

    async handleSave () {
      this.processing = true
      this.processingSave = true

      const { blocks } = this.report

      // Remove dataframes from elements before saving report
      const report = {
        ...this.report,
        blocks: blocks.map(block => {
          block.elements = block.elements.map(element => {
            delete element.dataframes
            return element
          })
          return block
        }),
      }

      // If new then create otherwise update
      if (this.isNew) {
        return this.$SystemAPI.reportCreate(report)
          .then(report => {
            this.report = new system.Report(report)
            this.initialReportState = this.report.clone()
            this.detectStateChange = false
            this.toastSuccess(this.$t('notification:report.created'))
            this.$router.push({ name: 'report.edit', params: { reportID: report.reportID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:report.createFailed')))
          .finally(() => {
            this.processing = false
            this.processingSave = false
          })
      } else {
        return this.$SystemAPI.reportUpdate(report)
          .then(report => {
            this.report = new system.Report(report)
            this.initialReportState = this.report.clone()
            this.detectStateChange = false
            this.toastSuccess(this.$t('notification:report.updated'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:report.updateFailed')))
          .finally(() => {
            this.processing = false
            this.processingSave = false
          })
      }
    },

    handleDelete () {
      this.processing = true
      this.processingDelete = true

      return this.$SystemAPI.reportDelete(this.report)
        .then(() => {
          this.report.deletedAt = new Date()
          this.toastSuccess(this.$t('notification:report.delete'))
          this.$router.push({ name: 'report.list' })
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.deleteFailed')))
        .finally(() => {
          this.processing = false
          this.processingDelete = false
        })
    },

    handleClone (report) {
      let { handle, meta, sources, blocks, scenarios, labels } = report
      handle = ''
      meta.name = `${meta.name} (${this.$t('general:cloneSuffix')})`
      this.processingClone = true
      return this.$SystemAPI.reportCreate({ handle, meta, sources, blocks, scenarios, labels })
        .then(report => {
          report = new system.Report(report)
          this.toastSuccess(this.$t('notification:report.created'))
          return report
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.createFailed')))
    },
  },
}

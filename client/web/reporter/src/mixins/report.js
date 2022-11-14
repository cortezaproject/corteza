import { system } from '@cortezaproject/corteza-js'

export default {
  data () {
    return {
      processing: false,
      report: undefined,
    }
  },

  methods: {
    async fetchReport (reportID) {
      this.processing = true

      return this.$SystemAPI.reportRead({ reportID })
        .then(report => {
          this.report = new system.Report(report)
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.fetchFailed')))
        .finally(() => {
          this.processing = false
        })
    },

    async handleSave () {
      this.processing = true

      // If new then create otherwise update
      if (this.isNew) {
        return this.$SystemAPI.reportCreate(this.report)
          .then(report => {
            this.report = new system.Report(report)
            this.toastSuccess(this.$t('notification:report.created'))
            this.$router.push({ name: 'report.edit', params: { reportID: report.reportID } })
          })
          .catch(this.toastErrorHandler(this.$t('notification:report.createFailed')))
          .finally(() => {
            this.processing = false
          })
      } else {
        return this.$SystemAPI.reportUpdate(this.report)
          .then(report => {
            this.report = new system.Report(report)
            this.toastSuccess(this.$t('notification:report.updated'))
          })
          .catch(this.toastErrorHandler(this.$t('notification:report.updateFailed')))
          .finally(() => {
            this.processing = false
          })
      }
    },

    async handleDelete () {
      this.processing = true

      return this.$SystemAPI.reportDelete(this.report)
        .then(() => {
          this.toastSuccess(this.$t('notification:report.delete'))
          this.$router.push({ name: 'report.list' })
        })
        .catch(this.toastErrorHandler(this.$t('notification:report.deleteFailed')))
        .finally(() => {
          this.processing = false
        })
    },
  },
}

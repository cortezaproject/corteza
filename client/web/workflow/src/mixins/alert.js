const success = { title: 'Success', variant: 'success', countdown: 3000 }
const danger = { title: 'Error', variant: 'danger', countdown: 10000 }
const warning = { title: 'Warning', variant: 'warning', countdown: 5000 }
const info = { title: 'Info', variant: 'info', countdown: 5000 }

export default {
  methods: {
    raiseSuccessAlert (message, title) {
      this.raiseAlert({
        ...success,
        message,
        title: title || success.title,
      })
    },

    raiseDangerAlert (message, title) {
      this.raiseAlert({
        ...danger,
        message,
        title: title || danger.title,
      })
    },

    raiseWarningAlert (message, title) {
      this.raiseAlert({
        ...warning,
        message,
        title: title || warning.title,
      })
    },

    raiseInfoAlert (message, title) {
      this.raiseAlert({
        ...info,
        message,
        title: title || info.title,
      })
    },

    raiseAlert (alert = {}) {
      this.$root.$emit('alert', alert)
    },

    defaultErrorHandler (title) {
      return (err = {}) => {
        /* eslint-disable no-console */
        err = err.message || title
        console.error(err)
        this.raiseWarningAlert(err, title)
      }
    },

    handleAlert (handler) {
      this.$root.$on('alert', handler)
    },
  },
}

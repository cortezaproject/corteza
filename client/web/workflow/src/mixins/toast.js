export default {
  methods: {
    toastSuccess (message, title = this.$t('notification:general.success')) {
      this.toast(message, { title, variant: 'success' })
    },

    toastWarning (message, title = this.$t('notification:general.warning')) {
      this.toast(message, { title, variant: 'warning' })
    },

    toastDanger (message, title = this.$t('notification:general.error')) {
      this.toast(message, { title, variant: 'danger' })
    },

    toastInfo (message, title = this.$t('notification:general.info')) {
      this.toast(message, { title, variant: 'info' })
    },

    toast (msg, opt = { variant: 'success' }) {
      this.$root.$bvToast.toast(msg, opt)
    },

    toastErrorHandler (opt) {
      if (typeof opt === 'string') {
        opt = { prefix: opt }
      }

      const { prefix, title } = opt

      return (err = {}) => {
        /* eslint-disable no-console */
        console.error(err)
        const msg = err.message ? (prefix + ': ' + err.message) : prefix
        this.toastDanger(msg, title)
      }
    },
  },
}

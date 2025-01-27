export default {
  methods: {
    toastSuccess (message, title = undefined) {
      if (title === undefined) {
        title = this.$t('notification:general.success')
      }

      this.toast(message, { title, variant: 'success' })
    },

    toastWarning (message, title = undefined) {
      if (title === undefined) {
        title = this.$t('notification:general.warning')
      }

      this.toast(message, { title, variant: 'warning' })
    },

    toastInfo (message, title = undefined) {
      if (title === undefined) {
        title = this.$t('notification:general.info')
      }

      this.toast(message, { title, variant: 'info' })
    },

    toastDanger (message, title = undefined) {
      if (title === undefined) {
        title = this.$t('notification:general.error')
      }

      this.toast(message, { title, variant: 'danger' })
    },

    toast (msg, opt = { variant: 'success' }) {
      this.$root.$bvToast.toast(msg, opt)
    },

    getToastMessage (err) {
      if (err.message && err.message.startsWith('notification')) {
        return this.$t(`notification:${err.message.substring('notification.'.length)}`)
      }

      return err.message
    },

    toastErrorHandler (opt) {
      if (typeof opt === 'string') {
        opt = { prefix: opt }
      }

      const { prefix, title } = opt

      return (err = {}) => {
        err.message = this.getToastMessage(err)

        // all other messages should be shown as they are
        const msg = err.message ? `${prefix}: ${err.message}` : prefix
        this.toastDanger(msg, title)

        return err.message
      }
    },
  },
}

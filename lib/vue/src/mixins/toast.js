export default {
  methods: {
    toastSuccess (message, title) {
      title = title || this.$t('notification:general.success')
      this.toast(message, { title, variant: 'success' })
    },

    toastWarning (message, title) {
      title = title || this.$t('notification:general.warning')
      this.toast(message, { title, variant: 'warning' })
    },

    toastInfo (message, title) {
      title = title || this.$t('notification:general.info')
      this.toast(message, { title, variant: 'info' })
    },

    toastDanger (message, title) {
      title = title || this.$t('notification:general.error')
      this.toast(message, { title, variant: 'danger' })
    },

    toast (msg, opt = { variant: 'success' }) {
      const toaster = this.$Settings?.get('ui.notifications.toastPosition', 'b-toaster-bottom-right')
      this.$root.$bvToast.toast(msg, { toaster, ...opt })
    },

    getToastMessage (err) {
      if (err.message && err.message.startsWith('notification')) {
        return this.$t(`notification:${err.message.substring('notification.'.length)}`)
      }

      return err.message
    },

    toastErrorHandler (opt) {
      if (typeof opt === 'string') {
        opt = { title: opt }
      }

      const { prefix, title } = opt

      return (err = {}) => {
        let toastMsg = ''
        let toastTitle = title

        err.message = this.getToastMessage(err)

        if (err.message) {
          toastMsg = err.message
        } else {
          toastMsg = title
          toastTitle = ''
        }

        if (prefix) {
          toastMsg = `${prefix}: ${toastMsg}`
        }

        toastMsg = toastTitle ? toastMsg.charAt(0).toUpperCase() + toastMsg.slice(1) : toastMsg
        toastTitle = toastTitle ? toastTitle.charAt(0).toUpperCase() + toastTitle.slice(1) : toastTitle

        this.toastDanger(toastMsg, toastTitle)

        return err.message
      }
    },
  },
}

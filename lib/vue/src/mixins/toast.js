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
      this.$root.$bvToast.toast(msg, opt)
    },

    getToastMessage (err) {
      if (err.message && err.message.startsWith('notification')) {
        return this.$t(`notification:${err.message.substring('notification.'.length)}`)
      }

      return err.message
    },

    // extractWorkflowErrorMeta looks for the namespaced "workflow.error.*"
    // meta keys attached by the backend workflow error step. It probes the
    // handful of shapes the error can arrive in (direct reject, axios wrap,
    // nested response.data.error) and returns a normalised object, or null
    // when the error is not a workflow error step error.
    //
    // This is purely additive: when it returns null, toastErrorHandler
    // falls through to the existing plain-toast behaviour.
    extractWorkflowErrorMeta (err) {
      if (!err || typeof err !== 'object') return null

      // candidate meta bags, in priority order
      const candidates = [
        err.meta,
        err.response && err.response.data && err.response.data.error && err.response.data.error.meta,
        err.data && err.data.error && err.data.error.meta,
        err.error && err.error.meta,
      ]

      for (const meta of candidates) {
        if (!meta || typeof meta !== 'object') continue
        const title = meta['workflow.error.title']
        const severity = meta['workflow.error.severity']
        if (title || severity) {
          return {
            title: title || '',
            severity: severity || 'error',
          }
        }
      }

      return null
    },

    toastErrorHandler (opt) {
      if (typeof opt === 'string') {
        opt = { title: opt }
      }

      const { prefix, title } = opt

      return (err = {}) => {
        // Rich workflow error step path. Guarded by a try/catch so that
        // any unexpected toast library or context issue silently degrades
        // to the legacy plain-toast rendering below.
        try {
          const wfErr = this.extractWorkflowErrorMeta(err)
          if (wfErr) {
            const variant = wfErr.severity === 'warning'
              ? 'warning'
              : wfErr.severity === 'info'
                ? 'info'
                : 'danger'

            // The toast body is always the flat error message (which
            // is equal to the author-configured `message` argument).
            //
            // NOTE: message/title originate from workflow-author-controlled
            // meta fields. BootstrapVue's $bvToast.toast(content, opts)
            // renders both content and title as plain text (no v-html),
            // which is what we want here — do not switch to an HTML
            // variant without sanitising, as that would expose the
            // workflow editor as a stored-XSS surface.
            const body = err.message || ''
            const toastTitle = wfErr.title || title || ''

            this.toast(body, { title: toastTitle, variant, solid: true })
            return err.message
          }
        } catch (_) {
          // fall through to legacy rendering
        }

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

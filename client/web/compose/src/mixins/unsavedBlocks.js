import { mapGetters, mapActions } from 'vuex'

/**
 * Warns before leaving a page while a block (record list in inline-edit
 * mode) still holds unsaved changes; covers in-app navigation through
 * confirmUnsavedBlocks() and the browser's reload/close via beforeunload
 */
export default {
  computed: {
    ...mapGetters({
      hasUnsavedBlocks: 'ui/hasUnsavedBlocks',
    }),
  },

  watch: {
    hasUnsavedBlocks: {
      immediate: true,
      handler (unsaved) {
        if (unsaved) {
          window.addEventListener('beforeunload', this.warnBeforeUnload)
        } else {
          window.removeEventListener('beforeunload', this.warnBeforeUnload)
        }
      },
    },
  },

  beforeDestroy () {
    window.removeEventListener('beforeunload', this.warnBeforeUnload)
    this.clearUnsavedBlocks()
  },

  methods: {
    ...mapActions({
      clearUnsavedBlocks: 'ui/clearUnsavedBlocks',
    }),

    warnBeforeUnload (e) {
      e.preventDefault()
      // legacy browsers need a value to show the prompt
      e.returnValue = ''
    },

    /**
     * Returns true when leaving is fine: nothing unsaved, or the user agreed
     */
    confirmUnsavedBlocks () {
      if (!this.hasUnsavedBlocks) {
        return true
      }

      if (!window.confirm(this.$t('general:record.unsavedChanges'))) {
        return false
      }

      this.clearUnsavedBlocks()
      return true
    },
  },
}

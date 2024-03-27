import { compose } from '@cortezaproject/corteza-js'
import { mapGetters } from 'vuex'

export default {
  props: {
    namespace: {
      // via router-view
      type: compose.Namespace,
      required: true,
    },
    page: {
      // via route-view
      type: compose.Page,
      required: true,
    },
    // We're using recordID to check if we need to display router-view or grid component
    recordID: {
      type: String,
      default: '',
    },
  },
  computed: {
    ...mapGetters({
      getPageLayouts: 'pageLayout/getByPageID',
    }),

    isRecordPage () {
      return this.recordID || this.$route.name === 'page.record.create'
    },

    expressionVariables () {
      return {
        user: this.$auth.user,
        record: this.record ? this.record.serialize() : {},
        screen: {
          width: window.innerWidth,
          height: window.innerHeight,
          userAgent: navigator.userAgent,
          breakpoint: this.getBreakpoint(), // This is from a global mixin uiHelpers
        },
        oldLayout: this.layout,
        layout: undefined,
        ...(this.isRecordPage && {
          isView: !this.inEditing && !this.inCreating,
          isCreate: this.inCreating,
          isEdit: this.inEditing && !this.inCreating,
        }),
      }
    },
  },
  methods: {
    evaluateBlocksExpressions (variables = {}) {
      const expressions = {}
      variables = {
        ...this.expressionVariables,
        ...variables,
      }

      this.layout.blocks.forEach(block => {
        const { visibility } = block.meta
        if (!(visibility || {}).expression) return

        expressions[block.blockID] = visibility.expression
      })

      return this.$SystemAPI.expressionEvaluate({ variables, expressions }).catch(e => {
        this.toastErrorHandler(this.$t('notification:evaluate.failed'))(e)
        Object.keys(expressions).forEach(key => (expressions[key] = false))

        return expressions
      })
    },
  },
}

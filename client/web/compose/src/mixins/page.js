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

  data () {
    return {
      layouts: [],
      layout: undefined,
      blocks: undefined,
    }
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
          isView: !this.edit && !this.isNew,
          isCreate: this.isNew,
          isEdit: this.edit && !this.isNew,
        }),
      }
    },
  },

  mounted () {
    this.createEvents()
  },

  methods: {
    evaluateLayoutExpressions (variables = {}) {
      const expressions = {}
      variables = {
        ...this.expressionVariables,
        ...variables,
      }

      this.layouts.forEach(layout => {
        const { config = {} } = layout
        if (!config.visibility.expression) return

        variables.layout = layout

        expressions[layout.pageLayoutID] = config.visibility.expression
      })

      return this.$SystemAPI.expressionEvaluate({ variables, expressions }).catch(e => {
        this.toastErrorHandler(this.$t('notification:evaluate.failed'))(e)
        Object.keys(expressions).forEach(key => (expressions[key] = false))

        return expressions
      })
    },

    async determineLayout (pageLayoutID, variables = {}, redirectOnFail = true) {
      // Clear stored records so they can be refetched with latest values
      this.clearRecordSet()

      if (this.isRecordPage) {
        this.resetErrors()
      }

      let expressions = {}

      // Only evaluate if one of the layouts has an expressions variable
      if (this.layouts.some(({ config = {} }) => config.visibility.expression)) {
        expressions = await this.evaluateLayoutExpressions(variables)
      }

      // Check layouts for expressions/roles and find the first one that fits
      const matchedLayout = this.layouts.find(l => {
        if (pageLayoutID && l.pageLayoutID !== pageLayoutID) return

        const { expression, roles = [] } = l.config.visibility

        if (expression && !expressions[l.pageLayoutID]) return false

        if (!roles.length) return true

        return this.$auth.user.roles.some(roleID => roles.includes(roleID))
      })

      this.processing = false

      if (!matchedLayout) {
        this.toastWarning(this.$t('notification:page.page-layout.notFound.view'))

        if (redirectOnFail) {
          this.$router.go(-1)
        }

        return this.$router.go(-1)
      }

      if (this.isRecordPage) {
        this.inEditing = this.edit
      }

      if (this.layout && matchedLayout.pageLayoutID === this.layout.pageLayoutID) {
        return
      }

      this.layout = matchedLayout

      if (this.isRecordPage) {
        this.handleRecordButtons()
      } else {
        const { handle, meta = {} } = this.layout || {}
        const title = meta.title || this.page.title
        this.pageTitle = title || handle || this.$t('navigation:noPageTitle')
        document.title = [title, this.namespace.name, this.$t('general:label.app-name.public')].filter(v => v).join(' | ')
      }

      await this.updateBlocks(variables)
    },

    async updateBlocks (variables = {}) {
      const tempBlocks = []
      const { blocks = [] } = this.layout || {}

      let blocksExpressions = {}

      if (blocks.some(({ meta = {} }) => (meta.visibility || {}).expression)) {
        blocksExpressions = await this.evaluateBlocksExpressions(variables)
      }

      blocks.forEach(({ blockID, xywh, meta }) => {
        const block = this.page.blocks.find(b => b.blockID === blockID)
        const { roles = [], expression = '' } = meta.visibility || {}

        if (block && (!expression || blocksExpressions[blockID])) {
          block.xywh = xywh

          if (!roles.length || this.$auth.user.roles.some(roleID => roles.includes(roleID))) {
            tempBlocks.push(block)
          }
        }
      })

      this.blocks = tempBlocks
    },

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

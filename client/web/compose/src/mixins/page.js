import { compose } from '@cortezaproject/corteza-js'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import { mapActions, mapGetters } from 'vuex'

export default {
  components: {
    PageTranslator,
  },

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

      tempRecord: undefined,
    }
  },

  computed: {
    ...mapGetters({
      getPageLayouts: 'pageLayout/getByPageID',
    }),

    isRecordPage () {
      return this.recordID || this.$route.name === 'page.record.create'
    },

    trPage: {
      get () {
        return this.page.clone()
      },

      set (v) {
        this.updatePageSet(v)
      },
    },

    pageEditor () {
      return { name: 'admin.pages.edit', params: { pageID: this.page.pageID } }
    },

    pageBuilder () {
      const { pageLayoutID } = this.layout || {}
      return { name: 'admin.pages.builder', params: { pageID: this.page.pageID }, query: { layoutID: pageLayoutID } }
    },

    moduleEditor () {
      if (!this.module) return undefined

      return { name: 'admin.modules.edit', params: { moduleID: this.module.moduleID } }
    },
  },

  methods: {
    ...mapActions({
      clearRecordSet: 'record/clearSet',
      updatePageSet: 'page/updateSet',
      popPreviousPages: 'ui/popPreviousPages',
      popModalPreviousPage: 'ui/popModalPreviousPage',
    }),

    expressionVariables () {
      const record = this.tempRecord || this.record

      return {
        user: this.$auth.user,
        record: record ? record.serialize() : {},
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

    async evaluateLayoutExpressions () {
      const expressions = {}
      const variables = this.expressionVariables()

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

    async determineLayout ({ pageLayoutID, redirectOnFail = true } = {}) {
      if (this.isRecordPage) {
        this.resetErrors()
      }

      let expressions = {}

      // Only evaluate if one of the layouts has an expressions variable
      if (this.layouts.some(({ config = {} }) => config.visibility.expression)) {
        expressions = await this.evaluateLayoutExpressions()
      }

      // Check layouts for expressions/roles and find the first one that fits
      const matchedLayout = this.layouts.find(l => {
        if (pageLayoutID && l.pageLayoutID !== pageLayoutID) return false

        const { expression, roles = [] } = l.config.visibility

        if (expression && !expressions[l.pageLayoutID]) return false

        if (!roles.length) return true

        return this.$auth.user.roles.some(roleID => roles.includes(roleID))
      })

      if (!matchedLayout) {
        this.toastWarning(this.$t('notification:page.page-layout.notFound.view'))
        if (redirectOnFail) {
          this.$router.go(-1)
        }

        return
      }

      if (this.isRecordPage) {
        this.inEditing = this.edit
      }

      if (this.layout && matchedLayout.pageLayoutID === this.layout.pageLayoutID) {
        this.evaluateBlocks()
        return
      }

      this.layout = matchedLayout

      if (this.isRecordPage) {
        this.handleRecordButtons()
      } else {
        const { handle, meta = {} } = this.layout || {}

        this.pageTitle = (meta.title || this.page.title) || handle || this.$t('navigation:noPageTitle')
        document.title = this.pageTitle
      }

      return this.prepareBlocks()
    },

    async prepareBlocks () {
      this.blocks = undefined

      const tempBlocks = []
      const { blocks = [] } = this.layout || {}

      blocks.forEach(({ blockID, xywh }) => {
        const block = this.page.blocks.find(b => b.blockID === blockID)

        if (block) {
          block.xywh = xywh
          tempBlocks.push(block)
        }
      })

      return this.evaluateBlocks(tempBlocks)
    },

    async evaluateBlocks (blocks = this.page.blocks) {
      const layoutBlocks = this.layout.blocks
      let layoutBlocksExpressions = {}

      // Only evaluate expressions if any blocks have visibility expressions
      if (layoutBlocks.some(({ meta = {} }) => (meta.visibility || {}).expression)) {
        layoutBlocksExpressions = await this.evaluateBlocksExpressions()
      }

      blocks.forEach(block => {
        const { blockID, meta = {} } = block
        const visibility = meta.visibility || {}
        const { roles = [] } = visibility

        // Determine if block should be shown based on expression and roles
        const validExpression = !visibility.expression || layoutBlocksExpressions[blockID]
        const validRole = !roles.length || this.$auth.user.roles.some(roleID => roles.includes(roleID))
        const showBlock = block && validExpression && validRole

        // Update invisible status based on visibility evaluation
        block.meta.invisible = !showBlock
      })

      return blocks
    },

    async evaluateBlocksExpressions () {
      const expressions = {}
      const variables = this.expressionVariables()

      this.layout.blocks.forEach(block => {
        const { visibility } = block.meta
        if (!(visibility || {}).expression) return

        expressions[block.blockID] = visibility.expression
      })

      return this.$SystemAPI.expressionEvaluate({ variables, expressions }).catch(e => {
        this.toastErrorHandler(this.$t('notification:evaluate.failed'))(e)
        Object.keys(expressions).forEach(key => (expressions[`${key}`] = false))

        return expressions
      })
    },
  },
}

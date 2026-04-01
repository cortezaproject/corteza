import { compose } from '@cortezaproject/corteza-js'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'
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
      layoutRequiredFields: 'ui/layoutRequiredFields',
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
      const tabbedIDs = new Set()

      blocks.forEach(({ blockID, xywh }) => {
        const block = this.page.blocks.find(b => b.blockID === blockID)

        if (block) {
          block.xywh = xywh
          tempBlocks.push(block)

          if (block.kind === 'Tabs') {
            const { tabs = [] } = block.options
            tabs.forEach(t => {
              if (t.blockID && !blocks.some(b => b.blockID === t.blockID)) {
                tabbedIDs.add(t.blockID)
              }
            })
          }
        }
      })

      // Include blocks that are only in tabs
      this.page.blocks.forEach(block => {
        if (tabbedIDs.has(block.blockID)) {
          tempBlocks.push(block)
        }
      })

      return this.evaluateBlocks(tempBlocks)
    },

    async evaluateBlocks (blocks = this.page.blocks) {
      let layoutBlocksExpressions = {}

      // Only evaluate expressions if any blocks have visibility expressions
      if (blocks.some(({ meta = {} }) => (meta.visibility || {}).expression)) {
        layoutBlocksExpressions = await this.evaluateBlocksExpressions(blocks)
      }

      blocks.forEach(block => {
        const { meta = {} } = block
        const blockID = fetchID(block)
        const visibility = meta.visibility || {}
        const { roles = [] } = visibility

        // Determine if block should be shown based on expression and roles
        const validExpression = !visibility.expression || layoutBlocksExpressions[blockID]
        const validRole = !roles.length || this.$auth.user.roles.some(roleID => roles.includes(roleID))
        const showBlock = block && validExpression && validRole

        // Update invisible status based on visibility evaluation
        if (!block.meta) {
          this.$set(block, 'meta', {})
        }
        this.$set(block.meta, 'invisible', !showBlock)
      })

      // Propagate invisibility to Tabs blocks
      blocks.forEach(block => {
        if (block.kind === 'Tabs' && !block.meta.invisible) {
          const { tabs = [] } = block.options
          const hasVisibleTab = tabs.some(t => {
            const b = blocks.find(b => fetchID(b) === t.blockID)
            return b ? !b.meta.invisible : !!t.title
          })

          if (!hasVisibleTab) {
            this.$set(block.meta, 'invisible', true)
          }
        }
      })

      // Evaluate layout required fields if this is a record page
      if (this.isRecordPage && this.evaluateLayoutRequiredFields) {
        await this.evaluateLayoutRequiredFields()
      }

      return blocks
    },

    async evaluateBlocksExpressions (blocks = this.page.blocks) {
      const expressions = {}
      const variables = this.expressionVariables()

      blocks.forEach(block => {
        const { visibility } = block.meta
        if (!(visibility || {}).expression) return

        expressions[fetchID(block)] = visibility.expression
      })

      return this.$SystemAPI.expressionEvaluate({ variables, expressions }).catch(e => {
        this.toastErrorHandler(this.$t('notification:evaluate.failed'))(e)
        Object.keys(expressions).forEach(key => (expressions[`${key}`] = false))

        return expressions
      })
    },

    async evaluateLayoutRequiredFields () {
      if (!this.layout) {
        this.$store.dispatch('ui/clearLayoutRequiredFields')
        return
      }

      const { config = {} } = this.layout || {}
      const { validation = {} } = config || {}
      const { requiredFields = [] } = validation || {}

      if (requiredFields.length === 0) {
        this.$store.dispatch('ui/clearLayoutRequiredFields')
        return
      }

      await new Promise(resolve => setTimeout(resolve, 300))

      const { expressions, variables } = this.prepareLayoutRequiredFieldsData()

      if (Object.keys(expressions).length === 0) {
        // No conditions to evaluate, mark all as required
        const fields = requiredFields
          .filter(rf => !rf.condition || rf.condition.trim() === '')
          .map(rf => rf.field)
        this.$store.dispatch('ui/setLayoutRequiredFields', fields)
        return
      }

      return this.$SystemAPI
        .expressionEvaluate({ variables, expressions })
        .then(res => {
          // Update required fields based on evaluation results
          const fields = []

          requiredFields.forEach(({ field, condition }) => {
            // If no condition, always required
            if (!condition || condition.trim() === '') {
              fields.push(field)
            } else if (res[field]) {
              // Condition evaluated to true, field is required
              fields.push(field)
            }
          })

          this.$store.dispatch('ui/setLayoutRequiredFields', fields)
        }).catch(this.toastErrorHandler(this.$t('notification:record.requiredFields.failed')))
    },

    prepareLayoutRequiredFieldsData () {
      const expressions = {}
      const variables = this.expressionVariables()

      const { requiredFields = [] } = this.layout.config.validation

      requiredFields.forEach(({ field, condition }) => {
        if (field && condition && condition.trim() !== '') {
          expressions[field] = condition
        }
      })

      return { expressions, variables }
    },
  },
}

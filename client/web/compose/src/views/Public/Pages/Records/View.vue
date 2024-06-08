<template>
  <div class="d-flex flex-column flex-grow-1 w-100 h-100">
    <b-alert
      v-if="isDeleted"
      show
      variant="warning"
      class="mb-2 mx-2"
    >
      {{ $t('record.recordDeleted') }}
    </b-alert>

    <portal
      :to="portalTopbarTitle"
    >
      {{ title }}
    </portal>

    <div
      v-if="!layout"
      class="d-flex align-items-center justify-content-center w-100 h-100"
    >
      <b-spinner />
    </div>

    <grid
      v-else-if="blocks"
      v-bind="$props"
      :errors="errors"
      :record="record"
      :blocks="blocks"
      :mode="inEditing ? 'editor' : 'base'"
      class="h-100"
    />

    <portal
      :to="portalRecordToolbar"
    >
      <record-toolbar
        v-if="layout"
        :module="module"
        :record="record"
        :labels="recordToolbarLabels"
        :processing="processing"
        :processing-submit="processingSubmit"
        :processing-delete="processingDelete"
        :processing-undelete="processingUndelete"
        :in-editing="inEditing"
        :show-record-modal="showRecordModal"
        :record-navigation="recordNavigation"
        :hide-back="!layoutButtons.has('back')"
        :hide-delete="!layoutButtons.has('delete')"
        :hide-new="!layoutButtons.has('new')"
        :hide-clone="!layoutButtons.has('clone')"
        :hide-edit="!layoutButtons.has('edit')"
        :hide-submit="!layoutButtons.has('submit')"
        :has-back="viewHasBack"
        @add="handleAdd()"
        @clone="handleClone()"
        @edit="handleEdit()"
        @view="handleView()"
        @delete="handleDelete()"
        @undelete="handleUndelete()"
        @back="handleBack()"
        @submit="handleFormSubmit('page.record')"
        @update-navigation="handleRedirectToPrevOrNext"
      >
        <template #start-actions>
          <b-button
            v-for="(action, index) in layoutActions.filter(a => a.placement === 'start')"
            :key="index"
            :variant="action.meta.style.variant"
            :disabled="processing"
            size="lg"
            class="text-nowrap"
            @click.prevent="handleAction(action)"
          >
            {{ action.meta.label }}
          </b-button>
        </template>

        <template #center-actions>
          <b-button
            v-for="(action, index) in layoutActions.filter(a => a.placement === 'center')"
            :key="index"
            :variant="action.meta.style.variant"
            :disabled="processing"
            size="lg"
            class="text-nowrap"
            @click.prevent="handleAction(action)"
          >
            {{ action.meta.label }}
          </b-button>
        </template>

        <template #end-actions>
          <b-button
            v-for="(action, index) in layoutActions.filter(a => a.placement === 'end')"
            :key="index"
            :variant="action.meta.style.variant"
            :disabled="processing"
            size="lg"
            class="text-nowrap"
            @click.prevent="handleAction(action)"
          >
            {{ action.meta.label }}
          </b-button>
        </template>
      </record-toolbar>
    </portal>
  </div>
</template>

<script>
import axios from 'axios'
import { isEqual } from 'lodash'
import { mapGetters, mapActions } from 'vuex'
import Grid from 'corteza-webapp-compose/src/components/Public/Page/Grid'
import RecordToolbar from 'corteza-webapp-compose/src/components/Common/RecordToolbar'
import record from 'corteza-webapp-compose/src/mixins/record'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'ViewRecord',

  components: {
    Grid,
    RecordToolbar,
  },

  mixins: [
    // The record mixin contains all of the logic for creating/editing/deleting/undeleting the record
    record,
  ],

  beforeRouteLeave (to, from, next) {
    next(this.checkUnsavedChanges())
  },

  beforeRouteUpdate (to, from, next) {
    const areParamsChanged = JSON.stringify(to.params) !== JSON.stringify(from.params)

    // If the route params have changed, we need to check for unsaved changes
    // We do this to avoid maginfy block to raise the unsaved changes prompt
    if (!areParamsChanged) {
      next()
      return
    }

    next(this.checkUnsavedChanges())
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    module: {
      type: compose.Module,
      required: false,
      default: () => ({}),
    },

    page: {
      type: compose.Page,
      required: true,
    },

    recordID: {
      type: String,
      required: false,
      default: '',
    },

    // When creating from related record blocks
    refRecord: {
      type: compose.Record,
      required: false,
      default: () => ({}),
    },

    // If component was called (via router) with some pre-seed values
    values: {
      type: Object,
      required: false,
      default: () => ({}),
    },

    // Open record in a modal
    showRecordModal: {
      type: Boolean,
      required: false,
    },

    edit: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      inEditing: this.edit,

      layouts: [],
      layout: undefined,
      layoutButtons: new Set(),
      blocks: undefined,

      recordNavigation: {
        prev: undefined,
        next: undefined,
      },

      abortableRequests: [],
    }
  },

  computed: {
    ...mapGetters({
      getNextAndPrevRecord: 'ui/getNextAndPrevRecord',
      getPageLayouts: 'pageLayout/getByPageID',
      previousPages: 'ui/previousPages',
      modalPreviousPages: 'ui/modalPreviousPages',
    }),

    isNew () {
      return !this.recordID || this.recordID === NoID
    },

    portalTopbarTitle () {
      return this.showRecordModal ? 'record-modal-header' : 'topbar-title'
    },

    portalRecordToolbar () {
      return this.showRecordModal ? 'record-modal-footer' : 'toolbar'
    },

    getUiEventResourceType () {
      return 'record-page'
    },

    recordToolbarLabels () {
      // Use an intermediate object so we can reflect all changes in one go;
      const aux = {}
      const { config = {} } = this.layout || {}
      const { buttons = {} } = config

      Object.entries(buttons).forEach(([key, { label = '' }]) => {
        aux[key] = label
      })
      return aux
    },

    layoutActions () {
      const { config = {} } = this.layout || {}
      const { actions = [] } = config

      return actions.filter(({ enabled }) => enabled)
    },

    title () {
      if (!this.layout) {
        return ''
      }

      const { config = {}, meta = {} } = this.layout || {}
      const { useTitle = false } = config

      if (useTitle) {
        try {
          return evaluatePrefilter(meta.title, {
            record: this.record,
            user: this.$auth.user || {},
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          })
        } catch (e) {
          return ''
        }
      }

      const { name, handle } = this.module

      const titlePrefix = this.isNew ? 'create' : this.inEditing ? 'edit' : 'view'

      return this.$t(`page:public.record.${titlePrefix}.title`, { name: name || handle, interpolation: { escapeValue: false } })
    },

    currentRecordNavigation () {
      const { recordID } = this.record || {}
      return this.getNextAndPrevRecord(recordID)
    },

    viewHasBack () {
      if (this.showRecordModal) {
        return this.modalPreviousPages.length > 1
      }

      return this.previousPages.length > 0
    },

    uniqueID () {
      return [(this.page || {}).pageID, this.recordID, this.edit]
    },
  },

  watch: {
    uniqueID: {
      immediate: true,
      handler (value = [], oldValue = []) {
        const [pageID = '', recordID = '', edit = false] = value
        const [oldPageID = '', oldRecordID = '', oldEdit = false] = oldValue

        // If page changed, get layouts
        if (pageID && pageID !== NoID && pageID !== oldPageID) {
          this.layout = undefined
          this.layouts = this.getPageLayouts(this.page.pageID)
        }

        // If page or record changed refresh record and determine layout
        if (pageID !== oldPageID || recordID !== oldRecordID) {
          this.record = undefined
          this.initialRecordState = undefined

          this.refresh()
        } else if (edit !== oldEdit) {
          this.determineLayout()
        }
      },
    },

    currentRecordNavigation: {
      handler (rn, oldRn) {
        // To prevent hiding and then showing the record navigation
        // We use the old value if its valid and the current one isn't
        if (rn.prev || rn.next) {
          this.recordNavigation = rn
        } else if (this.recordID !== NoID && (oldRn.prev || oldRn.next)) {
          this.recordNavigation = oldRn
        } else {
          this.recordNavigation = {
            prev: undefined,
            next: undefined,
          }
        }
      },
    },
  },

  mounted () {
    this.$root.$on('refetch-record-blocks', this.refetchRecordBlocks)

    if (this.showRecordModal) {
      this.$root.$on('bv::modal::hide', this.checkUnsavedChanges)
    }
  },

  beforeDestroy () {
    this.abortRequests()
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      popPreviousPages: 'ui/popPreviousPages',
      clearRecordSet: 'record/clearSet',
      popModalPreviousPage: 'ui/popModalPreviousPage',
    }),

    async loadRecord (recordID = this.recordID) {
      if (!this.page) {
        return
      }

      const { namespaceID, moduleID } = this.page

      if (moduleID !== NoID) {
        const module = Object.freeze(this.getModuleByID(moduleID).clone())

        if (recordID && recordID !== NoID) {
          const { response, cancel } = this.$ComposeAPI
            .recordReadCancellable({ namespaceID, moduleID, recordID })

          this.abortableRequests.push(cancel)

          return response()
            .then(record => {
              return new Promise(resolve => setTimeout(resolve, 300)).then(() => {
                return new compose.Record(module, record)
              })
            })
            .catch(e => {
              if (!axios.isCancel(e)) {
                this.toastErrorHandler(this.$t('notification:record.loadFailed'))(e)
                this.handleBack()
              }
            })
        } else {
          if (this.refRecord) {
            // Record create form called from a related records block,
            // we'll try to find an appropriate fields and cross-link this new record to ref

            this.module.fields.filter(f => f.kind === 'Record' && f.options.moduleID === this.refRecord.moduleID).forEach(f => {
              this.values[f.name] = this.refRecord.recordID
            })
          }

          return new compose.Record(module, { values: this.values })
        }
      }
    },

    async handleBack () {
      /**
       * Not the best way since we can not always know where we
       * came from (and "where" is back).
      */
      if (this.showRecordModal) {
        if (this.checkUnsavedChanges()) {
          this.popModalPreviousPage().then(({ recordID, recordPageID, edit }) => {
            this.$emit('on-modal-back', { recordID, recordPageID, pushModalPreviousPage: false, edit })
          })
        }
      } else {
        const previousPage = await this.popPreviousPages()
        const extraPop = !this.isNew

        this.$router.push(previousPage || { name: 'pages', params: { slug: this.namespace.slug || this.namespace.namespaceID } })
        // Pop an additional time so that the route we went back to isn't added to previousPages
        if (extraPop) {
          this.popPreviousPages()
        }
      }
    },

    handleAdd () {
      this.processing = true

      if (this.showRecordModal) {
        if (this.checkUnsavedChanges()) {
          this.$emit('handle-record-redirect', { recordID: NoID, recordPageID: this.page.pageID, edit: true })
        }
      } else {
        this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, edit: true } })
      }
    },

    handleClone () {
      this.processing = true

      if (this.showRecordModal) {
        if (this.checkUnsavedChanges()) {
          this.$emit('handle-record-redirect', { recordID: NoID, recordPageID: this.page.pageID, values: this.record.values, edit: true })
        }
      } else {
        this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, values: this.record.values, edit: true } })
      }
    },

    handleEdit () {
      this.processing = true

      if (this.showRecordModal) {
        this.$emit('handle-record-redirect', { recordID: this.recordID, recordPageID: this.page.pageID, edit: true })
      } else {
        this.$router.push({ name: 'page.record.edit', params: { recordID: this.recordID, pageID: this.page.pageID, edit: true } })
      }
    },

    handleView () {
      this.processing = true

      if (this.showRecordModal) {
        if (this.checkUnsavedChanges()) {
          this.$emit('handle-record-redirect', { recordID: this.recordID, recordPageID: this.page.pageID, edit: false })
        }
      } else {
        this.$router.push({ name: 'page.record', params: { recordID: this.recordID, pageID: this.page.pageID, edit: false } })
      }
    },

    handleRedirectToPrevOrNext (recordID) {
      if (!recordID) return

      this.processing = true

      if (this.showRecordModal) {
        if (this.checkUnsavedChanges()) {
          this.$emit('handle-record-redirect', { recordID, recordPageID: this.page.pageID })
        }
      } else {
        this.$router.push({
          params: { ...this.$route.params, recordID },
        })
        this.popPreviousPages()
      }
    },

    evaluateLayoutExpressions (variables = {}) {
      const expressions = {}
      variables = {
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
        isView: !this.edit && !this.isNew,
        isCreate: this.isNew,
        isEdit: this.edit && !this.isNew,
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
      this.resetErrors()
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

        return
      }

      this.inEditing = this.edit

      if (this.layout && matchedLayout.pageLayoutID === this.layout.pageLayoutID) {
        return
      }

      this.layout = matchedLayout

      const { blocks = [], config = {} } = this.layout || {}
      const { buttons = [] } = config

      this.layoutButtons = Object.entries(buttons).reduce((acc, [key, value]) => {
        if (value.enabled) {
          acc.add(key)
        }
        return acc
      }, new Set())

      const tempBlocks = []

      blocks.forEach(({ blockID, xywh }) => {
        const block = this.page.blocks.find(b => b.blockID === blockID)

        if (block) {
          block.xywh = xywh
          tempBlocks.push(block)
        }
      })

      this.blocks = tempBlocks
    },

    refetchRecordBlocks () {
      // Don't refresh when creating and prompt user before refreshing when editing
      if (this.isNew || (this.edit && this.compareRecordValues() && !window.confirm(this.$t('notification:record.staleDataRefresh')))) {
        return
      }

      // Refetch the record and other page blocks that use records
      this.loadRecord().then(record => {
        this.record = record
        this.initialRecordState = record.clone()
      })
      this.$root.$emit(`refetch-non-record-blocks:${this.page.pageID}`)
    },

    async refresh (variables = {}) {
      this.processing = true

      return this.loadRecord().then(record => {
        // Set record so the latest one can be used in the layout expressions
        variables.record = record ? record.serialize() : {}

        return this.determineLayout(undefined, variables).then(() => {
          this.record = record
          this.initialRecordState = record.clone()
        })
      }).finally(() => {
        this.processing = false
      })
    },

    handleAction ({ kind, params = {} }) {
      if (kind === 'toLayout') {
        this.processing = true

        this.determineLayout(params.pageLayoutID, {}, false).finally(() => {
          this.processing = false
        })
      } else if (kind === 'toURL') {
        window.open(params.url, params.openIn === 'newTab' ? '_blank' : '_self')
      }
    },

    setDefaultValues () {
      this.inEditing = false
      this.layouts = []
      this.layout = undefined
      this.layoutButtons.clear()
      this.blocks = undefined
      this.recordNavigation = {}
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    destroyEvents () {
      this.$root.$off('refetch-record-blocks', this.refetchRecordBlocks)

      if (this.showRecordModal) {
        this.$root.$off('bv::modal::hide', this.checkUnsavedChanges)
      }
    },

    compareRecordValues () {
      const recordValues = JSON.parse(JSON.stringify(this.record ? this.record.values : {}))
      const initialRecordState = JSON.parse(JSON.stringify(this.initialRecordState ? this.initialRecordState.values : {}))

      return !isEqual(recordValues, initialRecordState)
    },

    checkUnsavedChanges (bvEvent, modalId) {
      if ((bvEvent && modalId !== 'record-modal') || !this.edit) return true

      const recordStateChange = this.compareRecordValues() ? window.confirm(this.$t('general:record.unsavedChanges')) : true

      if (!recordStateChange) {
        this.processing = false

        if (bvEvent) {
          bvEvent.preventDefault()
        }
      } else {
        this.record = this.initialRecordState ? this.initialRecordState.clone() : undefined
      }

      return recordStateChange
    },
  },
}
</script>

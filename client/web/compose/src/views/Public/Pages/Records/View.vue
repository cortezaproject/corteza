<template>
  <div class="d-flex flex-grow-1 w-100 h-100">
    <b-alert
      v-if="isDeleted"
      show
      variant="warning"
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
      @reload="loadRecord()"
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
import { isEqual } from 'lodash'
import { mapGetters, mapActions } from 'vuex'
import Grid from 'corteza-webapp-compose/src/components/Public/Page/Grid'
import RecordToolbar from 'corteza-webapp-compose/src/components/Common/RecordToolbar'
import record from 'corteza-webapp-compose/src/mixins/record'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import axios from 'axios'

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
      inEditing: false,
      inCreating: false,

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

    portalTopbarTitle () {
      return this.showRecordModal ? 'record-modal-header' : 'topbar-title'
    },

    portalRecordToolbar () {
      return this.showRecordModal ? 'record-modal-footer' : 'toolbar'
    },

    newRouteParams () {
      // Remove recordID and values from route params
      const { recordID, values, ...params } = this.$route.params
      return params
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
            recordID: (this.record || {}).recordID || NoID,
            ownerID: (this.record || {}).ownedBy || NoID,
            userID: (this.$auth.user || {}).userID || NoID,
          })
        } catch (e) {
          return ''
        }
      }

      const { name, handle } = this.module

      const titlePrefix = this.inCreating ? 'create' : this.inEditing ? 'edit' : 'view'

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
  },

  watch: {
    recordID: {
      immediate: true,
      handler () {
        this.record = undefined
        this.initialRecordState = undefined
        this.refresh()
        this.loadRecord().then(() => {
          this.determineLayout()
        })
      },
    },

    'page.pageID': {
      immediate: true,
      handler (pageID) {
        if (pageID === NoID) return

        this.layouts = this.getPageLayouts(this.page.pageID)
        this.layout = undefined
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

    edit: {
      immediate: true,
      handler (value) {
        if (value) {
          this.inEditing = true
          this.inCreating = true
        }
      },
    },
  },

  mounted () {
    this.$root.$on('refetch-record-blocks', this.refetchRecordBlocks)

    if (this.showRecordModal) {
      this.$root.$on('bv::modal::hide', this.checkUnsavedChangesOnModal)
    }
  },

  beforeDestroy () {
    this.abortRequests()
    this.destroyEvents()
    this.setDefaultValues()
  },

  // Destroy event before route leave to ensure it doesn't destroy the newly created one
  beforeRouteLeave (to, from, next) {
    this.$root.$off('refetch-record-blocks', this.refetchRecordBlocks)
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
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
                this.record = new compose.Record(module, record)
                this.initialRecordState = this.record.clone()
              })
            })
            .catch(e => {
              if (!axios.isCancel(e)) {
                this.toastErrorHandler(this.$t('notification:record.loadFailed'))(e)
                this.handleBack()
              }
            })
        } else {
          this.record = new compose.Record(module, { values: this.values })
          this.initialRecordState = this.record.clone()

          this.inEditing = true
          this.inCreating = true
        }

        if (this.refRecord) {
          // Record create form called from a related records block,
          // we'll try to find an appropriate field and cross-link this new record to ref
          const recRefField = this.module.fields.find(f => f.kind === 'Record' && f.options.moduleID === this.refRecord.moduleID)
          if (recRefField) {
            this.record.values[recRefField.name] = this.refRecord.recordID
          }
        }
      }
    },

    async handleBack () {
      /**
       * Not the best way since we can not always know where we
       * came from (and "where" is back).
      */
      if (this.showRecordModal) {
        this.popModalPreviousPage().then(({ recordID, recordPageID }) => {
          this.$emit('on-modal-back', { recordID, recordPageID, pushModalPreviousPage: false })
          this.inCreating = recordID === NoID
          this.inEditing = false
        })

        return
      }

      const previousPage = await this.popPreviousPages()
      const extraPop = !this.inCreating
      this.$router.push(previousPage || { name: 'pages', params: { slug: this.namespace.slug || this.namespace.namespaceID } })
      // Pop an additional time so that the route we went back to isn't added to previousPages
      if (extraPop) {
        this.popPreviousPages()
      }
    },

    handleAdd () {
      if (!this.showRecordModal) {
        this.$router.push({ name: 'page.record.create', params: this.newRouteParams })
        return
      }

      this.inEditing = true
      this.inCreating = true
      this.record = new compose.Record(this.module, { values: this.values })
      this.initialRecordState = this.record.clone()
      this.$emit('handle-record-redirect', { recordID: NoID, recordPageID: this.page.pageID })
    },

    handleClone () {
      if (!this.showRecordModal) {
        this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, values: this.record.values } })
        return
      }

      this.inEditing = true
      this.inCreating = true
      this.record = new compose.Record(this.module, { values: this.record.values })
      this.initialRecordState = this.record.clone()
      this.$emit('handle-record-redirect', { recordID: NoID, recordPageID: this.page.pageID, values: this.record.values })
    },

    handleEdit () {
      this.processing = true

      this.refresh({ isView: false, isEdit: true, isCreate: false }).then(() => {
        this.inEditing = true
        this.inCreating = false

        this.$root.$emit(`refetch-non-record-blocks:${this.page.pageID}`)
      }).finally(() => {
        this.processing = false
      })
    },

    handleView () {
      this.processing = true

      this.refresh({ isView: true, isEdit: false, isCreate: false }).then(() => {
        this.inEditing = false
        this.inCreating = false

        this.$root.$emit(`refetch-non-record-blocks:${this.page.pageID}`)
      }).finally(() => {
        this.processing = false
      })
    },

    handleRedirectToPrevOrNext (recordID) {
      if (!recordID) return

      if (this.showRecordModal) {
        this.$emit('handle-record-redirect', { recordID, recordPageID: this.page.pageID })
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
        isView: !this.inEditing && !this.inCreating,
        isCreate: this.inCreating,
        isEdit: this.inEditing,
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

    async determineLayout (pageLayoutID, variables = {}) {
      // Clear stored records so they can be refetched with latest values
      this.clearRecordSet()
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

      if (!matchedLayout) {
        this.toastWarning(this.$t('notification:page.page-layout.notFound.view'))
        return this.$router.go(-1)
      }

      if (this.layout && matchedLayout.pageLayoutID === this.layout.pageLayoutID) {
        return
      }

      this.layout = matchedLayout

      const { config = {} } = this.layout
      const { buttons = [] } = config

      this.layoutButtons = Object.entries(buttons).reduce((acc, [key, value]) => {
        if (value.enabled) {
          acc.add(key)
        }
        return acc
      }, new Set())

      this.blocks = (this.layout || {}).blocks.map(({ blockID, xywh }) => {
        const block = this.page.blocks.find(b => b.blockID === blockID)
        block.xywh = xywh
        return block
      })
    },

    refetchRecordBlocks () {
      // Don't refresh when creating and prompt user before refreshing when editing
      if (this.inCreating || (this.inEditing && !window.confirm(this.$t('notification:record.staleDataRefresh')))) {
        return
      }

      // Refetch the record and other page blocks that use records
      this.loadRecord()
      this.$root.$emit(`refetch-non-record-blocks:${this.page.pageID}`)
    },

    async refresh (variables = {}) {
      return this.loadRecord().then(() => {
        return this.determineLayout(undefined, variables)
      })
    },

    handleAction ({ kind, params = {} }) {
      if (kind === 'toLayout') {
        this.determineLayout(params.pageLayoutID)
      } else if (kind === 'toURL') {
        window.open(params.url, params.openIn === 'newTab' ? '_blank' : '_self')
      }
    },

    setDefaultValues () {
      this.inEditing = false
      this.inCreating = false
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
        this.$root.$off('bv::modal::hide', this.checkUnsavedChangesOnModal)
      }
    },

    compareRecordValues () {
      const recordValues = JSON.parse(JSON.stringify(this.record ? this.record.values : {}))
      const initialRecordState = JSON.parse(JSON.stringify(this.initialRecordState ? this.initialRecordState.values : {}))

      return !isEqual(recordValues, initialRecordState)
    },

    checkUnsavedChanges (next, to) {
      if (this.inCreating) {
        next(true)
      } else {
        next(this.compareRecordValues() ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },

    checkUnsavedChangesOnModal (bvEvent, modalId) {
      if (modalId === 'record-modal' && !this.inCreating) {
        const recordStateChange = this.compareRecordValues() ? window.confirm(this.$t('general:editor.unsavedChanges')) : true

        if (!recordStateChange) {
          bvEvent.preventDefault()
        }
      }
    },
  },
}
</script>

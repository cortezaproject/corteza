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
import page from 'corteza-webapp-compose/src/mixins/page'
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
    page,
  ],

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteUpdate (to, from, next) {
    const areParamsChanged = JSON.stringify(to.params) !== JSON.stringify(from.params)

    // If the route params have changed, we need to check for unsaved changes
    // We do this to avoid maginfy block to raise the unsaved changes prompt
    if (!areParamsChanged) {
      next()
      return
    }

    this.checkUnsavedChanges(next, to)
  },

  props: {
    module: {
      type: compose.Module,
      required: false,
      default: () => ({}),
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

      layoutButtons: new Set(),

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

    createEvents () {
      this.$root.$on('refetch-record-blocks', this.refetchRecordBlocks)
      this.$root.$on('record-field-change', this.validateBlocksVisibilityCondition)

      if (this.showRecordModal) {
        this.$root.$on('bv::modal::hide', this.checkUnsavedChangesOnModal)
      }
    },

    validateBlocksVisibilityCondition ({ fieldName }) {
      const { blocks = [] } = this.page

      if (blocks.some(({ meta = {} }) => ((meta.visibility || {}).expression).includes(fieldName))) {
        this.updateBlocks()
      }
    },

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

    handleRecordButtons () {
      const { config = {} } = this.layout
      const { buttons = [] } = config

      this.layoutButtons = Object.entries(buttons).reduce((acc, [key, value]) => {
        if (value.enabled) {
          acc.add(key)
        }
        return acc
      }, new Set())
    },

    refetchRecordBlocks () {
      // Don't refresh when creating and prompt user before refreshing when editing
      if (this.inCreating || (this.inEditing && this.compareRecordValues() && !window.confirm(this.$t('notification:record.staleDataRefresh')))) {
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
      this.$root.$off('record-field-change', this.validateBlocksVisibilityCondition)

      if (this.showRecordModal) {
        this.$root.$off('bv::modal::hide', this.checkUnsavedChangesOnModal)
      }
    },

    compareRecordValues () {
      const recordValues = JSON.parse(JSON.stringify(this.record ? this.record.values : {}))
      const initialRecordState = JSON.parse(JSON.stringify(this.initialRecordState ? this.initialRecordState.values : {}))

      return !isEqual(recordValues, initialRecordState)
    },

    checkUnsavedChanges (next) {
      if (!this.inEditing) {
        next(true)
        return
      }

      next(this.compareRecordValues() ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
    },

    checkUnsavedChangesOnModal (bvEvent, modalId) {
      if (modalId !== 'record-modal' || !this.inEditing) return

      const recordStateChange = this.compareRecordValues() ? window.confirm(this.$t('general:editor.unsavedChanges')) : true

      if (!recordStateChange) {
        bvEvent.preventDefault()
      }
    },
  },
}
</script>

<template>
  <div class="d-flex flex-column flex-grow-1 w-100 h-100">
    <portal
      :to="portalTopbarTitle"
    >
      {{ title }}
    </portal>

    <portal to="topbar-tools">
      <b-button
        v-if="page && isRecordPage && page.canUpdatePage"
        variant="primary"
        :disabled="!moduleEditor"
        :to="moduleEditor"
        size="sm"
        class="d-flex align-items-center mr-2"
      >
        {{ $t('navigation.editModule') }}
        <font-awesome-icon
          :icon="['far', 'edit']"
          class="ml-2"
        />
      </b-button>

      <b-button-group
        v-if="page && page.canUpdatePage"
        size="sm"
        class="text-nowrap"
      >
        <b-button
          data-test-id="button-page-builder"
          variant="primary"
          class="d-flex align-items-center"
          :to="pageBuilder"
        >
          {{ $t('general:label.pageBuilder') }}
          <font-awesome-icon
            :icon="['fas', 'tools']"
            class="ml-2"
          />
        </b-button>

        <b-button
          v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.edit.page'), boundary: 'body' }"
          data-test-id="button-page-edit"
          :to="pageEditor"
          variant="primary"
          class="d-flex align-items-center"
          style="margin-left:2px;"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>

        <page-translator
          v-if="trPage"
          data-test-id="button-page-translations"
          :page.sync="trPage"
          :page-layout.sync="layout"
          button-variant="primary"
          style="margin-left:2px;"
        />
      </b-button-group>
    </portal>

    <b-alert
      v-if="isDeleted"
      show
      variant="warning"
      class="mx-1 my-2"
    >
      {{ $t('block:record.recordDeleted') }}
    </b-alert>

    <b-alert
      v-if="isDraft"
      show
      variant="light"
      class="d-flex align-items-center mx-1 my-2"
    >
      {{ $t('drafts:activeDraft') }}
    </b-alert>

    <div
      v-if="isLoading"
      class="d-flex align-items-center justify-content-center w-100 h-100"
    >
      <b-spinner />
    </div>

    <grid
      v-else
      v-bind="$props"
      :errors="errors"
      :record="record"
      :loading-record="loadingRecord"
      :blocks="blocks"
      :mode="inEditing ? 'editor' : 'base'"
      class="h-100"
    />

    <portal
      :to="portalRecordToolbar"
    >
      <record-toolbar
        :module="module"
        :record="record"
        :labels="recordToolbarLabels"
        :processing="processing"
        :processing-action="processingAction"
        :in-editing="inEditing"
        :in-modal="inModal"
        :is-created="!isNew"
        :record-navigation="recordNavigation"
        :is-draft="isDraft"
        :is-new="isNew"
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
            :to="generateActionLink(action)"
            :href="generateActionHref(action)"
            :target="generateActionTarget(action)"
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
            :to="generateActionLink(action)"
            :href="generateActionHref(action)"
            :target="generateActionTarget(action)"
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
            :to="generateActionLink(action)"
            :href="generateActionHref(action)"
            :target="generateActionTarget(action)"
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
import { isEqual, throttle } from 'lodash'
import { mapGetters, mapActions } from 'vuex'
import Grid from 'corteza-webapp-compose/src/components/Public/Page/Grid'
import RecordToolbar from 'corteza-webapp-compose/src/components/Common/RecordToolbar'
import record from 'corteza-webapp-compose/src/mixins/record'
import page from 'corteza-webapp-compose/src/mixins/page'
import { compose, system, NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  i18nOptions: {
    namespaces: 'page',
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
    next(this.checkUnsavedChanges())
  },

  beforeRouteUpdate (to, from, next) {
    const areParamsChanged = JSON.stringify(to.params) !== JSON.stringify(from.params)

    // If the route params have changed, we need to check for unsaved changes
    // We do this to avoid magnify block to raise the unsaved changes prompt
    if (!areParamsChanged) {
      next()
      return
    }

    next(this.checkUnsavedChanges())
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
      default: undefined,
    },

    // If component was called (via router) with some pre-seed values
    values: {
      type: Object,
      required: false,
      default: () => ({}),
    },

    // Open record in a modal
    inModal: {
      type: Boolean,
      default: false,
    },

    edit: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      inEditing: this.edit,

      loading: false,

      layoutButtons: new Set(),

      recordNavigation: {
        prev: undefined,
        next: undefined,
      },

      abortableRequests: [],

      loadingRecord: false,

      // Used to identify, load, and delete the correct draft
      activeDraftKey: null,
    }
  },

  computed: {
    ...mapGetters({
      getNextAndPrevRecord: 'ui/getNextAndPrevRecord',
      previousPages: 'ui/previousPages',
      modalPreviousPages: 'ui/modalPreviousPages',
    }),

    isNew () {
      return !this.recordID || this.recordID === NoID
    },

    isLoading () {
      return this.loading || !this.layout || !this.blocks
    },

    portalTopbarTitle () {
      return this.inModal ? 'record-modal-header' : 'topbar-title'
    },

    portalRecordToolbar () {
      return this.inModal ? 'record-modal-footer' : 'toolbar'
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

      if (this.isDraft) {
        aux.delete = this.$t('drafts:deleteDraft')
        aux.clone = this.$t('drafts:saveAsNewDraft')
        aux.new = this.$t('drafts:saveAsNewDraft')
        aux.submit = this.isNew ? this.$t('general:label.createRecordFromDraft') : this.$t('drafts:applyDraft')
        aux.edit = this.$t('drafts:viewRecord')
      }

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
          return e
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
      if (this.inModal) {
        return this.modalPreviousPages.length > 1
      }

      return this.previousPages.length > 0
    },

    uniqueID () {
      return [(this.page || {}).pageID, this.$route.query.layoutID, this.$route.query.modalLayoutID, this.recordID, this.edit, this.$route.query.draftID]
    },

    showDrafts () {
      return this.$Settings.get('ui.topbar.showDrafts', false)
    },

    isDraft () {
      return !!this.$route.query.draftID
    },
  },

  watch: {
    uniqueID: {
      immediate: true,
      handler (value = [], oldValue = []) {
        const [pageID = '', pageLayoutID = '', modalPageLayoutID = '', recordID = '', edit = '', draftID = ''] = value
        const [oldPageID = '', oldPageLayoutID = '', oldModalPageLayoutID = '', oldRecordID = '', oldEdit = '', oldDraftID = ''] = oldValue

        if (!pageID || pageID === NoID) return

        // If page changed, get layouts
        if (pageID !== oldPageID) {
          this.loading = true
          this.layouts = this.getPageLayouts(this.page.pageID)
        }

        // Only refresh if the record ID has changed or the edit state has changed
        if ((recordID === NoID && recordID !== oldRecordID) || recordID !== oldRecordID || edit !== oldEdit || pageID !== oldPageID || draftID !== oldDraftID) {
          this.refresh()
          return
        }

        // Determine which layout ID to use based on modal state
        const currentLayoutID = this.inModal ? modalPageLayoutID : pageLayoutID
        const oldLayoutID = this.inModal ? oldModalPageLayoutID : oldPageLayoutID

        // Only update layout if it has actually changed
        if (currentLayoutID !== oldLayoutID) {
          this.determineLayout({ pageLayoutID: currentLayoutID })
            .then(blocks => {
              if (blocks) {
                this.blocks = blocks
              }
            })
            .finally(() => {
              this.processing = false
            })
        }
      },
    },

    'layout.handle': {
      immediate: true,
      handler (handle, oldHandle) {
        if (handle !== oldHandle) {
          this.inModal ? this.setModalLayoutHandle(handle) : this.setLayoutHandle(handle)
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

    title: {
      immediate: true,
      handler (title) {
        if (title && !this.inModal) {
          document.title = title
        }
      },
    },
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.abortRequests()
    this.destroyEvents()
    this.cleanupDraftSync()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      popPreviousPages: 'ui/popPreviousPages',
      popModalPreviousPage: 'ui/popModalPreviousPage',
      setLayoutHandle: 'ui/setLayoutHandle',
      setModalLayoutHandle: 'ui/setModalLayoutHandle',
      updateRecordSet: 'record/updateRecords',
    }),

    createEvents () {
      this.$root.$on('refetch-records', this.refetchRecords)
      this.$root.$on('record-field-change', this.evaluateLayoutConditions)
      this.$root.$on('record-field-change', this.saveDraftRevision)

      if (this.inModal) {
        this.$root.$on('bv::modal::hide', this.checkUnsavedChanges)
      }
    },

    evaluateLayoutConditions () {
      this.evaluateBlocks()
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
              record = new compose.Record(module, record)
              this.updateRecordSet(record)

              return new Promise(resolve => setTimeout(resolve, 300)).then(() => {
                return record
              })
            })
            .catch(e => {
              if (!axios.isCancel(e)) {
                this.toastErrorHandler(this.$t('notification:record.loadFailed'))(e)
                this.handleBack()
              }
            })
        } else {
          if (this.refRecord && this.refRecord.recordID && this.refRecord.recordID !== NoID) {
            this.updateRecordSet(this.refRecord)

            // Record create form called from a related records block,
            // we'll try to find an appropriate fields and cross-link this new record to ref

            this.module.fields.filter(f => f.kind === 'Record' && f.options.moduleID === this.refRecord.moduleID).forEach(f => {
              if (f.isMulti) {
                this.values[f.name] = [this.refRecord.recordID]
              } else {
                this.values[f.name] = this.refRecord.recordID
              }
            })
          }

          const { userID } = this.$auth.user

          await new Promise(resolve => setTimeout(resolve, 300))
          // Prefill ownedBy field with current user
          return new compose.Record(module, { ownedBy: userID, values: this.values })
        }
      }
    },

    async handleBack () {
      /**
       * Not the best way since we can not always know where we
       * came from (and "where" is back).
      */
      if (this.inModal) {
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

      if (this.isDraft) {
        const currentDraftID = this.activeDraftKey
        this.activeDraftKey = this.generateChangeID()
        this.saveDraft()
        this.activeDraftKey = currentDraftID
        this.toastSuccess(this.$t('drafts:saved'))
        this.processing = false
        return
      }

      if (this.inModal) {
        if (this.checkUnsavedChanges()) {
          this.$emit('handle-record-redirect', { recordID: NoID, recordPageID: this.page.pageID, edit: true })
        }
      } else {
        this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, edit: true } })
      }
    },

    handleClone () {
      this.processing = true

      if (this.isDraft) {
        const currentDraftID = this.activeDraftKey
        this.activeDraftKey = this.generateChangeID()
        this.saveDraft()
        this.activeDraftKey = currentDraftID
        this.toastSuccess(this.$t('drafts:saved'))
        this.processing = false
        return
      }

      if (this.inModal) {
        if (this.checkUnsavedChanges()) {
          this.$emit('handle-record-redirect', { recordID: NoID, recordPageID: this.page.pageID, values: this.record.values, edit: true })
        }
      } else {
        this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, values: this.record.values, edit: true } })
      }
    },

    handleEdit () {
      this.processing = true

      if (this.inModal) {
        this.$emit('handle-record-redirect', { recordID: this.recordID, recordPageID: this.page.pageID, edit: true })
      } else {
        this.$router.push({ name: 'page.record.edit', params: { recordID: this.recordID, pageID: this.page.pageID, edit: true } })
      }
    },

    handleView () {
      this.processing = true

      if (this.isDraft) {
        this.openDraftRecord()
        this.processing = false
        return
      }

      if (this.inModal) {
        if (this.checkUnsavedChanges()) {
          this.$emit('handle-record-redirect', { recordID: this.recordID, recordPageID: this.page.pageID, edit: false })
        }
      } else {
        this.$router.push({ name: 'page.record', params: { recordID: this.recordID, pageID: this.page.pageID, edit: false } })
      }
    },

    handleDelete: throttle(function () {
      this.processing = true
      this.processingAction = 'delete'

      if (this.isDraft) {
        this.$store.dispatch('drafts/removeDraft', { changeID: this.activeDraftKey })
        this.activeDraftKey = null
        this.toastSuccess(this.$t('drafts:deleted'))
        this.openDraftRecord()
        this.processing = false
        return
      }

      return this.dispatchUiEvent('beforeDelete')
        .then(() => this.$ComposeAPI.recordDelete(this.record))
        .then(this.dispatchUiEvent('afterDelete'))
        .then(() => {
          // Clear draft revision on successful delete
          if (this.activeDraftKey && this.$store) {
            this.$store.dispatch('drafts/removeDraft', { changeID: this.activeDraftKey })
            this.activeDraftKey = null
          }

          this.$root.$emit('refetch-records')
          this.toastSuccess(this.$t('notification:record.deleteSuccess'))
        }).catch(e => {
          this.processing = false
          this.toastErrorHandler(this.$t('notification:record.deleteFailed'))(e)
        })
    }, 500),

    handleRedirectToPrevOrNext (recordID) {
      if (!recordID) return

      this.processing = true

      if (this.inModal) {
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

    openDraftRecord () {
      const { draftID, ...query } = this.$route.query
      this.activeDraftKey = null

      if (this.inModal) {
        this.$emit('handle-record-redirect', { recordID: this.recordID, recordPageID: this.page.pageID, edit: false })
      } else if (!this.isNew) {
        this.$router.push({ name: 'page.record', params: { ...this.$route.params, recordID: this.recordID }, query })
      } else {
        this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, edit: true } })
      }
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

    refetchRecords ({ recordID } = {}) {
      if ((recordID && recordID === this.recordID && this.inModal)) {
        return
      }

      // Don't refresh when creating and prompt user before refreshing when editing
      if (this.isNew || (this.edit && this.compareRecordValues() && !window.confirm(this.$t('notification:record.staleDataRefresh')))) {
        return
      }

      this.refresh()
    },

    async refresh () {
      this.processing = true
      this.loadingRecord = true

      if (this.isRecordPage) {
        this.inEditing = this.edit
      }

      return this.loadRecord().then(record => {
        this.tempRecord = record

        const pageLayoutID = this.$route.query.layoutID

        return this.determineLayout({ pageLayoutID }).then(blocks => {
          if (blocks) {
            this.blocks = blocks
          }

          this.record = this.tempRecord
          this.initialRecordState = this.record.clone()

          this.getRecordDraft()
        })
      }).finally(() => {
        this.tempRecord = undefined

        this.processing = false
        this.loading = false
        this.loadingRecord = false
      })
    },

    generateActionLink (action) {
      const { kind, params = {} } = action

      if (kind === 'toLayout') {
        const pageLayoutID = params.pageLayoutID

        if (pageLayoutID) {
          if (this.inModal) {
            return {
              ...this.$route,
              query: {
                ...this.$route.query,
                modalLayoutID: pageLayoutID,
              },
            }
          } else {
            return {
              ...this.$route,
              query: {
                ...this.$route.query,
                layoutID: pageLayoutID,
              },
            }
          }
        }
      }

      return undefined
    },

    generateActionHref (action) {
      const { kind, params = {} } = action

      if (kind === 'toURL') {
        return params.url
      }

      return undefined
    },

    generateActionTarget (action) {
      const { kind, params = {} } = action

      if (kind === 'toURL') {
        return params.openIn === 'newTab' ? '_blank' : '_self'
      }

      return undefined
    },

    setDefaultValues () {
      this.inEditing = false
      this.layoutButtons.clear()
      this.blocks = undefined
      this.recordNavigation = {
        prev: undefined,
        next: undefined,
      }
      this.abortableRequests = []
      this.loadingRecord = false
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    destroyEvents () {
      this.$root.$off('refetch-records', this.refetchRecords)
      this.$root.$off('record-field-change', this.evaluateLayoutConditions)
      this.$root.$off('record-field-change', this.saveDraftRevision)

      if (this.inModal) {
        this.$root.$off('bv::modal::hide', this.checkUnsavedChanges)
      }
    },

    compareRecordValues () {
      const recordValues = JSON.parse(JSON.stringify(this.record ? this.record.values : {}))
      const initialRecordState = JSON.parse(JSON.stringify(this.initialRecordState ? this.initialRecordState.values : {}))

      return !isEqual(recordValues, initialRecordState)
    },

    checkUnsavedChanges (bvEvent, modalId) {
      if ((bvEvent && modalId !== 'record-modal') || !this.edit || this.isDraft) return true

      let recordStateChange = true

      if (this.compareRecordValues()) {
        const message = this.showDrafts ? this.$t('general:record.unsavedChangesDraft') : this.$t('general:record.unsavedChanges')
        recordStateChange = window.confirm(message)
      }

      if (!recordStateChange) {
        this.processing = false

        if (bvEvent) {
          bvEvent.preventDefault()
        }
      } else if (this.record) {
        this.initialRecordState = this.record.clone()
      }

      return recordStateChange
    },

    getRecordDraft () {
      this.activeDraftKey = null

      if (!this.record || !this.inEditing) return

      const { draftID } = this.$route.query
      if (!draftID) return

      const draft = this.$store.getters['drafts/getDraft'](draftID)
      if (!draft) return

      const { revision } = draft
      if (revision.record) {
        const revRecord = new compose.Record(this.module, revision.record)
        this.record.values = revRecord.values
      } else {
        // Fallback to changes
        revision.changes.forEach(({ key, new: values }) => {
          if (this.module.fields.find(f => f.name === key)) {
            this.record.values[key] = this.module.makeValue(key, values)
          }
        })
      }

      this.activeDraftKey = revision.changeID
    },

    cleanupDraftSync () {
      this.activeDraftKey = null
    },

    saveDraftRevision: throttle(function () {
      this.saveDraft()
    }, 2000),

    saveDraft () {
      if (!this.record || !this.inEditing || !this.showDrafts) return

      const changes = this.computeChanges()
      if (changes.length === 0) return

      if (!this.activeDraftKey) {
        this.activeDraftKey = this.generateChangeID()
      }

      const resourceID = this.isNew ? '0' : String(this.record.recordID)
      const resource = `compose:record/${this.namespace.namespaceID}/${this.module.moduleID}/${resourceID}`

      const revision = new system.Revision({
        changeID: this.activeDraftKey,
        timestamp: new Date().toISOString(),
        resource,
        revision: this.record.revision || 0,
        operation: this.isNew ? 'created' : 'updated',
        status: 'draft',
        userID: String(this.$auth.user?.userID || '0'),
        changes,
        comment: '',
        record: this.record.clone(),
      })

      this.$store.dispatch('drafts/saveDraft', { revision })
    },

    computeChanges () {
      const changes = []
      const current = this.record?.values || {}
      const initial = this.initialRecordState?.values || {}
      const allKeys = new Set([...Object.keys(current), ...Object.keys(initial)])

      for (const key of allKeys) {
        const oldVal = initial[key]
        const newVal = current[key]

        if (JSON.stringify(oldVal) === JSON.stringify(newVal)) continue

        changes.push({
          key,
          old: this.valueToArray(oldVal),
          new: this.valueToArray(newVal),
        })
      }

      return changes
    },

    valueToArray (value) {
      if (value === undefined || value === null) return []
      if (Array.isArray(value)) return value.map(v => v ?? '')
      return [value]
    },

    generateChangeID () {
      return `local-${Date.now()}-${Math.random().toString(36).substring(2, 10)}`
    },
  },
}
</script>

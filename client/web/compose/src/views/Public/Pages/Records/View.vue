<template>
  <div
    :class="{ 'flex-column': showRecordModal }"
    class="d-flex flex-grow-1 w-100 h-100"
  >
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

    <grid
      v-bind="$props"
      :errors="errors"
      :record="record"
      :blocks="blocks"
      :mode="inEditing ? 'editor' : 'base'"
      @reload="loadRecord()"
    />

    <portal
      :to="portalRecordToolbar"
    >
      <record-toolbar
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
        @add="handleAdd()"
        @clone="handleClone()"
        @edit="handleEdit()"
        @delete="handleDelete()"
        @undelete="handleUndelete()"
        @back="handleBack()"
        @submit="handleFormSubmit('page.record')"
        @update-navigation="handleRedirectToPrevOrNext"
      >
        <template #actions>
          <b-button
            v-for="(action, index) in layoutActions"
            :key="index"
            variant="primary"
            size="lg"
            class="ml-2"
            @click.prevent="determineLayout(action.params.pageLayoutID)"
          >
            {{ action.Meta.label }}
          </b-button>
        </template>
      </record-toolbar>
    </portal>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Grid from 'corteza-webapp-compose/src/components/Public/Page/Grid'
import RecordToolbar from 'corteza-webapp-compose/src/components/Common/RecordToolbar'
import record from 'corteza-webapp-compose/src/mixins/record'
import { compose, NoID } from '@cortezaproject/corteza-js'

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

    // Open record in a modal
    showRecordModal: {
      type: Boolean,
      required: false,
    },
  },

  data () {
    return {
      inEditing: false,
      inCreating: false,

      layouts: [],
      layout: undefined,
      layoutButtons: new Set(),
      blocks: [],
    }
  },

  computed: {
    ...mapGetters({
      getNextAndPrevRecord: 'ui/getNextAndPrevRecord',
      getPageLayouts: 'pageLayout/getByPageID',
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

      return actions
    },

    title () {
      const { name, handle } = this.module
      const titlePrefix = this.inCreating ? 'create' : this.inEditing ? 'edit' : 'view'

      return this.$t(`page:public.record.${titlePrefix}.title`, { name: name || handle, interpolation: { escapeValue: false } })
    },

    recordNavigation () {
      const { recordID } = this.record || {}
      return this.getNextAndPrevRecord(recordID)
    },
  },

  watch: {
    recordID: {
      immediate: true,
      handler () {
        this.record = undefined
        this.loadRecord()
      },
    },

    'page.pageID': {
      immediate: true,
      handler () {
        this.layouts = this.getPageLayouts(this.page.pageID)
        this.layout = undefined
        this.blocks = []

        this.determineLayout()
      },
    },
  },

  created () {
    this.$root.$on('refetch-record-blocks', () => {
      // Don*t refresh when creating and prompt user before refreshing when editing
      if (this.inCreating || (this.inEditing && !window.confirm(this.$t('notification:record.staleDataRefresh')))) {
        return
      }

      // Refetch the record and other page blocks that use records
      this.loadRecord()
      this.$root.$emit(`refetch-non-record-blocks:${this.page.pageID}`)
    })
  },

  // Destroy event before route leave to ensure it doesn't destroy the newly created one
  beforeRouteLeave (to, from, next) {
    this.$root.$off('refetch-record-blocks')
    next()
  },

  methods: {
    async loadRecord () {
      this.record = undefined

      if (!this.page) {
        return
      }

      const { namespaceID, moduleID } = this.page

      if (moduleID !== NoID) {
        const module = Object.freeze(this.getModuleByID(moduleID).clone())

        if (this.recordID && this.recordID !== NoID) {
          await this.$ComposeAPI
            .recordRead({ namespaceID, moduleID, recordID: this.recordID })
            .then(record => {
              setTimeout(() => {
                this.record = new compose.Record(module, record)
              }, 300)
            })
            .catch(this.toastErrorHandler(this.$t('notification:record.loadFailed')))
        } else {
          this.record = new compose.Record(module, {})
        }
      }
    },

    handleBack () {
      /**
       * Not the best way since we can not always know where we
       * came from (and "were" is back).
       */
      if (this.showRecordModal) {
        this.inEditing = false
        this.inCreating = false
        this.$bvModal.hide('record-modal')
      } else {
        this.$router.back()
      }
    },

    handleAdd () {
      if (this.showRecordModal) {
        this.inEditing = true
        this.inCreating = true
        this.record = new compose.Record(this.module, { values: this.values })
      } else {
        this.$router.push({ name: 'page.record.create', params: this.newRouteParams })
      }
    },

    handleClone () {
      if (this.showRecordModal) {
        this.inEditing = true
        this.inCreating = true
        this.record = new compose.Record(this.module, { values: this.record.values })
      } else {
        this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, values: this.record.values } })
      }
    },

    handleEdit () {
      if (this.showRecordModal) {
        this.inCreating = false
        this.inEditing = true
      } else {
        this.$router.push({ name: 'page.record.edit', params: this.$route.params })
      }
    },

    handleRedirectToPrevOrNext (recordID) {
      if (!recordID) return

      if (this.showRecordModal) {
        this.$router.push({
          query: { ...this.$route.query, recordID },
        })
      } else {
        this.$router.push({
          params: { ...this.$route.params, recordID },
        })
      }
    },

    evaluateLayoutExpressions () {
      const expressions = {}
      const variables = {
        user: this.$auth.user,
        record: this.record || {},
        screen: {
          width: window.innerWidth,
          height: window.innerHeight,
          userAgent: navigator.userAgent,
          breakpoint: this.getBreakpoint(), // This is from a global mixin uiHelpers
        },
        layout: this.layout || {},
      }

      this.layouts.forEach(({ pageLayoutID, config }) => {
        if (!config.visibility.expression) return

        expressions[pageLayoutID] = config.visibility.expression
      })

      return this.$SystemAPI.expressionEvaluate({ variables, expressions }).catch(() => {
        Object.keys(expressions).forEach(key => (expressions[key] = false))
        return expressions
      })
    },

    async determineLayout (pageLayoutID) {
      const expressions = await this.evaluateLayoutExpressions()

      // Check layouts for expressions/roles and find the first one that fits
      this.layout = this.layouts.find(l => {
        if (pageLayoutID && l.pageLayoutID !== pageLayoutID) return

        const { expression, roles = [] } = l.config.visibility

        if (expression && !expressions[l.pageLayoutID]) return false

        if (!roles.length) return true

        return this.$auth.user.roles.some(roleID => roles.includes(roleID))
      })

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
  },
}
</script>

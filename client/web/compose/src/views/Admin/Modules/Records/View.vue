<template>
  <div
    class="overflow-auto p-2"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="modulePage || allRecords"
        size="sm"
      >
        <b-button
          variant="primary"
          :disabled="!modulePage"
          :to="modulePage"
          class="d-flex align-items-center"
        >
          {{ $t('edit.edit') }}
          <font-awesome-icon
            :icon="['far', 'edit']"
            size="sm"
            class="ml-2"
          />
        </b-button>

        <b-button
          v-if="allRecords"
          variant="primary"
          :disabled="!allRecords"
          :to="allRecords"
          style="margin-left:2px;"
          class="d-flex align-items-center"
        >
          <font-awesome-icon
            :icon="['fas', 'columns']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <b-alert
      v-if="isDeleted"
      show
      variant="warning"
      class="mb-2 mx-2"
    >
      {{ $t('block.record.recordDeleted') }}
    </b-alert>

    <b-row
      v-if="module"
      no-gutters
    >
      <b-col
        v-for="(block, index) in blocks"
        :key="index"
        cols="12"
        lg="3"
        style="max-height: 650px; height: 650px;"
      >
        <component
          :is="getRecordComponent"
          :errors="errors"
          v-bind="{ namespace, page, module, block, record }"
          class="p-2"
        />
      </b-col>
    </b-row>

    <portal to="admin-toolbar">
      <record-toolbar
        :module="module"
        :record="record"
        :processing="processing"
        :processing-submit="processingSubmit"
        :processing-delete="processingDelete"
        :processing-undelete="processingUndelete"
        :in-editing="inEditing"
        :record-navigation="recordNavigation"
        :hide-back="false"
        :hide-delete="false"
        :hide-new="false"
        :hide-clone="false"
        :hide-edit="false"
        :hide-submit="false"
        @add="handleAdd()"
        @clone="handleClone()"
        @edit="handleEdit()"
        @view="handleView()"
        @delete="handleDelete()"
        @undelete="handleUndelete()"
        @back="handleBack()"
        @submit="handleFormSubmitSimple('admin.modules.record.view')"
        @update-navigation="handleRedirectToPrevOrNext"
      />
    </portal>
  </div>
</template>

<script>
import axios from 'axios'
import { isEqual } from 'lodash'
import { mapGetters } from 'vuex'
import RecordToolbar from 'corteza-webapp-compose/src/components/Common/RecordToolbar'
import record from 'corteza-webapp-compose/src/mixins/record'
import { compose, NoID } from '@cortezaproject/corteza-js'
import RecordBase from 'corteza-webapp-compose/src/components/PageBlocks/RecordBase'
import RecordEditor from 'corteza-webapp-compose/src/components/PageBlocks/RecordEditor'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  name: 'ViewRecord',

  components: {
    RecordToolbar,
    RecordBase,
    RecordEditor,
  },

  mixins: [
    // The record mixin contains all of the logic for creating/editing/deleting the record
    record,
  ],

  beforeRouteLeave (to, from, next) {
    next(this.checkUnsavedChanges())
  },

  beforeRouteUpdate (to, from, next) {
    next(this.checkUnsavedChanges())
  },

  props: {
    namespace: {
      type: Object,
      required: false,
      default: undefined,
    },

    moduleID: {
      type: String,
      required: false,
      default: '',
    },

    recordID: {
      type: String,
      required: false,
      default: '',
    },

    edit: {
      type: Boolean,
      default: false,
    },

    // If component was called (via router) with some pre-seed values
    values: {
      type: Object,
      required: false,
      default: () => ({}),
    },
  },

  data () {
    return {
      blocks: [],

      page: new compose.Page(),

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
    }),

    isNew () {
      return !this.recordID || this.recordID === NoID
    },

    title () {
      const { name, handle } = this.module
      const titlePrefix = this.isNew ? 'create' : this.inEditing ? 'edit' : 'view'

      return this.$t(`allRecords.${titlePrefix}.title`, { name: name || handle, interpolation: { escapeValue: false } })
    },

    module () {
      if (this.moduleID) {
        return this.getModuleByID(this.moduleID)
      } else {
        return undefined
      }
    },

    fields () {
      if (!this.module) {
        // No module, no fields
        return []
      }

      const fields = []
      const fieldSetSize = 8

      let i, j
      for (i = 0, j = this.module.fields.length; i < j; i += fieldSetSize) {
        fields.push(this.module.fields.slice(i, i + fieldSetSize))
      }

      fields.push(this.module.systemFields())

      return fields
    },

    getUiEventResourceType () {
      return 'admin-record-page'
    },

    getRecordComponent () {
      return this.inEditing ? RecordEditor : RecordBase
    },

    modulePage () {
      if (this.module) {
        return { name: 'admin.modules.edit', params: { moduleID: this.module.moduleID }, query: null }
      }

      return undefined
    },

    allRecords () {
      if (this.module.moduleID) {
        return { name: 'admin.modules.record.list', params: { moduleID: this.module.moduleID } }
      }

      return undefined
    },

    currentRecordNavigation () {
      const { recordID } = this.record || {}
      return this.getNextAndPrevRecord(recordID)
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

    edit: {
      immediate: true,
      handler (edit) {
        this.inEditing = edit
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

  created () {
    this.createBlocks()
  },

  beforeDestroy () {
    this.abortRequests()
    this.setDefaultValues()
  },

  methods: {
    createBlocks () {
      this.fields.forEach(f => {
        const options = {
          moduleID: this.moduleID,
          fields: f,
        }
        this.blocks.push(new compose.PageBlockRecord({ options }))
      })
    },

    refresh () {
      return this.loadRecord()
    },

    loadRecord () {
      const { moduleID = NoID, recordID = NoID } = this

      if (!moduleID || moduleID === NoID) return
      const module = Object.freeze(this.getModuleByID(moduleID).clone())

      if (recordID && recordID !== NoID) {
        const { namespaceID } = this.namespace

        const { response, cancel } = this.$ComposeAPI
          .recordReadCancellable({ namespaceID, moduleID, recordID })

        this.abortableRequests.push(cancel)

        response()
          .then(record => {
            this.record = new compose.Record(module, record)
            this.initialRecordState = this.record.clone()
          })
          .catch((e) => {
            if (!axios.isCancel(e)) {
              this.toastErrorHandler(this.$t('notification:record.loadFailed'))(e)
            }
          })
      } else {
        this.record = new compose.Record(module, { values: this.values })
        this.initialRecordState = undefined
      }
    },

    handleBack () {
      this.$router.push({ name: 'admin.modules.record.list', params: { moduleID: this.module.moduleID } })
    },

    handleAdd () {
      this.$router.push({ name: 'admin.modules.record.create', params: { moduleID: this.module.moduleID, edit: true } })
    },

    handleClone () {
      this.$router.push({ name: 'admin.modules.record.create', params: { moduleID: this.module.moduleID, values: this.record.values, edit: true } })
    },

    handleEdit () {
      this.$router.push({ name: 'admin.modules.record.edit', params: { moduleID: this.module.moduleID, edit: true } })
    },

    handleView () {
      this.$router.push({ name: 'admin.modules.record.view', params: { moduleID: this.module.moduleID, edit: false } })
    },

    handleRedirectToPrevOrNext (recordID) {
      if (!recordID) return

      this.$router.push({
        params: { ...this.$route.params, recordID },
      })
    },

    setDefaultValues () {
      this.blocks = []
      this.bindParams = {}
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    compareRecordValues () {
      const recordValues = JSON.parse(JSON.stringify(this.record ? this.record.values : {}))
      const initialRecordState = JSON.parse(JSON.stringify(this.initialRecordState ? this.initialRecordState.values : {}))

      return !isEqual(recordValues, initialRecordState)
    },

    checkUnsavedChanges () {
      if (!this.edit) return true

      const recordStateChange = this.compareRecordValues() ? window.confirm(this.$t('general:record.unsavedChanges')) : true

      if (!recordStateChange) {
        this.processing = false
      } else {
        this.record = this.initialRecordState ? this.initialRecordState.clone() : undefined
      }

      return recordStateChange
    },
  },
}
</script>

<style lang="scss">
.value {
  min-height: 1.2rem;
}
</style>

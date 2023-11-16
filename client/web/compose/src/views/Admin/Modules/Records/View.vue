<template>
  <div
    class="overflow-auto px-2"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="modulePage || allRecords"
        size="sm"
        class="mr-1"
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
      variant="info"
      class="mb-0"
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
        md="3"
        cols="12"
        style="max-height: 650px; height: 650px;"
      >
        <component
          :is="getRecordComponent"
          :errors="errors"
          v-bind="{ ...bindParams, module, block, record }"
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

  props: {
    // If component was called (via router) with some pre-seed values
    values: {
      type: Object,
      required: false,
      default: () => ({}),
    },
  },

  data () {
    return {
      inEditing: false,
      inCreating: false,

      blocks: [],

      bindParams: {
        page: new compose.Page(),
        namespace: this.$attrs.namespace,
      },

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
      return this.record.recordID === NoID
    },

    title () {
      const { name, handle } = this.module
      const titlePrefix = this.inCreating ? 'create' : this.inEditing ? 'edit' : 'view'

      return this.$t(`allRecords.${titlePrefix}.title`, { name: name || handle, interpolation: { escapeValue: false } })
    },

    module () {
      if (this.$attrs.moduleID) {
        return this.getModuleByID(this.$attrs.moduleID)
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
    '$attrs.recordID': {
      immediate: true,
      handler () {
        this.record = undefined
        this.initialRecordState = undefined
        this.loadRecord()
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

  beforeRouteLeave (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  beforeRouteUpdate (to, from, next) {
    this.checkUnsavedChanges(next, to)
  },

  methods: {
    createBlocks () {
      this.fields.forEach(f => {
        const options = {
          moduleID: this.$attrs.moduleID,
          fields: f,
        }
        this.blocks.push(new compose.PageBlockRecord({ options }))
      })
    },

    loadRecord () {
      const { moduleID = NoID, recordID = NoID } = this.$attrs

      if (!moduleID || moduleID === NoID) return
      const module = Object.freeze(this.getModuleByID(moduleID).clone())

      if (recordID && recordID !== NoID) {
        const { namespaceID } = this.$attrs.namespace

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
        this.inEditing = true
        this.inCreating = true
      }
    },

    handleBack () {
      this.$router.push({ name: 'admin.modules.record.list', params: { moduleID: this.module.moduleID } })
    },

    handleAdd () {
      this.$router.push({ name: 'admin.modules.record.create', params: { moduleID: this.module.moduleID } })
    },

    handleClone () {
      this.$router.push({ name: 'admin.modules.record.create', params: { moduleID: this.module.moduleID, values: this.record.values } })
    },

    handleEdit () {
      this.inEditing = true
      this.inCreating = false
    },

    handleView () {
      this.inEditing = false
      this.inCreating = false
    },

    handleRedirectToPrevOrNext (recordID) {
      if (!recordID) return

      this.$router.push({
        params: { ...this.$route.params, recordID },
      })
    },

    setDefaultValues () {
      this.inEditing = false
      this.blocks = []
      this.bindParams = {}
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },

    checkUnsavedChanges (next, to) {
      if (this.isNew) {
        next(true)
      } else {
        const recordValues = JSON.parse(JSON.stringify(this.record.values))
        const initialRecordState = JSON.parse(JSON.stringify(this.initialRecordState.values))

        next(!isEqual(recordValues, initialRecordState) ? window.confirm(this.$t('general:editor.unsavedChanges')) : true)
      }
    },
  },
}
</script>

<style lang="scss">
.value {
  min-height: 1.2rem;
}
</style>

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
      v-if="module && record"
      no-gutters
    >
      <b-col
        v-for="(block, index) in blocks"
        :key="index"
        md="3"
        cols="12"
      >
        <component
          :is="getRecordComponent"
          :errors="errors"
          v-bind="{ ...bindParams, module, block, record }"
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
        :is-deleted="isDeleted"
        :in-editing="inEditing"
        @add="handleAdd()"
        @clone="handleClone()"
        @edit="handleEdit()"
        @delete="handleDelete()"
        @back="handleBack()"
        @submit="handleFormSubmitSimple('admin.modules.record.view')"
      />
    </portal>
  </div>
</template>

<script>
import RecordToolbar from 'corteza-webapp-compose/src/components/Common/RecordToolbar'
import users from 'corteza-webapp-compose/src/mixins/users'
import record from 'corteza-webapp-compose/src/mixins/record'
import { compose } from '@cortezaproject/corteza-js'
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
    users,
  ],

  data () {
    return {
      blocks: [],
      bindParams: {
        page: new compose.Page(),
        boundingRect: {},
        namespace: this.$attrs.namespace,
      },
    }
  },
  computed: {
    title () {
      const { name, handle } = this.module
      return this.$t('allRecords.view.title', { name: name || handle, interpolation: { escapeValue: false } })
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
      const moduleFields = this.module.fields.slice().sort((a, b) => a.label.localeCompare(b.label))

      let i, j
      for (i = 0, j = moduleFields.length; i < j; i += fieldSetSize) {
        fields.push(moduleFields.slice(i, i + fieldSetSize))
      }

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
  },

  watch: {
    '$attrs.recordID': {
      immediate: true,
      handler () {
        this.loadRecord()
      },
    },
  },

  created () {
    this.createBlocks()
  },

  methods: {
    createBlocks () {
      this.fields.forEach(f => {
        const block = new compose.PageBlockRecord()
        const options = {
          moduleID: this.$attrs.moduleID,
          fields: f,
        }
        block.options = options
        this.blocks.push(block)
      })
    },

    loadRecord () {
      if (this.$attrs.recordID && this.$attrs.moduleID) {
        const { namespaceID } = this.$attrs.namespace
        const { moduleID, recordID } = this.$attrs
        const module = Object.freeze(this.getModuleByID(moduleID).clone())
        this.$ComposeAPI
          .recordRead({ namespaceID, moduleID, recordID })
          .then(record => {
            this.record = new compose.Record(module, record)
            this.fetchUsers(this.module.fields, [this.record])
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.loadFailed')))
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
      this.$router.push({ name: 'admin.modules.record.edit', params: this.$route.params })
    },
  },
}
</script>

<style lang="scss">
.value {
  min-height: 1.2rem;
}
</style>

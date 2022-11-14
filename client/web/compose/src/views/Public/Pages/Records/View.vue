<template>
  <div class="d-flex flex-grow-1 w-100">
    <b-alert
      v-if="isDeleted"
      show
      variant="warning"
    >
      {{ $t('record.recordDeleted') }}
    </b-alert>

    <portal to="topbar-title">
      {{ title }}
    </portal>

    <grid
      v-bind="$props"
      :errors="errors"
      :record="record"
      :mode="inEditing ? 'editor' : 'base'"
      @reload="loadRecord()"
    />

    <portal to="toolbar">
      <record-toolbar
        :module="module"
        :record="record"
        :labels="recordToolbarLabels"
        :processing="processing"
        :processing-submit="processingSubmit"
        :processing-delete="processingDelete"
        :is-deleted="isDeleted"
        :in-editing="inEditing"
        :hide-clone="inCreating"
        :hide-add="inCreating"
        @add="handleAdd()"
        @clone="handleClone()"
        @edit="handleEdit()"
        @delete="handleDelete()"
        @back="handleBack()"
        @submit="handleFormSubmit('page.record')"
      />
    </portal>
  </div>
</template>
<script>
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
    // The record mixin contains all of the logic for creating/editing/deleting the record
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
  },

  data () {
    return {
      inEditing: false,
      inCreating: false,
    }
  },

  computed: {
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
      const { buttons = {} } = this.page.config || {}
      Object.entries(buttons).forEach(([key, { label = '' }]) => {
        aux[key] = label
      })
      return aux
    },

    title () {
      const { name, handle } = this.module
      return this.$t('page:public.record.view.title', { name: name || handle, interpolation: { escapeValue: false } })
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
              this.record = new compose.Record(module, record)
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
      this.$router.back()
    },

    handleAdd () {
      this.$router.push({ name: 'page.record.create', params: this.newRouteParams })
    },

    handleClone () {
      this.$router.push({ name: 'page.record.create', params: { pageID: this.page.pageID, values: this.record.values } })
    },

    handleEdit () {
      this.$router.push({ name: 'page.record.edit', params: this.$route.params })
    },
  },
}
</script>

<template>
  <div
    class="h-100 px-2"
  >
    <portal to="topbar-title">
      {{ title }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="modulePage"
        size="sm"
        class="mr-1"
      >
        <b-button
          variant="primary"
          :disabled="!modulePage"
          :to="modulePage"
          style="margin-right:2px;"
          class="d-flex align-items-center"
        >
          {{ $t('edit.edit') }}
          <font-awesome-icon
            :icon="['far', 'edit']"
            size="sm"
            class="ml-2"
          />
        </b-button>
      </b-button-group>
    </portal>

    <record-list-base
      v-if="block && page"
      :block="block"
      :page="page"
      :module="module"
      :namespace="namespace"
      :block-index="0"
      class="p-2"
      @save-fields="handleFieldsSave"
    />
  </div>
</template>

<script>

import { mapGetters, mapActions } from 'vuex'
import { compose } from '@cortezaproject/corteza-js'
import RecordListBase from 'corteza-webapp-compose/src/components/PageBlocks/RecordListBase'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  components: {
    RecordListBase,
  },

  data () {
    return {
      block: undefined,
      namespace: this.$attrs.namespace,
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      recordPaginationUsable: 'ui/recordPaginationUsable',
    }),

    title () {
      const { name, handle } = this.module
      return this.$t('allRecords.list.title', { name: name || handle, interpolation: { escapeValue: false } })
    },

    module () {
      if (this.$route.params.moduleID) {
        return this.getModuleByID(this.$route.params.moduleID)
      } else {
        return undefined
      }
    },

    modulePage () {
      if (this.module) {
        return { name: 'admin.modules.edit', params: { moduleID: this.module.moduleID }, query: null }
      }

      return undefined
    },

    page () {
      if (!this.module) {
        return undefined
      }

      // Fake the pageID so record list uniqueID can be properly made
      const { moduleID } = this.module
      return new compose.Page({ pageID: moduleID })
    },
  },

  watch: {
    module: {
      handler (module) {
        if (module) {
          const { meta = { ui: {} }, moduleID } = module || {}

          let fields = ((meta.ui || {}).admin || {}).fields || []
          fields = fields.length ? fields : module.fields

          this.block.options.moduleID = moduleID
          this.block.options.fields = fields
        }
      },
    },
  },

  created () {
    const { meta = { ui: {} }, moduleID } = this.module || {}

    let fields = ((meta.ui || {}).admin || {}).fields || []
    fields = fields.length ? fields : this.module.fields

    // Init block
    const block = new compose.PageBlockRecordList({
      blockIndex: 0,
      options: {
        moduleID,
        fields,
        hideRecordReminderButton: true,
        hideRecordViewButton: true,
        hideRecordCloneButton: false,
        hideRecordPermissionsButton: false,
        selectable: true,
        allowExport: true,
        perPage: 14,
        fullPageNavigation: true,
        showTotalCount: true,
        showDeletedRecordsOption: true,
        presort: 'createdAt DESC',
        enableRecordPageNavigation: true,
        hideConfigureFieldsButton: false,
        inlineRecordEditEnabled: true,
        customFilterPresets: true,
      },
    })

    block.options = {
      ...block.options,
      allRecords: true,
      rowViewUrl: 'admin.modules.record.view',
      rowEditUrl: 'admin.modules.record.edit',
      rowCreateUrl: 'admin.modules.record.create',
    }

    this.block = block

    // If the page changed we need to clear the record pagination since its not relevant anymore
    if (this.recordPaginationUsable) {
      this.setRecordPaginationUsable(false)
    } else {
      this.clearRecordIDs()
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      updateModule: 'module/update',
      setRecordPaginationUsable: 'ui/setRecordPaginationUsable',
      clearRecordIDs: 'ui/clearRecordIDs',
    }),

    handleFieldsSave (fields = []) {
      fields = fields.map((f) => f.fieldID)

      if (!this.module.meta.ui) {
        this.module.meta.ui = { admin: { fields } }
      } else {
        this.module.meta.ui.admin = { ...(this.module.meta.ui.admin || {}), fields }
      }

      this.updateModule(this.module).then(() => {
        this.toastSuccess(this.$t('notification:module.columns.saved'))
      }).catch(this.toastErrorHandler(this.$t('notification:module.columns.saveFailed')))
    },

    setDefaultValues () {
      this.block = undefined
      this.namespace = {}
    },
  },
}
</script>

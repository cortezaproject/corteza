<template>
  <div v-if="namespace.canManageNamespace">
    <b-dropdown
      v-if="recordPage"
      :size="size"
      variant="primary"
      :text="$t('related-pages')"
      :boundary="boundary"
      :disabled="!module.name"
      class="related-pages-dropdown flex-fill"
    >
      <b-dropdown-item
        v-if="recordPage"
        data-test-id="dropdown-link-record-page-edit"
        :disabled="!namespace.canManageNamespace"
        :to="{ name: 'admin.pages.builder', params: { pageID: recordPage.pageID } }"
      >
        {{ $t('recordPage.edit') }}
      </b-dropdown-item>

      <b-dropdown-item
        v-if="recordListPage"
        :to="{ name: 'admin.pages.builder', params: { pageID: recordListPage.pageID } }"
      >
        {{ $t('recordListPage.edit') }}
      </b-dropdown-item>

      <b-dropdown-item-button
        v-else
        data-test-id="dropdown-link-record-list-page-create"
        :disabled="processing"
        @click.prevent.stop="handleRecordListPageCreation"
      >
        {{ $t('recordListPage.create') }}
      </b-dropdown-item-button>
    </b-dropdown>

    <b-button
      v-else
      data-test-id="button-record-page-create"
      variant="primary"
      :size="size"
      :disabled="processing || !module.name"
      @click.stop.prevent="handleRecordPageCreation"
    >
      {{ $t('recordPage.create') }}
    </b-button>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    size: {
      type: String,
      default: 'md',
    },

    boundary: {
      type: String,
      default: 'viewport',
    },
  },

  data () {
    return {
      processing: false,
    }
  },

  computed: {
    ...mapGetters({
      pages: 'page/set',
    }),

    recordPage () {
      return this.pages.find(p => p.moduleID === this.module.moduleID)
    },

    recordListPage () {
      return this.pages.find(p => {
        return p.blocks.find(b => b.kind === 'RecordList' && b.options.moduleID === this.module.moduleID)
      })
    },
  },

  methods: {
    ...mapActions({
      createPage: 'page/create',
      updatePage: 'page/update',
      createPageLayout: 'pageLayout/create',
    }),

    handleRecordPageCreation () {
      this.processing = true

      const { name, moduleID } = this.module
      const { namespaceID } = this.namespace

      // A simple record block w/o preselected fields
      const blocks = [new compose.PageBlockRecord({ xywh: [0, 0, 48, 82] })]
      const selfID = (this.recordListPage || {}).pageID || NoID

      const page = new compose.Page({
        namespaceID,
        moduleID,
        selfID,
        title: this.$t('forModule.recordPage', { name, interpolation: { escapeValue: false } }),
        blocks,
      })

      this.createPage(page).then(({ pageID, title, blocks }) => {
        const pageLayout = new compose.PageLayout({ namespaceID, pageID, handle: 'primary', blocks, meta: { title } })
        return this.createPageLayout(pageLayout)
      }).catch(this.toastErrorHandler(this.$t('notification:module.recordPage.createFailed')))
        .finally(() => {
          this.processing = false
        })
    },

    handleRecordListPageCreation () {
      this.processing = true

      const { namespaceID } = this.namespace
      const { name, moduleID } = this.module

      const blocks = [new compose.PageBlockRecordList({
        xywh: [0, 0, 48, 82],
        options: {
          moduleID,
          fields: [],
        },
      })]

      const page = new compose.Page({
        title: name,
        namespaceID,
        blocks,
        visible: true,
      })

      this.createPage(page)
        .then(({ pageID, title, blocks }) => {
          const pageLayout = new compose.PageLayout({ namespaceID, pageID, handle: 'primary', blocks, meta: { title } })
          return Promise.all([
            this.updatePage({ ...this.recordPage, selfID: pageID }),
            this.createPageLayout(pageLayout),
          ])
        })
        .catch(this.toastErrorHandler(this.$t('notification:module.recordPage.createFailed')))
        .finally(() => {
          this.processing = false
        })
    },
  },
}
</script>

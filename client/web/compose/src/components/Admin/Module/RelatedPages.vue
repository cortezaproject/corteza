<template>
  <div
    class="d-inline-block"
  >
    <b-dropdown
      v-if="recordPage"
      :size="size"
      variant="primary"
      :text="$t('related-pages')"
      :boundary="boundary"
      class="related-pages-dropdown"
    >
      <b-dropdown-item>
        <b-button
          data-test-id="dropdown-link-record-page-edit"
          :disabled="!namespace.canManageNamespace"
          :to="{ name: 'admin.pages.builder', params: { pageID: recordPage.pageID } }"
          variant="link"
          class="text-dark text-decoration-none"
        >
          {{ $t('recordPage.edit') }}
        </b-button>
      </b-dropdown-item>

      <b-dropdown-item>
        <b-button
          v-if="recordListPage"
          data-test-id="dropdown-link-record-list-page-edit"
          :disabled="!namespace.canManageNamespace"
          :to="{ name: 'admin.pages.builder', params: { pageID: recordListPage.pageID } }"
          variant="link"
          class="text-dark text-decoration-none"
        >
          {{ $t('recordListPage.edit') }}
        </b-button>
        <b-button
          v-else
          data-test-id="dropdown-link-record-list-page-create"
          variant="link"
          href=""
          :disabled="processing"
          class="text-dark text-decoration-none"
          @click.stop.prevent="handleRecordListPageCreation"
        >
          {{ $t('recordListPage.create') }}
        </b-button>
      </b-dropdown-item>
    </b-dropdown>

    <b-button
      v-else
      data-test-id="button-record-page-create"
      variant="primary"
      :size="size"
      :disabled="processing"
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
        title: `${this.$t('forModule.recordPage')} "${name || moduleID}"`,
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
          perPage: 14,
          fullPageNavigation: false,
          showTotalCount: false,
        },
      })]

      const page = new compose.Page({
        title: `${this.$t('forModule.recordList')} "${name || moduleID}"`,
        namespaceID,
        blocks,
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

<style lang="scss">
.related-pages-dropdown {
  .dropdown-item {
    padding: 0;
  }
}
</style>

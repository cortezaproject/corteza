<template>
  <div
    v-if="!!page"
    class="d-flex w-100 overflow-hidden"
  >
    <portal to="topbar-title">
      {{ pageTitle }}
    </portal>

    <portal to="topbar-tools">
      <b-button-group
        v-if="page && page.canUpdatePage"
        size="sm"
        class="mr-1"
      >
        <b-button
          data-test-id="button-page-builder"
          variant="primary"
          class="d-flex align-items-center"
          :to="pageBuilder"
        >
          {{ $t('general:label.pageBuilder') }}
          <font-awesome-icon
            :icon="['fas', 'cogs']"
            class="ml-2"
          />
        </b-button>
        <page-translator
          v-if="trPage"
          data-test-id="button-page-translations"
          :page.sync="trPage"
          style="margin-left:2px;"
        />
        <b-button
          data-test-id="button-page-edit"
          :title="$t('tooltip.edit.page')"
          :to="pageEditor"
          variant="primary"
          class="d-flex align-items-center"
          style="margin-left:2px;"
        >
          <font-awesome-icon
            :icon="['far', 'edit']"
          />
        </b-button>
      </b-button-group>
    </portal>

    <div
      class="flex-grow-1 overflow-auto d-flex px-2 w-100"
    >
      <router-view
        v-if="recordID || isRecordCreatePage"
        :namespace="namespace"
        :module="module"
        :page="page"
        class="flex-grow-1 overflow-auto d-flex flex-column"
      />

      <grid
        v-else
        :namespace="namespace"
        :module="module"
        :page="page"
      />
    </div>

    <attachment-modal />
  </div>
</template>
<script>
import { mapActions } from 'vuex'
import Grid from 'corteza-webapp-compose/src/components/Public/Page/Grid'
import AttachmentModal from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Modal'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  components: {
    Grid,
    AttachmentModal,
    PageTranslator,
  },

  props: {
    namespace: { // via router-view
      type: compose.Namespace,
      required: true,
    },

    page: { // via route-view
      type: compose.Page,
      required: true,
    },

    // We're using recordID to check if we need to display router-view or grid component
    recordID: {
      type: String,
      default: '',
    },
  },

  computed: {
    isRecordCreatePage () {
      return this.$route.name === 'page.record.create'
    },

    module () {
      if (this.page.moduleID && this.page.moduleID !== NoID) {
        return this.$store.getters['module/getByID'](this.page.moduleID)
      }

      return undefined
    },

    trPage: {
      get () {
        return this.page.clone()
      },
      set (v) {
        this.updatePageSet(v)
      },
    },

    pageTitle () {
      if (this.page.pageID !== NoID) {
        const { title = '', handle = '' } = this.page
        return title || handle || this.$t('navigation:noPageTitle')
      }

      return ''
    },

    pageEditor () {
      return { name: 'admin.pages.edit', params: { pageID: this.page.pageID } }
    },

    pageBuilder () {
      return { name: 'admin.pages.builder', params: { pageID: this.page.pageID } }
    },
  },

  watch: {
    'page.title': {
      immediate: true,
      handler (title) {
        document.title = [title, this.namespace.name, this.$t('general:label.app-name.public')].filter(v => v).join(' | ')
      },
    },
  },

  created () {
    this.$root.$on('refetch-records', () => {
      // If on a record page, let it take care of events else just refetch non record-blocks (that use records)
      this.$root.$emit(this.page.moduleID !== NoID ? 'refetch-record-blocks' : `refetch-non-record-blocks:${this.page.pageID}`)
    })
  },

  beforeDestroy () {
    this.$root.$off('refetch-records')
  },

  methods: {
    ...mapActions({
      updatePageSet: 'page/updateSet',
    }),
  },
}
</script>

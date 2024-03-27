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
        <page-translator
          v-if="trPage"
          data-test-id="button-page-translations"
          :page.sync="trPage"
          :page-layout.sync="layout"
          button-variant="primary"
          style="margin-left:2px;"
        />
        <b-button
          v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.edit.page'), container: '#body' }"
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
      </b-button-group>
    </portal>

    <div
      class="flex-grow-1 overflow-auto d-flex p-2 w-100"
    >
      <router-view
        v-if="isRecordPage"
        :namespace="namespace"
        :module="module"
        :page="page"
      />

      <div
        v-else-if="!layout"
        class="d-flex align-items-center justify-content-center w-100"
      >
        <b-spinner />
      </div>

      <grid
        v-else-if="blocks"
        :namespace="namespace"
        :module="module"
        :page="page"
        :blocks="blocks"
      />
    </div>

    <record-modal
      :namespace="namespace"
    />

    <magnification-modal
      :namespace="namespace"
    />
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import Grid from 'corteza-webapp-compose/src/components/Public/Page/Grid'
import RecordModal from 'corteza-webapp-compose/src/components/Public/Record/Modal'
import MagnificationModal from 'corteza-webapp-compose/src/components/Public/Page/Block/Modal'
import PageTranslator from 'corteza-webapp-compose/src/components/Admin/Page/PageTranslator'
import { NoID } from '@cortezaproject/corteza-js'
import page from 'corteza-webapp-compose/src/mixins/page'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  components: {
    Grid,
    RecordModal,
    PageTranslator,
    MagnificationModal,
  },

  mixins: [
    page,
  ],

  beforeRouteLeave (to, from, next) {
    this.setPreviousPages([])
    next()
  },

  beforeRouteUpdate (to, from, next) {
    const { recordID: toRecordID } = to.params
    const { recordID: fromRecordID } = from.params

    // Update either if coming from a record page and going to another record page and if the record isn't yet in previous pages to (avoid loop)
    const fromToRecordPage = fromRecordID && toRecordID !== fromRecordID
    // or if going from normal to record page
    const fromNormalToRecordPage = from.name === 'page' && to.name !== 'page'

    if (fromNormalToRecordPage || fromToRecordPage) {
      this.pushPreviousPages(from)
    }

    next()
  },

  data () {
    return {
      pageTitle: '',
    }
  },

  computed: {
    ...mapGetters({
      recordPaginationUsable: 'ui/recordPaginationUsable',
    }),

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

    pageEditor () {
      return { name: 'admin.pages.edit', params: { pageID: this.page.pageID } }
    },

    pageBuilder () {
      const { pageLayoutID } = this.layout || {}
      return { name: 'admin.pages.builder', params: { pageID: this.page.pageID }, query: { layoutID: pageLayoutID } }
    },
  },

  watch: {
    'page.pageID': {
      immediate: true,
      handler (pageID) {
        if (pageID === NoID) return

        this.layouts = this.getPageLayouts(this.page.pageID)
        this.layout = undefined
        this.pageTitle = this.page.title

        if (!this.isRecordPage) {
          this.determineLayout()
        } else {
          this.blocks = []
        }

        // If the page changed we need to clear the record pagination since its not relevant anymore
        if (this.recordPaginationUsable) {
          this.setRecordPaginationUsable(false)
        } else {
          this.clearRecordIDs()
        }
      },
    },
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      updatePageSet: 'page/updateSet',
      setRecordPaginationUsable: 'ui/setRecordPaginationUsable',
      clearRecordIDs: 'ui/clearRecordIDs',
      setPreviousPages: 'ui/setPreviousPages',
      pushPreviousPages: 'ui/pushPreviousPages',
      clearRecordSet: 'record/clearSet',
    }),

    createEvents () {
      this.$root.$on('refetch-records', this.refetchRecords)
    },

    refetchRecords () {
      // If on a record page, let it take care of events else just refetch non record-blocks (that use records)
      this.$root.$emit(this.page.moduleID !== NoID ? 'refetch-record-blocks' : `refetch-non-record-blocks:${this.page.pageID}`)
    },

    setDefaultValues () {
      this.layouts = []
      this.layout = undefined
      this.blocks = undefined
      this.pageTitle = ''
    },

    destroyEvents () {
      this.$root.$off('refetch-records', this.refetchRecords)
    },
  },
}
</script>

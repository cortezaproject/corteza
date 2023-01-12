<template>
  <b-modal
    v-model="showModal"
    scrollable
    body-class="card p-0"
    footer-class="p-0"
    :content-class="contentClass"
    :dialog-class="dialogClass"
    hide-header
    hide-footer
    size="xl"
    @hidden="hideModal"
  >
    <page-block
      v-if="showModal"
      :block="block"
      :record="record"
      :page="page"
      v-bind="$props"
      v-on="$listeners"
    />
  </b-modal>
</template>

<script>
import { mapGetters } from 'vuex'
import { compose } from '@cortezaproject/corteza-js'
import PageBlock from 'corteza-webapp-compose/src/components/PageBlocks'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'MagnificationModal',

  components: {
    PageBlock,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },
  },

  data () {
    return {
      showModal: false,
      block: undefined,
      record: undefined,
      page: undefined,
    }
  },

  computed: {
    ...mapGetters({
      getPageByID: 'page/getByID',
    }),

    dialogClass () {
      return this.block && this.block.options.magnifyOption === 'fullscreen' ? 'h-100 mw-100 m-0 mh-100' : 'h-100 mw-90'
    },

    contentClass () {
      return `${this.block && this.block.options.magnifyOption === 'fullscreen' ? 'mh-100' : ''} position-initial`
    },
  },

  watch: {
    '$route.query.blockID': {
      immediate: true,
      handler (blockID, oldBlockID) {
        if (!blockID) {
          this.showModal = false
          return
        }

        if (this.showModal && (blockID !== oldBlockID)) {
          this.showModal = false

          setTimeout(() => {
            this.$router.push({ query: { ...this.$route.query, blockID } })
          }, 300)

          return
        }

        setTimeout(() => {
          this.loadModal(blockID)
        }, 100)
      },
    },
  },

  created () {
    this.$root.$on('magnify-page-block', blockID => {
      this.$router.push({ query: { ...this.$route.query, blockID } })
    })
  },

  beforeDestroy () {
    this.$root.$off('magnify-page-block')
  },

  methods: {
    loadModal (blockID) {
      // Get data from route
      const { recordID: paramsRecordID, pageID } = this.$route.params
      const { recordID: queryRecordID, recordPageID } = this.$route.query

      // Get page that we should display
      this.page = this.getPageByID(recordPageID || pageID)

      if (!this.page) {
        return
      }

      this.block = this.page.blocks.find(block => block.blockID === blockID)
      const { namespaceID, moduleID } = this.page
      const recordID = paramsRecordID || queryRecordID

      if (recordID) {
        this.$ComposeAPI
          .recordRead({ namespaceID, moduleID, recordID })
          .then(record => {
            this.record = record
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.loadFailed')))
      }

      this.showModal = !!(this.block || {}).blockID
    },

    hideModal () {
      this.$router.push({ query: { ...this.$route.query, blockID: undefined } })
    },
  },
}

</script>

<style lang="scss">
.position-initial {
  position: initial;
}
.mw-90 {
  max-width: 90vw;
}
</style>

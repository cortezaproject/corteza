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
    @hidden="onHidden"
  >
    <page-block
      v-if="showModal"
      :block="block"
      :blocks="page.blocks"
      :module="module"
      :record="record"
      :page="page"
      magnified
      v-bind="$props"
      v-on="$listeners"
    />
  </b-modal>
</template>

<script>
import { mapGetters } from 'vuex'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { fetchID } from 'corteza-webapp-compose/src/lib/block'
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

      // Used if you want to display a specific block in the modal
      // Otherwise its retrieved based on the page and blockID
      customBlock: undefined,
    }
  },

  computed: {
    ...mapGetters({
      getPageByID: 'page/getByID',
      getModuleByID: 'module/getByID',
    }),

    dialogClass () {
      return this.block && this.block.options.magnifyOption === 'fullscreen' ? 'h-100 mw-100 m-0 mh-100' : 'h-100 mw-90'
    },

    contentClass () {
      return `${this.block && this.block.options.magnifyOption === 'fullscreen' ? 'mh-100 rounded-0' : 'card'} position-initial`
    },
  },

  mounted () {
    this.$root.$on('magnify-page-block', this.magnifyPageBlock)

    const { magnifiedBlockID } = this.$route.query

    if (magnifiedBlockID) {
      this.magnifyPageBlock({ blockID: magnifiedBlockID })
    }
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    magnifyPageBlock ({ blockID, block } = {}) {
      this.customBlock = block
      const magnifiedBlockID = blockID || (block || {}).blockID
      this.loadModal(magnifiedBlockID)

      setTimeout(() => {
        this.$router.push({
          query: {
            ...this.$route.query,
            magnifiedBlockID,
          },
        })
      }, 300)
    },

    loadModal (blockID) {
      // Get data from route
      const { recordID: paramsRecordID, pageID } = this.$route.params
      const { recordID: queryRecordID, recordPageID } = this.$route.query

      // Get page that we should display
      this.page = this.getPageByID(recordPageID || pageID)

      if (!this.page) {
        return
      }

      const { namespaceID, moduleID } = this.page
      const recordID = paramsRecordID || queryRecordID
      this.block = this.customBlock || this.page.blocks.find(block => fetchID(block) === blockID)
      this.module = moduleID !== NoID ? this.getModuleByID(moduleID) : undefined
      this.showModal = !!(this.block || {}).blockID

      if (recordID) {
        this.$ComposeAPI
          .recordRead({ namespaceID, moduleID, recordID })
          .then(record => {
            this.record = new compose.Record(this.module, record)
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.loadFailed')))
      } else if (this.module) {
        this.record = new compose.Record(this.module, {})
      }
    },

    onHidden () {
      setTimeout(() => {
        this.$router.replace({
          query: {
            ...this.$route.query,
            magnifiedBlockID: undefined,
          },
        })
      }, 300)
    },

    setDefaultValues () {
      this.showModal = false
      this.block = undefined
      this.record = undefined
      this.page = undefined
      this.customBlock = undefined
    },

    destroyEvents () {
      this.$root.$off('magnify-page-block', this.magnifyPageBlock)
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

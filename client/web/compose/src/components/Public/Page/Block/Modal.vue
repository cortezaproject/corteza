<template>
  <b-modal
    id="page-block-modal"
    v-model="showModal"
    scrollable
    body-class="p-0"
    :content-class="contentClass"
    :dialog-class="dialogClass"
    hide-header
    hide-footer
    size="xl"
    no-fade
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
import { mapGetters, mapActions } from 'vuex'
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
      return this.block && this.block.options.magnifyOption === 'fullscreen' ? 'h-100 mw-100 m-0 mh-100' : 'h-100 modal-max-width'
    },

    contentClass () {
      return `${this.block && this.block.options.magnifyOption === 'fullscreen' ? 'mh-100 rounded-0' : ''} position-initial`
    },

    magnifiedBlockID () {
      return this.$route.query.magnifiedBlockID
    },
  },

  watch: {
    magnifiedBlockID: {
      immediate: true,
      handler (magnifiedBlockID) {
        if (!magnifiedBlockID) {
          this.setDefaultValues()
        } else if (magnifiedBlockID !== fetchID(this.block)) {
          this.loadModal(magnifiedBlockID)
        }
      },
    },
  },

  mounted () {
    this.$root.$on('magnify-page-block', this.magnifyPageBlock)
  },

  beforeDestroy () {
    this.destroyEvents()
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      updateRecordSet: 'record/updateRecords',
    }),

    magnifyPageBlock ({ blockID, block } = {}) {
      this.customBlock = block
      const magnifiedBlockID = blockID || (block || {}).blockID
      this.loadModal(magnifiedBlockID)

      this.$router.push({
        query: {
          ...this.$route.query,
          magnifiedBlockID,
        },
      })
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
            this.updateRecordSet(this.record)
          })
          .catch(this.toastErrorHandler(this.$t('notification:record.loadFailed')))
      } else if (this.module) {
        this.record = new compose.Record(this.module, {})
      }
    },

    onHidden () {
      if (this.$route.query.magnifiedBlockID !== undefined) {
        this.$router.push({
          query: {
            ...this.$route.query,
            magnifiedBlockID: undefined,
          },
        })
      }
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

.modal-max-width {
  max-width: 97vw;
}
</style>

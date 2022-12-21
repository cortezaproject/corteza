<template>
  <b-modal
    v-model="showBlockModal"
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
      v-if="showBlockModal"
      :block="block"
      v-bind="$props"
      v-on="$listeners"
    />
  </b-modal>
</template>

<script>
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

    page: {
      type: compose.Page,
      required: true,
    },
  },

  data () {
    return {
      showBlockModal: false,
      block: undefined,
    }
  },

  computed: {
    dialogClass () {
      return this.block && this.block.options.magnifyOption === 'fullscreen' ? 'h-100 mw-100 m-0 mh-100' : 'h-100'
    },

    contentClass () {
      return `${this.block && this.block.options.magnifyOption === 'fullscreen' ? 'mh-100' : ''} position-initial`
    },
  },

  watch: {
    '$route.query.blockID': {
      immediate: true,
      handler (blockID, oldBlockID) {
        if (blockID && !oldBlockID) {
          this.showModal(blockID)
        } else if (!blockID && oldBlockID) {
          this.hideModal()
        }
      },
    },
  },

  created () {
    this.$root.$on('magnify-page-block', this.showModal)
  },

  beforeDestroy () {
    this.$root.$off('magnify-page-block')
  },

  methods: {
    showModal (blockID) {
      if (blockID) {
        this.block = this.page.blocks.find(block => block.blockID === blockID)
        this.showBlockModal = true
        this.$router.push({ query: { blockID } })
      } else {
        this.hideModal()
      }
    },

    hideModal () {
      this.showBlockModal = false
      this.$router.push({ query: {} })
      this.block = undefined
    },
  },
}

</script>

<style>
.position-initial {
  position: initial;
}
</style>

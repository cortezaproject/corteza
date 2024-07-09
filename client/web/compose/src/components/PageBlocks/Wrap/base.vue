<script>
import { compose, NoID } from '@cortezaproject/corteza-js'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    block: {
      type: compose.PageBlock,
      required: true,
    },

    record: {
      type: compose.Record,
      required: false,
      default: undefined,
    },

    scrollableBody: {
      type: Boolean,
      required: false,
      default: true,
    },

    cardClass: {
      type: String,
      required: false,
      default: '',
    },

    bodyClass: {
      type: String,
      required: false,
      default: '',
    },

    headerClass: {
      type: String,
      required: false,
      default: '',
    },

    magnified: {
      type: Boolean,
      required: false,
      default: false,
    },
  },

  computed: {
    blockID () {
      const { blockID, meta } = this.block || {}
      return meta.customID || blockID
    },

    customCSSClass () {
      const { meta } = this.block || {}
      return meta.customCSSClass
    },

    blockClass () {
      return [
        'block',
        { border: this.block.style.border.enabled },
        this.block.kind,
      ]
    },

    isBlockMagnified () {
      const { magnifiedBlockID } = this.$route.query
      return this.magnified && magnifiedBlockID === this.block.blockID
    },

    headerSet () {
      return !!this.$scopedSlots.header
    },

    toolbarSet () {
      return !!this.$scopedSlots.toolbar
    },

    footerSet () {
      return !!this.$scopedSlots.footer
    },

    showHeader () {
      return [
        this.headerSet,
        this.block.title,
        this.block.description,
        this.block.options.magnifyOption,
        this.block.options.showRefresh,
      ].some(c => !!c)
    },

    showOptions () {
      return [
        this.block.options.magnifyOption,
        this.block.options.showRefresh,
        this.isBlockMagnified,
      ].some(c => !!c)
    },

    magnifyParams () {
      const params = this.block.blockID === NoID ? { block: this.block } : { blockID: this.block.blockID }
      return this.isBlockMagnified ? undefined : params
    },

    blockTitle () {
      return evaluatePrefilter(this.block.title, {
        record: this.record,
        user: this.$auth.user || {},
        recordID: (this.record || {}).recordID || NoID,
        ownerID: (this.record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      })
    },

    blockDescription () {
      return evaluatePrefilter(this.block.description, {
        record: this.record,
        user: this.$auth.user || {},
        recordID: (this.record || {}).recordID || NoID,
        ownerID: (this.record || {}).ownedBy || NoID,
        userID: (this.$auth.user || {}).userID || NoID,
      })
    },
  },
}
</script>

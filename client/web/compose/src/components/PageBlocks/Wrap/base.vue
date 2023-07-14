<script>
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    block: {
      type: compose.PageBlock,
      required: true,
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

    magnified: {
      type: Boolean,
      required: false,
      default: false,
    },
  },

  computed: {
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
      ].some(c => !!c)
    },

    magnifyParams () {
      const params = this.block.blockID === NoID ? { block: this.block } : { blockID: this.block.blockID }
      return this.isBlockMagnified ? undefined : params
    },
  },
}
</script>

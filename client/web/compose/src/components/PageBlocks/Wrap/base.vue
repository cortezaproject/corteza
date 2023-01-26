<script>
import { compose } from '@cortezaproject/corteza-js'
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
      default: () => true,
    },
  },

  computed: {
    blockClass () {
      return [
        'block',
        this.block.kind,
      ]
    },

    isBlockOpened () {
      return this.block.blockID === this.$route.query.blockID
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
  },
}
</script>

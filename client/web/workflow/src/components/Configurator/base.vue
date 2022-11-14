<script>
export default {
  props: {
    item: {
      type: Object,
      default: () => {},
    },

    edges: {
      type: Object,
      default: () => {},
    },

    outEdges: {
      type: Number,
      default: 0,
    },
  },

  computed: {
    // getSelectedItemConfigJSON: {
    //   get () {
    //     let config = {}
    //     if (this.item) {
    //       config = this.item.config
    //     }
    //     return JSON.stringify(config, undefined, 2)
    //   },

    //   set (value) {
    //     this.item.config = JSON.parse(value)
    //   }
    // },
  },
}
</script>

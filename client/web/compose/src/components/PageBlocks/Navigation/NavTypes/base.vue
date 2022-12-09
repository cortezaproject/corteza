<script>
export default {
  i18nOptions: {
    namespaces: 'block',
  },

  props: {
    item: {
      type: Object,
      required: true,
    },
  },

  computed: {
    options: {
      get () {
        return this.item.options
      },
      set (options) {
        this.item.options = options
      },
    },
  },
}
</script>

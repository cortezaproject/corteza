<script>
export default {
  props: {
    displayElementOptions: {
      type: Object,
      default: () => ({}),
    },

    columns: {
      type: Array,
      required: true,
    },

    datasource: {
      type: Object,
      default: undefined,
    },
  },

  computed: {
    options: {
      get () {
        return this.displayElementOptions || undefined
      },

      set (options = {}) {
        this.$emit('update:displayElementOptions', options)
      },
    },
  },
}
</script>

<script>
export default {
  props: {
    displayElement: {
      type: Object,
      required: true,
    },

    labels: {
      type: Object,
      default: () => ({})
    },
  },

  computed: {
    dataframes () {
      return this.displayElement.dataframes || []
    },

    options () {
      return this.displayElement.options || undefined
    },
  },
}
</script>

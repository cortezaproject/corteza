<script>
/**
 * Provides base for each feed configurator.
 * Define common props, methods, ... in here.
 */
export default {
  filters: {
    /**
     * Prepares a set of select options
     * @param {Array} ff Raw options
     * @returns {Array}
     */
    optionizeFields (ff) {
      return ff.map(f => {
        return {
          text: f.label || f.name,
          value: f.name,
        }
      })
    },
  },

  props: {
    feed: {
      type: Object,
      required: true,
      default: () => ({}),
    },
    modules: {
      type: Array,
      required: false,
      default: () => [],
    },
  },
}
</script>

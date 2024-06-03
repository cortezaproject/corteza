<script>
/**
 * Provides base for each feed configurator.
 * Define common props, methods, ... in here.
 */
export default {
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

<script>
/**
 * Defines common props, methods between different toolbar item types
 */
export default {
  props: {
    format: {
      type: Object,
      required: true,
      default: () => ({}),
    },
    isActive: {
      type: Object,
      required: false,
      default: () => ({}),
    },
    getMarkAttrs: {
      type: Function,
      required: false,
      default: () => {},
    },
    currentValue: {
      type: String,
      required: false,
    },
  },
  methods: {
    /**
     * Helper method to emit format selection
     * @param {String} type Format's type
     * @param {Object} attrs Format's extra attributes
     */
    onClick (type, attrs) {
      this.$emit('click', { type, attrs })
    },

    /**
     * Helper method to determine if given format is active or not.
     * When attrs is provided, it will check for that exact match
     * @param {Object|undefined} attrs Format's extra attributes
     * @returns {Array|undefined}
     */
    activeClasses (attrs) {
      if ((this.isActive[this.format.type])(attrs)) {
        return ['text-success']
      }
      return undefined
    },
  },
}
</script>

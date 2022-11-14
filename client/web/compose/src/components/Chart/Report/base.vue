<script>
export default {
  props: {
    report: {
      type: Object,
      required: false,
      default: undefined,
    },

    modules: {
      type: Array,
      required: true,
    },

    supportedMetrics: {
      type: Number,
      required: false,
      default: -1,
    },

    // Specifies what field kinds are supported as a dimension field
    dimensionFieldKind: {
      type: Array,
      required: false,
      default: () => [
        'DateTime',
        'Select',
        'Number',
        'Bool',
        'String',
        'Record',
        'User',
      ],
    },

    usesDimensionsField: {
      type: Boolean,
      default: true,
    },

    unSkippable: {
      type: Boolean,
      required: false,
      default: false,
    },
  },

  computed: {
    editReport: {
      get () {
        return this.report
      },
      set (v) {
        this.$emit('update:report', v)
      },
    },
  },
}
</script>

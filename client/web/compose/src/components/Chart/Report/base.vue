<script>
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  props: {
    report: {
      type: Object,
      required: false,
      default: undefined,
    },

    chart: {
      type: compose.Chart,
      default: () => new compose.Chart(),
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

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
    }
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

    isNew () {
      return this.chart.chartID === NoID
    },
  },
}
</script>

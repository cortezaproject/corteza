<script>
import { pType, pVal, variantToSeverity } from '../utils'

export default {
  props: {
    loading: {
      type: Boolean,
      default: false,
    },

    payload: {
      type: Object,
      default: () => ({}),
    },
  },

  computed: {
    message() {
      return this.pVal('message', '')
    },

    label() {
      return this.pVal('label', '')
    },
  },

  methods: {
    pVal(k, def = undefined) {
      return pVal(this.payload, k, def)
    },

    pType(k, def = undefined) {
      return pType(this.payload, k, def)
    },

    pRaw(k, defValue = undefined, defType = undefined) {
      if (k && this.payload && this.payload[k] !== undefined) {
        return this.payload[k]
      }

      return { '@type': defType, '@value': defValue }
    },

    tF(key, fallback) {
      return this.$te(key) ? this.$t(key) : fallback
    },

    vS(k, def = 'primary') {
      return variantToSeverity(this.pVal(k, def))
    },
  },
}
</script>

<script lang="js">
import { pVal, pType } from '../utils.ts'

export default {
  props: {
    loading: {
      type: Boolean,
      default: () => false,
    },

    payload: {
      type: Object,
    },
  },

  computed: {
    message () {
      return this.pVal('message', '<i>Default prompt message.</i>')
    },
  },

  methods: {
    pVal (k, def = undefined) {
      return pVal(this.payload, k, def)
    },

    pType (k, def = undefined) {
      return pType(this.payload, k, def)
    },

    pRaw (k, defValue = undefined, defType = undefined) {
      if (this.payload && this.payload[k] && this.payload[k] !== undefined) {
        return this.payload[k]
      }

      return { '@type': defType, '@value': defValue}
    }
  },
}
</script>

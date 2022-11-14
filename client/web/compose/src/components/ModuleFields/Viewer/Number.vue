<script>
import base from './base'

export default {
  extends: base,

  computed: {
    formatted () {
      if (this.value) {
        if (this.field.isMulti) {
          return this.value.map(v => this.field.formatValue(v)).join(this.field.options.multiDelimiter)
        } else {
          return this.field.formatValue(this.value)
        }
      }
      return undefined
    },
  },
}
</script>

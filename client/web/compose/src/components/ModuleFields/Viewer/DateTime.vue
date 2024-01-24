<script>
import base from './base'

export default {
  extends: base,
  computed: {
    formatted () {
      if (this.field.isMulti && this.value) {
        return this.value.map(v => this.field.formatValue(v)).join(this.field.options.multiDelimiter)
      }
      return this.value ? this.field.formatValue(this.value) : null
    },
  },
}
</script>

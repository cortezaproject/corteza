<script>
import base from './base'

export default {
  extends: base,

  computed: {
    formatted () {
      const { trueLabel, falseLabel } = this.field.options

      if (this.value === '1') {
        return trueLabel || this.$t('general:label.yes')
      } else if (this.value === '0') {
        return falseLabel || this.$t('general:label.no')
      }

      return ''
    },
  },
}
</script>

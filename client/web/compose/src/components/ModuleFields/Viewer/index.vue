<script>
import base from './base'
import * as Viewers from './loader'

// Renders one of the field kind components
export default {
  i18nOptions: {
    namespaces: 'notification',
  },

  extends: base,

  computed: {
    component () {
      const kind = this.field.kind.toLocaleLowerCase()
      const keys = Object.keys(Viewers)
      const i = keys.map(c => c.toLocaleLowerCase()).findIndex(c => c === kind)

      if (i >= 0) {
        return Viewers[keys[i]]
      } else {
        return null
      }
    },
  },

  render (createElement) {
    const cmp = this.component
    if (cmp) {
      return createElement(cmp, {
        props: this.$props,
      })
    } else {
      return createElement('code', this.$t('field.unknownFieldKind', { kind: this.field.kind }))
    }
  },
}
</script>

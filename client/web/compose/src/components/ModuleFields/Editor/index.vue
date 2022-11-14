<script>
import base from './base'
import * as Editors from './loader'

export default {
  i18nOptions: {
    namespaces: 'notification',
  },

  components: {
    ...Editors,
  },

  extends: base,

  computed: {
    component () {
      const kind = this.field.kind.toLocaleLowerCase()
      const keys = Object.keys(this.$options.components)
      const i = keys.map(c => c.toLocaleLowerCase()).findIndex(c => c === kind)

      if (i >= 0) {
        return Editors[keys[i]]
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
        on: this.$listeners,
      })
    } else {
      return createElement('code', this.$t('field.unknownFieldKind', { kind: this.field.kind }))
    }
  },
}
</script>

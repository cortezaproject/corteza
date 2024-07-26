<script>
import base from './base'
import { mapGetters } from 'vuex'
import { NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  computed: {
    ...mapGetters({
      findByID: 'user/findByID',
    }),

    formatted () {
      const f = (u) => {
        if (u === NoID) {
          return
        }

        if (typeof u === 'string') {
          u = this.findByID(u)
        }

        if (!u) {
          return
        }

        return this.field.formatter(u)
      }

      if (this.field.isMulti) {
        return this.value.map(v => f(v) || this.$t('kind.user.na')).join(this.field.options.multiDelimiter)
      }

      return f(this.value) || this.$t('kind.user.na')
    },
  },

  created () {
    // Prefill value with current user only when creating a new record and field is ownedBy (since BE prefills it anyway if its undefined)

    const isNewRecord = this.record && this.record.recordID === NoID

    if (isNewRecord && this.field.name === 'ownedBy') {
      const { userID } = this.$auth.user
      this.record.ownedBy = userID
    }
  },
}
</script>

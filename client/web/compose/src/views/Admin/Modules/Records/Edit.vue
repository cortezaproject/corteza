<script>
import ViewRecord from './View'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  name: 'EditRecord',

  extends: ViewRecord,

  data () {
    return {
      inEditing: true,
    }
  },

  computed: {
    title () {
      const { name, handle } = this.module
      return this.$t('allRecords.edit.title', { name: name || handle, interpolation: { escapeValue: false } })
    },
  },

  methods: {
    handleBack () {
      this.$router.push({ name: 'admin.modules.record.view' })
    },
  },
}
</script>

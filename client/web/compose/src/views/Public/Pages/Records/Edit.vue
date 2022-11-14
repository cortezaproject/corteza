<script>
import ViewRecord from './View'

export default {
  name: 'EditRecord',

  extends: ViewRecord,

  data () {
    return {
      inEditing: true,
      inCreating: false,
    }
  },

  computed: {
    title () {
      const { name, handle } = this.module
      return this.$t('page:public.record.edit.title', { name: name || handle, interpolation: { escapeValue: false } })
    },
  },

  methods: {
    handleBack () {
      this.$router.push({ name: 'page.record', params: { recordID: this.recordID, pageID: this.page.pageID } })
    },
  },
}
</script>

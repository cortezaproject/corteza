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

  methods: {
    handleBack () {
      this.$router.push({ name: 'page.record', params: { recordID: this.recordID, pageID: this.page.pageID } })
    },
  },
}
</script>

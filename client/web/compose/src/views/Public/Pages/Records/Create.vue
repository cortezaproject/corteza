<script>
import ViewRecord from './View'
export default {
  name: 'CreateRecord',

  extends: ViewRecord,

  data () {
    return {
      inEditing: true,
      inCreating: true,
    }
  },
}
</script>

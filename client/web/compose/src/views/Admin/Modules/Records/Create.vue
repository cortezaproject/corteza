<script>
import ViewRecord from './View'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

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

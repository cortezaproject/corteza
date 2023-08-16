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
      inCreating: false,
    }
  },
}
</script>

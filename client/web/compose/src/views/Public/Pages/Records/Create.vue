<script>
import ViewRecord from './View'
import { compose } from '@cortezaproject/corteza-js'

export default {
  name: 'CreateRecord',

  extends: ViewRecord,

  props: {
    // When creating from related record blocks
    refRecord: {
      type: compose.Record,
      required: false,
      default: () => ({}),
    },

    // If component was called (via router) with some pre-seed values
    values: {
      type: Object,
      required: false,
      default: () => ({}),
    },
  },

  data () {
    return {
      inEditing: true,
      inCreating: true,
    }
  },

  created () {
    this.record = new compose.Record(this.module, { values: this.values })

    if (this.refRecord) {
      // Record create form called from a related records block,
      // we'll try to find an appropriate field and cross-link this new record to ref
      const recRefField = this.module.fields.find(f => f.kind === 'Record' && f.options.moduleID === this.refRecord.moduleID)
      if (recRefField) {
        this.record.values[recRefField.name] = this.refRecord.recordID
      }
    }
  },
}
</script>

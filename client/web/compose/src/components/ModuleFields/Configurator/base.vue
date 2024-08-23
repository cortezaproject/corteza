<script>
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    field: {
      type: compose.ModuleField,
      required: true,
    },

    hasRecords: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    f: {
      get () {
        return this.field
      },

      set (f) {
        this.$emit('update:field', f)
      },
    },

    isNew () {
      return this.module.moduleID === NoID || this.field.fieldID === NoID
    },

    hasData () {
      return !this.isNew && this.hasRecords
    },
  },
}
</script>

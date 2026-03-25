<script>
import { inject } from 'vue'

export default {
  props: {
    item: {
      type: Object,
      default: () => {},
    },

    edges: {
      type: Object,
      default: () => {},
    },

    outEdges: {
      type: Number,
      default: 0,
    },

    isSubworkflow: {
      type: Boolean,
      default: false,
    },
  },

  setup () {
    const $toast = inject('$toast')
    return {
      toastSuccess: $toast.toastSuccess,
      toastWarning: $toast.toastWarning,
      toastInfo: $toast.toastInfo,
      toastDanger: $toast.toastDanger,
      toastErrorHandler: $toast.toastErrorHandler,
    }
  },
}
</script>

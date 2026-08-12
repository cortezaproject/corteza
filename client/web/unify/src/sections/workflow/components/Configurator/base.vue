<script>
import { useToast } from 'primevue/usetoast'

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

  setup() {
    const toast = useToast()
    return { toast }
  },
}
</script>

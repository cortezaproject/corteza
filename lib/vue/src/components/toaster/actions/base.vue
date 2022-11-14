<script lang="js">
export default {
  props: {
    options: {
      type: Object,
      required: false,
      default: () => ({}),
    },

    label: {
      type: String,
      required: false,
      default: undefined,
    },
  },
}
</script>

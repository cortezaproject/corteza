<script>
export default {
  props: {
    componentID: {
      type: String,
      required: true,
    },

    accessToken: {
      type: String,
      default: undefined,
    },
  },
}
</script>

<script>
export default {
  props: {
    step: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

  computed: {
    kind () {
      return Object.keys(this.step)
    },
  },
}
</script>

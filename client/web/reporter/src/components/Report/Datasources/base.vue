<script>
export default {
  props: {
    index: {
      type: Number,
      required: true,
    },

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

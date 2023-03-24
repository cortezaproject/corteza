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

    datasources: {
      type: Array,
      required: false,
      default: () => [],
    },

    creating: {
      type: Boolean,
      default: true,
    },
  },

  computed: {
    kind () {
      return Object.keys(this.step)
    },

    nameState () {
      const name = this.step[this.kind].name
      const isDuplicate = this.datasources.some(({ step }, index) => {
        return index !== this.index && step[Object.keys(step)].name === name && name
      })

      return name.length > 0 && !isDuplicate ? null : false
    },
  },
}
</script>

export default {
  data () {
    return {
      visible: false,
    }
  },

  methods: {
    toggleNav (visible) {
      this.visible = visible
      this.$emit('toggleNav', visible)
    },
  },
}

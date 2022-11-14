<script>
export default {
  props: {
    request: {
      type: Object,
      required: true,
    },
  },

  computed: {
    payload () {
      const { payload = [] } = this.request || {}
      return payload[0]
    },
  },
}
</script>

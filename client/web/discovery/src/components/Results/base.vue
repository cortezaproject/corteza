<script>
import TextHighlight from 'vue-text-highlight'

export default {
  components: {
    TextHighlight,
  },

  props: {
    index: {
      type: Number,
      required: true,
    },
    hit: {
      type: Object,
      required: true,
    },
    showMap: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      defaultBlacklistedFields: ['deleted', 'created', 'updated', 'security'],
    }
  },

  computed: {
    blacklistedFields () {
      return this.defaultBlacklistedFields
    },

    query () {
      return this.$route.query.query || ''
    },
  },

  methods: {
    limitData () {
      const out = {}
      const limit = 5

      if (this.hit.value) {
        let counter = 0

        for (const key in this.hit.value) {
          const value = this.hit.value[key]

          if (counter < limit && !!value && this.blacklistedFields.indexOf(key) < 0) {
            out[key] = value
            counter++
          }
        }
      }

      return out
    },
  },
}
</script>

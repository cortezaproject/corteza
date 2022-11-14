<script lang="js">
import { getComponent } from './index.js'

export default {
  inheritAttrs: false,

  props: {
    inline: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    previewType () {
      const component = getComponent({ type: this.$attrs.mime, src: this.$attrs.src, name: this.$attrs.name })
      if (!component) {
        this.$emit('previewUnavailable', { type: this.$attrs.mime })
      }
      return component
    },
  },
}
</script>

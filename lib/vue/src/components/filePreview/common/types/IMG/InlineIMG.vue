<template>
  <div class="d-flex inline h-100 mb-1">
    <img
      ref="image"
      :key="src"
      :src="src"
      :title="title"
      :alt="alt"
      :class="getClass"
      :style="previewStyle"
      :width="getWidth"
      :height="getHeight"
      @click.stop="$emit('openPreview', {})"
      @error.once="reloadBrokenImage"
      @load="loaded=true"
    >
  </div>
</template>

<script lang="js">
import base from '../base.vue'

export default {
  extends: base,

  props: {
    alt: {
      type: String,
      default: null,
    },

    title: {
      type: String,
      default: null,
    },
  },

  data () {
    return {
      loaded: false,
    }
  },

  computed: {
    isSvg () {
      return this.mime === 'image/svg+xml' ||
        (this.name && this.name.toLowerCase().endsWith('.svg'))
    },

    getClass () {
      const rtr = [...this.previewClass]
      if (this.$listeners.click) {
        rtr.push('clickable')
      }
      if (this.loaded) {
        rtr.push('loaded')
      }
      return rtr
    },

    getWidth () {
      if (this.isSvg) return undefined
      return this.meta.preview.image.width
    },
    getHeight () {
      if (this.isSvg) return undefined
      return this.meta.preview.image.height
    },
  },

  methods: {
    reloadBrokenImage (ev) {
      if (ev.target && ev.target.src) {
        window.setTimeout(() => {
          if (!ev.target && !ev.target.src) return

          // This forces Vue to re-try image download
          // eslint-disable-next-line
          ev.target.src = ev.target.src
        }, 500)
      }
    },
  },
}
</script>

<style scoped lang="scss">
div {
  object-fit: contain;

  img {
    &.loaded {
      width: auto;
      height: auto;
      display: block;
    }
  }

  &.inline {
    img {
      cursor: zoom-in;
    }

    img.disable-zoom-cursor {
      cursor: default;
    }
  }
}
</style>

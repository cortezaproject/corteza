<template>
  <div class="popup-img-preview">
    <photo-swipe
      :is-open="true"
      :items="items"
      :options="options"
      @close="() => $emit('close')"
    />
  </div>
</template>

<script lang="js">
import { PhotoSwipe } from 'v-photoswipe'
import base from '../base.vue'

export default {
  components: {
    PhotoSwipe,
  },

  extends: base,

  data () {
    return {
      options: {
        index: 0,
        bgOpacity: 0,
        closeOnScroll: false,
        escKey: false,
        history: false,
        arrowKeys: false,
        modal: false,

        closeEl: false,
        captionEl: false,
        fullscreenEl: false,
        zoomEl: false,
        shareEl: false,
        counterEl: false,
        arrowEl: false,
        preloaderEl: false,

        clickToCloseNonZoomable: false,
      },
    }
  },

  computed: {
    items () {
      const { original, preview } = this.meta
      const image = (original || preview || {}).image
      if (!image) {
        this.$emit('close')
        return []
      }

      return [{
        src: this.src,
        w: image.width,
        h: image.height,
      }]
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    setDefaultValues() {
      this.options = {}
    },
  },
}
</script>

<style lang="scss">
.popup-img-preview {
  .pswp {
    pointer-events: none;
    .pswp__img {
      pointer-events: all;
    }
  }
  .pswp__top-bar {
    display: none!important;
  }
}
</style>

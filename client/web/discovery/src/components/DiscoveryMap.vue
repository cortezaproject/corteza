<template>
  <c-map
    style="height: calc(100vh - 64px);"
    :map="{
      zoom,
      center,
    }"
    :markers="makerValues"
    @on-marker-click="onMarkerClick"
    @on-map-click="clearClickedMarker"
  />
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'

const { CMap } = components

export default {
  components: {
    CMap,
  },

  props: {
    markers: {
      type: Array,
      required: true,
    },
    hoverIndex: {
      type: String,
      default: undefined,
    },
  },

  data () {
    return {
      zoom: 3,
      center: [30, 30],
      rotation: 0,
      attribution: '&copy; <a target="_blank" rel="noopener noreferrer" href="http://osm.org/copyright">OpenStreetMap</a>',
      clickedMarker: undefined,
    }
  },

  computed: {
    makerValues () {
      return this.markers.map((marker, i) => {
        return {
          id: marker.id,
          value: marker.coordinates,
          opacity: [this.hoverIndex, this.clickedMarker].includes(marker.id) ? 1.0 : 0.6,
        }
      })
    },
  },

  watch: {
    markers: {
      immediate: true,
      handler (markers = []) {
        if (markers.length) {
          const { coordinates = [30, 30] } = markers[0]
          this.center = coordinates
        }
      },
    },

    hoverIndex: {
      handler (hoverIndex) {
        if (hoverIndex) {
          const { coordinates } = this.markers.find(({ id }) => id === hoverIndex) || {}
          if (coordinates) {
            this.center = coordinates
          }
        }
      },
    },
  },

  methods: {
    onMarkerClick ({ index }) {
      this.clickedMarker = index
      this.$emit('hover', this.clickedMarker)
    },

    clearClickedMarker () {
      this.clickedMarker = undefined
      this.$emit('hover', this.clickedMarker)
    },
  },
}
</script>

<style lang="scss">
.vl-style-text {
  color: var(--white);
}
</style>

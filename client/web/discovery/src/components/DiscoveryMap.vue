<template>
  <l-map
    :zoom="zoom"
    :center="center"
    style="height: calc(100vh - 64px);"
    @click="clearClickedMarker()"
  >
    <l-tile-layer
      url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      :attribution="attribution"
    />
    <l-marker
      v-for="(marker, i) in markers"
      :key="i"
      :lat-lng="getLatLng(marker.coordinates)"
      :opacity="[hoverIndex, clickedMarker].includes(marker.id) ? 1.0 : 0.6"
      @click="onMarkerClick(marker.id)"
    />
  </l-map>
</template>

<script>
import { isNumber } from 'lodash'
import { latLng } from 'leaflet'

export default {
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

  watch: {
    markers: {
      immediate: true,
      handler (markers = []) {
        if (markers.length) {
          const { coordinates = [30, 30] } = markers[0]
          this.center = this.getLatLng(coordinates)
        }
      },
    },

    hoverIndex: {
      handler (hoverIndex) {
        if (hoverIndex) {
          const { coordinates } = this.markers.find(({ id }) => id === hoverIndex) || {}
          if (coordinates) {
            this.center = this.getLatLng(coordinates)
          }
        }
      },
    },
  },

  methods: {
    getLatLng (coordinates = [0, 0]) {
      const [lat, lng] = coordinates

      if (isNumber(lat) && isNumber(lng)) {
        return latLng(lat, lng)
      }
    },

    onMarkerClick (ID) {
      this.clickedMarker = ID
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
  color: $white;
}
</style>

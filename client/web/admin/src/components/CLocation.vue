<template>
  <div>
    <b-button
      variant="link"
      rounded
      class="p-0"
      @click="openMap"
    >
      <font-awesome-icon
        :icon="['fas', 'map-marked-alt']"
      />
    </b-button>

    <b-modal
      v-model="map.show"
      size="lg"
      title="Map"
      body-class="p-0"
      hide-header
      hide-footer
    >
      <l-map
        ref="map"
        :zoom="map.zoom"
        :center="map.center"
        style="height: 75vh; width: 100%;"
        :style="{ 'cursor': editable ? 'pointer' : 'grab'}"
        @click="placeMarker"
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="map.attribution"
        />
        <l-marker
          v-if="value && value.length"
          :lat-lng="getLatLng(value)"
          @click="removeMarker"
        />
      </l-map>
    </b-modal>
  </div>
</template>

<script>
import { latLng } from 'leaflet'

export default {
  props: {
    value: {
      type: Array,
      required: true,
    },

    editable: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      map: {
        show: false,
        zoom: 3,
        center: [30, 30],
        rotation: 0,
        attribution: '&copy; <a target="_blank" rel="noopener noreferrer" href="http://osm.org/copyright">OpenStreetMap</a>',
      },
    }
  },

  methods: {
    openMap () {
      this.map.show = true

      setTimeout(() => {
        this.$refs.map.mapObject.invalidateSize()
      }, 100)
    },

    getLatLng (coordinates = [0, 0]) {
      return latLng(coordinates[0], coordinates[1])
    },

    placeMarker ({ latlng = {} }) {
      const { lat = 0, lng = 0 } = latlng
      this.$emit('input', [lat, lng])
    },

    removeMarker () {
      this.$emit('input', [])
    },
  },
}
</script>

<style lang="scss">

</style>

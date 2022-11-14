<template>
  <l-map
    :zoom="zoom"
  >
    <l-tile-layer
      url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      :attribution="attribution"
    />

    <l-marker
      v-for="connection in validMarkers"
      :key="connection.name"
      :lat-lng="getLocationCoordinates(connection)"
    >
      <l-tooltip
        :options="{
          offset: [-14, 10],
          direction: 'bottom',
        }"
      >
        <h5
          class="text-primary"
        >
          {{ $t('server-details') }}
        </h5>
        <b-form-group
          :label="$t('name')"
          label-class="text-primary"
        >
          {{ connection.meta.name }}
        </b-form-group>

        <b-form-group
          :label="$t('location')"
          label-class="text-primary"
        >
          {{ getLocationName(connection) }}
        </b-form-group>
      </l-tooltip>
    </l-marker>
  </l-map>
</template>

<script>
import { latLng } from 'leaflet'
import { LMap, LTileLayer, LMarker, LTooltip } from 'vue2-leaflet'

export default {
  i18nOptions: {
    namespaces: 'map',
  },

  components: {
    LMap,
    LTileLayer,
    LMarker,
    LTooltip,
  },

  props: {
    connections: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
      zoom: 2,
      center: [47.313220, -1.319482],
      rotation: 0,
      attribution: '&copy; <a target="_blank" rel="noopener noreferrer" href="http://osm.org/copyright">OpenStreetMap</a>',
    }
  },

  computed: {
    validMarkers () {
      return this.connections.filter(({ meta = {} }) => {
        const { location = {} } = meta
        const { geometry = {} } = location
        const { coordinates = [] } = geometry

        return coordinates && !!coordinates.length
      })
    },
  },

  methods: {
    getLocationCoordinates ({ meta = {} }) {
      const { location = {} } = meta
      const { geometry = {} } = location
      return this.getLatLng(geometry.coordinates)
    },

    getLocationName (connection) {
      return connection.meta.location.properties.name || this.$t('unnamed-location')
    },

    getLatLng (coordinates = [0, 0]) {
      return latLng(coordinates[0], coordinates[1])
    },
  },
}
</script>

<style lang="scss">
.vl-style-text {
  color: white;
}
</style>

<template>
  <div>
    <div class="my-2">
      <l-map
        ref="map"
        :zoom="options.zoomStarting"
        :min-zoom="options.zoomMin"
        :max-zoom="options.zoomMax"
        :center="options.center"
        :bounds="bounds"
        :max-bounds="options.bounds"
        class="w-100 cursor-pointer"
        style="height: 45vh;"
        @update:zoom="zoomUpdated"
        @update:center="updateCenter"
        @update:bounds="boundsUpdated"
        @locationfound="onLocationFound"
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="map.attribution"
        />
        <l-control class="leaflet-bar">
          <a
            :title="$t('geometry.tooltip.goToCurrentLocation')"
            role="button"
            class="d-flex justify-content-center align-items-center"
            @click="goToCurrentLocation"
          >
            <font-awesome-icon
              :icon="['fas', 'location-arrow']"
              class="text-primary"
            />
          </a>
        </l-control>
      </l-map>
      <b-form-text id="password-help-block">
        {{ $t('geometry.mapHelpText') }}
      </b-form-text>
    </div>
    <hr>

    <b-row
      class="mb-2 mt-4"
    >
      <b-col
        sm="12"
        md="4"
      >
        <b-form-group
          :label="$t('geometry.zoom.zoomStartingLabel')"
          class="rounded-left"
        >
          <b-form-input
            v-model="options.zoomStarting"
            number
            readonly
            type="number"
          />
        </b-form-group>
      </b-col>

      <b-col
        sm="12"
        md="4"
      >
        <b-form-group
          :label="$t('geometry.zoom.zoomMinLabel')"
          :description="`${options.zoomMin}`"
          class="rounded-0"
        >
          <b-form-input
            v-model="options.zoomMin"
            number
            :min="1"
            :max="18"
            type="range"
          />
        </b-form-group>
      </b-col>

      <b-col
        sm="12"
        md="4"
      >
        <b-form-group
          :label="$t('geometry.zoom.zoomMaxLabel')"
          :description="`${options.zoomMax}`"
        >
          <b-form-input
            v-model="options.zoomMax"
            number
            :min="1"
            :max="18"
            type="range"
          />
        </b-form-group>
      </b-col>

      <b-col
        sm="12"
        md="4"
      >
        <b-form-group
          :label="$t('geometry.bounds.lockBounds')"
          class="rounded-left"
        >
          <b-form-checkbox
            v-model="options.lockBounds"
            name="lock-bounds"
            switch
            size="lg"
            @change="updateBounds"
          />
        </b-form-group>
      </b-col>

      <b-col
        sm="12"
        md="4"
      >
        <b-form-group
          :label="$t('geometry.onMarkerClick')"
          label-class="text-primary"
        >
          <b-form-select
            v-model="options.displayOption"
            :options="displayOptions"
          />
        </b-form-group>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import base from '../base'
import { latLng } from 'leaflet'
import { LControl } from 'vue2-leaflet'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    LControl,
  },

  extends: base,

  data () {
    return {
      map: {
        show: false,
        zoom: 3,
        center: [30, 30],
        rotation: 0,
        attribution: '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a>',
      },

      localValue: { coordinates: [] },
      center: [],
      bounds: null,
    }
  },

  computed: {
    displayOptions () {
      return [
        { value: 'sameTab', text: this.$t('geometry.openInSameTab') },
        { value: 'newTab', text: this.$t('geometry.openInNewTab') },
        { value: 'modal', text: this.$t('geometry.openInModal') },
      ]
    },
  },

  methods: {
    getLatLng (coordinates = [undefined, undefined]) {
      const [lat, lng] = coordinates

      if (lat && lng) {
        return latLng(lat, lng)
      }
    },

    updateCenter (coordinates) {
      let { lat = 0, lng = 0 } = coordinates || {}

      lat = Math.round(lat * 1e7) / 1e7
      lng = Math.round(lng * 1e7) / 1e7

      this.options.center = [lat, lng]
    },

    boundsUpdated (coordinates) {
      this.bounds = coordinates

      this.updateBounds(this.options.lockBounds)
    },

    zoomUpdated (zoom) {
      this.options.zoomStarting = zoom
    },

    updateBounds (value) {
      if (value) {
        const bounds = this.bounds || this.$refs.map.mapObject.getBounds()
        const { _northEast, _southWest } = bounds

        this.options.bounds = [Object.values(_northEast), Object.values(_southWest)]
      } else {
        this.options.bounds = null
      }
    },

    goToCurrentLocation () {
      this.$refs.map.mapObject.locate()
    },

    onLocationFound ({ latitude, longitude }) {
      const zoom = this.$refs.map.mapObject._zoom >= 13 ? this.$refs.map.mapObject._zoom : 13
      this.$refs.map.mapObject.flyTo([latitude, longitude], zoom)
    },
  },
}
</script>

<style></style>

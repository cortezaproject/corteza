<template>
  <div class="d-flex flex-column">
    <b-form-group>
      <b-form-checkbox v-model="f.options.prefillWithCurrentLocation">
        {{ $t('prefillWithCurrentLocation') }}
      </b-form-checkbox>

      <b-form-checkbox v-model="f.options.hideCurrentLocationButton">
        {{ $t('hideCurrentLocationButton') }}
      </b-form-checkbox>

      <b-form-checkbox v-model="f.options.hideGeoSearch">
        {{ $t('hideGeoSearch') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      :label="$t('initialZoomAndPosition')"
      label-class="text-primary"
      class="mb-0"
    >
      <l-map
        ref="map"
        :zoom="zoom"
        :center="center"
        class="w-100"
        style="height: 50vh;"
        @update:zoom="f.options.zoom = $event"
        @update:center="f.options.center = $event"
        @locationfound="onLocationFound"
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="attribution"
        />
        <l-control class="leaflet-bar">
          <a
            :title="$t('tooltip.goToCurrentLocation')"
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
    </b-form-group>
  </div>
</template>

<script>
import base from './base'
import { LControl } from 'vue2-leaflet'

export default {
  i18nOptions: {
    namespaces: 'field',
    keyPrefix: 'kind.geometry',
  },

  components: {
    LControl,
  },

  extends: base,

  data () {
    return {
      attribution: '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a>',
    }
  },

  computed: {
    center () {
      return this.f.options.center || [30, 30]
    },

    zoom () {
      return this.f.options.zoom || 3
    },
  },

  methods: {
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

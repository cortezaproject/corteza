<template>
  <div>
    <span
      v-for="(c, index) of localValue"
      :key="index"
      :class="{ 'd-block': field.options.multiDelimiter === '\n' }"
    >
      <a
        class="text-primary pointer"
        @click.stop="openMap(index)"
      >
        {{ c.lat }}, {{ c.lng }}
        <font-awesome-icon
          :icon="['fas', 'map-marked-alt']"
        />
        {{ index !== localValue.length - 1 ? field.options.multiDelimiter : '' }}
      </a>
    </span>

    <b-modal
      v-model="map.show"
      :title="field.label || field.name"
      size="lg"
      body-class="p-0"
      hide-footer
    >
      <div
        v-if="!field.options.hideGeoSearch"
        class="geosearch-container"
      >
        <c-input-search
          v-model="geoSearch.query"
          :placeholder="$t('geosearchInputPlaceholder')"
          :autocomplete="'off'"
          :debounce="300"
          @input="onGeoSearch"
        />

        <div class="geosearch-results">
          <div
            v-for="(result, idx) in geoSearch.results"
            :key="idx"
            class="geosearch-result"
            @click="placeGeoSearchMarker(result)"
          >
            {{ result.label }}
          </div>
        </div>
      </div>

      <l-map
        ref="map"
        :zoom="map.zoom"
        :center="map.center"
        style="height: 75vh; width: 100%;"
        @locationfound="onLocationFound"
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="map.attribution"
        />
        <l-marker
          v-for="(marker, i) in localValue"
          :key="i"
          :lat-lng="marker"
          :opacity="localValueIndex === undefined || i == localValueIndex ? 1.0 : 0.6"
        />
        <l-control class="leaflet-bar">
          <a
            v-if="!field.options.hideCurrentLocationButton"
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
    </b-modal>

    <errors :errors="errors" />
  </div>
</template>
<script>
import base from './base'
import { latLng } from 'leaflet'
import { LControl } from 'vue2-leaflet'
import { OpenStreetMapProvider } from 'leaflet-geosearch'
import { components } from '@cortezaproject/corteza-vue'
const { CInputSearch } = components

export default {
  i18nOptions: {
    namespaces: 'field',
    keyPrefix: 'kind.geometry',
  },

  components: {
    LControl,
    CInputSearch,
  },

  extends: base,

  data () {
    return {
      map: {
        show: false,
        zoom: 14,
        center: [30, 30],
        rotation: 0,
        attribution: '&copy; <a target="_blank" rel="noopener noreferrer" href="http://osm.org/copyright">OpenStreetMap</a>',
      },
      localValueIndex: undefined,

      geoSearch: {
        query: '',
        provider: new OpenStreetMapProvider(),
        results: [],
      },
    }
  },

  computed: {
    localValue () {
      if (this.field.isMulti) {
        return this.value.map(v => {
          return this.getLatLng(JSON.parse(v || '{"coordinates":[]}').coordinates || [])
        }).filter(c => c)
      } else {
        return [this.getLatLng(JSON.parse(this.value || '{"coordinates":[]}').coordinates || [])].filter(c => c)
      }
    },
  },

  methods: {
    openMap (index) {
      this.localValueIndex = index
      this.map.center = this.localValue[index] || this.field.options.center
      this.map.zoom = index >= 0 ? 13 : this.field.options.zoom
      this.map.show = true

      setTimeout(() => {
        this.$refs.map.mapObject.invalidateSize()
      }, 100)
    },

    getLatLng (coordinates = [undefined, undefined]) {
      const [lat, lng] = coordinates

      if (lat && lng) {
        return latLng(lat, lng)
      }
    },

    goToCurrentLocation () {
      this.$refs.map.mapObject.locate()
    },

    onLocationFound ({ latitude, longitude }) {
      const zoom = this.$refs.map.mapObject._zoom >= 13 ? this.$refs.map.mapObject._zoom : 13
      this.$refs.map.mapObject.flyTo([latitude, longitude], zoom)
    },

    placeGeoSearchMarker (result) {
      const zoom = this.$refs.map.mapObject._zoom >= 15 ? this.$refs.map.mapObject._zoom : 15
      this.$refs.map.mapObject.flyTo([result.latlng.lat, result.latlng.lng], zoom, { animate: false })
      this.geoSearch.results = []
    },

    onGeoSearch (query) {
      if (!query) {
        this.geoSearch.results = []
        return
      }

      this.geoSearch.provider.search({ query }).then(results => {
        this.geoSearch.results = results.map(result => ({
          ...result,
          latlng: {
            lat: result.raw.lat,
            lng: result.raw.lon,
          },
        }))
      }).catch(() => {
        this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.locationSearchFailed'))()
      })
    },
  },
}
</script>

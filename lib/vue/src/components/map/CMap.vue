<template>
  <div
    class="position-relative"
  >
    <div
      v-if="!hideGeoSearch"
      class="geosearch-container"
      @mouseover="disableMap"
      @mouseleave="enableMap"
    >
      <c-input-search
        v-model="geoSearch.query"
        :placeholder="labels.geosearchInputPlaceholder"
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
      :zoom="mapOptions.zoom"
      :center="mapOptions.center"
      :min-zoom="mapOptions.minZoom"
      :max-zoom="mapOptions.maxZoom"
      :bounds="mapOptions.bounds"
      :max-bounds="mapOptions.maxBounds"
      class="w-100 h-100"
      @click="onMapClick"
      @locationfound="onLocationFound"
      @update:zoom="onZoom"
      @update:center="onCenter"
      @update:bounds="onBoundsUpdate"
    >
      <l-tile-layer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        :attribution="mapOptions.attribution"
      />

      <l-polygon
        v-for="(geometry, i) in polygons"
        :key="`polygon-${i}`"
        :lat-lngs="geometry.map(value => value.geometry)"
        :color="(geometry.find(g => g) || {}).color"
      />

      <l-marker
        v-for="(marker, i) in markerValues"
        :key="`marker-${i}`"
        :lat-lng="marker.value"
        :icon="getIcon(marker.color)"
        :opacity="marker.opacity || 1.0"
        @click="onMarkerClick(i, marker)"
      >
        <l-tooltip
          v-if="$scopedSlots['marker-tooltip'] || marker.title"
          :options="{
            offset: [-1, 5],
            direction: 'bottom',
          }"
        >
          <slot
            name="marker-tooltip"
            :marker="marker"
          >
            {{ marker.title }}
          </slot>
        </l-tooltip>
      </l-marker>

      <l-marker
        v-if="geoSearch.marker"
        :lat-lng="geoSearch.marker.latlng"
        :icon="getIcon(getCSSVariable('--secondary'))"
      >
        <l-tooltip
          :options="{
            offset: [-1, 5],
            direction: 'bottom',
          }"
        >
          {{ geoSearch.marker.title }}
        </l-tooltip>
      </l-marker>

      <l-control class="leaflet-bar">
        <a
          v-if="!hideCurrentLocationButton"
          v-b-tooltip.noninteractive.hover="{ title: labels.tooltip && labels.tooltip.goToCurrentLocation, boundary: 'body' }"
          role="button"
          class="d-flex justify-content-center align-items-center"
          @click="goToCurrentLocation"
        >
          <font-awesome-icon
            :icon="['fas', 'location-crosshairs']"
          />
        </a>
      </l-control>
    </l-map>
  </div>
</template>

<script>
import { divIcon, latLng, latLngBounds } from 'leaflet'
import {
  OpenStreetMapProvider,
  OpenCageProvider,
  EsriProvider,
  GeoapifyProvider,
  GeocodeEarthProvider,
  GoogleProvider,
  LocationIQProvider,
} from 'leaflet-geosearch'
import { isNumber } from 'lodash'
import { LControl, LMap, LMarker, LPolygon, LTileLayer, LTooltip } from 'vue2-leaflet'
import CInputSearch from '../input/CInputSearch.vue'

import 'leaflet/dist/leaflet.css'

export default {
  components: {
    LControl,
    CInputSearch,
    LPolygon,
    LTileLayer,
    LMarker,
    LMap,
    LTooltip,
  },

  props: {
    hideCurrentLocationButton: {
      type: Boolean,
      default: false,
    },

    labels: {
      type: Object,
      default: () => ({}),
    },

    map: {
      type: Object,
      default: () => ({}),
    },

    markers: {
      type: Array,
      default: () => ([]),
    },

    polygons: {
      type: Array,
      default: () => ([]),
    },

    hideGeoSearch: {
      type: Boolean,
      default: false,
    },

    disabled: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      geoSearch: {
        query: '',
        results: [],
        marker: null,
      },
    }
  },

  computed: {
    markerValues () {
      return this.markers.map((m) => {
        return {
          ...m,
          value: this.getLatLng(m.value),
        }
      }).filter(c => c.value) || []
    },

    mapOptions () {
      const map = { ...this.map }
      const defaultOptions = {
        zoom: 3,
        center: [30, 30],
        rotation: 0,
        attribution: '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a>',
      }

      map.bounds = map.bounds ? latLngBounds(map.bounds) : null

      return {
        ...defaultOptions,
        ...map,
      }
    },

    geoSearchApiKey () {
      return this.$Settings.get('ui.location.geoSearchApiKey', '')
    },

    geoSearchProviderName () {
      return this.$Settings.get('ui.location.geoSearchProvider', '')
    },

    geoSearchProvider () {
      const providerName = this.geoSearchProviderName.toLowerCase() || 'openstreetmap'
      const apiKey = this.geoSearchApiKey

      // Map of available providers
      const providers = {
        openstreetmap: () => new OpenStreetMapProvider(),
        opencage: () => new OpenCageProvider({
          params: { key: apiKey },
        }),
        esri: () => new EsriProvider(),
        geoapify: () => new GeoapifyProvider({
          params: { apiKey },
        }),
        geocodeearth: () => new GeocodeEarthProvider({
          params: { api_key: apiKey },
        }),
        google: () => new GoogleProvider({
          apiKey,
        }),
        locationiq: () => new LocationIQProvider({
          params: { key: apiKey },
        }),
      }

      if (providers[providerName]) {
        return providers[providerName]()
      }

      console.warn(`Unknown geosearch provider: ${providerName}, falling back to OpenStreetMap`)
      return new OpenStreetMapProvider()
    },

    tileProviderUrl () {
      return ''
    },

    tileProvider () {
      if (this.tileProviderUrl) {
        return this.tileProviderUrl
      }

      return 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
    },
  },

  mounted () {
    if (this.$refs.map && this.$refs.map.mapObject) {
      this.onBoundsUpdate(this.$refs.map.mapObject.getBounds())
    }
  },

  methods: {
    onGeoSearch (query) {
      if (!query || !this.geoSearchProvider) {
        this.geoSearch.results = []
        this.geoSearch.marker = null
        return
      }

      this.geoSearchProvider.search({ query }).then(results => {
        this.geoSearch.results = results.map(result => {
          // Different providers return coordinates in different formats
          let lat, lng

          if (result.y !== undefined && result.x !== undefined) {
            // ESRI and some others use x, y
            lat = result.y
            lng = result.x
          } else if (result.raw) {
            // Most providers (OpenStreetMap, OpenCage, etc.) use raw.lat/lon or raw.latitude/longitude
            lat = result.raw.lat || result.raw.latitude
            lng = result.raw.lon || result.raw.lng || result.raw.longitude
          }

          return {
            ...result,
            latlng: {
              lat,
              lng,
            },
          }
        })
      }).catch(() => {
        this.$emit('on-geosearch-error')
      })
    },

    placeGeoSearchMarker (result) {
      const zoom = this.$refs.map.mapObject._zoom >= 15 ? this.$refs.map.mapObject._zoom : 15
      this.$refs.map.mapObject.flyTo([result.latlng.lat, result.latlng.lng], zoom, { animate: false })
      this.geoSearch.marker = { title: result.label, latlng: result.latlng }
      this.geoSearch.results = []
      this.onMapClick(result)
    },

    getLatLng (coordinates = [undefined, undefined]) {
      const [lat, lng] = coordinates

      if (isNumber(lat) && isNumber(lng)) {
        return latLng(lat, lng)
      }
    },

    onLocationFound ({ latitude, longitude }) {
      const zoom = this.$refs.map.mapObject._zoom >= 15 ? this.$refs.map.mapObject._zoom : 15
      this.$refs.map.mapObject.flyTo([latitude, longitude], zoom)
      this.$emit('location-found', { latlng: { lat: latitude, lng: longitude } })
    },

    disableMap () {
      if (this.disabled) this.$refs.map.mapObject._handlers.forEach(handler => handler.disable())
    },

    enableMap () {
      if (this.disabled) this.$refs.map.mapObject._handlers.forEach(handler => handler.enable())
    },

    onMarkerClick (index, marker) {
      this.$emit('on-marker-click', { index, marker })
    },

    goToCurrentLocation () {
      this.$refs.map.mapObject.locate()
    },

    onMapClick (e) {
      this.$emit('on-map-click', e)
    },

    getCSSVariable (variable) {
      return getComputedStyle(document.documentElement).getPropertyValue(variable)
    },

    getIcon (markerColor = this.getCSSVariable('--primary')) {
      const markerIconHtml = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 34.892337" height="60" width="40" style="margin-top: -40px;margin-left: -15px;height: 35px;">
          <g transform="translate(-814.59595,-274.38623)">
            <g transform="matrix(1.1855854,0,0,1.1855854,-151.17715,-57.3976)">
              <path d="m 817.11249,282.97118 c -1.25816,1.34277 -2.04623,3.29881 -2.01563,5.13867 0.0639,3.84476 1.79693,5.3002 4.56836,10.59179 0.99832,2.32851 2.04027,4.79237 3.03125,8.87305 0.13772,0.60193 0.27203,1.16104 0.33416,1.20948 0.0621,0.0485 0.19644,-0.51262 0.33416,-1.11455 0.99098,-4.08068 2.03293,-6.54258 3.03125,-8.87109 2.77143,-5.29159 4.50444,-6.74704 4.56836,-10.5918 0.0306,-1.83986 -0.75942,-3.79785 -2.01758,-5.14062 -1.43724,-1.53389 -3.60504,-2.66908 -5.91619,-2.71655 -2.31115,-0.0475 -4.4809,1.08773 -5.91814,2.62162 z" style="fill:${markerColor};stroke:${markerColor};"/>
              <circle r="3.0355" cy="288.25278" cx="823.03064" id="path3049" style="display:inline;fill:#FFFFFF;"/>
            </g>
          </g>
        </svg>`

      return divIcon({
        className: 'marker-pin',
        html: markerIconHtml,
      })
    },

    onZoom (e) {
      this.$emit('on-zoom', e)
    },

    onCenter (e) {
      this.$emit('on-center', e)
    },

    onBoundsUpdate (value) {
      this.$nextTick(() => {
        setTimeout(() => {
          this.$refs.map.mapObject.invalidateSize()
        }, 100)
      })

      value = value || this.$refs.map.mapObject.getBounds()

      this.$emit('on-bounds-update', value)
    },
  },
}
</script>

<style lang="scss">
.leaflet-touch .leaflet-bar {
  border: 1px solid transparent;
  border-radius: 0.3rem;
}

.leaflet-bar a {
  background-color: var(--white) !important;
  color: var(--primary) !important;
  text-decoration: none !important;

  &:hover {
    background-color: var(--white) !important;
    transition: background-color 0.15s ease;
  }
}

.geosearch-result {
  &:hover {
    background-color: var(--light) !important;
    color: var(--black);
  }

  &:active {
    color: var(--white) !important;
    background-color: var(--primary) !important;
  }
}
</style>

<style scoped>
.geosearch-container {
  position: absolute;
  display: block;
  height: auto;
  width: 50%;
  max-width: 50%;
  cursor: auto;
  z-index: 1030;
  left: 50%;
  transform: translateX(-50%);
  top: 10px;
}

.geosearch-results {
  margin: 1px;
  border-radius: 2px;
  background-color: var(--white);
  max-height: 50%;
  overflow: auto;
}

.geosearch-result {
  border-radius: 2px;
  line-height: 32px;
  padding: 0 8px;
  font-size: 12px;
  white-space: nowrap;
}

.geosearch-result:hover {
  background-color: var(--gray-200);
  cursor: pointer;
}
</style>

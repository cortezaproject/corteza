<template>
  <b-form-group
    label-class="text-primary"
    :state="state"
    :class="formGroupStyleClasses"
  >
    <template
      v-if="!valueOnly"
      #label
    >
      <div
        class="d-flex align-items-top"
      >
        <label
          class="mb-0"
        >
          {{ label }}
        </label>

        <hint
          :id="field.fieldID"
          :text="hint"
        />
      </div>

      <small
        class="form-text font-weight-light text-muted"
      >
        {{ description }}
      </small>
    </template>

    <div class="d-flex w-100">
      <b-button
        v-if="field.isMulti"
        :title="$t('tooltip.openMap')"
        variant="light"
        class="w-100"
        @click="openMap()"
      >
        <font-awesome-icon
          :icon="['fas', 'map-marked-alt']"
          class="text-primary"
        />
      </b-button>
    </div>

    <multi
      v-if="field.isMulti"
      v-slot="ctx"
      :value.sync="localValue"
      :errors="errors"
      single-input
    >
      <b-input-group>
        <b-form-input
          v-model="localValue[ctx.index].coordinates[0]"
          type="number"
          step="0.000001"
          number
          :placeholder="$t('latitude')"
        />
        <b-form-input
          v-model="localValue[ctx.index].coordinates[1]"
          type="number"
          step="0.000001"
          number
          :placeholder="$t('longitude')"
        />
        <b-input-group-append>
          <b-button
            :title="$t('tooltip.openMap')"
            variant="light"
            class="d-flex align-items-center"
            @click="openMap(ctx.index)"
          >
            <font-awesome-icon
              :icon="['fas', 'map-marked-alt']"
              class="text-primary"
            />
          </b-button>

          <b-button
            v-if="!field.options.hideCurrentLocationButton"
            :title="$t('tooltip.useCurrentLocation')"
            variant="light"
            class="d-flex align-items-center"
            @click="useCurrentLocation(ctx.index)"
          >
            <font-awesome-icon
              :icon="['fas', 'location-arrow']"
              class="text-primary"
            />
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </multi>

    <template v-else>
      <b-input-group>
        <b-form-input
          v-model="localValue.coordinates[0]"
          type="number"
          step="0.000001"
          number
          :placeholder="$t('latitude')"
        />
        <b-form-input
          v-model="localValue.coordinates[1]"
          type="number"
          step="0.000001"
          number
          :placeholder="$t('longitude')"
        />
        <b-input-group-append>
          <b-button
            :title="$t('tooltip.openMap')"
            variant="light"
            class="d-flex align-items-center"
            @click="openMap()"
          >
            <font-awesome-icon
              :icon="['fas', 'map-marked-alt']"
              class="text-primary"
            />
          </b-button>

          <b-button
            v-if="!field.options.hideCurrentLocationButton"
            :title="$t('tooltip.useCurrentLocation')"
            variant="light"
            class="d-flex align-items-center"
            @click="useCurrentLocation()"
          >
            <font-awesome-icon
              :icon="['fas', 'location-arrow']"
              class="text-primary"
            />
          </b-button>
        </b-input-group-append>
      </b-input-group>

      <errors :errors="errors" />
    </template>

    <b-modal
      v-model="map.show"
      :title="field.label || field.name"
      size="lg"
      body-class="p-0"
    >
      <template #modal-footer>
        {{ $t('clickToPlaceMarker') }}
      </template>

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
        style="height: 75vh; width: 100%; cursor: pointer;"
        @click="placeMarker"
        @locationfound="onLocationFound"
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="map.attribution"
        />
        <l-marker
          v-for="(marker, i) in markers"
          :key="i"
          :lat-lng="marker"
          :opacity="localValueIndex === undefined || i == localValueIndex ? 1.0 : 0.6"
          @click="removeMarker(i)"
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
  </b-form-group>
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
      localValue: undefined,
      localValueIndex: undefined,

      map: {
        show: false,
        zoom: 3,
        center: [30, 30],
        rotation: 0,
        attribution: '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a>',
      },

      geoSearch: {
        query: '',
        provider: new OpenStreetMapProvider(),
        results: [],
      },
    }
  },

  computed: {
    markers () {
      let markers = [this.localValue.coordinates]

      if (this.field.isMulti) {
        markers = this.localValue.map(({ coordinates }) => coordinates && coordinates.length ? coordinates : undefined)
      }

      return markers.map(this.getLatLng).filter(c => c)
    },
  },

  watch: {
    localValue: {
      deep: true,
      handler (value) {
        this.value = this.field.isMulti ? value.filter(v => (v || {}).coordinates).map(v => JSON.stringify(v)) : JSON.stringify(value)
      },
    },
  },

  created () {
    if (this.field.isMulti) {
      this.localValue = this.value.map(v => {
        return JSON.parse(v || '{"coordinates":[]}')
      })
    } else {
      this.localValue = JSON.parse(this.value || '{"coordinates":[]}')
    }

    if (this.field.options.prefillWithCurrentLocation) {
      this.useCurrentLocation()
    }
  },

  methods: {
    openMap (index) {
      this.localValueIndex = index
      const firstCoordinates = (index ? this.localValue[index] : (this.field.isMulti ? this.localValue[0] : this.localValue)) || {}
      this.map.center = firstCoordinates.coordinates && firstCoordinates.coordinates.length ? firstCoordinates.coordinates : this.field.options.center
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

    placeMarker (e, index) {
      const { lat = 0, lng = 0 } = e.latlng || {}
      const coords = {
        coordinates: [
          Math.round(lat * 1e7) / 1e7,
          Math.round(lng * 1e7) / 1e7,
        ],
      }

      if (this.field.isMulti) {
        if (index >= 0) {
          this.localValue.splice(index, 1, coords)
        } else {
          this.localValue.push(coords)
        }
      } else {
        this.localValue = coords
      }
    },

    removeMarker (i) {
      if (this.field.isMulti) {
        this.localValue.splice(i, 1)
      } else {
        this.localValue = { coordinates: [] }
      }
    },

    useCurrentLocation (index) {
      try {
        if (!navigator.geolocation) {
          this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.notSupported'))()
        }

        navigator.geolocation.getCurrentPosition(
          ({ coords }) => {
            const latlng = { lat: coords.latitude, lng: coords.longitude }
            this.placeMarker({ latlng }, index)
          },
          error => {
            switch (error.code) {
              case error.PERMISSION_DENIED:
                this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.permissionDenied'))()
                break
              case error.POSITION_UNAVAILABLE:
                this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.positionUnavailable'))()
                break
              case error.TIMEOUT:
                this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.timeout'))()
                break
              default:
                this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.unknownError'))()
                break
            }
          },
        )
      } catch (error) {
        this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.errorOccurred'))()
      }
    },

    goToCurrentLocation () {
      this.$refs.map.mapObject.locate()
    },

    onLocationFound ({ latitude, longitude }) {
      const zoom = this.$refs.map.mapObject._zoom >= 15 ? this.$refs.map.mapObject._zoom : 15
      const latlng = { lat: latitude, lng: longitude }
      this.placeMarker({ latlng })
      this.$refs.map.mapObject.flyTo([latitude, longitude], zoom)
    },

    placeGeoSearchMarker (result) {
      const zoom = this.$refs.map.mapObject._zoom >= 15 ? this.$refs.map.mapObject._zoom : 15
      this.$refs.map.mapObject.flyTo([result.latlng.lat, result.latlng.lng], zoom, { animate: false })
      this.placeMarker(result)
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

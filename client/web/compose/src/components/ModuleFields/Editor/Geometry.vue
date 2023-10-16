<template>
  <b-form-group
    :label-cols-md="horizontal && '5'"
    :label-cols-xl="horizontal && '4'"
    :content-cols-md="horizontal && '7'"
    :content-cols-xl="horizontal && '8'"
    :state="state"
    :class="formGroupStyleClasses"
  >
    <template
      #label
    >
      <div
        v-if="!valueOnly"
        class="d-flex align-items-center text-primary p-0"
      >
        <span
          :title="label"
          class="d-inline-block mw-100 py-1"
        >
          {{ label }}
        </span>

        <c-hint :tooltip="hint" />

        <slot name="tools" />
      </div>
      <div
        class="small text-muted"
        :class="{ 'mb-1': description }"
      >
        {{ description }}
      </div>
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

      <c-map
        :map="map"
        :hide-geo-search="field.options.hideGeoSearch"
        :hide-current-location-button="field.options.hideCurrentLocationButton"
        :markers="markers"
        :labels="{
          tooltip: { 'goToCurrentLocation': $t('tooltip.goToCurrentLocation') }
        }"
        style="height: 75vh; width: 100%; cursor: pointer;"
        @on-map-click="placeMarker"
        @on-marker-click="removeMarker"
        @on-geosearch-error="onGeoSearchError"
      />
    </b-modal>
  </b-form-group>
</template>
<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
import { isNumber } from 'lodash'
const { CMap } = components

export default {
  i18nOptions: {
    namespaces: 'field',
    keyPrefix: 'kind.geometry',
  },

  components: {
    CMap,
  },

  extends: base,

  data () {
    return {
      localValue: undefined,
      localValueIndex: undefined,

      map: {
        show: false,
      },
    }
  },

  computed: {
    markers () {
      let markers = [{ value: this.localValue.coordinates, opacity: 1.0 }]

      if (this.field.isMulti) {
        markers = this.localValue.map(({ coordinates }, i) => ({
          value: coordinates && coordinates.length ? coordinates : undefined,
          opacity: this.localValueIndex === undefined || i === this.localValueIndex ? 1.0 : 0.6,
        }))
      }

      return markers
    },
  },

  watch: {
    localValue: {
      deep: true,
      handler (value) {
        this.value = this.field.isMulti ? value.filter(v => (v || {}).coordinates).map(v => JSON.stringify(v)) : JSON.stringify(value)
      },
    },

    'field.isMulti': {
      immediate: true,
      handler (value) {
        if (this.field.isMulti) {
          this.localValue = this.value.map(v => {
            return JSON.parse(v || '{"coordinates":[]}')
          })
        } else {
          this.localValue = JSON.parse(this.value || '{"coordinates":[]}')
        }
      },
    },

    'field.options.prefillWithCurrentLocation': {
      immediate: true,
      handler (value) {
        if (value) {
          this.useCurrentLocation()
        }
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    openMap (index) {
      this.localValueIndex = index
      const firstCoordinates = (index >= 0 ? this.localValue[index] : this.localValue) || {}
      firstCoordinates.coordinates = firstCoordinates.coordinates ? [...firstCoordinates.coordinates] : []

      this.map.center = firstCoordinates.coordinates &&
                        firstCoordinates.coordinates.length === 2 &&
                        firstCoordinates.coordinates.every(isNumber) ? firstCoordinates.coordinates : this.field.options.center

      this.map.zoom = index >= 0 ? 13 : this.field.options.zoom
      this.map.show = true
    },

    placeMarker (e, index = this.localValueIndex) {
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

    removeMarker ({ index }) {
      if (this.field.isMulti) {
        this.localValue.splice(index, 1)
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

    onGeoSearchError () {
      this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.locationSearchFailed'))()
    },

    setDefaultValues () {
      this.localValue = undefined
      this.localValueIndex = undefined
      this.map = {}
    },
  },
}
</script>

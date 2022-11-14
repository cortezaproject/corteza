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
        variant="primary"
        rounded
        class="w-100"
        @click="openMap"
      >
        <font-awesome-icon
          :icon="['fas', 'map-marked-alt']"
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
          number
          :placeholder="$t('latitude')"
        />
        <b-form-input
          v-model="localValue[ctx.index].coordinates[1]"
          type="number"
          number
          :placeholder="$t('longitude')"
        />
      </b-input-group>
    </multi>

    <template v-else>
      <b-input-group>
        <b-form-input
          v-model="localValue.coordinates[0]"
          type="number"
          number
          :placeholder="$t('latitude')"
        />
        <b-form-input
          v-model="localValue.coordinates[1]"
          type="number"
          number
          :placeholder="$t('longitude')"
        />
        <b-input-group-append>
          <b-button
            variant="light"
            rounded
            @click="openMap"
          >
            <font-awesome-icon
              :icon="['fas', 'map-marked-alt']"
            />
          </b-button>
        </b-input-group-append>
      </b-input-group>

      <errors :errors="errors" />
    </template>

    <b-modal
      v-model="map.show"
      size="lg"
      title="Map"
      body-class="p-0"
      hide-header
    >
      <template #modal-footer>
        <h6
          class="w-100"
        >
          {{ $t('clickToPlaceMarker') }}
        </h6>
      </template>

      <l-map
        ref="map"
        :zoom="map.zoom"
        :center="map.center"
        style="height: 75vh; width: 100%; cursor: pointer;"
        @click="placeMarker"
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="map.attribution"
        />
        <l-marker
          v-for="(marker, i) in markers"
          :key="i"
          :lat-lng="marker"
          @click="removeMarker(i)"
        />
      </l-map>
    </b-modal>
  </b-form-group>
</template>
<script>
import base from './base'
import { latLng } from 'leaflet'

export default {
  i18nOptions: {
    namespaces: 'field',
    keyPrefix: 'kind.geometry',
  },

  extends: base,

  data () {
    return {
      localValue: undefined,

      map: {
        show: false,
        zoom: 3,
        center: [30, 30],
        rotation: 0,
        attribution: '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a>',
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
  },

  methods: {
    openMap () {
      const firstCoordinates = (this.field.isMulti ? this.localValue[0] : this.localValue) || {}
      this.map.center = firstCoordinates.coordinates && firstCoordinates.coordinates.length ? firstCoordinates.coordinates : this.field.options.center
      this.map.zoom = this.field.options.zoom
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

    placeMarker (e) {
      let { lat = 0, lng = 0 } = e.latlng || {}
      lat = Math.round(lat * 1e7) / 1e7
      lng = Math.round(lng * 1e7) / 1e7

      if (this.field.isMulti) {
        this.localValue.push({ coordinates: [lat, lng] })
      } else {
        this.localValue = { coordinates: [lat, lng] }
      }
    },

    removeMarker (i) {
      if (this.field.isMulti) {
        this.localValue.splice(i, 1)
      } else {
        this.localValue = { coordinates: [] }
      }
    },
  },
}
</script>

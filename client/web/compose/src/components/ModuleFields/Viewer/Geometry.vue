<template>
  <div>
    <span
      v-for="(c, index) of localValue"
      :key="index"
      :class="{ 'd-block': field.options.multiDelimiter === '\n' }"
    >
      <a
        class="text-primary pointer"
        @click.stop="openMap"
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
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="map.attribution"
        />
        <l-marker
          v-for="(marker, i) in localValue"
          :key="i"
          :lat-lng="marker"
        />
      </l-map>
    </b-modal>

    <errors :errors="errors" />
  </div>
</template>
<script>
import base from './base'
import { latLng } from 'leaflet'

export default {
  extends: base,

  i18nOptions: {
    namespaces: 'field',
    keyPrefix: 'kind.geometry',
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
    openMap () {
      const firstCoordinates = this.localValue[0]
      this.map.center = firstCoordinates && firstCoordinates.length ? firstCoordinates : this.field.options.center
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
  },
}
</script>

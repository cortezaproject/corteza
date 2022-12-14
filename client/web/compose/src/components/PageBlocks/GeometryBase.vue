<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>
    <template v-else>
      <div
        class="w-100 h-100"
        @mouseover="disableMap"
        @mouseleave="enableMap"
      >
        <l-map
          v-if="map"
          ref="map"
          :zoom="map.zoom"
          :center="map.center"
          :min-zoom="map.zoomMin"
          :max-zoom="map.zoomMax"
          :bounds="map.bounds"
          :max-bounds="map.bounds"
          class="w-100 h-100"
        >
          <l-tile-layer
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            :attribution="map.attribution"
          />
          <l-polygon
            v-for="(geometry, i) in geometries"
            :key="`polygon-${i}`"
            :lat-lngs="geometry.map(value => value.geometry)"
            :color="colors[i]"
          />

          <l-marker
            v-for="(marker, i) in localValue"
            :key="`marker-${i}`"
            :lat-lng="marker.value"
            :icon="getIcon(marker)"
          />
        </l-map>
      </div>
    </template>
  </wrap>
</template>

<script>
import { divIcon, latLng, latLngBounds } from 'leaflet'
import {
  LPolygon,
} from 'vue2-leaflet'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { mapGetters, mapActions } from 'vuex'
import { evaluatePrefilter } from 'corteza-webapp-compose/src/lib/record-filter'
import base from './base'

export default {
  components: { LPolygon },

  extends: base,

  data () {
    return {
      map: undefined,

      processing: false,
      show: false,

      geometries: [],
      colors: [],
      markers: [],
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
    }),

    localValue () {
      const values = []

      this.geometries.forEach((geo) => {
        geo.forEach((value) => {
          if (value.displayMarker) {
            value.markers.map(subValue => {
              values.push({ value: this.getLatLng(subValue), color: value.color })
            })
          }
        })
      })

      return values
    },
  },

  watch: {
    'record.recordID': {
      immediate: true,
      handler () {
        this.loadEvents()
      },
    },
    options: {
      deep: true,
      handler () {
        this.loadEvents()
      },
    },
    boundingRect () {
      this.loadEvents()
    },
  },

  created () {
    this.bounds = this.options.bounds
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
    }),

    loadEvents () {
      this.geometries = []

      this.processing = true

      this.colors = this.options.feeds.map(feed => feed.options.color)

      const {
        center,
        zoomStarting,
        zoomMin,
        zoomMax,
      } = this.options

      this.map = {
        bounds: this.options.bounds ? latLngBounds(bounds) : null,
        center,
        zoom: zoomStarting,
        zoomMin,
        zoomMax,
      }

      Promise.all(this.options.feeds.map((feed, idx) => {
        return this.findModuleByID({ namespace: this.namespace, moduleID: feed.options.moduleID })
          .then(module => {
            // Interpolate prefilter variables
            if (feed.options.prefilter) {
              feed.options.prefilter = evaluatePrefilter(feed.options.prefilter, {
                record: this.record,
                recordID: (this.record || {}).recordID || NoID,
                ownerID: (this.record || {}).ownedBy || NoID,
                userID: (this.$auth.user || {}).userID || NoID,
              })
            }

            return compose.PageBlockGeometry.RecordFeed(this.$ComposeAPI, module, this.namespace, feed)
              .then(records => {
                const mapModuleField = module.fields.find(f => f.name === feed.geometryField)

                if (mapModuleField) {
                  this.geometries[idx] = records.map(e => {
                    let geometry = e.values[feed.geometryField]
                    let markers = []

                    if (mapModuleField.isMulti) {
                      geometry = geometry.map(value => this.parseGeometryField(value))
                      markers = geometry
                    } else {
                      geometry = this.parseGeometryField(geometry)
                      markers = [geometry]
                    }

                    return ({
                      title: e.values[feed.titleField],
                      geometry: feed.displayPolygon ? geometry : [],
                      markers,
                      color: feed.options.color,
                      displayMarker: feed.displayMarker,
                    })
                  })
                }
              })
          })
      })).finally(() => {
        this.processing = false
      })
    },

    getIcon (item) {
      item.circleColor = '#ffffff'

      return divIcon({
        className: 'marker-pin',
        html: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 34.892337" height="60" width="40" style="margin-top: -40px;margin-left: -15px;height: 35px;">
          <g transform="translate(-814.59595,-274.38623)">
            <g transform="matrix(1.1855854,0,0,1.1855854,-151.17715,-57.3976)">
              <path d="m 817.11249,282.97118 c -1.25816,1.34277 -2.04623,3.29881 -2.01563,5.13867 0.0639,3.84476 1.79693,5.3002 4.56836,10.59179 0.99832,2.32851 2.04027,4.79237 3.03125,8.87305 0.13772,0.60193 0.27203,1.16104 0.33416,1.20948 0.0621,0.0485 0.19644,-0.51262 0.33416,-1.11455 0.99098,-4.08068 2.03293,-6.54258 3.03125,-8.87109 2.77143,-5.29159 4.50444,-6.74704 4.56836,-10.5918 0.0306,-1.83986 -0.75942,-3.79785 -2.01758,-5.14062 -1.43724,-1.53389 -3.60504,-2.66908 -5.91619,-2.71655 -2.31115,-0.0475 -4.4809,1.08773 -5.91814,2.62162 z" style="fill:${item.color};stroke:${item.color};"/>
              <circle r="3.0355" cy="288.25278" cx="823.03064" id="path3049" style="display:inline;fill:${item.circleColor};"/>
            </g>
          </g>
        </svg>`,
      })
    },

    parseGeometryField (value) {
      return JSON.parse(value || '{"coordinates":[]}').coordinates || []
    },

    getLatLng (coordinates = [undefined, undefined]) {
      const [lat, lng] = coordinates

      if (lat && lng) {
        return latLng(lat, lng)
      }
    },
    disableMap () {
      if (this.editable) this.$refs.map.mapObject._handlers.forEach(handler => handler.disable())
    },
    enableMap () {
      if (this.editable) this.$refs.map.mapObject._handlers.forEach(handler => handler.enable())
    },
  },
}
</script>

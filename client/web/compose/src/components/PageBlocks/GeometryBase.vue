<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
    @refreshBlock="refresh"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <div
      v-else
      class="w-100 h-100"
    >
      <c-map
        v-if="map"
        :map="{
          ...map,
          maxBounds: map.bounds
        }"
        :labels="{
          tooltip: { 'goToCurrentLocation': $t('geometry.tooltip.goToCurrentLocation') }
        }"
        :markers="localValue"
        :disabled="editable"
        hide-geo-search
        :polygons="geometries"
        class="w-100 h-100"
        @on-marker-click="onMarkerCLick"
        @on-geosearch-error="onGeoSearchError"
      />
    </div>
  </wrap>
</template>

<script>
import axios from 'axios'
import { compose, NoID } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import { mapGetters, mapActions } from 'vuex'
import { evaluatePrefilter, isFieldInFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { isNumber } from 'lodash'

import base from './base'

const { CMap } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    CMap,
  },

  extends: base,

  data () {
    return {
      map: undefined,

      processing: false,
      show: false,

      geometries: [],
      colors: [],
      markers: [],

      cancelTokenSource: axios.CancelToken.source(),
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      pages: 'page/set',
    }),

    localValue () {
      const values = []

      this.geometries.forEach((geo) => {
        geo.forEach((value) => {
          if (value.displayMarker) {
            value.markers.map(subValue => {
              if (subValue) {
                values.push({
                  value: subValue || {},
                  color: value.color,
                  recordID: value.recordID,
                  moduleID: value.moduleID,
                })
              }
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
  },

  created () {
    this.bounds = this.options.bounds
    this.refreshBlock(this.refresh)
  },

  mounted () {
    this.createEvents()
  },

  beforeDestroy () {
    this.setDefaultValues()
    this.abortRequests()
    this.destroyEvents()
  },

  methods: {
    ...mapActions({
      findModuleByID: 'module/findByID',
    }),

    createEvents () {
      this.$root.$on('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$on('record-field-change', this.refetchOnPrefilterValueChange)
    },

    refetchOnPrefilterValueChange ({ fieldName }) {
      const { feeds } = this.options

      if (feeds.some(({ options }) => isFieldInFilter(fieldName, options.prefilter))) {
        this.refresh()
      }
    },

    refreshOnRelatedRecordsUpdate ({ moduleID, notPageID }) {
      this.options.feeds.forEach((feed) => {
        if (feed.options.moduleID === moduleID && this.page.pageID !== notPageID) {
          this.refresh()
        }
      })
    },

    loadEvents () {
      this.geometries = []

      this.processing = true

      this.colors = this.options.feeds.map(feed => feed.options.color)

      const {
        bounds,
        center,
        zoomStarting,
        zoomMin,
        zoomMax,
      } = this.options

      this.map = {
        bounds,
        center,
        zoom: zoomStarting,
        zoomMin,
        zoomMax,
      }

      Promise.all(this.options.feeds.filter(f => f.isValid()).map((feed, idx) => {
        return this.findModuleByID({ namespace: this.namespace, moduleID: feed.options.moduleID })
          .then(module => {
            const f = { ...feed } // Clone feed, so we dont modify the original

            // Interpolate prefilter variables
            if (f.options.prefilter) {
              f.options.prefilter = evaluatePrefilter(f.options.prefilter, {
                record: this.record,
                user: this.$auth.user || {},
                recordID: (this.record || {}).recordID || NoID,
                ownerID: (this.record || {}).ownedBy || NoID,
                userID: (this.$auth.user || {}).userID || NoID,
              })
            }

            return compose.PageBlockGeometry.RecordFeed(this.$ComposeAPI, module, this.namespace, f, { cancelToken: this.cancelTokenSource.token })
              .then(records => {
                const mapModuleField = module.fields.find(f => f.name === f.geometryField)

                if (mapModuleField) {
                  this.geometries[idx] = records.map(record => {
                    let geometry = record.values[f.geometryField]
                    let markers = []

                    if (mapModuleField.isMulti) {
                      geometry = geometry.map(value => this.parseGeometryField(value))
                      markers = geometry
                    } else {
                      geometry = this.parseGeometryField(geometry)
                      markers = [geometry]
                    }

                    if (geometry.length && geometry.length === 2) {
                      return ({
                        title: record.values[f.titleField],
                        geometry: f.displayPolygon ? geometry : [],
                        markers,
                        color: f.options.color,
                        displayMarker: f.displayMarker,
                        recordID: record.recordID,
                        moduleID: record.moduleID,
                      })
                    }
                  }).filter(g => g)
                }
              })
          })
      })).finally(() => {
        this.processing = false
      })
    },

    parseGeometryField (value) {
      value = JSON.parse(value || '{"coordinates":[]}').coordinates || []
      return value.every(isNumber) ? value : []
    },

    onMarkerCLick ({ marker: { recordID, moduleID } }) {
      const page = this.pages.find(p => p.moduleID === moduleID)
      if (!page) {
        return
      }

      const route = { name: 'page.record', params: { recordID, pageID: page.pageID } }

      if (this.inModal) {
        this.$root.$emit('show-record-modal', {
          recordID,
          recordPageID: page.pageID,
        })
      } else if (this.options.displayOption === 'newTab') {
        window.open(this.$router.resolve(route).href)
      } else if (this.options.displayOption === 'modal') {
        this.$root.$emit('show-record-modal', {
          recordID,
          recordPageID: page.pageID,
        })
      } else {
        this.$router.push(route)
      }
    },

    refresh () {
      this.loadEvents()
    },

    onGeoSearchError () {
      this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.locationSearchFailed'))()
    },

    setDefaultValues () {
      this.map = undefined
      this.processing = false
      this.show = false
      this.geometries = []
      this.colors = []
      this.markers = []
    },

    abortRequests () {
      this.cancelTokenSource.cancel(`abort-request-${this.block.blockID}`)
    },

    destroyEvents () {
      this.$root.$off('module-records-updated', this.refreshOnRelatedRecordsUpdate)
      this.$root.$off('record-field-change', this.refetchOnPrefilterValueChange)
    },
  },
}
</script>

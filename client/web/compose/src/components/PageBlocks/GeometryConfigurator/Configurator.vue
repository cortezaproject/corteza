<template>
  <div>
    <div>
      <c-map
        :map="mapOptions"
        :labels="{
          tooltip: { 'goToCurrentLocation': $t('geometry.tooltip.goToCurrentLocation') },
        }"
        hide-geo-search
        class="w-100 cursor-pointer"
        style="height: 45vh;"
        @on-bounds-update="boundsUpdated"
        @on-center="updateCenter"
        @on-zoom="options.zoomStarting = $event"
      />

      <b-form-text id="password-help-block">
        {{ $t('geometry.mapHelpText') }}
      </b-form-text>
    </div>

    <hr>

    <b-row>
      <b-col
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('geometry.zoom.zoomStartingLabel')"
          label-class="text-primary"
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
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('geometry.zoom.zoomMinLabel')"
          :description="`${options.zoomMin}`"
          label-class="text-primary"
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
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('geometry.zoom.zoomMaxLabel')"
          :description="`${options.zoomMax}`"
          label-class="text-primary"
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
        cols="12"
        lg="4"
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

      <b-col
        cols="12"
        lg="4"
      >
        <b-form-group
          :label="$t('geometry.bounds.lockBounds')"
          label-class="text-primary"
          class="rounded-left"
        >
          <c-input-checkbox
            v-model="options.lockBounds"
            switch
            :labels="checkboxLabel"
          />
        </b-form-group>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import base from '../base'
import { components } from '@cortezaproject/corteza-vue'
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
      map: {},
      localValue: { coordinates: [] },
      center: [],
      bounds: null,

      checkboxLabel: {
        on: this.$t('general:label.yes'),
        off: this.$t('general:label.no'),
      },
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

    mapOptions: {
      get () {
        return {
          zoom: this.options.zoomStarting,
          minZoom: this.options.zoomMin,
          maxZoom: this.options.zoomMax,
          center: this.options.center,
          bounds: this.bounds,
          maxBounds: this.options.bounds,
        }
      },

      set (options) {
        this.options.zoomStarting = options.zoom
        this.options.zoomMin = options.minZoom
        this.options.zoomMax = options.maxZoom
        this.options.center = options.center
        this.options.bounds = options.center
        this.bounds = options.maxBounds
      },
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
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
        const bounds = this.bounds || {}
        const { _northEast, _southWest } = bounds

        this.options.bounds = [Object.values(_northEast), Object.values(_southWest)]
      } else {
        this.options.bounds = null
      }
    },

    setDefaultValues () {
      this.map = {}
      this.localValue = {}
      this.center = []
      this.bounds = null
    },
  },
}
</script>

<style></style>

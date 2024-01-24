<template>
  <div :class="classes">
    <span
      v-for="(c, index) of localValue"
      :key="index"
      :class="{ 'd-block': field.options.multiDelimiter === '\n' }"
    >
      <a
        class="text-nowrap text-primary pointer"
        @click.stop="openMap(index)"
      >
        {{ c.value[0] }}, {{ c.value[1] }}
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
      <c-map
        :map="map"
        :markers="localValue"
        :hide-current-location-button="field.options.hideCurrentLocationButton"
        :hide-geo-search="field.options.hideGeoSearch"
        :labels="{
          tooltip: { 'goToCurrentLocation': $t('tooltip.goToCurrentLocation') }
        }"
        style="height: 75vh; width: 100%;"
        @on-geosearch-error="onGeoSearchError"
      />
    </b-modal>

    <errors :errors="errors" />
  </div>
</template>
<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
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
      map: {
        show: false,
        zoom: 14,
        center: [30, 30],
        rotation: 0,
        attribution: '&copy; <a target="_blank" rel="noopener noreferrer" href="http://osm.org/copyright">OpenStreetMap</a>',
      },
      localValueIndex: undefined,
    }
  },

  computed: {
    localValue () {
      if (this.field.isMulti) {
        return this.value.map((v, i) => {
          return {
            value: JSON.parse(v || '{"coordinates":[]}').coordinates || [],
            opacity: this.localValueIndex === undefined || i === this.localValueIndex ? 1.0 : 0.6,
          }
        }).filter(c => c && c.value && c.value.length)
      } else {
        return [{ value: JSON.parse(this.value || '{"coordinates":[]}').coordinates || [] }].filter(c => c && c.value && c.value.length)
      }
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    openMap (index) {
      this.localValueIndex = index

      const { value } = this.localValue[index] || {}

      this.map.center = value || this.field.options.center
      this.map.zoom = index >= 0 ? 13 : this.field.options.zoom
      this.map.show = true
    },

    onGeoSearchError () {
      this.toastErrorHandler(this.$t('notification:field-geometry.geolocationErrors.locationSearchFailed'))()
    },

    setDefaultValues () {
      this.map = {}
      this.localValueIndex = undefined
      this.geoSearch = {}
    },
  },
}
</script>

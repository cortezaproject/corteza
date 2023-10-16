<template>
  <div class="d-flex flex-column">
    <b-form-group>
      <b-form-checkbox v-model="f.options.prefillWithCurrentLocation">
        {{ $t('prefillWithCurrentLocation') }}
      </b-form-checkbox>

      <b-form-checkbox v-model="f.options.hideCurrentLocationButton">
        {{ $t('hideCurrentLocationButton') }}
      </b-form-checkbox>

      <b-form-checkbox v-model="f.options.hideGeoSearch">
        {{ $t('hideGeoSearch') }}
      </b-form-checkbox>
    </b-form-group>

    <b-form-group
      :label="$t('initialZoomAndPosition')"
      label-class="text-primary"
      class="mb-0"
    >
      <c-map
        :map="{
          zoom,
          center,
        }"
        :labels="{
          tooltip: { 'goToCurrentLocation': $t('tooltip.goToCurrentLocation') }
        }"
        style="height: 50vh;"
        @on-zoom="f.options.zoom = $event"
        @on-center="f.options.center = $event"
      />
    </b-form-group>
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

  computed: {
    center () {
      return this.f.options.center || [30, 30]
    },

    zoom () {
      return this.f.options.zoom || 3
    },
  },
}
</script>

<template>
  <div>
    <b-form-group
      :label="$t('kind.geometry.initialZoomAndPosition')"
      class="mb-0"
    >
      <l-map
        :zoom="zoom"
        :center="center"
        class="w-100"
        style="height: 60vh;"
        @update:zoom="f.options.zoom = $event"
        @update:center="f.options.center = $event"
      >
        <l-tile-layer
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          :attribution="attribution"
        />
      </l-map>
    </b-form-group>
  </div>
</template>

<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'field',
  },

  extends: base,

  data () {
    return {
      attribution: '&copy; <a target="_blank" href="http://osm.org/copyright">OpenStreetMap</a>',
    }
  },

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

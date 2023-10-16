<template>
  <div>
    <b-button
      variant="link"
      rounded
      class="p-0"
      @click="openMap"
    >
      <font-awesome-icon
        :icon="['fas', 'map-marked-alt']"
      />
    </b-button>

    <b-modal
      v-model="map.show"
      size="lg"
      title="Map"
      body-class="p-0"
      hide-header
      hide-footer
    >
      <c-map
        :map="map"
        :markers="[{ value }]"
        style="height: 75vh; width: 100%;"
        @on-marker-click="removeMarker"
        @on-map-click="placeMarker"
      />
    </b-modal>
  </div>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'

const { CMap } = components

export default {
  components: {
    CMap,
  },

  props: {
    value: {
      type: Array,
      required: true,
    },

    editable: {
      type: Boolean,
      default: false,
    },
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

  methods: {
    openMap () {
      this.map.show = true
    },

    placeMarker ({ latlng = {} }) {
      const { lat = 0, lng = 0 } = latlng
      this.$emit('input', [lat, lng])
    },

    removeMarker () {
      this.$emit('input', [])
    },
  },
}
</script>

<style lang="scss">

</style>

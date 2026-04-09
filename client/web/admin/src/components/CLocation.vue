<template>
  <div>
    <b-button
      v-if="label || editable"
      :variant="label ? 'link' : 'primary'"
      rounded
      :class="label ? 'p-0 border-0' : ''"
      @click="openMap"
    >
      <span v-if="label">
        {{ label }}
      </span>

      <span v-else>
        <font-awesome-icon
          :icon="['fas', 'map-marked-alt']"
        />
        {{ $t('openMap') }}
      </span>
    </b-button>

    <b-modal
      v-model="map.show"
      size="xl"
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
  i18nOptions: {
    namespaces: 'general.label',
  },

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

    label: {
      type: String,
      default: '',
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

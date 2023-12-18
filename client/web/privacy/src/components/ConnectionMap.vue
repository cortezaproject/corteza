<template>
  <c-map
    :map="{
      zoom,
    }"
    hide-geo-search
    hide-current-location-button
    :markers="validMarkerValues"
    style="min-height: 400px; height: 100% !important;"
  >
    <template #marker-tooltip="{ marker }">
      <h5
        class="text-primary"
      >
        {{ $t('server-details') }}
      </h5>
      <b-form-group
        :label="$t('name')"
        label-class="text-primary"
      >
        {{ marker.meta.name }}
      </b-form-group>

      <b-form-group
        :label="$t('location')"
        label-class="text-primary"
      >
        {{ getLocationName(marker) }}
      </b-form-group>
    </template>
  </c-map>
</template>

<script>
import { components } from '@cortezaproject/corteza-vue'

const { CMap } = components

export default {
  i18nOptions: {
    namespaces: 'map',
  },

  components: {
    CMap,
  },

  props: {
    connections: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
      zoom: 2,
    }
  },

  computed: {
    validMarkerValues () {
      return this.connections
        .filter(({ meta = {} }) => {
          const { location = {} } = meta
          const { geometry = {} } = location
          const { coordinates = [] } = geometry

          return coordinates && !!coordinates.length
        })
        .map((connection) => {
          return {
            id: connection.id,
            value: this.getLocationCoordinates(connection),
            ...connection,
          }
        })
    },
  },

  methods: {
    getLocationCoordinates ({ meta = {} }) {
      const { location = {} } = meta
      const { geometry = {} } = location
      return geometry.coordinates
    },

    getLocationName (connection) {
      return connection.meta.location.properties.name || this.$t('unnamed-location')
    },
  },
}
</script>

<style lang="scss">
.vl-style-text {
  color: var(--white);
}
</style>

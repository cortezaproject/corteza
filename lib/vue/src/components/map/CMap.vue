<template>
  <div ref="rootRef" class="c-map relative">
    <!-- Geo Search -->
    <div v-if="!hideGeoSearch" class="geo-search-container">
      <InputText
        v-model="geoSearchQuery"
        :placeholder="$t('field.kind.geometry.geosearchInputPlaceholder')"
        class="w-full"
        @input="onGeoSearch"
      />
      <div v-if="geoSearchResults.length" class="geo-search-results">
        <div
          v-for="(result, idx) in geoSearchResults"
          :key="idx"
          class="geo-search-result"
          @click="placeGeoSearchResult(result)"
        >
          {{ result.label }}
        </div>
      </div>
    </div>

    <!-- Leaflet Map -->
    <LMap
      ref="mapRef"
      :zoom="effectiveZoom"
      :center="effectiveCenter"
      :min-zoom="minZoom || undefined"
      :max-zoom="maxZoom || undefined"
      :max-bounds="maxBounds || undefined"
      :use-global-leaflet="false"
      class="w-full h-full"
      @click="onMapClick"
      @update:center="onCenterUpdate"
      @update:zoom="onZoomUpdate"
      @update:bounds="onBoundsUpdate"
    >
      <LTileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution="&copy; <a target='_blank' href='http://osm.org/copyright'>OpenStreetMap</a>"
      />

      <!-- Polygons -->
      <LPolygon
        v-for="(polygon, i) in validPolygons"
        :key="`polygon-${i}`"
        :lat-lngs="polygon.latLngs"
        :color="polygon.color || '#09344E'"
        :fill-opacity="polygon.fillOpacity ?? 0.2"
      />

      <!-- Markers -->
      <template v-for="(marker, i) in validMarkers" :key="`marker-${i}`">
        <LCircleMarker
          v-if="marker.color"
          :lat-lng="marker.latLng"
          :radius="8"
          :color="marker.color"
          :fill-color="marker.color"
          :fill-opacity="0.8"
          @click="onMarkerClick(i, marker)"
        >
          <LTooltip v-if="marker.title">{{ marker.title }}</LTooltip>
        </LCircleMarker>
        <LMarker
          v-else
          :lat-lng="marker.latLng"
          @click="onMarkerClick(i, marker)"
        >
          <LTooltip v-if="marker.title">{{ marker.title }}</LTooltip>
        </LMarker>
      </template>

      <!-- Geo search marker -->
      <LMarker v-if="geoSearchMarker" :lat-lng="geoSearchMarker" />
    </LMap>

    <!-- Current location button -->
    <button
      v-if="!hideCurrentLocationButton"
      v-tooltip.top="$t('field.kind.geometry.tooltip.goToCurrentLocation')"
      class="current-location-btn"
      @click="goToCurrentLocation"
    >
      <FontAwesomeIcon :icon="faLocationCrosshairs" />
    </button>
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import 'leaflet/dist/leaflet.css'
import { LMap, LTileLayer, LMarker, LCircleMarker, LPolygon, LTooltip } from '@vue-leaflet/vue-leaflet'
import { OpenStreetMapProvider } from 'leaflet-geosearch'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faLocationCrosshairs } from '@fortawesome/free-solid-svg-icons'

const { t: $t } = useI18n()

const props = defineProps({
  center: {
    type: Array,
    default: () => [30, 30],
  },
  zoom: {
    type: Number,
    default: 3,
  },
  markers: {
    type: Array,
    default: () => [],
  },
  hideGeoSearch: {
    type: Boolean,
    default: false,
  },
  hideCurrentLocationButton: {
    type: Boolean,
    default: false,
  },
  minZoom: {
    type: Number,
    default: 0,
  },
  maxZoom: {
    type: Number,
    default: 0,
  },
  maxBounds: {
    type: Array,
    default: null,
  },
  polygons: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits([
  'map-click',
  'marker-click',
  'location-found',
  'update:center',
  'update:zoom',
  'update:bounds',
])

const mapRef = ref(null)
const rootRef = ref(null)

// Geo search state
const geoSearchQuery = ref('')
const geoSearchResults = ref([])
const geoSearchMarker = ref(null)
let searchTimeout = null
let resizeObserver = null

const provider = new OpenStreetMapProvider()

// Map center/zoom can be overridden from props
const effectiveCenter = computed(() => {
  if (Array.isArray(props.center) && props.center.length === 2) {
    const [lat, lng] = props.center
    if (typeof lat === 'number' && typeof lng === 'number' && !isNaN(lat) && !isNaN(lng)) {
      return props.center
    }
  }
  return [30, 30]
})

const effectiveZoom = computed(() => props.zoom || 3)

// Filter markers to only valid lat/lng pairs
const validMarkers = computed(() =>
  props.markers
    .filter(m => m.value && Array.isArray(m.value) && m.value.length === 2)
    .filter(m => typeof m.value[0] === 'number' && typeof m.value[1] === 'number')
    .map(m => ({
      ...m,
      latLng: m.value,
    })),
)

// Polygons: each entry { latLngs: [[lat,lng], ...], color?, fillOpacity? }
const validPolygons = computed(() =>
  (props.polygons || [])
    .filter(p => Array.isArray(p?.latLngs) && p.latLngs.length >= 3)
    .map(p => ({
      ...p,
      latLngs: p.latLngs.filter(
        pt => Array.isArray(pt) && pt.length === 2 &&
          typeof pt[0] === 'number' && typeof pt[1] === 'number',
      ),
    }))
    .filter(p => p.latLngs.length >= 3),
)

function onMapClick(e) {
  emit('map-click', e)
}

function normalizeLatLng(center) {
  if (Array.isArray(center) && center.length === 2) return [center[0], center[1]]
  if (center && typeof center.lat === 'number') return [center.lat, center.lng]
  return null
}

function onCenterUpdate(center) {
  const next = normalizeLatLng(center)
  if (!next) return
  const cur = normalizeLatLng(props.center)
  if (cur && cur[0] === next[0] && cur[1] === next[1]) return
  emit('update:center', next)
}

function onZoomUpdate(zoom) {
  if (zoom === props.zoom) return
  emit('update:zoom', zoom)
}

function onBoundsUpdate(bounds) {
  if (!bounds) return
  const sw = typeof bounds.getSouthWest === 'function' ? bounds.getSouthWest() : bounds._southWest
  const ne = typeof bounds.getNorthEast === 'function' ? bounds.getNorthEast() : bounds._northEast
  if (!sw || !ne) return
  emit('update:bounds', [[sw.lat, sw.lng], [ne.lat, ne.lng]])
}

function onMarkerClick(index, marker) {
  emit('marker-click', { index, marker })
}

function onGeoSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  if (!geoSearchQuery.value) {
    geoSearchResults.value = []
    geoSearchMarker.value = null
    return
  }
  searchTimeout = setTimeout(async () => {
    try {
      const results = await provider.search({ query: geoSearchQuery.value })
      geoSearchResults.value = results.map(r => ({
        label: r.label,
        lat: r.y,
        lng: r.x,
      }))
    } catch {
      geoSearchResults.value = []
    }
  }, 300)
}

function placeGeoSearchResult(result) {
  geoSearchMarker.value = [result.lat, result.lng]
  geoSearchResults.value = []
  geoSearchQuery.value = result.label

  // Fly to the result
  const map = mapRef.value?.leafletObject
  if (map) {
    map.flyTo([result.lat, result.lng], 15, { animate: true })
  }

  // Emit as a click so the editor can place a marker
  emit('map-click', { latlng: { lat: result.lat, lng: result.lng } })
}

function goToCurrentLocation() {
  if (!navigator.geolocation) return

  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      const map = mapRef.value?.leafletObject
      if (map) {
        map.flyTo([coords.latitude, coords.longitude], 15)
      }
      emit('location-found', { latlng: { lat: coords.latitude, lng: coords.longitude } })
    },
    () => {
      // Geolocation error — silent fail
    },
  )
}

// Allow parent to trigger invalidateSize (e.g. after dialog open)
function invalidateSize() {
  const map = mapRef.value?.leafletObject
  if (map) {
    map.invalidateSize()
  }
}

function fitBounds (bounds, options = {}) {
  const map = mapRef.value?.leafletObject
  if (map && Array.isArray(bounds) && bounds.length === 2) {
    map.fitBounds(bounds, options)
  }
}

defineExpose({ invalidateSize, fitBounds })

// Auto-detect visibility/size changes via ResizeObserver.
// When a map is inside a hidden container (tab, dialog, modal) and becomes
// visible, the container resizes from 0×0 → actual size. We catch that and
// call invalidateSize() so Leaflet recalculates its viewport.
onMounted(() => {
  if (rootRef.value && typeof ResizeObserver !== 'undefined') {
    resizeObserver = new ResizeObserver(() => {
      nextTick(() => invalidateSize())
    })
    resizeObserver.observe(rootRef.value)
  }
})

onBeforeUnmount(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
})

// Watch center changes to re-center the map
watch(
  () => props.center,
  newCenter => {
    if (Array.isArray(newCenter) && newCenter.length === 2) {
      const map = mapRef.value?.leafletObject
      if (map) {
        map.panTo(newCenter)
      }
    }
  },
)
</script>

<style>
.c-map {
  overflow: hidden;
}

.geo-search-container {
  position: absolute;
  z-index: 1000;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
  width: 50%;
  max-width: 400px;
}

.geo-search-results {
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: 0 0 0.375rem 0.375rem;
  max-height: 200px;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.geo-search-result {
  padding: 0.5rem 0.75rem;
  font-size: 0.85rem;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.geo-search-result:hover {
  background: var(--p-content-hover-background);
}

.current-location-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  padding: 0;
  border: 2px solid rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  background: #fff;
  color: #333;
  cursor: pointer;
  background-clip: padding-box;
}

.current-location-btn:hover {
  background: #f4f4f4;
}
</style>

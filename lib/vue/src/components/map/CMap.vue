<template>
  <div ref="rootRef" class="c-map relative">
    <!-- Geo Search -->
    <div v-if="!hideGeoSearch" class="geo-search-container">
      <InputText
        v-model="geoSearchQuery"
        :placeholder="$t('field.kind.geometry.geosearchInputPlaceholder')"
        class="w-full"
        autocomplete="off"
        name="c-map-geosearch"
        spellcheck="false"
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

    <!-- Leaflet Map (rendering only — markers/polygons & events managed below) -->
    <!-- use-global-leaflet must be true so vue-leaflet and our direct L.*
         calls below share one leaflet module instance (else bounds math
         operates across two prototype chains and throws). -->
    <LMap
      ref="mapRef"
      :zoom="effectiveZoom"
      :center="effectiveCenter"
      :min-zoom="minZoom || undefined"
      :max-zoom="maxZoom || undefined"
      :max-bounds="maxBounds || undefined"
      :use-global-leaflet="true"
      class="w-full h-full"
      @ready="onMapReady"
    >
      <LTileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution="&copy; <a target='_blank' href='http://osm.org/copyright'>OpenStreetMap</a>"
      />
    </LMap>

    <!-- Floating map controls (zoom + current location) -->
    <div class="map-controls">
      <div class="map-controls-group">
        <Button
          v-tooltip.left="$t('field.kind.geometry.tooltip.zoomIn')"
          icon="pi pi-plus"
          severity="secondary"
          size="small"
          aria-label="Zoom in"
          @click="onZoomIn"
        />
        <Button
          v-tooltip.left="$t('field.kind.geometry.tooltip.zoomOut')"
          icon="pi pi-minus"
          severity="secondary"
          size="small"
          aria-label="Zoom out"
          @click="onZoomOut"
        />
      </div>

      <Button
        v-if="!hideCurrentLocationButton"
        v-tooltip.left="$t('field.kind.geometry.tooltip.goToCurrentLocation')"
        icon="pi pi-map-marker"
        severity="secondary"
        size="small"
        aria-label="Go to current location"
        @click="goToCurrentLocation"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import { LMap, LTileLayer } from '@vue-leaflet/vue-leaflet'
import { OpenStreetMapProvider } from 'leaflet-geosearch'

const { t: $t } = useI18n()

const props = defineProps({
  center: { type: Array, default: () => [30, 30] },
  zoom: { type: Number, default: 3 },
  markers: { type: Array, default: () => [] },
  polygons: { type: Array, default: () => [] },
  hideGeoSearch: { type: Boolean, default: false },
  hideCurrentLocationButton: { type: Boolean, default: false },
  minZoom: { type: Number, default: 0 },
  maxZoom: { type: Number, default: 0 },
  maxBounds: { type: Array, default: null },
})

const emit = defineEmits([
  'ready',
  'map-click',
  'marker-click',
  'location-found',
  'update:center',
  'update:zoom',
  'update:bounds',
])

const mapRef = ref(null)
const rootRef = ref(null)
const leafletMap = ref(null)

// Persistent layer that holds whatever markers/polygons CMap is responsible for
let renderLayer = null
let resizeObserver = null

// Geo search state
const geoSearchQuery = ref('')
const geoSearchResults = ref([])
let searchTimeout = null
const provider = new OpenStreetMapProvider()

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

const validMarkers = computed(() =>
  (props.markers || [])
    .filter(m => m?.value && Array.isArray(m.value) && m.value.length === 2)
    .filter(m => typeof m.value[0] === 'number' && typeof m.value[1] === 'number'),
)

const validPolygons = computed(() =>
  (props.polygons || [])
    .filter(p => Array.isArray(p?.latLngs) && p.latLngs.length >= 3)
    .map(p => ({
      ...p,
      latLngs: p.latLngs.filter(
        pt =>
          Array.isArray(pt) &&
          pt.length === 2 &&
          typeof pt[0] === 'number' &&
          typeof pt[1] === 'number',
      ),
    }))
    .filter(p => p.latLngs.length >= 3),
)

function onMapReady(map) {
  const lmap = map || mapRef.value?.leafletObject
  if (!lmap) return
  leafletMap.value = lmap

  // Hide leaflet's default +/- zoom control — we render PrimeVue buttons
  // overlay-style instead so they respect the active theme (dark/light).
  if (lmap.zoomControl) {
    lmap.removeControl(lmap.zoomControl)
  }

  renderLayer = L.layerGroup().addTo(lmap)

  // vue-leaflet 0.10 does not reliably forward leaflet events to Vue,
  // so bind them directly on the leaflet map.
  lmap.on('click', handleMapClick)
  lmap.on('moveend', handleMoveEnd)

  renderAll()
  emit('ready', lmap)
}

function onZoomIn() {
  const lmap = leafletMap.value
  if (lmap) lmap.zoomIn()
}

function onZoomOut() {
  const lmap = leafletMap.value
  if (lmap) lmap.zoomOut()
}

function handleMapClick(e) {
  // If propagation was stopped (e.g. by a marker click), bail out.
  if (e.originalEvent?.defaultPrevented) return
  emit('map-click', e)
}

function handleMoveEnd() {
  const lmap = leafletMap.value
  if (!lmap) return
  emit('update:zoom', lmap.getZoom())
  emit('update:center', [lmap.getCenter().lat, lmap.getCenter().lng])
  const bounds = lmap.getBounds()
  emit('update:bounds', [
    [bounds.getSouthWest().lat, bounds.getSouthWest().lng],
    [bounds.getNorthEast().lat, bounds.getNorthEast().lng],
  ])
}

function clearRenderLayer() {
  if (renderLayer) renderLayer.clearLayers()
}

// Build a colored teardrop pin via divIcon. Same silhouette as leaflet's
// default marker, just recolored. Used when a consumer supplies `color`.
function teardropIcon(color) {
  const safeColor = String(color || '#3b82f6').replace(/"/g, '')
  const html = `
    <svg xmlns="http://www.w3.org/2000/svg" width="25" height="41" viewBox="0 0 25 41">
      <path d="M12.5 0.5 C5.87 0.5 0.5 5.87 0.5 12.5 C0.5 22 12.5 40 12.5 40 C12.5 40 24.5 22 24.5 12.5 C24.5 5.87 19.13 0.5 12.5 0.5 Z"
            fill="${safeColor}" stroke="#ffffff" stroke-width="1.5" />
      <circle cx="12.5" cy="12.5" r="4" fill="#ffffff" />
    </svg>`
  return L.divIcon({
    html,
    className: 'c-map-pin',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    tooltipAnchor: [12, -28],
  })
}

function renderMarkers() {
  if (!renderLayer) return
  validMarkers.value.forEach((marker, index) => {
    const latLng = marker.value
    const layer = marker.color
      ? L.marker(latLng, { icon: teardropIcon(marker.color) })
      : L.marker(latLng)
    if (marker.title) {
      layer.bindTooltip(marker.title)
    }
    layer.on('click', e => {
      // Prevent the click from also triggering map-click underneath.
      L.DomEvent.stopPropagation(e)
      if (e.originalEvent) e.originalEvent.preventDefault()
      emit('marker-click', { index, marker, event: e })
    })
    layer.addTo(renderLayer)
  })
}

function renderPolygons() {
  if (!renderLayer) return
  validPolygons.value.forEach(polygon => {
    L.polygon(polygon.latLngs, {
      color: polygon.color || '#09344E',
      fillOpacity: polygon.fillOpacity ?? 0.2,
    }).addTo(renderLayer)
  })
}

function renderAll() {
  clearRenderLayer()
  renderMarkers()
  renderPolygons()
}

watch(validMarkers, renderAll, { deep: true })
watch(validPolygons, renderAll, { deep: true })

function onGeoSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  if (!geoSearchQuery.value) {
    geoSearchResults.value = []
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
  const lmap = leafletMap.value
  if (!lmap) return

  geoSearchResults.value = []
  geoSearchQuery.value = result.label
  lmap.flyTo([result.lat, result.lng], 15, { animate: true })
  // Drop the consumer's marker at the picked location. Goes through the
  // normal map-click path — for CInputLocation this only updates draftCoords
  // (in-dialog), the lat/lng value is still committed only on Save.
  emit('map-click', { latlng: { lat: result.lat, lng: result.lng } })
}

function goToCurrentLocation() {
  if (!navigator.geolocation) return

  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      const lmap = leafletMap.value
      if (lmap) lmap.flyTo([coords.latitude, coords.longitude], 15)
      emit('location-found', { latlng: { lat: coords.latitude, lng: coords.longitude } })
    },
    () => {
      // Geolocation error — silent fail
    },
  )
}

function invalidateSize() {
  const lmap = leafletMap.value
  if (lmap) lmap.invalidateSize()
}

function fitBounds(bounds, options = {}) {
  const lmap = leafletMap.value
  if (lmap && Array.isArray(bounds) && bounds.length === 2) {
    lmap.fitBounds(bounds, options)
  }
}

defineExpose({ invalidateSize, fitBounds })

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
  const lmap = leafletMap.value
  if (lmap) {
    lmap.off('click', handleMapClick)
    lmap.off('moveend', handleMoveEnd)
  }
  if (renderLayer) {
    renderLayer.clearLayers()
    renderLayer = null
  }
})

watch(
  () => props.center,
  newCenter => {
    const lmap = leafletMap.value
    if (lmap && Array.isArray(newCenter) && newCenter.length === 2) {
      lmap.panTo(newCenter)
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

/* Floating overlay containing zoom +/- and current-location buttons. The
   PrimeVue buttons inside use theme tokens so they respect dark mode. */
.map-controls {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
}

.map-controls-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

/* Colored teardrop pin (divIcon). Leaflet adds a transparent box by default; we
   strip it so only the SVG paints. */
.c-map-pin {
  background: transparent;
  border: 0;
}
.c-map-pin svg {
  display: block;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.35));
}
</style>

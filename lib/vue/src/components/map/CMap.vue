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

    <!-- Leaflet renders into this element. Layers, events and controls are all
         driven imperatively below: leaflet owns the viewport, props only ever
         nudge it, and every emit is de-duplicated so a parent that writes back
         what we emitted cannot start a pan/emit loop. -->
    <div ref="mapEl" class="c-map-canvas w-full h-full" />

    <!-- Floating map controls (zoom + current location) -->
    <div class="map-controls">
      <div class="map-controls-group">
        <Button
          v-tooltip.left="$t('field.kind.geometry.tooltip.zoomIn')"
          icon="pi pi-plus"
          severity="secondary"
          size="small"
          :disabled="!canZoomIn"
          aria-label="Zoom in"
          @click="onZoomIn"
        />
        <Button
          v-tooltip.left="$t('field.kind.geometry.tooltip.zoomOut')"
          icon="pi pi-minus"
          severity="secondary"
          size="small"
          :disabled="!canZoomOut"
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
import { computed, ref, shallowRef, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import { OpenStreetMapProvider } from 'leaflet-geosearch'
import { parseBounds, parseLatLng, sameBounds, sameLatLng } from './geo'

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
  // Pins the viewport centre: no dragging, no box zoom, no keyboard panning,
  // and zoom anchors on the centre instead of the pointer. Configurators turn
  // this on once bounds are locked — the saved area is then exactly what the
  // preview shows, and nothing but unlocking can move it.
  disablePan: { type: Boolean, default: false },
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

const rootRef = ref(null)
const mapEl = ref(null)
// shallowRef: leaflet keeps identity-sensitive internal references, so the map
// must never be wrapped in a deep reactive proxy.
const leafletMap = shallowRef(null)

// Persistent layer that holds whatever markers/polygons CMap is responsible for
let renderLayer = null
let resizeObserver = null

// Last view we told the parent about. Also seeded before a prop-driven move, so
// the moveend that move causes is recognised as our own doing.
let lastCenter = null
let lastZoom = null
let lastBounds = null
// What we last handed to setMaxBounds — re-applying an unchanged box would
// re-pan the view for nothing.
let appliedMaxBounds = null

const currentZoom = ref(0)
// Read off the map rather than guessed from props: with no maxZoom set, the
// ceiling comes from the tile layer.
const zoomFloor = ref(0)
const zoomCeil = ref(18)

// Geo search state
const geoSearchQuery = ref('')
const geoSearchResults = ref([])
let searchTimeout = null
const provider = new OpenStreetMapProvider()

const effectiveCenter = computed(() => parseLatLng(props.center) || [30, 30])
const effectiveZoom = computed(() => props.zoom || 3)
const effectiveMaxBounds = computed(() => parseBounds(props.maxBounds))

const canZoomIn = computed(() => currentZoom.value < zoomCeil.value)
const canZoomOut = computed(() => currentZoom.value > zoomFloor.value)

const validMarkers = computed(() =>
  (props.markers || []).filter(m => parseLatLng(m?.value) !== null),
)

const validPolygons = computed(() =>
  (props.polygons || [])
    .filter(p => Array.isArray(p?.latLngs))
    .map(p => ({ ...p, latLngs: p.latLngs.filter(pt => parseLatLng(pt) !== null) }))
    .filter(p => p.latLngs.length >= 3),
)

function initMap() {
  if (!mapEl.value || leafletMap.value) return

  const map = L.map(mapEl.value, {
    center: effectiveCenter.value,
    zoom: effectiveZoom.value,
    // Our own themed buttons stand in for leaflet's default control, which
    // ignores the active theme.
    zoomControl: false,
    minZoom: props.minZoom || undefined,
    maxZoom: props.maxZoom || undefined,
    maxBounds: effectiveMaxBounds.value || undefined,
  })

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: "&copy; <a target='_blank' href='http://osm.org/copyright'>OpenStreetMap</a>",
  }).addTo(map)

  renderLayer = L.layerGroup().addTo(map)

  map.on('click', handleMapClick)
  map.on('moveend zoomend', handleViewChange)

  leafletMap.value = map
  appliedMaxBounds = effectiveMaxBounds.value
  snapshotView()
  applyInteraction()
  renderAll()
  // The viewport is a function of the container size, so only the map knows it.
  // Reporting it up front means a consumer that captures bounds has them before
  // the first gesture rather than only after one.
  emit('update:bounds', lastBounds)
  emit('ready', map)
}

/** Read the map's view without emitting — used to seed the echo guard. */
function snapshotView() {
  const map = leafletMap.value
  if (!map) return
  const c = map.getCenter()
  lastCenter = [c.lat, c.lng]
  lastZoom = map.getZoom()
  lastBounds = readBounds()
  currentZoom.value = lastZoom
  syncZoomLimits()
}

function syncZoomLimits() {
  const map = leafletMap.value
  if (!map) return
  const floor = map.getMinZoom()
  const ceil = map.getMaxZoom()
  zoomFloor.value = Number.isFinite(floor) ? floor : 0
  zoomCeil.value = Number.isFinite(ceil) ? ceil : 18
}

function readBounds() {
  const map = leafletMap.value
  if (!map) return null
  const b = map.getBounds()
  return [
    [b.getSouthWest().lat, b.getSouthWest().lng],
    [b.getNorthEast().lat, b.getNorthEast().lng],
  ]
}

/**
 * The single place the map reports its viewport. Each value is emitted only
 * when it actually differs from the last one the parent was given, so a parent
 * that stores what we emit and hands it straight back settles after one round
 * instead of oscillating.
 */
function handleViewChange() {
  const map = leafletMap.value
  if (!map) return

  const c = map.getCenter()
  const center = [c.lat, c.lng]
  if (!sameLatLng(center, lastCenter)) {
    lastCenter = center
    emit('update:center', center)
  }

  const zoom = map.getZoom()
  if (zoom !== lastZoom) {
    lastZoom = zoom
    currentZoom.value = zoom
    emit('update:zoom', zoom)
  }

  const bounds = readBounds()
  if (!sameBounds(bounds, lastBounds)) {
    lastBounds = bounds
    emit('update:bounds', bounds)
  }
}

/**
 * Panning is what a locked view must not allow; zooming still may. A zoom
 * anchored on the pointer shifts the centre, which is a pan by another name, so
 * while pan is disabled every zoom gesture anchors on the centre instead.
 */
function applyInteraction() {
  const map = leafletMap.value
  if (!map) return

  const locked = props.disablePan
  const anchor = locked ? 'center' : true

  if (locked) {
    map.dragging.disable()
    map.boxZoom.disable()
    map.keyboard.disable()
  } else {
    map.dragging.enable()
    map.boxZoom.enable()
    map.keyboard.enable()
  }

  // Handlers read these at gesture time, so flipping them is enough.
  map.options.scrollWheelZoom = anchor
  map.options.doubleClickZoom = anchor
  map.options.touchZoom = anchor
}

function onZoomIn() {
  leafletMap.value?.zoomIn()
}

function onZoomOut() {
  leafletMap.value?.zoomOut()
}

function handleMapClick(e) {
  // If propagation was stopped (e.g. by a marker click), bail out.
  if (e.originalEvent?.defaultPrevented) return
  emit('map-click', e)
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
      weight: polygon.weight ?? 3,
      dashArray: polygon.dashArray || undefined,
      fillOpacity: polygon.fillOpacity ?? 0.2,
      interactive: polygon.interactive ?? true,
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
  const map = leafletMap.value
  if (!map) return

  geoSearchResults.value = []
  geoSearchQuery.value = result.label
  // A locked view stays where it is; the pick is still reported so a consumer
  // placing a marker gets it.
  if (!props.disablePan) {
    map.flyTo([result.lat, result.lng], 15, { animate: true })
  }
  // Drop the consumer's marker at the picked location. Goes through the
  // normal map-click path — for CInputLocation this only updates draftCoords
  // (in-dialog), the lat/lng value is still committed only on Save.
  emit('map-click', { latlng: { lat: result.lat, lng: result.lng } })
}

function goToCurrentLocation() {
  if (!navigator.geolocation) return

  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      const map = leafletMap.value
      if (map && !props.disablePan) {
        map.flyTo([coords.latitude, coords.longitude], 15)
      }
      emit('location-found', { latlng: { lat: coords.latitude, lng: coords.longitude } })
    },
    () => {
      // Geolocation error — silent fail
    },
  )
}

function invalidateSize() {
  leafletMap.value?.invalidateSize()
}

function fitBounds(bounds, options = {}) {
  const map = leafletMap.value
  const box = parseBounds(bounds)
  if (map && box) {
    map.fitBounds(box, options)
  }
}

defineExpose({ invalidateSize, fitBounds })

onMounted(() => {
  initMap()

  if (rootRef.value && typeof ResizeObserver !== 'undefined') {
    resizeObserver = new ResizeObserver(() => {
      nextTick(() => invalidateSize())
    })
    resizeObserver.observe(rootRef.value)
  }
})

onBeforeUnmount(() => {
  if (searchTimeout) clearTimeout(searchTimeout)
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
  const map = leafletMap.value
  if (map) {
    map.off()
    map.remove()
  }
  leafletMap.value = null
  renderLayer = null
})

watch(
  () => props.center,
  center => {
    const map = leafletMap.value
    const target = parseLatLng(center)
    if (!map || !target) return
    const c = map.getCenter()
    // Ignore our own emitted centre coming back through the parent.
    if (sameLatLng([c.lat, c.lng], target)) return
    lastCenter = target
    map.panTo(target)
  },
)

watch(
  () => props.zoom,
  zoom => {
    const map = leafletMap.value
    if (!map || !zoom || zoom === map.getZoom()) return
    lastZoom = zoom
    map.setZoom(zoom)
  },
)

watch(
  () => [props.minZoom, props.maxZoom],
  ([min, max]) => {
    const map = leafletMap.value
    if (!map) return
    map.setMinZoom(min || undefined)
    map.setMaxZoom(max || undefined)
    syncZoomLimits()
  },
)

watch(effectiveMaxBounds, bounds => {
  const map = leafletMap.value
  // setMaxBounds re-pans the view into the new box, which is a real move worth
  // reporting — but re-applying an unchanged box would move it for nothing.
  if (!map || sameBounds(bounds, appliedMaxBounds)) return
  appliedMaxBounds = bounds
  map.setMaxBounds(bounds || null)
})

watch(() => props.disablePan, applyInteraction)
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

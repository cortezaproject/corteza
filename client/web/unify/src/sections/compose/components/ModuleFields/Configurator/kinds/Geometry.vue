<template>
  <div class="flex flex-col gap-3">
    <CFormGroup
      :label="$t('field.kind.geometry.startingView')"
      :description="
        lockBounds
          ? $t('field.kind.geometry.boundsLockedHint')
          : $t('field.kind.geometry.mapHelpText')
      "
    >
      <CMap
        :center="center"
        :zoom="zoom"
        :max-bounds="lockedBounds"
        :polygons="boundsOutline"
        :hide-geo-search="hideGeoSearch"
        :hide-current-location-button="hideCurrentLocationButton"
        style="height: 40vh"
        @update:center="onMapCenter"
        @update:zoom="onMapZoom"
        @update:bounds="onMapBounds"
      />
    </CFormGroup>

    <div class="grid grid-cols-2 gap-3 items-start">
      <CInputToggleCard
        v-model="prefillWithCurrentLocation"
        :label="$t('field.kind.geometry.prefillWithCurrentLocation')"
        :description="$t('field.kind.geometry.prefillWithCurrentLocationDescription')"
      />
      <CInputToggleCard
        v-model="hideCurrentLocationButton"
        :label="$t('field.kind.geometry.hideCurrentLocationButton')"
        :description="$t('field.kind.geometry.hideCurrentLocationButtonDescription')"
      />
      <CInputToggleCard
        v-model="hideGeoSearch"
        :label="$t('field.kind.geometry.hideGeoSearch')"
        :description="$t('field.kind.geometry.hideGeoSearchDescription')"
      />
      <CInputToggleCard
        :model-value="lockBounds"
        :label="$t('field.kind.geometry.lockBounds')"
        :description="$t('field.kind.geometry.lockBoundsDescription')"
        @update:model-value="onLockBoundsToggle"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { components, mapGeo } from '@planetcrust/human-vue'

const { CInputToggleCard, CMap } = components

const field = inject('fieldDraft')

function getPrimaryColor() {
  return (
    getComputedStyle(document.documentElement).getPropertyValue('--p-primary-color').trim() ||
    '#09344E'
  )
}

function updateOptions(patch) {
  Object.assign(field.value.options, patch)
}

const center = computed(() => mapGeo.parseLatLng(field.value.options?.center) || [30, 30])

const zoom = computed(() => field.value.options?.zoom || 3)

function onMapCenter(next) {
  const rounded = mapGeo.roundLatLng(next)
  if (rounded) updateOptions({ center: rounded })
}

function onMapZoom(z) {
  updateOptions({ zoom: z })
}

const prefillWithCurrentLocation = computed({
  get: () => !!field.value.options?.prefillWithCurrentLocation,
  set: v => updateOptions({ prefillWithCurrentLocation: v }),
})

const hideCurrentLocationButton = computed({
  get: () => !!field.value.options?.hideCurrentLocationButton,
  set: v => updateOptions({ hideCurrentLocationButton: v }),
})

const hideGeoSearch = computed({
  get: () => !!field.value.options?.hideGeoSearch,
  set: v => updateOptions({ hideGeoSearch: v }),
})

const lockBounds = computed(() => !!field.value.options?.lockBounds)

// options.bounds is [[swLat, swLng], [neLat, neLng]], and holds the locked area
// only — unlocking clears it.
const lockedBounds = computed(() =>
  lockBounds.value ? mapGeo.parseBounds(field.value.options?.bounds) : null,
)

// Outlines the locked area, which stays visible when zooming out past it.
const boundsOutline = computed(() => {
  const ring = mapGeo.boundsRing(lockedBounds.value)
  if (!ring) return []
  return [
    {
      latLngs: ring,
      color: getPrimaryColor(),
      weight: 2,
      dashArray: '6 4',
      fillOpacity: 0.05,
      interactive: false,
    },
  ]
})

// The map's live viewport, reported by CMap on ready and on every move.
const currentMapBounds = ref(null)

function onMapBounds(b) {
  currentMapBounds.value = b
}

// Locking captures the area on screen right now and freezes the preview on it;
// unlocking gives the map back and drops the area, since nothing is bounded any
// more.
function onLockBoundsToggle(v) {
  if (!v) {
    updateOptions({ lockBounds: false, bounds: null })
    return
  }

  const bounds = mapGeo.roundBounds(currentMapBounds.value)
  updateOptions({ lockBounds: true, ...(bounds ? { bounds } : {}) })
}

onMounted(() => {
  if (!field.value.options.center) {
    field.value.options.center = [30, 30]
  }
  if (!field.value.options.zoom) {
    field.value.options.zoom = 3
  }
})
</script>

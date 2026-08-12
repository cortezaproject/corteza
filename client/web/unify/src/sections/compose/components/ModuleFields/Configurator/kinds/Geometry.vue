<template>
  <div class="flex flex-col gap-3">
    <!-- Starting view: pan/zoom to set starting center and zoom -->
    <div class="flex flex-col gap-2">
      <label class="text-primary font-medium text-sm">
        {{ $t('field.kind.geometry.startingView') }}
      </label>
      <CMap
        :center="center"
        :zoom="zoom"
        :max-bounds="lockBounds ? lockedBounds : null"
        :hide-geo-search="hideGeoSearch"
        :hide-current-location-button="hideCurrentLocationButton"
        style="height: 40vh"
        @update:center="onMapCenter"
        @update:zoom="onMapZoom"
        @update:bounds="onMapBounds"
      />
    </div>

    <!-- Toggles -->
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
import { components } from '@planetcrust/human-vue'
import CMap from '@planetcrust/human-vue/src/components/map/CMap.vue'

const { CInputToggleCard } = components

const field = inject('fieldDraft')

const center = computed(() => {
  const c = field.value.options?.center
  return Array.isArray(c) && c.length === 2 ? c : [30, 30]
})

const zoom = computed(() => field.value.options?.zoom || 3)

function onMapCenter([lat, lng]) {
  field.value.options.center = [Math.round(lat * 1e6) / 1e6, Math.round(lng * 1e6) / 1e6]
}

function onMapZoom(z) {
  field.value.options.zoom = z
}

const prefillWithCurrentLocation = computed({
  get: () => !!field.value.options?.prefillWithCurrentLocation,
  set: v => {
    field.value.options.prefillWithCurrentLocation = v
  },
})

const hideCurrentLocationButton = computed({
  get: () => !!field.value.options?.hideCurrentLocationButton,
  set: v => {
    field.value.options.hideCurrentLocationButton = v
  },
})

const hideGeoSearch = computed({
  get: () => !!field.value.options?.hideGeoSearch,
  set: v => {
    field.value.options.hideGeoSearch = v
  },
})

const lockBounds = computed(() => !!field.value.options?.lockBounds)

const lockedBounds = computed(() => {
  const b = field.value.options?.bounds
  if (Array.isArray(b) && b.length === 2 && b.every(p => Array.isArray(p) && p.length === 2)) {
    return b
  }
  return null
})

const currentMapBounds = ref(null)

function onMapBounds(b) {
  currentMapBounds.value = b
}

function onLockBoundsToggle(v) {
  if (v) {
    const b = currentMapBounds.value || lockedBounds.value
    if (b) {
      field.value.options.bounds = b
    }
    field.value.options.lockBounds = true
  } else {
    field.value.options.lockBounds = false
  }
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

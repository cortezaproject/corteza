<template>
  <div class="flex flex-col gap-2">
    <!-- Lat/Lng manual input -->
    <div class="flex gap-2 items-center overflow-hidden">
      <InputGroup class="flex-1 min-w-0">
        <InputNumber
          :model-value="latitude"
          :disabled="disabled"
          :step="0.000001"
          :max-fraction-digits="7"
          :placeholder="$t('field.kind.geometry.latitude')"
          class="flex-1 min-w-0"
          @update:model-value="onLatChange"
        />
        <InputNumber
          :model-value="longitude"
          :disabled="disabled"
          :step="0.000001"
          :max-fraction-digits="7"
          :placeholder="$t('field.kind.geometry.longitude')"
          class="flex-1 min-w-0"
          @update:model-value="onLngChange"
        />
      </InputGroup>
      <Button
        v-tooltip.top="$t('field.kind.geometry.tooltip.openMap')"
        icon="pi pi-map"
        text
        severity="primary"
        class="flex-shrink-0"
        @click="showMapDialog = true"
      />
      <Button
        v-if="!field.options?.hideCurrentLocationButton"
        v-tooltip.top="$t('field.kind.geometry.tooltip.useCurrentLocation')"
        icon="pi pi-compass"
        text
        severity="primary"
        class="flex-shrink-0"
        @click="useCurrentLocation"
      />
    </div>

    <!-- Map dialog -->
    <Dialog
      v-model:visible="showMapDialog"
      :header="field.label || field.name"
      modal
      :style="{ width: '80vw', height: '80vh' }"
      :pt="{
        root: 'flex flex-col',
        content: 'flex flex-col flex-1 min-h-0 p-0',
        footer: 'p-3',
      }"
      @after-show="mapComponentRef?.invalidateSize()"
    >
      <CMap
        ref="mapComponentRef"
        :center="mapCenter"
        :zoom="mapZoom"
        :markers="mapMarkers"
        :max-bounds="lockedBounds"
        :hide-geo-search="field.options?.hideGeoSearch"
        :hide-current-location-button="field.options?.hideCurrentLocationButton"
        class="flex-1 min-h-0 w-full"
        @map-click="onMapClick"
        @location-found="onMapClick"
      />

      <template #footer>
        <div class="flex items-center text-muted-color text-sm mr-auto">
          {{ $t('field.kind.geometry.clickToPlaceMarker') }}
        </div>
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="showMapDialog = false"
        />
        <Button
          :label="$t('general.label.save')"
          severity="primary"
          size="small"
          @click="saveMapValue"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CMap from '../../map/CMap.vue'

const { t: $t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const showMapDialog = ref(false)
const tempMapCoords = ref(null)
const mapComponentRef = ref(null)

// Parse the stored JSON value: { "coordinates": [lat, lng] }
const parsedValue = computed(() => {
  try {
    const v = JSON.parse(props.modelValue || '{}')
    return v?.coordinates || []
  } catch {
    return []
  }
})

const latitude = computed(() => parsedValue.value[0] ?? null)
const longitude = computed(() => parsedValue.value[1] ?? null)

function emitCoords(lat, lng) {
  if (lat == null && lng == null) {
    emit('update:modelValue', '')
    return
  }
  emit('update:modelValue', JSON.stringify({ coordinates: [lat, lng] }))
}

function onLatChange(val) {
  emitCoords(val, longitude.value)
}

function onLngChange(val) {
  emitCoords(latitude.value, val)
}

// Map state
const mapCenter = computed(() => {
  if (tempMapCoords.value) return tempMapCoords.value
  if (latitude.value != null && longitude.value != null) {
    return [latitude.value, longitude.value]
  }
  return props.field.options?.center || [30, 30]
})

const mapZoom = computed(() => {
  if (latitude.value != null && longitude.value != null) return 13
  return props.field.options?.zoom || 3
})

const lockedBounds = computed(() => {
  if (!props.field.options?.lockBounds) return null
  const b = props.field.options?.bounds
  if (Array.isArray(b) && b.length === 2 && b.every(p => Array.isArray(p) && p.length === 2)) {
    return b
  }
  return null
})

const mapMarkers = computed(() => {
  if (tempMapCoords.value) {
    return [{ value: tempMapCoords.value }]
  }
  if (latitude.value != null && longitude.value != null) {
    return [{ value: [latitude.value, longitude.value] }]
  }
  return []
})

function onMapClick(e) {
  const { lat, lng } = e.latlng || {}
  if (lat != null && lng != null) {
    const roundedLat = Math.round(lat * 1e7) / 1e7
    const roundedLng = Math.round(lng * 1e7) / 1e7
    tempMapCoords.value = [roundedLat, roundedLng]
  }
}

function saveMapValue() {
  if (tempMapCoords.value) {
    emitCoords(tempMapCoords.value[0], tempMapCoords.value[1])
    tempMapCoords.value = null
  }
  showMapDialog.value = false
}

function useCurrentLocation() {
  if (!navigator.geolocation) return

  navigator.geolocation.getCurrentPosition(
    ({ coords }) => {
      const lat = Math.round(coords.latitude * 1e7) / 1e7
      const lng = Math.round(coords.longitude * 1e7) / 1e7
      emitCoords(lat, lng)
    },
    () => {
      // Geolocation error
    },
  )
}

// Prefill with current location for new records
onMounted(() => {
  if (props.field.options?.prefillWithCurrentLocation && !props.modelValue) {
    useCurrentLocation()
  }
})
</script>

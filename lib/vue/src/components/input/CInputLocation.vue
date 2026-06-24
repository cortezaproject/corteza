<template>
  <div class="flex flex-col gap-2">
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
        <InputGroupAddon v-if="!disabled && hasValue">
          <Button
            v-tooltip.top="$t('general.label.clear')"
            icon="pi pi-times"
            text
            severity="secondary"
            class="!p-2"
            @click="clearValue"
          />
        </InputGroupAddon>
      </InputGroup>
      <Button
        v-tooltip.top="$t('field.kind.geometry.tooltip.openMap')"
        icon="pi pi-map"
        text
        severity="primary"
        :disabled="disabled"
        class="flex-shrink-0"
        @click="openMap"
      />
      <Button
        v-if="!hideCurrentLocationButton"
        v-tooltip.top="$t('field.kind.geometry.tooltip.useCurrentLocation')"
        icon="pi pi-map-marker"
        text
        severity="primary"
        :disabled="disabled"
        class="flex-shrink-0"
        @click="useCurrentLocation"
      />
    </div>

    <Dialog
      v-model:visible="showMapDialog"
      :header="dialogHeader"
      modal
      :style="{ width: '80vw', height: '80vh' }"
      :pt="{
        root: 'flex flex-col',
        content: 'flex flex-col flex-1 min-h-0 p-0',
        footer: 'p-3',
      }"
      @after-show="mapComponentRef?.invalidateSize()"
      @hide="cancelMap"
    >
      <CMap
        ref="mapComponentRef"
        :center="dialogCenter"
        :zoom="dialogZoom"
        :markers="dialogMarkers"
        :max-bounds="bounds || null"
        :hide-geo-search="hideGeoSearch"
        :hide-current-location-button="hideCurrentLocationButton"
        class="flex-1 min-h-0 w-full"
        @map-click="onMapClick"
        @marker-click="onMarkerClick"
        @location-found="onLocationFound"
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
          @click="cancelMap"
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
import CMap from '../map/CMap.vue'

const { t: $t } = useI18n()

const props = defineProps({
  // GeoJSON Point — { type: 'Point', coordinates: [lng, lat] } — or null
  modelValue: { type: Object, default: null },
  disabled: { type: Boolean, default: false },
  hideGeoSearch: { type: Boolean, default: false },
  hideCurrentLocationButton: { type: Boolean, default: false },
  prefillWithCurrentLocation: { type: Boolean, default: false },
  // initialCenter is [lat, lng] (Leaflet convention) for the map dialog
  initialCenter: { type: Array, default: () => [30, 30] },
  initialZoom: { type: Number, default: 3 },
  bounds: { type: Array, default: null },
  dialogHeader: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const showMapDialog = ref(false)
const mapComponentRef = ref(null)

// Dialog-local working copy: GeoJSON coordinates [lng, lat] or null
const draftCoords = ref(null)

// View (center/zoom) is snapshotted on dialog open and NEVER updated from the
// draft afterwards. If we let it follow draftCoords, removing a marker would
// fall back to initialCenter/initialZoom and visibly jump/zoom the map.
const dialogCenter = ref([30, 30])
const dialogZoom = ref(3)

const longitude = computed(() => props.modelValue?.coordinates?.[0] ?? null)
const latitude = computed(() => props.modelValue?.coordinates?.[1] ?? null)
const hasValue = computed(() => latitude.value != null && longitude.value != null)

function round7(n) {
  return Math.round(n * 1e7) / 1e7
}

function buildPoint(lng, lat) {
  if (lng == null || lat == null || Number.isNaN(lng) || Number.isNaN(lat)) {
    return null
  }
  return { type: 'Point', coordinates: [round7(lng), round7(lat)] }
}

function emitValue(lng, lat) {
  emit('update:modelValue', buildPoint(lng, lat))
}

function onLatChange(val) {
  emitValue(longitude.value, val)
}

function onLngChange(val) {
  emitValue(val, latitude.value)
}

function clearValue() {
  emit('update:modelValue', null)
}

function openMap() {
  draftCoords.value = props.modelValue?.coordinates ? [...props.modelValue.coordinates] : null
  // Snapshot view once — based on the value at open time, not the live draft.
  if (props.modelValue?.coordinates) {
    const [lng, lat] = props.modelValue.coordinates
    dialogCenter.value = [lat, lng]
    dialogZoom.value = 13
  } else {
    dialogCenter.value = props.initialCenter
    dialogZoom.value = props.initialZoom
  }
  showMapDialog.value = true
}

const dialogMarkers = computed(() => {
  if (!draftCoords.value) return []
  const [lng, lat] = draftCoords.value
  // No `color` — falls through to L.marker so we get leaflet's default
  // teardrop pin (proper marker that stays visually anchored at any zoom).
  return [{ value: [lat, lng] }]
})

function onMapClick(e) {
  const { lat, lng } = e.latlng || {}
  if (lat == null || lng == null) return
  draftCoords.value = [round7(lng), round7(lat)]
}

function onLocationFound(e) {
  onMapClick(e)
}

function onMarkerClick() {
  draftCoords.value = null
}

function saveMapValue() {
  if (draftCoords.value) {
    const [lng, lat] = draftCoords.value
    emitValue(lng, lat)
  } else {
    emit('update:modelValue', null)
  }
  showMapDialog.value = false
  draftCoords.value = null
}

function cancelMap() {
  draftCoords.value = null
  showMapDialog.value = false
}

function useCurrentLocation() {
  if (!navigator.geolocation) return

  navigator.geolocation.getCurrentPosition(
    ({ coords }) => emitValue(coords.longitude, coords.latitude),
    () => {
      // silent fail
    },
  )
}

onMounted(() => {
  if (props.prefillWithCurrentLocation && !hasValue.value) {
    useCurrentLocation()
  }
})
</script>

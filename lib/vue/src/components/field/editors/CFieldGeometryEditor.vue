<template>
  <CInputLocation
    :model-value="geometryPoint"
    :disabled="disabled"
    :hide-geo-search="!!field.options?.hideGeoSearch"
    :hide-current-location-button="!!field.options?.hideCurrentLocationButton"
    :prefill-with-current-location="!!field.options?.prefillWithCurrentLocation && !modelValue"
    :initial-center="field.options?.center || [30, 30]"
    :initial-zoom="field.options?.zoom || 3"
    :bounds="lockedBounds"
    :dialog-header="field.label || field.name"
    @update:model-value="onPointUpdate"
  />
</template>

<script setup>
import { computed } from 'vue'
import CInputLocation from '../../input/CInputLocation.vue'

const props = defineProps({
  field: { type: Object, required: true },
  modelValue: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

// Legacy storage shape is JSON: { coordinates: [lat, lng] }
// CInputLocation speaks GeoJSON: { type: 'Point', coordinates: [lng, lat] }
const geometryPoint = computed(() => {
  try {
    const parsed = JSON.parse(props.modelValue || '{}')
    const coords = parsed?.coordinates
    if (!Array.isArray(coords) || coords.length !== 2) return null
    const [lat, lng] = coords
    if (typeof lat !== 'number' || typeof lng !== 'number') return null
    return { type: 'Point', coordinates: [lng, lat] }
  } catch {
    return null
  }
})

const lockedBounds = computed(() => {
  if (!props.field.options?.lockBounds) return null
  const b = props.field.options?.bounds
  if (Array.isArray(b) && b.length === 2 && b.every(p => Array.isArray(p) && p.length === 2)) {
    return b
  }
  return null
})

function onPointUpdate(point) {
  if (!point?.coordinates) {
    emit('update:modelValue', '')
    return
  }
  const [lng, lat] = point.coordinates
  emit('update:modelValue', JSON.stringify({ coordinates: [lat, lng] }))
}
</script>

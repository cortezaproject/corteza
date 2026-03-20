<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <Checkbox
        v-model="field.options.prefillWithCurrentLocation"
        inputId="prefillLocation"
        :binary="true"
      />
      <label for="prefillLocation" class="cursor-pointer">
        {{ $t('field.kind.geometry.prefillWithCurrentLocation') }}
      </label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox
        v-model="field.options.hideCurrentLocationButton"
        inputId="hideLocationBtn"
        :binary="true"
      />
      <label for="hideLocationBtn" class="cursor-pointer">
        {{ $t('field.kind.geometry.hideCurrentLocationButton') }}
      </label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox
        v-model="field.options.hideGeoSearch"
        inputId="hideGeoSearch"
        :binary="true"
      />
      <label for="hideGeoSearch" class="cursor-pointer">
        {{ $t('field.kind.geometry.hideGeoSearch') }}
      </label>
    </div>

    <!-- Map preview for setting initial zoom and position -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-muted-color text-sm">
        {{ $t('field.kind.geometry.initialZoomAndPosition') }}
      </label>
      <CMap
        :center="center"
        :zoom="zoom"
        hide-geo-search
        hide-current-location-button
        style="height: 50vh;"
        @map-click="onMapUpdate"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted } from 'vue'
import CMap from '@cortezaproject/corteza-vue-next/src/components/map/CMap.vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const center = computed(() => props.field.options?.center || [30, 30])
const zoom = computed(() => props.field.options?.zoom || 3)

function onMapUpdate(e) {
  // When loading the map, we can listen to zoom/center updates
  // but CMap doesn't emit separate zoom/center events yet,
  // so for now the map just lets users visualize the default position.
}

onMounted(() => {
  if (!props.field.options.center) {
    props.field.options.center = [30, 30]
  }
  if (!props.field.options.zoom) {
    props.field.options.zoom = 3
  }
})
</script>
